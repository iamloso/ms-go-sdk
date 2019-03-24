package aliyunms

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"sms-service/lib"
	"sms-service/log"
	"sort"
	"strings"
	"time"
)

type msgSend interface {
	send() error
}

//发送短信接口， 响应信息
type Response struct {
	RequestId string `json:"RequestId"`
	Code      string `json:"Code,omitempty"`
	Message   string `json:"Message,omitempty"`
	BizId     string `json:"BizId"`
}

type Sms struct {
	paras SmsConfig
}

type Vms struct {
	paras VmsConfig
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

	Smslog.InfoData("aliyunms.sms.send:(开启跟踪)短信发送sdk", params)

	SmsParams.PhoneNumbers = phoneNums
	SmsParams.TemplateParam = templateParam
	SmsParams.TemplateCode = templateCode
	SmsParams.SignatureNonce = Smslib.GenXid()
	SmsParams.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05Z")

	smsParams := baseParmas(SmsParams)
	baseParams := baseParmas(SmsParams.BaseConfig)
	for key, val := range baseParams {
		smsParams[key] = val
	}
	AccessSecret := smsParams["AccessSecret"]
	delete(smsParams, "BaseConfig")
	delete(smsParams, "AccessSecret")

	signUrl := createSignUrl(smsParams, SMSDOMAIN, AccessSecret)

	resp, err := http.Get(signUrl)
	if err != nil {
		params["system_error"] = err
		Smslog.ErrorData("aliyunms.sms.send:(系统错误)阿里通道发送短信失败!", params)
		return result, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		params["system_error"] = err
		Smslog.ErrorData("aliyunms.sms.send:(系统错误)阿里通道接受返回数据失败!", params)
		return result, err
	}
	result = Response{}

	if err := json.Unmarshal(body, &result); err != nil {
		params["system_error"] = err
		Smslog.ErrorData("aliyunms.sms.send:(系统错误)阿里通道解析返回数据失败! ", params)
		return result, err
	}

	if result.Code == "SignatureNonceUsed" {
		return s.Send(phoneNums, templateParam, templateCode)
	} else if result.Code != "OK" {
		params["system_error"] = result.Code + ":" + result.Message
		Smslog.ErrorData("aliyunms.sms.send:(系统错误)阿里通道发送短信失败!", params)
		return result, errors.New(result.Code)
	}

	return result, nil
}

/**
 * 创建请求地址签名
 */
func createSignUrl(params map[string]string, domain string, accessSecret string) string {
	var keys []string

	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var sortQueryString string

	for _, v := range keys {
		sortQueryString = fmt.Sprintf("%s&%s=%s", sortQueryString, replace(v), replace(params[v]))
	}

	querySign := fmt.Sprintf("GET&%s&%s", replace("/"), replace(sortQueryString[1:]))

	mac := hmac.New(sha1.New, []byte(fmt.Sprintf("%s&", accessSecret)))
	mac.Write([]byte(querySign))
	sign := replace(base64.StdEncoding.EncodeToString(mac.Sum(nil)))
	signUrl := fmt.Sprintf("http://%s/?Signature=%s%s", domain, sign, sortQueryString)

	return signUrl
}

//基础参数 返回map数据结构
func baseParmas(obj interface{}) map[string]string {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).String()
	}
	return data
}

/**
 * phoneNums string 手机号
 * templateParam string 模板变量, json格式键值对
 * templateCode  string 模板ID
 */
func (v *Vms) Send(phoneNums, templateParam, templateCode string) (result Response, err error) {
	var params = map[string]interface{}{"phone": phoneNums, "templateParam": templateParam, "templateCode": templateCode}

	Smslog.InfoData("aliyunms.vms.send:(开启跟踪)语言短信发送sdk", params)

	VmsParams.CalledNumber = phoneNums
	VmsParams.TtsParam = templateParam
	VmsParams.TtsCode = templateCode
	VmsParams.CalledShowNumber = "01086397932"
	VmsParams.SignatureNonce = fmt.Sprintf("%d", rand.Int63())
	VmsParams.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05Z")

	vmsParams := baseParmas(VmsParams)
	baseParams := baseParmas(VmsParams.BaseConfig)
	for key, val := range baseParams {
		vmsParams[key] = val
	}
	AccessSecret := vmsParams["AccessSecret"]
	delete(vmsParams, "BaseConfig")
	delete(vmsParams, "AccessSecret")

	signUrl := createSignUrl(vmsParams, VMSDOMAIN, AccessSecret)

	resp, err := http.Get(signUrl)
	if err != nil {
		params["system_error"] = err
		Smslog.ErrorData("aliyunms.vms.send:(系统错误)阿里语音通道发送短信失败!", params)
		return result, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		params["system_error"] = err
		Smslog.ErrorData("aliyunms.vms.send:(系统错误)阿里语音通道接受返回数据失败!", params)
		return result, err
	}

	result = Response{}

	if err := json.Unmarshal(body, &result); err != nil {
		params["system_error"] = err
		Smslog.ErrorData("aliyunms.vms.send:(系统错误)阿里语音通道解析返回数据失败!", params)
		return result, err
	}

	if result.Code == "SignatureNonceUsed" {
		return v.Send(phoneNums, templateParam, templateCode)
	} else if result.Code != "OK" {
		params["system_error"] = result.Code + ":" + result.Message
		Smslog.ErrorData("aliyunms.vms.send:(系统错误)阿里语音通道发送短信失败!", params)
		return result, errors.New(result.Code)
	}

	return result, nil
}
