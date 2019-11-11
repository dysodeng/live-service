package config

// 短信模版配置
type SmsConfig struct {
	// 短信模版
	SmsTemplate struct {
		Register struct {
			TemplateId string	`yaml:"template_id"`
			Name string			`yaml:"name"`
			Params uint8		`yaml:"params"`
		} `yaml:"register"`
	} `yaml:"sms_template"`

	// 短信验证码过期时间(分钟)
	ValidCodeExpire int64	`yaml:"valid_code_expire"`
}

