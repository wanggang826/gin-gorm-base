package douyin

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"gin-gorm-base/pkg/constants"
	"gin-gorm-base/pkg/gredis"
	"gin-gorm-base/pkg/logging"
	"gin-gorm-base/pkg/setting"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

type DyConfig struct {
	Appid           string `json:"appid"`
	Secret          string `json:"secret"`
	Salt            string `json:"salt"`
	Token           string `json:"token"`
	NotifyUrl       string `json:"notify_url"`
	RefundNotifyUrl string `json:"refund_notify_url"`
}

type DyService struct {
	*DyConfig
}

func NewDyService() *DyService {
	dyService := DyService{GetDyConfig()}
	return &dyService
}

// GetDyConfig 初始抖音配置
func GetDyConfig() *DyConfig {
	return &DyConfig{
		Appid:           setting.DouYinSetting.Appid,
		Secret:          setting.DouYinSetting.Secret,
		Salt:            setting.DouYinSetting.Salt,
		Token:           setting.DouYinSetting.Token,
		NotifyUrl:       setting.DouYinSetting.NotifyUrl,
		RefundNotifyUrl: setting.DouYinSetting.NotifyUrl,
	}
}

// GetAccessToken 获取AccessToken
func (dy *DyService) GetAccessToken() string {
	accessToken := ""
	gredis.Delete(constants.RedisAccessTokenKey)
	if gredis.Exists(constants.RedisAccessTokenKey) {
		accessToken, _ = gredis.GetString(constants.RedisAccessTokenKey)
	} else {
		accessTokenResp, err := dy.GetClientTokenFromDy()
		fmt.Println(accessTokenResp)
		fmt.Println(err)
		//if err != nil {
		//	return ""
		//}
		//if accessTokenResp != nil && accessTokenResp.AccessToken != "" {
		//	gredis.SetString(constants.RedisAccessTokenKey, accessTokenResp, accessTokenResp.ExpiresIn)
		//	accessToken = accessTokenResp
		//}
	}
	return accessToken
}

type AccessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// GetAccessTokenFromDy 从抖音平台获取client_token （新的AccessToken 暂未用）
func (dy *DyService) GetAccessTokenFromDy() (*AccessTokenResp, error) {
	url := "https://developer.toutiao.com/api/apps/v2/token"
	params := make(map[string]interface{})
	params["appid"] = dy.Appid
	params["secret"] = dy.Secret
	params["grant_type"] = "client_credential"
	jsonData, err := json.Marshal(params)
	logging.Debug("GetAccessTokenFromDy jsonData", string(jsonData))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logging.Error("GetAccessTokenFromDy  HTTPPost ", "err", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	type Result struct {
		ErrNo   int              `json:"err_no"`
		Data    *AccessTokenResp `json:"data"`
		ErrTips string           `json:"err_tips"`
	}

	result := &Result{}
	err = json.Unmarshal(body, result)
	if err != nil {
		logging.Error("GetAccessTokenFromDy  Result Unmarshal ", "err", err)
		return nil, err
	}
	if result.ErrNo != 0 {
		logging.Error("Code2Session  Result fail", "result.ErrTips", result.ErrTips)
		return nil, errors.New(result.ErrTips)
	}
	logging.Debug("GetAccessTokenFromDy ", "Result", err)
	return result.Data, nil
}

type Code2SessionResult struct {
	SessionKey      string `json:"session_key"`
	Openid          string `json:"openid"`
	AnonymousOpenid string `json:"anonymous_openid"`
	Unionid         string `json:"unionid"`
}

// Code2Session 登录凭证校验
func (dy *DyService) Code2Session(code, anonymousCode string) (*Code2SessionResult, error) {
	url := "https://developer.toutiao.com/api/apps/v2/jscode2session"
	params := make(map[string]interface{})
	params["appid"] = dy.Appid
	params["secret"] = dy.Secret
	params["code"] = code
	params["anonymous_code"] = anonymousCode
	jsonData, err := json.Marshal(params)
	logging.Debug("Code2Session jsonData", string(jsonData))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logging.Error("Code2Session  HTTPPost err", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	type Result struct {
		ErrNo   int                 `json:"err_no"`
		Data    *Code2SessionResult `json:"data"`
		ErrTips string              `json:"err_tips"`
	}

	result := &Result{}
	err = json.Unmarshal(body, result)
	if err != nil {
		logging.Error("Code2Session  Result Unmarshal err", err)
		return nil, err
	}
	if result.ErrNo != 0 {
		logging.Error("Code2Session  Result fail", "result.ErrTips", result.ErrTips)
		return nil, errors.New(result.ErrTips)
	}
	return result.Data, nil
}

// GetClientTokenFromDy 从抖音平台获取client_token （新的AccessToken 暂未用）
func (dy *DyService) GetClientTokenFromDy() (string, error) {
	url := "https://open.douyin.com/oauth/client_token/"
	params := make(map[string]interface{})
	params["client_key"] = dy.Appid
	params["client_secret"] = dy.Secret
	params["grant_type"] = "client_credential"
	jsonData, err := json.Marshal(params)
	logging.Debug("GetClientTokenFromDy jsonData", string(jsonData))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Content-Type", "multipart/form-data")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logging.Error("GetClientTokenFromDy  HTTPPost err", err)
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	type ClientTokenData struct {
		ExpiresIn   int    `json:"expires_in"`
		AccessToken string `json:"access_token"`
		Description string `json:"description"`
		ErrorCode   int    `json:"error_code"`
	}
	type Result struct {
		Data    ClientTokenData `json:"data"`
		Message string          `json:"err_tips"`
	}

	result := &Result{}
	err = json.Unmarshal(body, result)
	if err != nil {
		logging.Error("GetClientTokenFromDy  Result Unmarshal err", err)
		return "", err
	}
	fmt.Println(result)
	return result.Data.AccessToken, nil
}

type PayRequestParams struct {
	OutOrderNo  string `json:"out_order_no"` //开发者侧的订单号。 只能是数字、大小写字母_-*且在同一个app_id下唯一
	TotalAmount int    `json:"total_amount"` //支付价格。 单位为[分]
	Subject     string `json:"subject"`      //商品描述
	Body        string `json:"body"`         //商品详情
	ValidTime   int    `json:"valid_time"`   //订单过期时间(秒)。最小5分钟，最大2天，小于5分钟会被置为5分钟，大于2天会被置为2天
}

type PayResponseParams struct {
	OrderId    string `json:"order_id"`
	OrderToken string `json:"order_token"`
}

// CreateOrder 预下单
func (dy *DyService) CreateOrder(reqParams *PayRequestParams) (*PayResponseParams, error) {
	url := "https://developer.toutiao.com/api/apps/ecpay/v1/create_order"
	params := make(map[string]interface{})
	params["app_id"] = dy.Appid //小程序APPID
	params["out_order_no"] = reqParams.OutOrderNo
	params["total_amount"] = reqParams.TotalAmount
	params["subject"] = reqParams.Subject
	params["body"] = reqParams.Body
	params["valid_time"] = reqParams.ValidTime
	params["valid_time"] = reqParams.ValidTime
	params["thirdparty_id"] = ""
	params["notify_url"] = dy.NotifyUrl
	params["sign"] = dy.getSign(params) //签名
	jsonData, err := json.Marshal(params)
	logging.Debug("DyService Pay", "jsonData", string(jsonData))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logging.Error("DyService Pay  HTTPPost", "err", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	type Result struct {
		ErrNo   int                `json:"err_no"`
		Data    *PayResponseParams `json:"data"`
		ErrTips string             `json:"err_tips"`
	}

	result := &Result{}
	logging.Debug("DyService CreateOrder", "resp body", string(body))
	err = json.Unmarshal(body, result)
	if err != nil {
		logging.Error("CreateOrder  Result Unmarshal", "err", err)
		return nil, err
	}
	if result.ErrNo != 0 {
		logging.Error("CreateOrder  Result fail", "result.ErrTips", result.ErrTips)
		return nil, errors.New(result.ErrTips)
	}
	logging.Debug("DyService CreateOrder", "result", result)
	return result.Data, nil
}

type RefundRequestParams struct {
	OutOrderNo   string `json:"out_order_no"`
	OutRefundNo  string `json:"out_refund_no"`
	Reason       string `json:"reason"`
	RefundAmount int    `json:"refund_amount"`
}

func (dy *DyService) Refund(reqParams *RefundRequestParams) (string, error) {
	url := "https://developer.toutiao.com/api/apps/ecpay/v1/create_refund"
	params := make(map[string]interface{})
	params["app_id"] = dy.Appid //小程序APPID
	params["out_order_no"] = reqParams.OutOrderNo
	params["out_refund_no"] = reqParams.OutRefundNo
	params["reason"] = reqParams.Reason
	params["refund_amount"] = reqParams.RefundAmount
	params["notify_url"] = dy.RefundNotifyUrl
	params["sign"] = dy.getSign(params) //签名
	jsonData, err := json.Marshal(params)
	logging.Debug("DyService Refund", "jsonData", string(jsonData))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logging.Error("DyService Refund  HTTPPost", "err", err)
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	type Result struct {
		ErrNo    int    `json:"err_no"`
		RefundNo string `json:"refund_no"`
		ErrTips  string `json:"err_tips"`
	}

	result := &Result{}
	logging.Debug("DyService Refund", "resp body", string(body))
	err = json.Unmarshal(body, result)
	if err != nil {
		logging.Error("Refund  Result Unmarshal", "err", err)
		return "", err
	}
	if result.ErrNo != 0 {
		logging.Error("Refund  Result fail", "result.ErrTips", result.ErrTips)
		return "", errors.New(result.ErrTips)
	}
	logging.Debug("DyService Refund", "result", result)
	return "", nil
}

type PayNotifyMsg struct {
	Appid          string `json:"appid"`            //当前交易发起的小程序id
	CpOrderno      string `json:"cp_orderno"`       //开发者侧的订单号
	CpExtra        string `json:"cp_extra"`         //预下单时开发者传入字段
	Way            string `json:"way"`              //way 字段中标识了支付渠道： 1-微信支付，2-支付宝支付，10-抖音支付
	ChannelNo      string `json:"channel_no"`       //支付渠道侧单号(抖音平台请求下游渠道微信或支付宝时传入的单号)
	PaymentOrderNo string `json:"payment_order_no"` //支付渠道侧PC单号，支付页面可见(微信支付宝侧的订单号)
	TotalAmount    int    `json:"total_amount"`     //支付金额，单位为分
	Status         string `json:"status"`           //固定SUCCESS
	ItemId         string `json:"item_id"`          //订单来源视频对应视频 id
	SellerUid      string `json:"seller_uid"`       //该笔交易卖家商户号
	PaidAt         int    `json:"paid_at"`          //支付时间
	OrderId        string `json:"order_id"`         //抖音侧订单号
}

type PayNotifyParams struct {
	Timestamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
	Msg          string `json:"msg"`
	PaymentType  string `json:"type"`
	MsgSignature string `json:"msg_signature"`
}

// getSign 获取签名
func (dy *DyService) getSign(paramsMap map[string]interface{}) string {
	var paramsArr []string
	for k, v := range paramsMap {
		if k == "other_settle_params" {
			continue
		}
		value := strings.TrimSpace(fmt.Sprintf("%v", v))
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") && len(value) > 1 {
			value = value[1 : len(value)-1]
		}
		value = strings.TrimSpace(value)
		if value == "" || value == "null" {
			continue
		}
		switch k {
		// app_id, thirdparty_id, sign 字段用于标识身份，不参与签名
		case "app_id", "thirdparty_id", "sign":
		default:
			paramsArr = append(paramsArr, value)
		}
	}

	paramsArr = append(paramsArr, dy.Salt)
	sort.Strings(paramsArr)
	return fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(paramsArr, "&"))))
}

// CheckNotifySign  支付回调验签
func (dy *DyService) CheckNotifySign(params *PayNotifyParams) bool {
	sortedString := make([]string, 0)
	sortedString = append(sortedString, dy.Token)
	sortedString = append(sortedString, params.Timestamp)
	sortedString = append(sortedString, params.Nonce)
	sortedString = append(sortedString, params.Msg)
	sort.Strings(sortedString)
	h := sha1.New()
	h.Write([]byte(strings.Join(sortedString, "")))
	bs := h.Sum(nil)
	_signature := fmt.Sprintf("%x", bs)
	if _signature != params.MsgSignature {
		logging.Error("CheckNotifySign fail", "_signature", _signature, "params.MsgSignature", params.MsgSignature)
		return false
	}
	return true
}
