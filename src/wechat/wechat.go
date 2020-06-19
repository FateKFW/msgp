package wechat

import (
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
	accessToken string
}

//初始化相关参数
func (wc *WeChat) InitWeChatParams(){
	util.Timer(wc.reqAccessToken, EXPIRE)
}

//微信服务器请求接入
func (wc *WeChat) WxAccess(res http.ResponseWriter, req *http.Request) {
	//处理回传
	if wc.RecordReturn {
		//TODO:返回XML格式数据
		buff := make([]byte, 4096)
		i,_ := req.Body.Read(buff)
		wlog.Result("wechat return", string(buff[:i]))
	}
	//取出参数
	query := req.URL.Query()
	signature, timestamp, nonce, echostr :=
		query.Get("signature"), query.Get("timestamp"),
		query.Get("nonce"), query.Get("echostr")

	//按照微信接入规则校验参数
	param := []string{wc.Token, timestamp, nonce}
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