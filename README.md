sms
===

中国短信短信验证码发送接口golang sdk

### 使用样例
	// 打开sms.go填入APP_ID和APP_SECRET
	ats, err := sms.GetAccessToken()
	if err!=nil {
		
	}
	s,err:=sms.NewSms(ats.Access_token)
	if err!=nil {
		
	}
	rand.Seed(time.Now().UnixNano())
	randstr := fmt.Sprintf("%06d", rand.Intn(999999))
	// 验证码有效时间5分钟
	s.CustomSms("13888888888", randstr, "5")
