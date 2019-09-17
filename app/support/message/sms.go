package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	baseRedis "github.com/gomodule/redigo/redis"
	"live-service/app/util"
	"live-service/app/util/config"
	"live-service/app/util/database"
	"log"
	"time"
)

type SmsCode struct {
	Code string		`redis:"code"`
	Time int64		`redis:"time"`
	Expire int64	`redis:"expire"`
}

type SmsSender interface {
	// 发送短信
	Send() (bool, error)
}

// 阿里云短信
type AliYunSmsSender struct {
	phoneNumber string
	signName string
	templateCode string
	templateParam string
	accessKey string
	accessSecret string
}

func NewAliYunSms(phoneNumber string, templateCode string, templateParam map[string]string) SmsSender {
	conf,err := config.GetAppConfig()
	if err != nil {
		log.Fatalln(err)
	}

	sender := new(AliYunSmsSender)

	sender.accessKey = conf.App.Sms.AccessId
	sender.accessSecret = conf.App.Sms.AccessKey
	sender.phoneNumber = phoneNumber
	sender.signName = conf.App.Sms.SignName
	sender.templateCode = templateCode
	param, err := json.Marshal(templateParam)
	if err != nil {
		log.Println("NewAliYunSms: line: 42 message:"+err.Error())
		log.Fatalln(err)
	}
	sender.templateParam = string(param)

	return sender
}

func (aliYun *AliYunSmsSender) Send() (bool, error) {
	client, err := sdk.NewClientWithAccessKey("default", aliYun.accessKey, aliYun.accessSecret)
	if err != nil {
		return false, err
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https"
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"

	params := make(map[string]string)
	params["PhoneNumbers"] = aliYun.phoneNumber
	params["SignName"] = aliYun.signName
	params["TemplateCode"] = aliYun.templateCode
	params["TemplateParam"] = aliYun.templateParam

	request.QueryParams = params

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		log.Println("NewAliYunSms->Send: line: 74 message:"+err.Error())
		return false, err
	}

	r := response.GetHttpContentString()
	log.Println("send")
	log.Println(r)

	return true, nil
}

// 阿里大于短信
type AliTopSmsSender struct {
	phoneNumber string
	signName string
	templateCode string
	templateParam string

	accessKey string
	accessSecret string
}

func NewAliTopSms(phoneNumber string, templateCode string, templateParam map[string]string) SmsSender {
	conf,err := config.GetAppConfig()
	if err != nil {
		log.Fatalln(err)
	}

	sender := new(AliTopSmsSender)
	sender.accessKey = conf.App.Sms.AliTopAppKey
	sender.accessSecret = conf.App.Sms.AliTopSecretKey
	sender.phoneNumber = phoneNumber
	sender.signName = conf.App.Sms.SignName
	sender.templateCode = templateCode
	param, err := json.Marshal(templateParam)
	if err != nil {
		log.Println("NewAliTopSms: line: 110 message:"+err.Error())
		log.Fatalln(err)
	}
	sender.templateParam = string(param)

	return sender
}

func (top *AliTopSmsSender) Send() (bool, error) {
	client := NewTopClient(top.accessKey, top.accessSecret)
	req := NewAliBaBaAliQinFcSmsNumSendRequest()

	req.SmsFreeSignName = top.signName
	req.RecNum = top.phoneNumber
	req.SmsTemplateCode = top.templateCode
	req.SmsParam = top.templateParam

	response, err := client.Execute(req)
	if err != nil {
		log.Println("NewAliTopSms: line: 129 message:"+err.Error())
		return false, err
	}

	log.Println(response)

	return true, nil
}

// 发送验证码
func SendSmsCode(phoneNumber string, template string) error {
	appConf,err := config.GetAppConfig()
	if err != nil {
		log.Fatalln(err)
	}
	smsConf,err := config.GetSmsConfig()
	if err != nil {
		log.Fatalln(err)
	}

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
	redis,err := database.GetRedis()
	if err != nil {
		return errors.New("redis init error")
	}
	defer redis.Close()

	key := "sms_code_"+template+":"+phoneNumber

	_, _ = redis.Do("DEL", key)

	smsCode := SmsCode{
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

	var sender SmsSender

	switch appConf.App.Sms.SmsSender {
	case "aliyun":
		sender = NewAliYunSms(phoneNumber, templateCode, templateParam)
		break
	case "alitop":
		sender = NewAliTopSms(phoneNumber, templateCode, templateParam)
		break
	default:
		return errors.New("sms storage error:"+appConf.App.Sms.SmsSender)
	}

	_,err = sender.Send()
	if err != nil {
		return err
	}

	return nil
}

// 验证短信验证码
func ValidSmsCode(phoneNumber string, template string, smsCode string) error {
	// 验证码缓存
	redis,err := database.GetRedis()
	if err != nil {
		return errors.New("redis init error")
	}
	defer redis.Close()

	key := "sms_code_"+template+":"+phoneNumber

	value, err := baseRedis.Values(redis.Do("HGETALL", key))
	if err != nil {
		log.Println(err)
		return errors.New("验证码已过期，请重新获取")
	}

	code := &SmsCode{}
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
