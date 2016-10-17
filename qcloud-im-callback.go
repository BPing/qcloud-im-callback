package qcloud_im_callback

import (
	"net/http"
	"io"
	"io/ioutil"
	"encoding/json"
)

var (
	BodyMaxLen int64=2048 //请求body实体内容限制长度值
)

// 默认的
//   系统初始实例的的事件队列长度为50，主要消费线程为1，默认处理返回操作成功结构
//   用户可以通过RegisterDefaultCallbackHandler重新替换默认句柄的
var defaultCallbackHandler *CallbackHandler

func init(){
	defaultHandle:=func(event *CallbackEvent)interface{}{
		return &BaseResponse{ActionStatus:OkStatus,ErrorCode:0}
	}
	RegisterDefaultCallbackHandler(1,50,defaultHandle)
}

func RegisterDefaultCallbackHandler(masterNum,msgEventLen int,defaultHandle func(*CallbackEvent) interface{})(err error){
	temp,err:=NewCallbackHandler(masterNum,msgEventLen,CallbackHandle(defaultHandle))
	if err==nil{
		defaultCallbackHandler=temp
	}
	return
}

func RegisterDefaultHandle(defaultHandle func(*CallbackEvent) interface{}){
         defaultCallbackHandler.RegisterDefaultHandle(CallbackHandle(defaultHandle))
}

func RegisterBeforeHook(defaultHandle func(*CallbackEvent) interface{}){
         defaultCallbackHandler.RegisterBeforeHook(CallbackHandle(defaultHandle))
}

func RegisterRouterInfo(cc CallbackCommand, ri RouterInfo){
	defaultCallbackHandler.Register(cc,ri)
}

func Handle(cc CallbackCommand, up *URLParams, body []byte)interface{}{
        return defaultCallbackHandler.NewCallbackEvent(cc,up,body).Handle()
}

func CreateEvents(cc CallbackCommand, up *URLParams, body []byte)*CallbackEvent{
	return defaultCallbackHandler.NewCallbackEvent(cc,up,body)
}

func HandleEvents(event *CallbackEvent)interface{}{
	return event.Handle()

}

// for std（raw） http
//----------------------------------------------------------------------------------------------------------------------

// 原始http请求处理
//    对于body内容目前做长度限制，具体查看变量BodyMaxLen
func HandleEventsHttp(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	up:=&URLParams{
		SdkAppid:getParam("SdkAppid",r),
		CallbackCommand:CallbackCommand(getParam("CallbackCommand",r)),
		ContentType:getParam("contenttype",r),
		ClientIP:getParam("ClientIP",r),
		OptPlatform:OptPlatform(getParam("OptPlatform",r)),
	}
	resp:=HandleEvents(CreateEvents(up.CallbackCommand,up,readBody(r)))
	data,_:=json.Marshal(resp)
	w.Write(data)
}

// 获取url附带的参数
func getParam(key string,r *http.Request)string{
	if len(r.Form[key]) > 0 {
		return r.Form[key][0]
	}
      return ""
}

// 读取请求body内容
//  对于body内容目前做长度限制，具体查看变量BodyMaxLen
func readBody(r *http.Request)[]byte{
	safe := &io.LimitedReader{R: r.Body, N: BodyMaxLen}
	reqBody, _ := ioutil.ReadAll(safe)
	r.Body.Close()
        return reqBody
}

// Beego框架例子
//    因为不想耦合beego,所以就不提供beego定制化接口了。请自行参考此例子代码
//----------------------------------------------------------------------------------------------------------------------
//
//func (u *Controller) HandleEventsBeego() {
//	up:=&URLParams{
//		SdkAppid:u.GetString("SdkAppid"),
//		CallbackCommand:CallbackCommand(u.GetString("CallbackCommand")),
//		ContentType:u.GetString("contenttype"),
//		ClientIP:u.GetString("ClientIP"),
//		OptPlatform:OptPlatform(u.GetString("OptPlatform")),
//	}
//	resp:=HandleEvents(CreateEvents(up.CallbackCommand,up,u.Ctx.Input.RequestBody))
//	u.Data["json"] = resp
//	u.ServeJSON()
//}
