package sms

import (
	"errors"
	"fmt"
	baseRedis "github.com/gomodule/redigo/redis"
	"live-service/app/config"
	"live-service/app/support/sms/sender"
	"live-service/app/util"
	"live-service/app/util/database"
	"log"
	"time"
)

type Code struct {
	Code string		`redis:"code"`
	Time int64		`redis:"time"`
	Expire int64	`redis:"expire"`
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
			templateParam["time"] = fmt.Sprintf("%d", smsConf.ValidCodeExpire)+"分钟"
		}
		break
	default:
		return errors.New("sms storage error:"+appConf.App.Sms.SmsSender)
	}

	templateParam["code"] = util.GenValidateCode(6) // 验证码

	// 验证码缓存
	redis := database.GetRedis()
	defer redis.Close()

	key := "sms_code_"+template+":"+phoneNumber

	_, _ = redis.Do("DEL", key)

	smsCode := Code{
		Code: templateParam["code"],
		Time: time.Now().Unix(),
		Expire: smsConf.ValidCodeExpire,
	}
	log.Println(smsCode)
	log.Println(smsCode.Time + smsCode.Expire * 60)

	result, err := redis.Do("HMSET", baseRedis.Args{}.Add(key).AddFlat(smsCode)...)
	if err != nil {
		return err
	}
	_, _ = redis.Do("EXPIRE", key, smsConf.ValidCodeExpire*60)
	log.Println(result)

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
	// 验证码缓存
	redis := database.GetRedis()
	defer redis.Close()

	key := "sms_code_"+template+":"+phoneNumber

	value, err := baseRedis.Values(redis.Do("HGETALL", key))
	if err != nil {
		log.Println(err)
		return errors.New("验证码已过期，请重新获取")
	}

	code := &Code{}
	err = baseRedis.ScanStruct(value, code)
	if err != nil {
		log.Println(err)
		return errors.New("验证码已过期，请重新获取")
	}

	if code.Time + code.Expire * 60 > time.Now().Unix() {
		if code.Code != smsCode {
			return errors.New("验证码错误")
		}

		_, _ = redis.Do("DEL", key)
	} else {
		return errors.New("验证码已过期，请重新获取")
	}

	return nil
}
