package event

import (
	"context"

	"github.com/lin-snow/ech0/internal/async"
	"github.com/lin-snow/ech0/internal/fediverse"
)

// FediverseAgent 处理联邦相关的事件
type FediverseAgent struct {
	pool *async.WorkerPool // 任务池
	core *fediverse.FediverseCore
}

func NewFediverseAgent(core *fediverse.FediverseCore) *FediverseAgent {
	return &FediverseAgent{
		pool: async.NewWorkerPool(3, 3), // 假设最大并发数为 3，任务队列大小为 3
		core: core,
	}
}

func (fa *FediverseAgent) Handle(ctx context.Context, e *Event) error {
	// 处理事件，与联邦宇宙交互
	switch e.Type {
	case EventTypeEchoCreated:
		fa.HandleCreateEchoEvent(ctx, e)

	default:
		return nil // 忽略其他事件
	}

	return nil
}

func (fa *FediverseAgent) Wait() {
	fa.pool.Wait()
}

func (fa *FediverseAgent) HandleCreateEchoEvent(ctx context.Context, e *Event) error {
	// 将 Echo 推送到联邦宇宙

	return nil
}
