package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-agent/internal/worker/base"
	"go-agent/internal/worker/order"
)

// BackgroundWorker 后台工作进程
type BackgroundWorker struct {
	manager   *base.TaskManager
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	taskQueue chan *base.Task
	isRunning bool
}

// NewBackgroundWorker 创建新的后台工作进程
func NewBackgroundWorker(maxWorkers, queueSize int) *BackgroundWorker {
	ctx, cancel := context.WithCancel(context.Background())

	return &BackgroundWorker{
		manager:   base.NewTaskManager(maxWorkers, queueSize),
		ctx:       ctx,
		cancel:    cancel,
		taskQueue: make(chan *base.Task, queueSize),
		isRunning: false,
	}
}

// Start 启动后台工作进程
func (bw *BackgroundWorker) Start() {
	if bw.isRunning {
		return
	}

	bw.isRunning = true

	// 启动任务处理协程
	bw.wg.Add(1)
	go bw.processTasks()

	fmt.Println("后台工作进程已启动")
}

// Stop 停止后台工作进程
func (bw *BackgroundWorker) Stop() {
	if !bw.isRunning {
		return
	}

	bw.isRunning = false
	bw.cancel()
	close(bw.taskQueue)
	bw.wg.Wait()
	bw.manager.Shutdown()

	fmt.Println("后台工作进程已停止")
}

// SubmitTask 提交任务到后台工作进程
func (bw *BackgroundWorker) SubmitTask(task *base.Task) error {
	if !bw.isRunning {
		return fmt.Errorf("后台工作进程未运行")
	}

	select {
	case bw.taskQueue <- task:
		return nil
	case <-bw.ctx.Done():
		return fmt.Errorf("后台工作进程已停止")
	default:
		return fmt.Errorf("任务队列已满")
	}
}

// processTasks 处理任务队列中的任务
func (bw *BackgroundWorker) processTasks() {
	defer bw.wg.Done()

	for {
		select {
		case task, ok := <-bw.taskQueue:
			if !ok {
				// 通道已关闭
				return
			}

			// 提交任务到任务管理器
			if err := bw.manager.Submit(task); err != nil {
				fmt.Printf("提交任务失败: %v\n", err)
			}

		case <-bw.ctx.Done():
			// 上下文已取消
			return
		}
	}
}

// SubmitSampleTasks 提交示例任务
func (bw *BackgroundWorker) SubmitSampleTasks() {
	// 提交一些示例任务
	tasks := []*base.Task{
		order.ProcessOrderTask("bg-001"),
		order.ProcessOrderTask("bg-002"),
		order.ProcessRefundTask("bg-refund-001"),
	}

	for _, task := range tasks {
		if err := bw.SubmitTask(task); err != nil {
			fmt.Printf("提交示例任务失败 %s: %v\n", task.ID, err)
		} else {
			fmt.Printf("成功提交示例任务: %s\n", task.ID)
		}
	}

	// 定时提交更多任务
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		taskCounter := 3
		for {
			select {
			case <-ticker.C:
				taskCounter++
				task := order.ProcessOrderTask(fmt.Sprintf("bg-timed-%d", taskCounter))
				if err := bw.SubmitTask(task); err != nil {
					fmt.Printf("提交定时任务失败 %s: %v\n", task.ID, err)
				} else {
					fmt.Printf("成功提交定时任务: %s\n", task.ID)
				}
			case <-bw.ctx.Done():
				return
			}
		}
	}()
}
