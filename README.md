# qcloud_im_callback [![Build Status](https://travis-ci.org/BPing/qcloud-im-callback.svg?branch=dev)](https://travis-ci.org/BPing/qcloud-im-callback) [![Coverage Status](https://coveralls.io/repos/github/BPing/qcloud-im-callback/badge.svg?branch=dev)](https://coveralls.io/github/BPing/qcloud-im-callback?branch=dev)
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
	up:=&qcloud_im_callback.URLParams{
		CallbackCommand:qcloud_im_callback.CallbackBeforeSendMsgCommand,
		SdkAppid:"888888",
		ContentType:"json",
	}

	// http请求中的body内容
	// json格式
	str:=`{"CallbackCommand": "C2C.CallbackBeforeSendMsg", "From_Account": "jared", "To_Account": "Jonh", "MsgBody": [ {"MsgType": "TIMTextElem","MsgContent": {"Text": "red packet"}}]}`

	resp:=qcloud_im_callback.HandleEvents(qcloud_im_callback.CreateEvents(up.CallbackCommand,up,[]byte(str)))

	// 自行处理返回内容
	fmt.Println(*resp.(*qcloud_im_callback.BaseResponse))
}
	
```

# 详细

>**Note:** 使用前，记得初始化类库和注册各种事件处理程序

* 更改默认CallbackHandler
```go
    // @masterNum 主消费线程数目,必须大于等于1
    // @chanLen   消费信息（事件）队列长度
    qcloud_im_callback.RegisterDefaultCallbackHandler(masterNum,msgEventLen int,defaultHandle func(*CallbackEvent) interface{})
```
* 注册事件处理路由信息（处理程序）。具体事件类型查看文件type.go
```go
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
```

### http
已提供了标准http调用处理函数了，可以直接调用，如下：
```go
    http.HandleFunc("/im/callback",qcloud_im_callback.HandleEventsHttp)
```

### beego
因为不想耦合beego,所以就不提供beego定制化接口了。可以参考下面例子：
```go
    func (u *Controller) HandleEventsBeego() {
    	up:=&qcloud_im_callback.URLParams{
    		SdkAppid:u.GetString("SdkAppid"),
    		CallbackCommand:qcloud_im_callback.CallbackCommand(u.GetString("CallbackCommand")),
    		ContentType:u.GetString("contenttype"),
    		ClientIP:u.GetString("ClientIP"),
    		OptPlatform:qcloud_im_callback.OptPlatform(u.GetString("OptPlatform")),
    	}
    	resp:=qcloud_im_callback.HandleEvents(qcloud_im_callback.CreateEvents(up.CallbackCommand,up,u.Ctx.Input.RequestBody))
    	u.Data["json"] = resp
    	u.ServeJSON()
    }
```


### cache队列（如：reids）

如果你想用redis第三方缓存工具来传递信息，你可以通过以下方法初始化。如果没有指定cache队列类型，则统一默认为channel缓冲队列类型

* 更改默认

```go
    qcloud_im_callback.RegisterDefaultNewCallbackHandlerWithCache(masterNum, msgEventLen int, defaultHandle func(*CallbackEvent) interface{},cache qcloud_im_callback.ICache) (err error)   
```

* 新建

```go
    qcloud_im_callback.NewCallbackHandlerWithCache(masterNum, msgEventLen int, defaultHandle CallbackHandle,cache qcloud_im_callback.ICache) (*qcloud_im_callback.CallbackHandler, error)  
```

> 请注意，必须实现qcloud_im_callback.ICache接口,具体如下：

```go
    // 缓存Cache接口
   	type ICache interface {
   		// BLPOP key1 timeout(秒)
   		// 移出并获取列表的第一个元素，
   		// 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
   		BLPop(key string,timeout int64)(map[string]string, error)
   		// 在列表尾部中添加一个或多个值
   		RPush(key string,values ... interface{}) (int64, error)
   		// 获取列表长度
   		LLen(key string) (int64, error)
   	}
```

# 文档

https://godoc.org/github.com/BPing/qcloud-im-callback [![GoDoc](https://godoc.org/github.com/BPing/qcloud-im-callback?status.svg)](https://godoc.org/github.com/BPing/qcloud-im-callback)

# 依赖包
  github.com/BPing/Golib/producer_consumer
