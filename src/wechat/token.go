package wechat

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//请求微信服务器获取accessToken
func (wc *WeChat) reqAccessToken() {
	//TODO:开发调试使用
	if true {
		wc.accessToken = "34_a9xBLUObuPdZ1kRNh_ZSKJKEFkcOpqYzMiBbRbT8" +
			"-DSrjflyAt0qQHdmfegAX-mbPhs3oFwtvqfnE-Y1tt-vAVvZXIUZBAhMLl_Xqs9G2hHa2y5lz2XPGieifxeV0YbJnga9O6yYjouoRhdMSROeAIACBY"
		return
	}

	wlog.Info("Request access_token from WeChat server")
	//获取accessToken
	res,err := http.Get(fmt.Sprintf(ACCESS_TOKEN_URL, wc.AppID, wc.AppSecret))
	defer res.Body.Close()

	if err != nil {
		wlog.NError(err)
	}

	token := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		wlog.NError(err)
	}

	if _,ok := token["errcode"]; ok {
		wlog.NErrorf("errcode:%v errmsg:%v", token["errcode"], token["errmsg"])
	}

	//设置accessToken
	wc.accessToken = token["access_token"].(string)
}