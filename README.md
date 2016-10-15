# qcloud_im_callback [![Build Status](https://travis-ci.org/BPing/qcloud-im-callback.svg?branch=master)](https://travis-ci.org/BPing/qcloud-im-callback) [![Coverage Status](https://coveralls.io/repos/github/BPing/qcloud-im-callback/badge.svg)](https://coveralls.io/github/BPing/qcloud-im-callback)
腾讯IM回调处理服务


# quick start

```go

   // 默认的处理函数
   	defaultHandle:=func(event *CallbackEvent)interface{}{
                  return &BaseResponse{ActionStatus:OkStatus,ErrorCode:0}
   	}
   	// 新建处理句柄
   	callbackHandler,_:=NewCallbackHandler(1,30,CallbackHandle(defaultHandle))
   
   	
   	CallbackBeforeSendMsgHandle:=func(event *CallbackEvent)interface{}{
   		var sendMsgBody SendMsgBody
   		event.ToJSON(&sendMsgBody)
   		fmt.Println("CallbackBeforeSendMsgHandle",sendMsgBody.MsgBody)
   		return &BaseResponse{ActionStatus:OkStatus,ErrorCode:0}
   	}
   	
   	// 注册
   	callbackHandler.Register(CallbackBeforeSendMsgCommand,RouterInfo{
   		Async:false,
   		Handle:CallbackHandle(CallbackBeforeSendMsgHandle),
   
   	})
   
    // http请求中url后面参数，
    // 实际应用中自行获取参数构成此结构体
   	up:=URLParams{
   		SdkAppid:"888888",
   		ContentType:"json",
   	}
   	
   	// http请求中的body内容
   	// json格式
   	str:=`{"CallbackCommand": "C2C.CallbackBeforeSendMsg", "From_Account": "jared", "To_Account": "Jonh", "MsgBody": [ {"MsgType": "TIMTextElem","MsgContent": {"Text": "red packet"}}]}`
   	resp,ok:=callbackHandler.NewCallbackEvent(CallbackBeforeSendMsgCommand,up,[]byte(str)).Handle().(*BaseResponse)
   
   // 自行处理返回内容
  	fmt.Println(*resp.ok)
   	
```

# 文档

https://godoc.org/github.com/BPing/qcloud-im-callback [![GoDoc](https://godoc.org/github.com/BPing/qcloud-im-callback?status.svg)](https://godoc.org/github.com/BPing/qcloud-im-callback)

# 依赖包
  github.com/BPing/Golib/producer_consumer
