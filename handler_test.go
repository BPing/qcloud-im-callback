package qcloud_im_callback

import (
	"testing"
	"fmt"
	"time"
)

func TestHandler(t *testing.T){
	defaultHandle:=func(event *CallbackEvent)interface{}{
               return &BaseResponse{ActionStatus:OkStatus,ErrorCode:0}
	}
	callbackHandler,_:=NewCallbackHandler(0,30,CallbackHandle(defaultHandle))

	// 注册
	CallbackBeforeSendMsgHandle:=func(event *CallbackEvent)interface{}{
		var sendMsgBody SendMsgBody
		event.ToJSON(&sendMsgBody)
		fmt.Println("CallbackBeforeSendMsgHandle",sendMsgBody.MsgBody)
		return &BaseResponse{ActionStatus:OkStatus,ErrorCode:0}
	}
	callbackHandler.Register(CallbackBeforeSendMsgCommand,RouterInfo{
		Async:false,
		Handle:CallbackHandle(CallbackBeforeSendMsgHandle),

	})

	CallbackAfterSendMsgHandle:=func(event *CallbackEvent)interface{}{
		fmt.Println(CallbackAfterSendMsgCommand)
		return &BaseResponse{ActionStatus:OkStatus,ErrorCode:0}
	}
	callbackHandler.Register(CallbackAfterSendMsgCommand,RouterInfo{
		Async:true,
		AsyncResponse:&BaseResponse{ActionStatus:FAILStatus,ErrorCode:0},
		Handle:CallbackHandle(CallbackAfterSendMsgHandle),

	})

	up:=URLParams{
		SdkAppid:"888888",
		ContentType:"json",
	}

	str:=`{"CallbackCommand": "C2C.CallbackBeforeSendMsg", "From_Account": "jared", "To_Account": "Jonh", "MsgBody": [ {"MsgType": "TIMTextElem","MsgContent": {"Text": "red packet"}}]}`

	resp,ok:=callbackHandler.NewCallbackEvent(CallbackBeforeSendMsgCommand,up,[]byte(str)).Handle().(*BaseResponse)
	if !ok ||resp.ActionStatus!=OkStatus{
		t.Fatal(ok,CallbackBeforeSendMsgCommand,"handle fail")
	}

	resp,ok=callbackHandler.NewCallbackEvent(CallbackAfterSendMsgCommand,up,nil).Handle().(*BaseResponse)
	if !ok ||resp.ActionStatus!=FAILStatus{
		t.Fatal(ok,CallbackAfterSendMsgCommand,"handle fail")
	}

        // 没有注册
	resp,ok=callbackHandler.NewCallbackEvent(CallbackBeforeCreateGroupCommand,up,nil).Handle().(*BaseResponse)
	if !ok ||resp.ActionStatus!=OkStatus{
		t.Fatal(ok,CallbackBeforeCreateGroupCommand,"handle fail")
	}


	//RegisterDefaultHandle

	callbackHandler.RegisterDefaultHandle(CallbackHandle(func(event *CallbackEvent)interface{}{
		return &BaseResponse{ActionStatus:FAILStatus,ErrorCode:0}
	}))

	resp,ok=callbackHandler.NewCallbackEvent(CallbackBeforeCreateGroupCommand,up,nil).Handle().(*BaseResponse)
	if !ok ||resp.ActionStatus!=FAILStatus{
		t.Fatal(ok,CallbackBeforeCreateGroupCommand,"handle fail")
	}

	//Exist and   UnRegister

	if !callbackHandler.Exist(CallbackBeforeSendMsgCommand){
		t.Fatal(ok,CallbackBeforeSendMsgCommand,"Exist  fail")
	}

	callbackHandler.UnRegister(CallbackBeforeSendMsgCommand)

	if callbackHandler.Exist(CallbackBeforeSendMsgCommand){
		t.Fatal(ok,CallbackBeforeSendMsgCommand,"UnRegister Exist  fail")
	}

	resp,ok=callbackHandler.NewCallbackEvent(CallbackBeforeSendMsgCommand,up,[]byte(str)).Handle().(*BaseResponse)
	if !ok ||resp.ActionStatus!=FAILStatus{
		t.Fatal(ok,CallbackBeforeSendMsgCommand,"handle fail")
	}

	_,err:=NewCallbackHandler(0,-1,CallbackHandle(defaultHandle))

	if nil==err{
		t.Fatal("test NewCallbackHandler fail")
	}

	beforeHook:=func(event *CallbackEvent)interface{}{
		return &BaseResponse{ActionStatus:OkStatus,ErrorCode:-1}
	}
	callbackHandler.RegisterBeforeHook(CallbackHandle(beforeHook))

	resp,ok=callbackHandler.NewCallbackEvent(CallbackBeforeSendMsgCommand,up,[]byte(str)).Handle().(*BaseResponse)
	if !ok ||resp.ErrorCode!=-1{
		t.Fatal(ok,CallbackBeforeSendMsgCommand,"handle fail")
	}

	beforeHook=func(event *CallbackEvent)interface{}{
		return nil
	}
	callbackHandler.RegisterBeforeHook(CallbackHandle(beforeHook))

	resp,ok=callbackHandler.NewCallbackEvent(CallbackBeforeSendMsgCommand,up,[]byte(str)).Handle().(*BaseResponse)
	if !ok ||resp.ActionStatus!=FAILStatus{
		t.Fatal(ok,CallbackBeforeSendMsgCommand,"handle fail")
	}

	time.Sleep(time.Second)
}
