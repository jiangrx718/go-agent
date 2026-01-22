package order

import (
	"context"
	"fmt"
	"time"

	"go-agent/internal/worker/base"
)

// ProcessOrderTask 处理订单任务
func ProcessOrderTask(orderID string) *base.Task {
	return &base.Task{
		ID:   fmt.Sprintf("order_%s", orderID),
		Name: "处理订单任务",
		Function: func(ctx context.Context) error {
			// 模拟订单处理逻辑
			fmt.Printf("开始处理订单: %s\n", orderID)

			// 模拟一些处理时间
			time.Sleep(2 * time.Second)

			// 模拟处理成功
			fmt.Printf("订单处理完成: %s\n", orderID)

			return nil
		},
	}
}

// ProcessRefundTask 处理退款任务
func ProcessRefundTask(refundID string) *base.Task {
	return &base.Task{
		ID:   fmt.Sprintf("refund_%s", refundID),
		Name: "处理退款任务",
		Function: func(ctx context.Context) error {
			// 模拟退款处理逻辑
			fmt.Printf("开始处理退款: %s\n", refundID)

			// 模拟一些处理时间
			time.Sleep(3 * time.Second)

			// 模拟处理成功
			fmt.Printf("退款处理完成: %s\n", refundID)

			return nil
		},
	}
}
