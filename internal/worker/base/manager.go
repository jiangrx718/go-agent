package base

import (
	"go-agent/gopkg/log"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/alitto/pond/v2"
	"go.uber.org/zap"
)

// TaskFunc 任务函数类型
type TaskFunc func(context.Context) error

// Task 任务结构
type Task struct {
	ID       string
	Name     string
	Function TaskFunc
}

// TaskManager 任务管理器
type TaskManager struct {
	pool       pond.Pool
	ctx        context.Context
	cancel     context.CancelFunc
	logger     *zap.SugaredLogger
	taskStats  map[string]*TaskStats
	statsMutex sync.RWMutex
}

// TaskStats 任务统计信息
type TaskStats struct {
	TotalExecutions int64
	SuccessCount    int64
	ErrorCount      int64
	LastExecution   time.Time
}

// NewTaskManager 创建新的任务管理器
func NewTaskManager(maxWorkers int, queueSize int) *TaskManager {
	ctx, cancel := context.WithCancel(context.Background())

	manager := &TaskManager{
		pool:      pond.NewPool(maxWorkers, pond.WithQueueSize(queueSize)),
		ctx:       ctx,
		cancel:    cancel,
		logger:    log.Sugar(),
		taskStats: make(map[string]*TaskStats),
	}

	return manager
}

// Submit 提交任务到线程池
func (tm *TaskManager) Submit(task *Task) error {
	select {
	case <-tm.ctx.Done():
		return fmt.Errorf("task manager is shutdown")
	default:
		return tm.pool.Go(tm.wrapTask(task))
	}
}

// SubmitAsync 异步提交任务到线程池
func (tm *TaskManager) SubmitAsync(task *Task) (pond.Task, error) {
	select {
	case <-tm.ctx.Done():
		return nil, fmt.Errorf("task manager is shutdown")
	default:
		return tm.pool.Submit(tm.wrapTask(task)), nil
	}
}

// wrapTask 包装任务函数以添加统计和错误处理
func (tm *TaskManager) wrapTask(task *Task) func() {
	return func() {
		startTime := time.Now()
		tm.logger.Infof("开始执行任务: %s (ID: %s)", task.Name, task.ID)

		// 执行任务
		err := task.Function(tm.ctx)

		// 更新统计信息
		tm.updateStats(task.ID, err == nil)

		// 记录执行结果
		duration := time.Since(startTime)
		if err != nil {
			tm.logger.Errorf("任务执行失败: %s (ID: %s), 耗时: %v, 错误: %v", task.Name, task.ID, duration, err)
		} else {
			tm.logger.Infof("任务执行成功: %s (ID: %s), 耗时: %v", task.Name, task.ID, duration)
		}
	}
}

// updateStats 更新任务统计信息
func (tm *TaskManager) updateStats(taskID string, success bool) {
	tm.statsMutex.Lock()
	defer tm.statsMutex.Unlock()

	stats, exists := tm.taskStats[taskID]
	if !exists {
		stats = &TaskStats{}
		tm.taskStats[taskID] = stats
	}

	stats.TotalExecutions++
	stats.LastExecution = time.Now()

	if success {
		stats.SuccessCount++
	} else {
		stats.ErrorCount++
	}
}

// GetStats 获取任务统计信息
func (tm *TaskManager) GetStats(taskID string) *TaskStats {
	tm.statsMutex.RLock()
	defer tm.statsMutex.RUnlock()

	if stats, exists := tm.taskStats[taskID]; exists {
		// 返回副本以避免并发问题
		return &TaskStats{
			TotalExecutions: stats.TotalExecutions,
			SuccessCount:    stats.SuccessCount,
			ErrorCount:      stats.ErrorCount,
			LastExecution:   stats.LastExecution,
		}
	}

	return nil
}

// GetAllStats 获取所有任务统计信息
func (tm *TaskManager) GetAllStats() map[string]*TaskStats {
	tm.statsMutex.RLock()
	defer tm.statsMutex.RUnlock()

	result := make(map[string]*TaskStats)
	for id, stats := range tm.taskStats {
		result[id] = &TaskStats{
			TotalExecutions: stats.TotalExecutions,
			SuccessCount:    stats.SuccessCount,
			ErrorCount:      stats.ErrorCount,
			LastExecution:   stats.LastExecution,
		}
	}

	return result
}

// Shutdown 关闭任务管理器
func (tm *TaskManager) Shutdown() {
	tm.logger.Info("正在关闭任务管理器...")
	tm.cancel()
	// 注意：pond/v2 不需要显式停止池
	tm.logger.Info("任务管理器已关闭")
}

// IsClosed 检查任务管理器是否已关闭
func (tm *TaskManager) IsClosed() bool {
	select {
	case <-tm.ctx.Done():
		return true
	default:
		return false
	}
}
