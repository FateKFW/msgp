package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"msgp/util"
	"net/http"
	"strings"
	"time"
)

//启动模板消息发送任务
func (wx *WeChat) SendTemplateMessage(res http.ResponseWriter, req *http.Request) {
	openids, tid, data, url :=
		req.PostFormValue("openids"),
		req.PostFormValue("tid"),
		req.PostFormValue("data"),
		req.PostFormValue("url")

	go wx.handleTemplateMessage(openids, tid, data, url)

	res.Write([]byte("模板消息发送已执行，发送结果请查看XXXXX"))
}

//执行模板消息发送
func (wx *WeChat) handleTemplateMessage(openids string, tid string, data string, url string) {
	//需要发送的人员
	openidArr := strings.Split(openids, ",")
	//缓冲管道，缓冲中最多存储5000个待发
	ch := make(chan string, 5000)

	for _, obj := range openidArr {
		go (func(openid string){
			//组装请求参数
			content,err := util.JsonStr2Map(data)
			if err != nil {
				wlog.NError(err)
				ch <- openid + ",-1," + err.Error()
				return
			}

			content["touser"] = openid
			content["template_id"] = tid
			if url != "" {
				content["url"] = url
			}

			//发送模板消息
			buff,err := json.Marshal(content)
			if err != nil {
				wlog.NError(err)
				ch <- openid + ",-2," + err.Error()
				return
			}

			res,err := http.Post(fmt.Sprintf(SEND_TEMPLATE_URL, wx.accessToken), "", bytes.NewBuffer(buff))
			if err != nil {
				wlog.NError(err)
				ch <- openid + ",-3," + err.Error()
				return
			}

			var result map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&result)
			if err != nil {
				wlog.NError(err)
				ch <- openid + ",-4," + err.Error()
				return
			}

			ch <- fmt.Sprintf("%s,%.0f,%s", openid, result["errcode"], result["errmsg"])
		})(obj)
	}

	for i:=0; i<len(openidArr); i++  {
		result := strings.Split(<- ch, ",")
		//TODO
		wlog.Result("Send message execution result", result)
	}
}

//处理普通消息-文本消息
func (wx *WeChat) handleTextMessage(push WXPush) []byte{
	wlog.Infof("accept openid:%s message", push.Openid)
	msg := fmt.Sprintf("<xml>" +
		"<ToUserName><![CDATA[%s]]></ToUserName>" +
		"<FromUserName><![CDATA[%s]]></FromUserName>" +
		"<CreateTime>%v</CreateTime>" +
		"<MsgType><![CDATA[text]]></MsgType>" +
		"<Content><![CDATA[%s]]></Content>" +
		"</xml>",push.Openid, wx.fromUserName, time.Now().Unix(), getRobotMessage(push.Content))
	return []byte(msg)
}

func getRobotMessage(content string) string{
	content = strings.ReplaceAll(content, "吗", "")
	content = strings.ReplaceAll(content, "?", "")
	return content
}