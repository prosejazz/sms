package sms

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	APP_ID     string = ""
	APP_SECRET string = ""
)

type SmsResponse struct {
	Identifier string
	Create_at  string
}

type Sms struct {
	AccessToken string
	Code        float64
	Message     string
	Token       string
}

func NewSms(accesstoken string) (*Sms, error) {
	s := new(Sms)
	s.AccessToken = accesstoken
	err := s.getToken()
	if err != nil {
		return nil, err
	}
	return s, nil

}

func (s *Sms) CustomSms(phone, randcode, exp_time string) (*SmsResponse, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	querys := map[string]string{}
	querys["app_id"] = APP_ID
	querys["access_token"] = s.AccessToken
	querys["timestamp"] = timestamp
	querys["token"] = s.Token
	querys["phone"] = phone
	querys["randcode"] = randcode
	querys["exp_time"] = exp_time
	u := url.Values{}
	for k, v := range querys {
		u.Set(k, v)
	}
	q := createSignQuery(querys)
	// 生成签名
	u.Set("sign", createSign(q))
	// 发送请求
	response, err := http.PostForm("http://api.189.cn/v2/dm/randcode/sendSms", u)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	resStruct := new(SmsResponse)
	err = json.Unmarshal(result, resStruct)
	return resStruct, err
}

type tokenResult struct {
	Res_code int
	Token    string
}

func (s *Sms) getToken() error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	u, err := url.Parse("http://api.189.cn/v2/dm/randcode/token")
	if err != nil {
		return err
	}
	q := u.Query()
	querys := map[string]string{}
	querys["app_id"] = APP_ID
	querys["access_token"] = s.AccessToken
	querys["timestamp"] = timestamp
	for k, v := range querys {
		q.Set(k, v)
	}
	sq := createSignQuery(querys)
	// 签名
	q.Set("sign", createSign(sq))
	u.RawQuery = q.Encode()
	// 发送get请求
	response, err := http.Get(u.String())
	if err != nil {
		return err
	}
	// 读取返回结果
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	// 把结果解析到结构体
	tokenresult := new(tokenResult)
	err = json.Unmarshal(result, tokenresult)
	if err != nil {
		return err
	}
	if tokenresult.Res_code == 0 {
		s.Token = tokenresult.Token
		return nil
	} else {
		return errors.New(strconv.Itoa(tokenresult.Res_code) + "  " + tokenresult.Token)
	}
}

func getMapKeys(themap map[string]string) []string {
	keys := []string{}
	for k := range themap {
		keys = append(keys, k)
	}
	return keys
}

func createSignQuery(params map[string]string) string {
	keys := getMapKeys(params)
	sort.Strings(keys)
	q := ""
	for _, v := range keys {
		q += "&" + v + "=" + params[v]
	}
	return strings.TrimLeft(q, "&")
}

func createSign(querys string) string {
	return base64Encode(sha1Encode(querys))
}

func sha1Encode(str string) []byte {
	h := hmac.New(sha1.New, []byte(APP_SECRET))
	h.Write([]byte(str))
	return h.Sum(nil)
}

func base64Encode(str []byte) string {
	return base64.StdEncoding.EncodeToString(str)
}
