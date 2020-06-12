package sms

import (
	"encoding/json"
	"errors"
	"live-service/app/config"
	"live-service/app/support/redis"
	"live-service/app/support/sms/sender"
	"live-service/app/util"
	"log"
	"time"
)

type Code struct {
	Code string		`redis:"code";json:"code"`
	Time int64		`redis:"time";json:"time"`
	Expire int64	`redis:"expire";json:"expire"`
}

// 发送验证码
func SendSmsCode(phoneNumber string, template string) error {
	appConf := config.GetAppConfig()
	smsConf := config.GetSmsConfig()

	var templateCode string
	templateParam := make(map[string]string)
	switch template {
	case "register":
		templateCode = smsConf.SmsTemplate.Register.TemplateId
		if smsConf.SmsTemplate.Register.Params > 1 {
			templateParam["time"] = string(smsConf.ValidCodeExpire)
		}
		break
	default:
		return errors.New("sms storage error:"+appConf.App.Sms.SmsSender)
	}

	templateParam["code"] = util.GenValidateCode(6) // 验证码

	// 验证码缓存
	redisClient := redis.Client()

	key := "sms_code_"+template+":"+phoneNumber

	redisClient.Del(key)

	smsCode := Code{
		Code: templateParam["code"],
		Time: time.Now().Unix(),
		Expire: 10,
	}
	log.Println(smsCode)
	log.Println(smsCode.Time + smsCode.Expire * 60)

	code, err := json.Marshal(smsCode)
	redisClient.Set(key, code, time.Second * 60 * 60)

	if err != nil {
		log.Println("sms storage err: ", err)
		return errors.New("短信发送失败")
	}

	var smsSender sender.SmsSender

	switch appConf.App.Sms.SmsSender {
	case "aliyun":
		smsSender = sender.NewAliYunSms(phoneNumber, templateCode, templateParam)
		break
	case "alitop":
		smsSender = sender.NewAliTopSms(phoneNumber, templateCode, templateParam)
		break
	default:
		return errors.New("sms storage error:"+appConf.App.Sms.SmsSender)
	}

	_,err = smsSender.Send()
	if err != nil {
		return err
	}

	return nil
}

// 验证短信验证码
func ValidSmsCode(phoneNumber string, template string, smsCode string) error {

	redisClient := redis.Client()

	// 验证码缓存
	key := "sms_code_"+template+":"+phoneNumber

	value, err := redisClient.Get(key).Result()

	if err != nil {
		log.Println(err)
		return errors.New("验证码已过期，请重新获取")
	}

	code := &Code{}
	err = json.Unmarshal([]byte(value), code)
	if err != nil {
		log.Println(err)
		return errors.New("验证码已过期，请重新获取")
	}

	if code.Time + code.Expire * 60 > time.Now().Unix() {
		if code.Code != smsCode {
			return errors.New("验证码错误")
		}

		redisClient.Del(key)
	} else {
		return errors.New("验证码已过期，请重新获取")
	}

	return nil
}
