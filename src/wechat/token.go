package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var mu sync.Mutex

//缓存access_token(全局唯一接口调用凭据)
func (wc *WeChat) getAccessToken() (string, error){
	//TODO:开发调试时使用一个token
	/*if true {
		return "34_5ToBQtdOonHBr2E-I07e8jvyZJl9os-PcPaA07SFY6hGKsPMNxWcqVEK65PfKtEjcZak0sdpnPbrSWdWMaKyq0PmLHeRHwdp0WBjxwStzW3I3hpzA-BYA5fbK5SzEs2OE7pElzzKQ-weY_ZTMLIdAHATSP", nil
	}*/

	//是否第一次请求accesstoken
	if _,ok := wc.accessToken["accessToken"]; !ok {
		at,err := wc.reqAccessToken()
		if err != nil {
			wlog.NError(err)
			return "", err
		}
		return at, nil
	}

	//验证token是否过期
	nowTimeStamp := time.Now().Unix()
	tokenTimeStamp := wc.accessToken["time"].(int64)

	//缓存accessToken失效，重新获取
	if tokenTimeStamp - nowTimeStamp > EXPIRE {
		mu.Lock()
		defer mu.Unlock()

		nowTimeStamp = time.Now().Unix()
		tokenTimeStamp = wc.accessToken["time"].(int64)
		if tokenTimeStamp - nowTimeStamp > EXPIRE {
			at, err := wc.reqAccessToken()
			if err != nil {
				wlog.NError(err)
				return "", err
			}
			return at, nil
		}
	}

	return wc.accessToken["accessToken"].(string), nil
}

//请求微信服务器获取accessToken
func (wc *WeChat) reqAccessToken() (string, error){
	wlog.Info("Request access_token from WeChat server")
	//获取accessToken
	res,err := http.Get(fmt.Sprintf(ACCESS_TOKEN_URL, wc.AppID, wc.AppSecret))
	defer res.Body.Close()

	if err != nil {
		wlog.NError(err)
		return "", err
	}

	token := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		wlog.NError(err)
		return "", err
	}

	if _,ok := token["errcode"]; ok {
		wlog.Infof("errcode:%v errmsg:%v", token["errcode"], token["errmsg"])
		return "", errors.New(token["errmsg"].(string))
	}

	//设置accessToken
	wc.accessToken["time"] = time.Now().Unix()
	wc.accessToken["accessToken"] = token["access_token"]
	wc.accessToken["expiresIn"] = token["expires_in"]
	return token["access_token"].(string), nil
}