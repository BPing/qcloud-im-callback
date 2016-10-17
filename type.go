package qcloud_im_callback

//回调命令
//@link https://www.qcloud.com/doc/product/269/1523
type CallbackCommand string

// 客户端平台
type OptPlatform string

// 请求处理的结果，OK表示处理成功，FAIL表示失败。
type ActionStatus string

const (
	//-- 单发消息 --
	CallbackBeforeSendMsgCommand = CallbackCommand("C2C.CallbackBeforeSendMsg") // 发消息之前回调
	CallbackAfterSendMsgCommand  = CallbackCommand("C2C.CallbackAfterSendMsg")  // 发消息之后回调

	//-- 群组系统 --
	CallbackBeforeCreateGroupCommand     = CallbackCommand("Group.CallbackBeforeCreateGroup")     // 创建群组之前回调
	CallbackAfterCreateGroupCommand      = CallbackCommand("Group.CallbackAfterCreateGroup")      // 创建群组之后回调
	CallbackBeforeApplyJoinGroupCommand  = CallbackCommand("Group.CallbackBeforeApplyJoinGroup")  // 申请入群之前回调
	CallbackBeforeInviteJoinGroupCommand = CallbackCommand("Group.CallbackBeforeInviteJoinGroup") // 拉人入群之后回调
	CallbackAfterNewMemberJoinCommand    = CallbackCommand("Group.CallbackAfterNewMemberJoin")    // 新成员入群之后回调
	CallbackAfterMemberExitCommand       = CallbackCommand("Group.CallbackAfterMemberExit")       // 群成员离开之后回调
	GroupCallbackBeforeSendMsgCommand    = CallbackCommand("Group.CallbackBeforeSendMsg")         // 群内发言之前回调
	GroupCallbackAfterSendMsgCommand     = CallbackCommand("Group.CallbackAfterSendMsg")          // 群内发言之后回调
	CallbackAfterGroupFullCommand        = CallbackCommand("Group.CallbackAfterGroupFull")        // 群组满员之后回调
	CallbackAfterGroupDestroyedCommand   = CallbackCommand("Group.CallbackAfterGroupDestroyed")   // 群组解散之后回调
	CallbackAfterGroupInfoChangedCommand = CallbackCommand("Group.CallbackAfterGroupInfoChanged") // 群组资料修改之后回调

	//-- 关系链系统 --
	CallbackFriendAddCommand       = CallbackCommand("Sns.CallbackFriendAdd")       // 添加好友之后回调
	CallbackFriendDeleteCommand    = CallbackCommand("Sns.CallbackFriendDelete")    // 删除好友之后回调
	CallbackBlackListAddCommand    = CallbackCommand("Sns.CallbackBlackListAdd")    // 添加黑名单之后回调
	CallbackBlackListDeleteCommand = CallbackCommand("Sns.CallbackBlackListDelete") // 删除黑名单之后回调

	//-- 在线状态 --
	StateChangeCommand = CallbackCommand("State.StateChange") // 状态变更回调

	// 客户端平台
	RestAPIPlatform = OptPlatform("RESTAPI")
	WebPlatform     = OptPlatform("Web")
	AndroidPlatform = OptPlatform("Android")
	IOSPlatform     = OptPlatform("iOS")
	WindowsPlatform = OptPlatform("Windows")
	MacPlatform     = OptPlatform("Mac")
	UnkownPlatform  = OptPlatform("Unkown")

	// 请求处理结果
	OkStatus   = ActionStatus("OK")
	FAILStatus = ActionStatus("FAIL")
)

// 腾讯云在发起回调时，会在APP提供的URL之后增加如下几个参数：
// @link https://www.qcloud.com/doc/product/269/1522  回调协议
type URLParams struct {

	// APP在云通讯申请的Appid。
	SdkAppid string

	// 回调命令字。
	CallbackCommand CallbackCommand

	//固定为：json。对应：contenttype
	ContentType string

	// 客户端IP地址
	ClientIP string

	// 客户端平台。对应不同的平台类型，
	// 可能的取值有： RESTAPI（使用REST API发送请求）、
	//              Web（使用Web SDK发送请求）、
	//              Android、
	//              iOS、
	//              Windows、
	//              Mac、
	//              Unkown（使用未知类型的设备发送请求）。
	OptPlatform OptPlatform
}

func (up *URLParams) Clone() *URLParams{
	new_obj := (*up)
	return &new_obj
}

// 基本应答包字段
type BaseResponse struct {
	// 请求处理的结果，
	// OK表示处理成功，FAIL表示失败。
	ActionStatus ActionStatus

	// 错误码
	ErrorCode int

	//错误信息。
	ErrorInfo string
}

// 基本请求body内容字段
type BaseBody struct {
	CallbackCommand CallbackCommand
}

// body json
//----------------------------------------------------------------------------------------------------------------------

//发单聊消息
type SendMsgBody struct {
	BaseBody
	From_Account string
	To_Account   string
	MsgBody      []struct {
		MsgType    string
		MsgContent map[string]interface{}
	}
}
