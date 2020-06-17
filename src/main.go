package main

import (
	"flag"
	"msgplog"
	"net/http"
	"os"
	"wechat"
)

var mlog = msgplog.Logger

type GateWay struct {
	Port string
	OpenWX bool
	wx *wechat.WeChat
}

func bindParam() *GateWay{
	//命令行参数绑定(是否改成配置文件形式?)
	var port, appid, appsecret, token string
	var help, recordreturn, logfile, openwx bool
	var loglevel int

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
	mlog.InitMSGPLog(logfile, loglevel)

	//消息处理总接口初始
	mlog.Info("Initialize the gateway")
	gw := &GateWay{":"+port, openwx, nil}
	mlog.Info("Successfully Initialize the gateway")

	//微信初始
	if openwx {
		mlog.Info("Initialize WeChat module")
		wx := new(wechat.WeChat)
		wx.AppID = appid
		wx.AppSecret = appsecret
		wx.Token = token
		wx.RecordReturn = recordreturn
		wx.InitWeChatParams()
		gw.wx = wx
		mlog.Info("Successfully initialized WeChat module")
	}

	//TODO:后续接入模块绑定

	return gw
}

func (gw *GateWay) start() {
	if gw.OpenWX {
		//微信服务器接入
		http.HandleFunc("/wx/access", gw.wx.WxAccess)
		//发送模板消息
		http.HandleFunc("/wx/send", gw.wx.SendTemplateMessage)
	}

	//开启服务
	mlog.Info("Gateway service started successfully")
	err := http.ListenAndServe(gw.Port, nil)
	if err != nil {
		mlog.Error(err)
	}
}

func main() {
	gw := bindParam()
	gw.start()
}