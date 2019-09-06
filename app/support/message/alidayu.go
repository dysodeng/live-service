package message

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type TopClient struct {
	session string
	timestamp string
	format string
	v string
	partnerId string
	simplify bool

	appKey string
	secretKey string
	targetAppKey string
	gatewayUrl string
}

func NewTopClient(appKey string, secretKey string) *TopClient {

	top := new(TopClient)
	top.appKey = appKey
	top.secretKey = secretKey
	top.gatewayUrl = "https://eco.taobao.com/router/rest"
	top.format = "json"
	top.v = "2.0"

	return top
}

func (client *TopClient) Execute(req AliBaBaRequest) (map[string]interface{}, error) {
	params, err := client.buildParams(req)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.PostForm(client.gatewayUrl, params)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (client *TopClient) buildParams(req AliBaBaRequest) (url.Values, error) {
	if err := req.ParamsIsValid(); err != nil {
		return nil, err
	}
	if len(req.GetMethodName()) == 0 {
		return nil, errors.New("method is required")
	}
	if len(client.appKey) == 0 {
		return nil, errors.New("app_key is required")
	}

	paramsCommon := make(map[string]string)
	paramsCommon["method"] = req.GetMethodName()
	paramsCommon["app_key"] = client.appKey
	paramsCommon["target_app_key"] = client.targetAppKey
	paramsCommon["sign_method"] = "md5"
	paramsCommon["session"] = client.session
	paramsCommon["timestamp"] = client.timeStamp()
	paramsCommon["format"] = client.format
	paramsCommon["v"] = client.v
	if client.simplify {
		paramsCommon["simplify"] = "1"
	} else {
		paramsCommon["simplify"] = "0"
	}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	paramsApi := make(map[string]string)
	if err := json.Unmarshal(data, &paramsApi); err != nil {
		return nil, err
	}

	paramsCommon["sign"] = client.sign(paramsApi, paramsCommon)
	params := make(url.Values)
	for key, value := range paramsApi {
		params.Set(key, value)
	}
	for key, value := range paramsCommon {
		params.Set(key, value)
	}

	return params, nil
}

func (client *TopClient) timeStamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (client *TopClient) sign(paramsApi map[string]string, paramsCommon map[string]string) string {
	ks := make([]string, 0)
	params := make(map[string]string)

	for key, value := range paramsApi {
		ks = append(ks, key)
		params[key] = value
	}

	for key, value := range paramsCommon {
		ks = append(ks, key)
		params[key] = value
	}

	sort.Strings(ks)
	str := ""
	for _, k := range ks {
		str = fmt.Sprintf("%v%v%v", str, k, params[k])
	}

	str = fmt.Sprintf("%v%v%v", client.secretKey, str, client.secretKey)
	hash := md5.Sum([]byte(str))
	sign := fmt.Sprintf("%x", hash)

	return strings.ToUpper(sign)
}

type AliBaBaRequest interface {
	GetMethodName() string
	ParamsIsValid() error
}

type AliBaBaAliQinFcSmsNumSendRequest struct {
	Extend string	`json:"extend"`
	RecNum string	`json:"rec_num"`
	SmsFreeSignName string	`json:"sms_free_sign_name"`
	SmsParam string	`json:"sms_param"`
	SmsTemplateCode string	`json:"sms_template_code"`
	SmsType string	`json:"sms_type"`
}

func NewAliBaBaAliQinFcSmsNumSendRequest() *AliBaBaAliQinFcSmsNumSendRequest {
	req := new(AliBaBaAliQinFcSmsNumSendRequest)
	req.SmsType = "normal"
	return req
}

func (req *AliBaBaAliQinFcSmsNumSendRequest) GetMethodName() string {
	return "alibaba.aliqin.fc.sms.num.send"
}

func (req *AliBaBaAliQinFcSmsNumSendRequest) ParamsIsValid() error {
	if len(req.SmsType) == 0 {
		return errors.New("sms_type is required")
	}
	if len(req.SmsFreeSignName) == 0 {
		return errors.New("sms_free_sign_name is required")
	}
	if len(req.RecNum) == 0 {
		return errors.New("rec_num is required")
	}
	if len(req.SmsTemplateCode) == 0 {
		return errors.New("sms_template_code is required")
	}

	return nil
}