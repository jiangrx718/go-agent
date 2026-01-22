package worker

import (
	"fmt"
	"strings"
	"time"

	"go-agent/internal/worker/base"
	"go-agent/internal/worker/order"

	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "worker",
		Usage: "工作进程",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "type",
				Usage:   "指定要运行的worker类型 (product/order), 不指定则运行所有worker",
				Aliases: []string{"t"},
			},
		},
		Action: func(ctx *cli.Context) error {
			// 创建任务管理器 (最大并发数: 10, 队列大小: 100)
			taskManager := base.NewTaskManager(10, 100)
			defer taskManager.Shutdown()

			workerType := strings.ToLower(ctx.String("t"))

			switch workerType {
			case "product":
				fmt.Printf("这里需要处理 %s 类型数据\n", workerType)
				// 示例：提交一些产品相关任务
				for i := 1; i <= 5; i++ {
					task := order.ProcessOrderTask(fmt.Sprintf("product-%d", i))
					if err := taskManager.Submit(task); err != nil {
						fmt.Printf("提交任务失败: %v\n", err)
					}
				}
			case "order":
				fmt.Printf("这里需要处理 %s 类型数据\n", workerType)
				// 示例：提交一些订单相关任务
				for i := 1; i <= 3; i++ {
					task := order.ProcessOrderTask(fmt.Sprintf("order-%d", i))
					if err := taskManager.Submit(task); err != nil {
						fmt.Printf("提交任务失败: %v\n", err)
					}
				}
			case "":
				fmt.Println("这里需要处理所有数据")
				// 示例：提交不同类型的任务
				for i := 1; i <= 3; i++ {
					orderTask := order.ProcessOrderTask(fmt.Sprintf("all-%d", i))
					refundTask := order.ProcessRefundTask(fmt.Sprintf("refund-%d", i))

					if err := taskManager.Submit(orderTask); err != nil {
						fmt.Printf("提交订单任务失败: %v\n", err)
					}

					if err := taskManager.Submit(refundTask); err != nil {
						fmt.Printf("提交退款任务失败: %v\n", err)
					}
				}
			default:
				return cli.Exit("无效的 worker 类型。可用类型: product/order", 1)
			}

			// 等待一段时间以确保任务完成
			time.Sleep(5 * time.Second)

			// 显示任务统计信息
			fmt.Println("\n任务统计信息:")
			stats := taskManager.GetAllStats()
			for taskID, stat := range stats {
				fmt.Printf("任务ID: %s, 总执行次数: %d, 成功: %d, 失败: %d, 最后执行: %v\n",
					taskID, stat.TotalExecutions, stat.SuccessCount, stat.ErrorCount, stat.LastExecution)
			}

			return nil
		},
	}
}
