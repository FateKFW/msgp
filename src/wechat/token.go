package wechat

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//请求微信服务器获取accessToken
func (wc *WeChat) reqAccessToken() {
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