package qcloud_im_callback

import (
	"fmt"
	"testing"
	"time"
	"container/list"
	"errors"
)

func TestHandler(t *testing.T) {
	defaultHandle := func(event *CallbackEvent) interface{} {
		return &BaseResponse{ActionStatus: OkStatus, ErrorCode: 0}
	}
	callbackHandler, _ := NewCallbackHandler(0, 30, CallbackHandle(defaultHandle))

	// 注册
	CallbackBeforeSendMsgHandle := func(event *CallbackEvent) interface{} {
		var sendMsgBody SendMsgBody
		event.ToJSON(&sendMsgBody)
		fmt.Println("CallbackBeforeSendMsgHandle", sendMsgBody.MsgBody)
		return &BaseResponse{ActionStatus: OkStatus, ErrorCode: 0}
	}
	callbackHandler.Register(CallbackBeforeSendMsgCommand, RouterInfo{
		Async:  false,
		Handle: CallbackHandle(CallbackBeforeSendMsgHandle),
	})

	CallbackAfterSendMsgHandle := func(event *CallbackEvent) interface{} {
		fmt.Println(CallbackAfterSendMsgCommand)
		return &BaseResponse{ActionStatus: OkStatus, ErrorCode: 0}
	}
	callbackHandler.Register(CallbackAfterSendMsgCommand, RouterInfo{
		Async:         true,
		AsyncResponse: &BaseResponse{ActionStatus: FAILStatus, ErrorCode: 0},
		Handle:        CallbackHandle(CallbackAfterSendMsgHandle),
	})

	up := &URLParams{
		SdkAppid:    "888888",
		ContentType: "json",
	}

	str := `{"CallbackCommand": "C2C.CallbackBeforeSendMsg", "From_Account": "jared", "To_Account": "Jonh", "MsgBody": [ {"MsgType": "TIMTextElem","MsgContent": {"Text": "red packet"}}]}`

	resp, ok := callbackHandler.NewCallbackEvent(CallbackBeforeSendMsgCommand, up, []byte(str)).Handle().(*BaseResponse)
	if !ok || resp.ActionStatus != OkStatus {
		t.Fatal(ok, CallbackBeforeSendMsgCommand, "handle fail")
	}

	resp, ok = callbackHandler.NewCallbackEvent(CallbackAfterSendMsgCommand, up, nil).Handle().(*BaseResponse)
	if !ok || resp.ActionStatus != FAILStatus {
		t.Fatal(ok, CallbackAfterSendMsgCommand, "handle fail")
	}

	// test AsyncResponse  is nil
	callbackHandler.Register(CallbackAfterSendMsgCommand, RouterInfo{
		Async:         true,
		AsyncResponse: nil,
		Handle:        CallbackHandle(CallbackAfterSendMsgHandle),
	})
	resp, ok = callbackHandler.NewCallbackEvent(CallbackAfterSendMsgCommand, up, nil).Handle().(*BaseResponse)
	if !ok || resp.ActionStatus != OkStatus {
		t.Fatal(ok, CallbackAfterSendMsgCommand, "test AsyncResponse  is nil fail")
	}

	// 没有注册
	resp, ok = callbackHandler.NewCallbackEvent(CallbackBeforeCreateGroupCommand, up, nil).Handle().(*BaseResponse)
	if !ok || resp.ActionStatus != OkStatus {
		t.Fatal(ok, CallbackBeforeCreateGroupCommand, "handle fail")
	}

	//RegisterDefaultHandle

	callbackHandler.RegisterDefaultHandle(CallbackHandle(func(event *CallbackEvent) interface{} {
		return &BaseResponse{ActionStatus: FAILStatus, ErrorCode: 0}
	}))

	resp, ok = callbackHandler.NewCallbackEvent(CallbackBeforeCreateGroupCommand, up, nil).Handle().(*BaseResponse)
	if !ok || resp.ActionStatus != FAILStatus {
		t.Fatal(ok, CallbackBeforeCreateGroupCommand, "handle fail")
	}

	//Exist and   UnRegister

	if !callbackHandler.Exist(CallbackBeforeSendMsgCommand) {
		t.Fatal(ok, CallbackBeforeSendMsgCommand, "Exist  fail")
	}

	callbackHandler.UnRegister(CallbackBeforeSendMsgCommand)

	if callbackHandler.Exist(CallbackBeforeSendMsgCommand) {
		t.Fatal(ok, CallbackBeforeSendMsgCommand, "UnRegister Exist  fail")
	}

	resp, ok = callbackHandler.NewCallbackEvent(CallbackBeforeSendMsgCommand, up, []byte(str)).Handle().(*BaseResponse)
	if !ok || resp.ActionStatus != FAILStatus {
		t.Fatal(ok, CallbackBeforeSendMsgCommand, "handle fail")
	}

	_, err := NewCallbackHandler(0, -1, CallbackHandle(defaultHandle))

	if nil == err {
		t.Fatal("test NewCallbackHandler fail")
	}

	beforeHook := func(event *CallbackEvent) interface{} {
		return &BaseResponse{ActionStatus: OkStatus, ErrorCode: -1}
	}
	callbackHandler.RegisterBeforeHook(CallbackHandle(beforeHook))

	resp, ok = callbackHandler.NewCallbackEvent(CallbackBeforeSendMsgCommand, up, []byte(str)).Handle().(*BaseResponse)
	if !ok || resp.ErrorCode != -1 {
		t.Fatal(ok, CallbackBeforeSendMsgCommand, "handle fail")
	}

	beforeHook = func(event *CallbackEvent) interface{} {
		return nil
	}
	callbackHandler.RegisterBeforeHook(CallbackHandle(beforeHook))

	resp, ok = callbackHandler.NewCallbackEvent(CallbackBeforeSendMsgCommand, up, []byte(str)).Handle().(*BaseResponse)
	if !ok || resp.ActionStatus != FAILStatus {
		t.Fatal(ok, CallbackBeforeSendMsgCommand, "handle fail")
	}

	fmt.Println(callbackHandler.ConsumerNumGoroutine())

	time.Sleep(time.Second)
}

