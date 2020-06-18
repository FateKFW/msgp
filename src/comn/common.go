package comn

import (
	"flag"
	"msgp/log"
	"msgp/wechat"
	"net/http"
	"os"
)

var clog = msgplog.Logger

type GateWay struct {
	Port string
	OpenWX bool
	Wx *wechat.WeChat
}

func BindParam() *GateWay{
	//命令行参数绑定(是否改成配置文件形式?)
	var port, appid, appsecret, token string
	var help, recordreturn, logfile, openwx bool
	var loglevel = 0x01

	flag.BoolVar(&help, "h", false, "Parameter Description")
	flag.StringVar(&port, "p", "8080", "Gateway port number")
	flag.BoolVar(&openwx, "openwx", false, "Whether to open the WeChat interface")
	//日志相关
	flag.BoolVar(&logfile, "lf", false, "Whether the log is recorded to a file")
	flag.IntVar(&loglevel, "ll", 1, "Log output level")
	//微信相关
	flag.StringVar(&appid, "appid", "", "WeChat appID")
	flag.StringVar(&appsecret, "appsecret", "", "WeChat appsecret")
	flag.StringVar(&token, "token", "", "WeChat Token")
	flag.BoolVar(&recordreturn, "wxrr", false, "Whether to process the WeChat server callback message")

	flag.Parse()
	if help {
		flag.PrintDefaults()
		os.Exit(1)
	}

	//日志模块初始
	clog.InitMSGPLog(logfile, loglevel)

	//消息处理总接口初始
	clog.Info("Initialize the msgp")
	gw := &GateWay{":"+port, openwx, nil}
	clog.Info("Successfully Initialize the msgp")

	//微信初始
	if openwx {
		clog.Info("Initialize WeChat module")
		wx := new(wechat.WeChat)
		wx.AppID = appid
		wx.AppSecret = appsecret
		wx.Token = token
		wx.RecordReturn = recordreturn
		wx.InitWeChatParams()
		gw.Wx = wx
		clog.Info("Successfully initialized WeChat module")
	}

	//TODO:后续接入模块绑定
	return gw
}

func (gw *GateWay) Start() {
	if gw.OpenWX {
		//微信服务器接入
		http.HandleFunc("/wx/access", gw.Wx.WxAccess)
		//发送模板消息
		http.HandleFunc("/wx/send", gw.Wx.SendTemplateMessage)
	}

	//开启服务
	clog.Info("msgp service started successfully")
	err := http.ListenAndServe(gw.Port, nil)
	if err != nil {
		clog.Error(err)
	}
}