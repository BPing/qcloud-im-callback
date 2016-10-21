package qcloud_im_callback

import (
	"sync"

	"github.com/BPing/Golib/producer_consumer"
)

// 具体事件处理程序
type CallbackHandle func(*CallbackEvent) interface{}

//
// 回调事件处理句柄
//   为某个的事件注册相应的事件处理程序
//
type CallbackHandler struct {
	// 事件处理路由，CallbackCommand对应的的处理方式
	router map[CallbackCommand]RouterInfo

	// 默认的处理程序
	// 当Router里面没有注册的理由存在时候，
	// 将默认使用本程序处理事件
	defaultHandle CallbackHandle

	// 开始处理事件之前的钩子
	//  如果返回数据不为nil，代表结束处理事件，
	//  否则，继续处理事件
	beforeHook CallbackHandle

	//生产/消费 消费异步事件
	producerConsumer *producerConsumer.Container

	// CallbackEvent对象池，避免创建过多对象
	eventPool *sync.Pool
}

// 事件处理信息和程序
type RouterInfo struct {
	// 异步或者同步
	Async bool

	// 如果是异步处理，
	// 默认的返回的数据。
	// 一般为BaseResponse结构体即可
	AsyncResponse interface{}

	// 处理句柄
	//   事件具体处理的程序。
	//   如果同步处理，返回的数据将返回到客户端去，如果异步的话，将会忽略
	Handle CallbackHandle
}

// 新建回调事件处理句柄
func NewCallbackHandler(masterNum, msgEventLen int, defaultHandle CallbackHandle) (*CallbackHandler, error) {
	ch := &CallbackHandler{
		router:        make(map[CallbackCommand]RouterInfo),
		defaultHandle: defaultHandle,
		eventPool: &sync.Pool{New: func() interface{} {
			return NewCallbackEvent("", nil, nil)
		}}}
	err := ch.InitProducerConsumer(masterNum, msgEventLen)
	return ch, err
}

//
// @masterNum 主消费线程数目,必须大于等于1
// @chanLen   消费信息（事件）队列长度
func (ch *CallbackHandler) InitProducerConsumer(masterNum, msgEventLen int) error {
	pc, err := producerConsumer.NewContainerPC(msgEventLen, func(msg producerConsumer.IMessage) {
		// 处理异步延后处理消息（事件）
		event ,ok:=msg.(*CallbackEvent)
		if ok {
			event.handle()
		}

	})
	if err != nil {
		return err
	}

	ch.producerConsumer = pc

	if masterNum < 1 {
		masterNum = 1
	}

	for i := 0; i < masterNum; i++ {
		ch.producerConsumer.Consume()
	}
	return nil
}

// 注册
//   如果重复注册，新的将覆盖旧的
func (ch *CallbackHandler) Register(cc CallbackCommand, ri RouterInfo) *CallbackHandler {
	ch.router[cc] = ri
	return ch
}

// 注册默认处理程序
func (ch *CallbackHandler) RegisterDefaultHandle(callbackHandle CallbackHandle) *CallbackHandler {
	ch.defaultHandle = callbackHandle
	return ch
}

// 注册钩子
func (ch *CallbackHandler) RegisterBeforeHook(beforeHook CallbackHandle) *CallbackHandler {
	ch.beforeHook = beforeHook
	return ch
}

// 注销事件处理路由信息
func (ch *CallbackHandler) UnRegister(cc CallbackCommand) *CallbackHandler {
	delete(ch.router, cc)
	return ch
}

// 是否已注册
func (ch *CallbackHandler) Exist(cc CallbackCommand) bool {
	_, ok := ch.router[cc]
	return ok
}

// 获取事件处理路由信息
func (ch *CallbackHandler) Get(cc CallbackCommand) (RouterInfo, bool) {
	ri, ok := ch.router[cc]
	return ri, ok
}

// 事件队列处理协程数目情况
func (ch *CallbackHandler) ConsumerNumGoroutine() (master, assistActive int64) {
	if nil != ch.producerConsumer {
		master, assistActive = ch.producerConsumer.NumGoroutine()
	}
	return
}

// 处理事件
func (ch *CallbackHandler) Handle(ce *CallbackEvent) interface{} {
	if ch.beforeHook != nil {
		hr := ch.beforeHook(ce)
		if nil != hr {
			// 如果钩子有返回数据则代表结束
			return hr
		}
	}

	ri, ok := ch.Get(ce.CallbackCommand)
	if ok {
		if ri.Async {
			if nil != ch.producerConsumer {
				// 放进消费队列延后处理
				ch.producerConsumer.Produce(ce)
			}
			if nil != ri.AsyncResponse {
				return ri.AsyncResponse
			}
			return ch.defaultHandle(ce)
		} else {
			defer ch.eventPool.Put(ce)
			resp := ri.Handle(ce)
			return resp
		}

	}
	defer ch.eventPool.Put(ce)
	return ch.defaultHandle(ce)
}

// producerConsumer消费事件时调用
func (ch *CallbackHandler) handle(ce *CallbackEvent) {
	defer ch.eventPool.Put(ce)
	ri, ok := ch.Get(ce.CallbackCommand)
	if ok && ri.Async {
		ri.Handle(ce)
	}
}

// 新建事件
func (ch *CallbackHandler) NewCallbackEvent(cc CallbackCommand, up *URLParams, body []byte) *CallbackEvent {
	ce := ch.eventPool.Get().(*CallbackEvent)
	ce.CallbackCommand = cc
	ce.URLParams = up
	ce.Body = body
	ce.Handler = ch
	return ce
}