func TestCacheHandler(t *testing.T){
	cacheR:=&TestRedis{list.New()}

	defaultHandle := func(event *CallbackEvent) interface{} {
		return &BaseResponse{ActionStatus: OkStatus, ErrorCode: 0}
	}
	callbackHandler, _ := NewCallbackHandlerWithCache(0, 30, CallbackHandle(defaultHandle),cacheR)

	// 注册
	CallbackBeforeSendMsgHandle := func(event *CallbackEvent) interface{} {
		var sendMsgBody SendMsgBody
		event.ToJSON(&sendMsgBody)
		fmt.Println("CallbackBeforeSendMsgHandle", sendMsgBody.MsgBody)
		return &BaseResponse{ActionStatus: OkStatus, ErrorCode: 0}
	}
	callbackHandler.Register(CallbackBeforeSendMsgCommand, RouterInfo{
		Async:  false,
		Handle: CallbackHandle(CallbackBeforeSendMsgHandle),
	})

	CallbackAfterSendMsgHandle := func(event *CallbackEvent) interface{} {
		fmt.Println(CallbackAfterSendMsgCommand)
		return &BaseResponse{ActionStatus: OkStatus, ErrorCode: 0}
	}
	callbackHandler.Register(CallbackAfterSendMsgCommand, RouterInfo{
		Async:         true,
		AsyncResponse: &BaseResponse{ActionStatus: FAILStatus, ErrorCode: 0},
		Handle:        CallbackHandle(CallbackAfterSendMsgHandle),
	})

	up := &URLParams{
		SdkAppid:    "888888",
		ContentType: "json",
	}

	str := `{"CallbackCommand": "C2C.CallbackBeforeSendMsg", "From_Account": "jared", "To_Account": "Jonh", "MsgBody": [ {"MsgType": "TIMTextElem","MsgContent": {"Text": "red packet"}}]}`

	resp, ok := callbackHandler.NewCallbackEvent(CallbackBeforeSendMsgCommand, up, []byte(str)).Handle().(*BaseResponse)
	if !ok || resp.ActionStatus != OkStatus {
		t.Fatal(ok, CallbackBeforeSendMsgCommand, "Cache: handle fail")
	}

	resp, ok = callbackHandler.NewCallbackEvent(CallbackAfterSendMsgCommand, up, nil).Handle().(*BaseResponse)
	if !ok || resp.ActionStatus != FAILStatus {
		t.Fatal(ok, CallbackAfterSendMsgCommand, "Cache: handle fail")
	}

	// test AsyncResponse  is nil
	callbackHandler.Register(CallbackAfterSendMsgCommand, RouterInfo{
		Async:         true,
		AsyncResponse: nil,
		Handle:        CallbackHandle(CallbackAfterSendMsgHandle),
	})
	resp, ok = callbackHandler.NewCallbackEvent(CallbackAfterSendMsgCommand, up, nil).Handle().(*BaseResponse)
	if !ok || resp.ActionStatus != OkStatus {
		t.Fatal(ok, CallbackAfterSendMsgCommand, "Cache: test AsyncResponse  is nil fail")
	}

	// 没有注册
	resp, ok = callbackHandler.NewCallbackEvent(CallbackBeforeCreateGroupCommand, up, nil).Handle().(*BaseResponse)
	if !ok || resp.ActionStatus != OkStatus {
		t.Fatal(ok, CallbackBeforeCreateGroupCommand, "Cache: handle fail")
	}

	//RegisterDefaultHandle

	callbackHandler.RegisterDefaultHandle(CallbackHandle(func(event *CallbackEvent) interface{} {
		return &BaseResponse{ActionStatus: FAILStatus, ErrorCode: 0}
	}))

	resp, ok = callbackHandler.NewCallbackEvent(CallbackBeforeCreateGroupCommand, up, nil).Handle().(*BaseResponse)
	if !ok || resp.ActionStatus != FAILStatus {
		t.Fatal(ok, CallbackBeforeCreateGroupCommand, "Cache: handle fail")
	}

	//Exist and   UnRegister

	if !callbackHandler.Exist(CallbackBeforeSendMsgCommand) {
		t.Fatal(ok, CallbackBeforeSendMsgCommand, "Cache: Exist  fail")
	}

	callbackHandler.UnRegister(CallbackBeforeSendMsgCommand)

	if callbackHandler.Exist(CallbackBeforeSendMsgCommand) {
		t.Fatal(ok, CallbackBeforeSendMsgCommand, "Cache: UnRegister Exist  fail")
	}

	resp, ok = callbackHandler.NewCallbackEvent(CallbackBeforeSendMsgCommand, up, []byte(str)).Handle().(*BaseResponse)
	if !ok || resp.ActionStatus != FAILStatus {
		t.Fatal(ok, CallbackBeforeSendMsgCommand, "Cache: handle fail")
	}

	_, err := NewCallbackHandler(0, -1, CallbackHandle(defaultHandle))

	if nil == err {
		t.Fatal("Cache: test NewCallbackHandler fail")
	}

	beforeHook := func(event *CallbackEvent) interface{} {
		return &BaseResponse{ActionStatus: OkStatus, ErrorCode: -1}
	}
	callbackHandler.RegisterBeforeHook(CallbackHandle(beforeHook))

	resp, ok = callbackHandler.NewCallbackEvent(CallbackBeforeSendMsgCommand, up, []byte(str)).Handle().(*BaseResponse)
	if !ok || resp.ErrorCode != -1 {
		t.Fatal(ok, CallbackBeforeSendMsgCommand, "Cache: handle fail")
	}

	beforeHook = func(event *CallbackEvent) interface{} {
		return nil
	}
	callbackHandler.RegisterBeforeHook(CallbackHandle(beforeHook))

	resp, ok = callbackHandler.NewCallbackEvent(CallbackBeforeSendMsgCommand, up, []byte(str)).Handle().(*BaseResponse)
	if !ok || resp.ActionStatus != FAILStatus {
		t.Fatal(ok, CallbackBeforeSendMsgCommand, "Cache: handle fail")
	}

	fmt.Println(callbackHandler.ConsumerNumGoroutine())

	time.Sleep(time.Second)

}

type TestRedis struct {
	msgList *list.List
}

func (tr *TestRedis) BLPop(key string, timeout int64) (map[string]string, error) {
	//fmt.Println("BLPop-----------------------")
	resp := make(map[string]string)
	ele := tr.msgList.Front()
	if ele != nil {
		tr.msgList.Remove(ele)
		resp[key] = ele.Value.(string)
		return resp, nil
	}else{
		time.Sleep(time.Second * time.Duration(timeout))
	}

	return nil, errors.New("list is empty")
}

func (tr *TestRedis) RPush(key string, values ...interface{}) (int64, error) {
	//fmt.Println("RPush-----------------------")
	switch value := values[0].(type) {
	case []byte:
		tr.msgList.PushBack(string(value))
	case string:
		tr.msgList.PushBack(value)

	default:
		return 0, errors.New("type error ([]byte or string)")
	}
	return 0, nil
}

func (tr *TestRedis) LLen(key string) (int64, error) {
	//fmt.Println("LLen-----------------------")
	return int64(tr.msgList.Len()), nil
}