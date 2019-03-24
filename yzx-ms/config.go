package yzxms

import "reflect"

const SMSDOMAIN = "https://open.ucpaas.com/ol/sms/sendsms"

type Config interface {
	init()
}

type BaseConfig struct {
	AccountsId string `json:"sid"`
	Token      string `json:"token"`
	AppId      string `json:"appid"`
	Uid        string `json:"uid"`
}

//短信验证码配置
type SmsConfig struct {
	BaseConfig
	PhoneNumbers  string `json:"mobile"`
	TemplateParam string `json:"param"`
	TemplateCode  string `json:"templateid"`
}

var BaseParams = BaseConfig{
	AccountsId: "d5fd8f9118b9cd6cf9cfc764d111894e",
	Token:      "616bee00db429d890ce218ae1555b83c",
	AppId:      "5ad60d2827a44b23b4c53936143c2353",
	Uid:        "2d92c6132139467b989d087c84a365d8",
}

var SmsParams SmsConfig

//短信参数初始化函数
func (s SmsConfig) init() {
	t := reflect.TypeOf(BaseParams)
	value := reflect.ValueOf(BaseParams)
	smsParamsTable := reflect.ValueOf(&SmsParams).Elem()
	for k := 0; k < t.NumField(); k++ {

		smsParamsTable.FieldByName(t.Field(k).Name).SetString(value.Field(k).String())
	}
}

func init() {
	SmsParams.init()
}
