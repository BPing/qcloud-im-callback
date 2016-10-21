package qcloud_im_callback

import (
	"testing"
)

func TestEvents(t *testing.T) {
	up := &URLParams{
		SdkAppid:    "888888",
		ContentType: "json",
	}

	str := `{"CallbackCommand": "C2C.CallbackBeforeSendMsg", "From_Account": "jared", "To_Account": "Jonh", "MsgBody": [ {"MsgType": "TIMTextElem","MsgContent": {"Text": "red packet"}}]}`
	var sendMsgBody SendMsgBody

	err := NewCallbackEvent(CallbackBeforeSendMsgCommand, up, []byte(str)).ToJSON(&sendMsgBody)

	if err != nil || sendMsgBody.From_Account != "jared" {
		t.Fatal("ToJSON fail")
	}

	err = NewCallbackEvent(CallbackBeforeSendMsgCommand, up, nil).ToJSON(&sendMsgBody)

	if err == nil || err.Error() != "body is nil" {
		t.Fatal("test nil ToJSON fail")
	}

	if NewCallbackEvent(CallbackBeforeSendMsgCommand, up, nil).Handle() != nil {
		t.Fatal("test (handler is nil) fial")
	}

	NewCallbackEvent(CallbackBeforeSendMsgCommand, up, nil).handle()

}
