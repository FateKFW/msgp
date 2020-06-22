package wechat

import (
	"encoding/xml"
	"msgp/log"
	"msgp/util"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	//accessToken过期时长(不能设置超过2小时)
	EXPIRE = 1*time.Hour + 55*time.Minute
	//获取accessToken的URL
	ACCESS_TOKEN_URL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	//发送模板消息的URL
	SEND_TEMPLATE_URL = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"
)

var wlog = msgplog.Logger

type WeChat struct {
	AppID string		`json:appid`
	AppSecret string	`json:appsecret`
	Token string		`json:token`
	RecordReturn bool	`json:recordreturn`
	fromUserName string
	accessToken string
}

type WXPush struct {
	XMLName xml.Name	`xml:xml`
	ToUserName string	`xml:ToUserName`
	FromUserName string	`xml:FromUserName`
	CreateTime int64	`xml:CreateTime`
	MsgType string		`xml:MsgType`
	MsgId string		`xml:MsgId`
	//普通消息使用
	Content string		`xml:Content`
	Openid string
}


//初始化相关参数
func (wx *WeChat) InitWeChatParams(){
	wx.fromUserName = "gh_1e29e748bead"
	util.Timer(wx.reqAccessToken, EXPIRE)
}

//微信服务器请求接入
func (wx *WeChat) WxAccess(res http.ResponseWriter, req *http.Request) {
	//取出参数
	query := req.URL.Query()
	signature, timestamp, nonce, echostr,openid :=
		query.Get("signature"),
		query.Get("timestamp"),
		query.Get("nonce"),
		query.Get("echostr"),
		query.Get("openid")

	if openid != "" {	//微信服务器回传
		//回传处理是否开启
		if wx.RecordReturn {
			var push WXPush
			err := xml.NewDecoder(req.Body).Decode(&push)
			if err != nil {
				wlog.NError(err)
			}
			push.Openid = openid
			res.Write(wx.handlePush(push))
		}
	} else {			//微信服务器对接
		//按照微信接入规则校验参数
		param := []string{wx.Token, timestamp, nonce}
		sort.Strings(param)
		cry := util.SHA1(strings.Join(param, ""))

		//验证成功
		if signature == cry {
			res.Write([]byte(echostr))
			wlog.Info("WeChat server access verification succeeded")
		} else {
			res.Write([]byte("failed"))
			wlog.Info("WeChat server access verification failed")
		}
	}
}

//处理微信服务器回传
func (wx *WeChat) handlePush(push WXPush) []byte{
	if "text" == push.MsgType {
		return wx.handleTextMessage(push)
	} else {
		wlog.Info(push)
		wlog.Info("No support for processing messages")
		return []byte("(✺ω✺)")
	}
}