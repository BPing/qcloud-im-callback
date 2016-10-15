package qcloud_im_callback

import (
	"testing"
	"fmt"
)

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

}