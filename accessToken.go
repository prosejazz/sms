package sms

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type accessTokenResponse struct {
	Access_token string
	Expires_in   int
	Res_code     string
	Res_message  string
}

func GetAccessToken(state string) (*accessTokenResponse, error) {
	u := url.Values{}
	u.Add("grant_type", "client_credentials")
	u.Add("app_id", APP_ID)
	u.Add("app_secret", APP_SECRET)
	u.Add("state", state)
	u.Add("scope", "")

	client := http.Client{}
	req, err := http.NewRequest("POST", "https://oauth.api.189.cn/emp/oauth2/v3/access_token", strings.NewReader(u.Encode()))
	if err != nil {
		return nil, err
	}
	// 设置请求头，表示为表单提交
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	r, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	// 解析到结构体中
	atr := new(accessTokenResponse)
	err = json.Unmarshal(r, atr)
	if err != nil {
		return nil, err
	}
	return atr, nil
}
