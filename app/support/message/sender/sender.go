package sender

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"live-service/app/config"
	"log"
)

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
	conf := config.GetAppConfig()

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	sender := new(AliYunSmsSender)

	sender.accessKey = conf.App.Sms.AccessId
	sender.accessSecret = conf.App.Sms.AccessKey
	sender.phoneNumber = phoneNumber
	sender.signName = conf.App.Sms.SignName
	sender.templateCode = templateCode
	param, err := json.Marshal(templateParam)
	if err != nil {
		log.Println("NewAliYunSms: message:"+err.Error())
		panic(err)
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
		log.Println("NewAliYunSms->Send: message:"+err.Error())
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
	conf := config.GetAppConfig()

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	sender := new(AliTopSmsSender)
	sender.accessKey = conf.App.Sms.AliTopAppKey
	sender.accessSecret = conf.App.Sms.AliTopSecretKey
	sender.phoneNumber = phoneNumber
	sender.signName = conf.App.Sms.SignName
	sender.templateCode = templateCode
	param, err := json.Marshal(templateParam)
	if err != nil {
		log.Println("NewAliTopSms: line: 110 message:"+err.Error())
		panic(err)
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
		log.Println("NewAliTopSms: message:"+err.Error())
		return false, err
	}

	log.Println(response)

	return true, nil
}
