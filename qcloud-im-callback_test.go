package qcloud_im_callback

import (
	"testing"
	"fmt"
	"net/http/httptest"
	"net/http"
	"encoding/json"
	"bytes"
)

// @depend TestFacade
func TestHttp(t *testing.T){
	RegisterBeforeHook(func(event *CallbackEvent)interface{}{
		return nil
	})

	callbackAfterCreateGroupHandle:=func(event *CallbackEvent)interface{}{
		var sendMsgBody SendMsgBody
		event.ToJSON(&sendMsgBody)
		fmt.Println("TestHttp Group.CallbackAfterCreateGroup",sendMsgBody.MsgBody)
		return &BaseResponse{ActionStatus:OkStatus,ErrorCode:0}
	}

	// 注册
	RegisterRouterInfo(CallbackAfterCreateGroupCommand,RouterInfo{
		Async:false,
		Handle:CallbackHandle(callbackAfterCreateGroupHandle),

	})

	str:=`{"CallbackCommand": "Group.CallbackAfterCreateGroup", "From_Account": "jared", "To_Account": "Jonh", "MsgBody": [ {"MsgType": "TIMTextElem","MsgContent": {"Text": "red packet"}}]}`
	url:="/callback?SdkAppid=888888&CallbackCommand=Group.CallbackAfterCreateGroup&contenttype=json&ClientIP=127&OptPlatform=window"
	writer:=httptest.NewRecorder()
	req,_:=http.NewRequest("POST",url,bytes.NewBuffer([]byte(str)))

	HandleEventsHttp(writer,req)
	var data BaseResponse
	json.Unmarshal(writer.Body.Bytes(),&data)

	if data.ActionStatus!=OkStatus{
		t.Fatal("test HandleEventsHttp")
	}
}

func TestFacade(t *testing.T){
	CallbackBeforeSendMsgHandle:=func(event *CallbackEvent)interface{}{
		var sendMsgBody SendMsgBody
		event.ToJSON(&sendMsgBody)
		fmt.Println("CallbackBeforeSendMsgHandle",sendMsgBody.MsgBody)
		return &BaseResponse{ActionStatus:OkStatus,ErrorCode:0}
	}

	// 注册
	RegisterRouterInfo(CallbackBeforeSendMsgCommand,RouterInfo{
		Async:false,
		Handle:CallbackHandle(CallbackBeforeSendMsgHandle),

	})

	// http请求中url后面参数，
	// 实际应用中自行获取参数构成此结构体
	up:=URLParams{
		CallbackCommand:CallbackBeforeSendMsgCommand,
		SdkAppid:"888888",
		ContentType:"json",
	}

	// http请求中的body内容
	// json格式
	str:=`{"CallbackCommand": "C2C.CallbackBeforeSendMsg", "From_Account": "jared", "To_Account": "Jonh", "MsgBody": [ {"MsgType": "TIMTextElem","MsgContent": {"Text": "red packet"}}]}`

	resp:=HandleEvents(CreateEvents(up.CallbackCommand,up,[]byte(str)))

	if resp.(*BaseResponse).ActionStatus!=OkStatus{
		t.Fatal("test facade")
	}

	resp=Handle(up.CallbackCommand,up,[]byte(str))

	if resp.(*BaseResponse).ActionStatus!=OkStatus{
		t.Fatal("test facade Handle")
	}
	
	resp=Handle(CallbackAfterSendMsgCommand,up,[]byte(str))

	if resp.(*BaseResponse).ActionStatus!=OkStatus{
		t.Fatal("test facade default Handle")
	}

	RegisterDefaultHandle(func(event *CallbackEvent)interface{}{
		return &BaseResponse{ActionStatus:FAILStatus,ErrorCode:0}
	})

	resp=Handle(CallbackAfterSendMsgCommand,up,[]byte(str))

	if resp.(*BaseResponse).ActionStatus!=FAILStatus{
		t.Fatal("test facade RegisterDefaultHandle")
	}

	//RegisterBeforeHook

	RegisterBeforeHook(func(event *CallbackEvent)interface{}{
		return nil
	})

	resp=Handle(CallbackAfterSendMsgCommand,up,[]byte(str))

	if resp.(*BaseResponse).ActionStatus!=FAILStatus{
		t.Fatal("test facade RegisterBeforeHook 1")
	}

	RegisterBeforeHook(func(event *CallbackEvent)interface{}{
		return &BaseResponse{ActionStatus:FAILStatus,ErrorCode:-1}
	})

	resp=Handle(CallbackAfterSendMsgCommand,up,[]byte(str))

	if resp.(*BaseResponse).ErrorCode!=-1{
		t.Fatal("test facade RegisterBeforeHook 2",resp.(*BaseResponse).ErrorCode)
	}

	RegisterBeforeHook(func(event *CallbackEvent)interface{}{
		return nil
	})
}

