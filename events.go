package qcloud_im_callback

import (
	"encoding/json"
	"errors"
)

// 事件实体
type CallbackEvent struct {
        // 事件标识命令
	CallbackCommand CallbackCommand

	// 参数
        URLParams *URLParams

	// 实体内容
	Body []byte

	// 句柄
	Handler *CallbackHandler
}


// 新建
func NewCallbackEvent(cc CallbackCommand,up *URLParams,body []byte )*CallbackEvent{
	return &CallbackEvent{CallbackCommand:cc,URLParams:up,Body:body}
}

// 将body字节内容以JSON格式转化
func (ce *CallbackEvent) ToJSON(v interface{}) error {
	if ce.Body == nil {
		return errors.New("body is nil")
	}
	return json.Unmarshal(ce.Body, v)
}

// @see CallbackHandler.Handle()
func (ce *CallbackEvent) Handle() interface{} {
	if nil!=ce.Handler {
		return ce.Handler.Handle(ce)
	}
	return nil
}

// @see CallbackHandler.handle()
func (ce *CallbackEvent) handle() {
	if nil!=ce.Handler {
		ce.Handler.handle(ce)
	}
}