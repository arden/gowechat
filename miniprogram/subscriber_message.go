package miniprogram

import (
    "encoding/json"
    "fmt"

    "github.com/arden/wechat/context"
    "github.com/arden/wechat/util"
)

const (
    subscribeMsgSendURL = "https://api.weixin.qq.com/cgi-bin/message/subscribe/send"
)

//Template 模板消息
type Subscriber struct {
    *context.Context
}

//NewTemplate 实例化
func NewSubscriber(context *context.Context) *Subscriber {
    tpl := new(Subscriber)
    tpl.Context = context
    return tpl
}

// MsgBody 消息内容体
type MsgBody map[string]map[string]string

// SubscribeMsg 小程序订阅消息
type SubscribeMsg struct {
    OpenID   string  // 接收者（用户）的 openid
    TplID    string  // 所需下发的订阅模板ID
    PagePath string  // 点击模板卡片后的跳转页面，仅限本小程序内的页面。支持带参数,（示例index?foo=bar）。该字段不填则模板无跳转
    Data     MsgBody // 模板内容，格式形如：{"key1": {"value": any}, "key2": {"value": any}}
    MPState  string  // 跳转小程序类型：developer为开发版；trial为体验版；formal为正式版；默认为正式版
    Lang     string  // 进入小程序查看”的语言类型，支持zh_CN(简体中文)、en_US(英文)、zh_HK(繁体中文)、zh_TW(繁体中文)，默认为zh_CN
}

//DataItem 模版内某个 .DATA 的值
type DataItem struct {
    Value string `json:"value"`
    Color string `json:"color,omitempty"`
}

type resSubscriberMsgSend struct {
    util.CommonError

    MsgID int64 `json:"msgid"`
}

//Send 发送订阅消息
func (sub *Subscriber) Send(msg *SubscribeMsg) (msgID int64, err error) {
    var accessToken string
    accessToken, err = sub.GetAccessToken()
    if err != nil {
        return
    }
    uri := fmt.Sprintf("%s?access_token=%s", subscribeMsgSendURL, accessToken)
    response, err := util.PostJSON(uri, msg)

    var result resSubscriberMsgSend
    err = json.Unmarshal(response, &result)
    if err != nil {
        return
    }
    if result.ErrCode != 0 {
        err = fmt.Errorf("subscriber msg send error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
        return
    }
    msgID = result.MsgID
    return
}
