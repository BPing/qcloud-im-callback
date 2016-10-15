# qcloud_im_callback [![Build Status](https://travis-ci.org/BPing/qcloud-im-callback.svg?branch=master)](https://travis-ci.org/BPing/qcloud-im-callback) [![Coverage Status](https://coveralls.io/repos/github/BPing/qcloud-im-callback/badge.svg)](https://coveralls.io/github/BPing/qcloud-im-callback)
腾讯IM回调处理服务


# 快速开始

```go

   package main

import (
	"fmt"
	"github.com/Bping/qcloud-im-callback"
)

func main(){

	CallbackBeforeSendMsgHandle:=func(event *qcloud_im_callback.CallbackEvent)interface{}{
		var sendMsgBody qcloud_im_callback.SendMsgBody
		event.ToJSON(&sendMsgBody)
		fmt.Println("CallbackBeforeSendMsgHandle",sendMsgBody.MsgBody)
		return &qcloud_im_callback.BaseResponse{ActionStatus:qcloud_im_callback.OkStatus,ErrorCode:0}
	}

	// 注册
	qcloud_im_callback.RegisterRouterInfo(qcloud_im_callback.CallbackBeforeSendMsgCommand,qcloud_im_callback.RouterInfo{
		Async:false,
		Handle:qcloud_im_callback.CallbackHandle(CallbackBeforeSendMsgHandle),

	})

	// http请求中url后面参数，
	// 实际应用中自行获取参数构成此结构体
	up:=qcloud_im_callback.URLParams{
		CallbackCommand:qcloud_im_callback.CallbackBeforeSendMsgCommand,
		SdkAppid:"888888",
		ContentType:"json",
	}

	// http请求中的body内容
	// json格式
	str:=`{"CallbackCommand": "C2C.CallbackBeforeSendMsg", "From_Account": "jared", "To_Account": "Jonh", "MsgBody": [ {"MsgType": "TIMTextElem","MsgContent": {"Text": "red packet"}}]}`

	resp:=qcloud_im_callback.HandleEvents(up.CallbackCommand,up,[]byte(str))

	// 自行处理返回内容
	fmt.Println(*resp.(*qcloud_im_callback.BaseResponse))
}
	
```

# 文档

https://godoc.org/github.com/BPing/qcloud-im-callback [![GoDoc](https://godoc.org/github.com/BPing/qcloud-im-callback?status.svg)](https://godoc.org/github.com/BPing/qcloud-im-callback)

# 依赖包
  github.com/BPing/Golib/producer_consumer
