//阿里云 云通信 go sdk 系统参数配置
//支持短信参数配置 语音参数配置
package aliyunms

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

const SMSDOMAIN = "dysmsapi.aliyuncs.com"
const VMSDOMAIN = "dyvmsapi.aliyuncs.com"

type Config interface {
	init()
}

type BaseConfig struct {
	AccessKeyId, AccessSecret, SignatureMethod, SignatureNonce string
	SignatureVersion, Timestamp, Format                        string
	Action, Version, RegionId, SignName                        string
}

//短信验证码配置
type SmsConfig struct {
	BaseConfig
	PhoneNumbers  string
	TemplateParam string
	TemplateCode  string
}

//语音验证码配置
type VmsConfig struct {
	BaseConfig
	CalledShowNumber string
	CalledNumber     string
	TtsParam         string
	TtsCode          string
	OutId            string
}

var BaseParams = BaseConfig{
	AccessKeyId:      "LTAIMz4fpJk26Xmg",
	AccessSecret:     "fs4swygV9gFdoXE0esAOTVuiu9mkvW",
	SignatureMethod:  "HMAC-SHA1",
	SignatureNonce:   fmt.Sprintf("%d", rand.Int63()),
	SignatureVersion: "1.0",
	Timestamp:        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
	Format:           "JSON",
	Action:           "SendSms",
	Version:          "2017-05-25",
	RegionId:         "cn-hangzhou",
	SignName:         "狸米数学",
}

var SmsParams SmsConfig

var VmsParams VmsConfig

//短信参数初始化函数
func (s SmsConfig) init() {
	t := reflect.TypeOf(BaseParams)
	value := reflect.ValueOf(BaseParams)
	smsParamsTable := reflect.ValueOf(&SmsParams).Elem()
	for k := 0; k < t.NumField(); k++ {

		smsParamsTable.FieldByName(t.Field(k).Name).SetString(value.Field(k).String())
	}
}

//语音参数初始化函数
func (v VmsConfig) init() {
	t := reflect.TypeOf(BaseParams)
	value := reflect.ValueOf(BaseParams)
	vmsParamsTable := reflect.ValueOf(&VmsParams).Elem()
	for k := 0; k < t.NumField(); k++ {

		vmsParamsTable.FieldByName(t.Field(k).Name).SetString(value.Field(k).String())
	}
	VmsParams.AccessKeyId = "LTAImgiMYlfO4QNG"
	VmsParams.AccessSecret = "mgNV0o8pr4ff0kChEeH3Io1HfxetC6"
	VmsParams.Action = "SingleCallByTts"
}

func init() {
	VmsParams.init()
	SmsParams.init()
}
