package yzxms

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"sms-service/log"
	"strings"
)

type msgSend interface {
	send() error
}

//发送短信接口， 响应信息
type Response struct {
	RequestId string `json:"RequestId"`
	Code      string `json:"code"`
	Message   string `json:"msg"`
	BizId     string `json:"BizId"`
}

type Sms struct {
	paras SmsConfig
}

func replace(in string) string {
	rep := strings.NewReplacer("+", "%20", "*", "%2A", "%7E", "~")
	return rep.Replace(url.QueryEscape(in))
}

/**
 * phoneNums string 手机号(支持批量，英文逗号分割)
 * templateParam string 模板变量, json格式键值对
 * templateCode  string 模板ID
 */
func (s *Sms) Send(phoneNums, templateParam, templateCode string) (result Response, err error) {
	var params = map[string]interface{}{"phone": phoneNums, "templateParam": templateParam, "templateCode": templateCode}

	Smslog.InfoData("yzx sms send:(开启跟踪)云之讯通道短信发送sdk", params)

	SmsParams.PhoneNumbers = phoneNums
	SmsParams.TemplateParam = templateParam
	SmsParams.TemplateCode = templateCode

	smsParams := baseParmas(SmsParams)

	baseParams := baseParmas(SmsParams.BaseConfig)

	for key, val := range smsParams {
		if key != "" {
			baseParams[key] = val
		}
	}

	jsonData, err := json.Marshal(baseParams)

	if err != nil {
		params["system_error"] = err
		Smslog.ErrorData("yzx sms send:(系统错误)参数转json失败！", params)
		return result, err
	}

	postReq, err := http.NewRequest("POST", SMSDOMAIN, strings.NewReader(string(jsonData)))

	if err != nil {
		params["system_error"] = err
		Smslog.ErrorData("yzx sms send:(系统错误)云之讯通道发送短信建立请求失败!", params)
		return result, err
	}

	postReq.Header.Set("Content-Type", "application/json; encoding=utf-8; charset=utf-8")
	postReq.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(postReq)

	if err != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		params["system_error"] = err
		Smslog.ErrorData("yzx sms send:(系统错误)云之讯通道发送短信失败!", params)
		return result, err
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			params["system_error"] = err
			Smslog.ErrorData("yzx sms send:(系统错误)云之讯通道接受返回数据失败!", params)
			return result, err
		}

		result = Response{}
		if err := json.Unmarshal(body, &result); err != nil {
			params["system_error"] = err
			Smslog.ErrorData("yzx sms send:(系统错误)云之讯通道解析返回数据失败!", params)
			return result, err
		}

		if result.Code != "000000" {
			params["system_error"] = err
			Smslog.ErrorData("yzx sms send:(系统错误)云之讯通道发送短信失败!", params)
			return result, errors.New(result.Code)
		}
	}

	return result, nil
}

//基础参数 返回map数据结构
func baseParmas(obj interface{}) map[string]string {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Tag.Get("json")] = v.Field(i).String()
	}
	return data
}
