package event

// EventHandlers 事件处理器集合
type EventHandlers struct {
	wbd *WebhookDispatcher  // webhook 事件处理器
	dlr *DeadLetterResolver // 死信处理器
	fa  *FediverseAgent     // 联邦事件处理器
}

// NewEventHandlers 创建一个新的事件处理器集合
func NewEventHandlers(wbd *WebhookDispatcher, dlr *DeadLetterResolver, fa *FediverseAgent) *EventHandlers {
	return &EventHandlers{wbd: wbd, dlr: dlr, fa: fa}
}

// EventRegistrar 事件注册器
type EventRegistrar struct {
	eb IEventBus      // 事件总线
	eh *EventHandlers // 事件处理器集合
}

// NewEventRegistry 创建一个新的事件注册表
func NewEventRegistry(ebp func() IEventBus, eh *EventHandlers) *EventRegistrar {
	return &EventRegistrar{eb: ebp(), eh: eh}
}

// Register 注册事件处理函数
func (er *EventRegistrar) Register() error {
	// 订阅死信事件
	er.eb.Subscribe(er.eh.dlr.Handle, EventTypeDeadLetterRetried) // 订阅死信事件，交给 DeadLetterResolver 处理
	er.eb.Subscribe(er.eh.fa.Handle, EventTypeEchoCreated)        // 订阅 EchoCreated 事件，交给 FediverseAgent 处理

	// 订阅所有事件，交给 WebhookDispatcher 处理
	er.eb.SubscribeAll(er.eh.wbd.Handle, EventTypeDeadLetterRetried) // 订阅所有事件，交给 WebhookDispatcher 处理,但是排除死信事件

	return nil
}

// Wait 等待所有事件处理完成
func (er *EventRegistrar) Wait() {
	er.eh.wbd.Wait()
	er.eh.fa.Wait()
}
