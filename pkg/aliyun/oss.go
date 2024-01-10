package aliyun

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gin-gorm-base/pkg/constants"
	"gin-gorm-base/pkg/gredis"
	"gin-gorm-base/pkg/logging"
	"gin-gorm-base/pkg/util"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"gin-gorm-base/pkg/setting"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const (
	// StsSignVersion sts sign version
	StsSignVersion = "1.0"
	// StsAPIVersion sts api version
	StsAPIVersion = "2015-04-01"
	// StsHost sts host
	StsHost = "https://sts.aliyuncs.com/"
	// TimeFormat time fomrat
	TimeFormat = "2006-01-02T15:04:05Z"
	// RespBodyFormat  respone body format
	RespBodyFormat = "JSON"
	// PercentEncode '/'
	PercentEncode = "%2F"
	// HTTPGet http get method
	HTTPGet = "GET"
)

var RedisKey = constants.RedisKeyOssSts

//AliOssVersion sdk 版本
func AliOssVersion() (version string) {
	return oss.Version
}

//InitServer 初始化oss服务
func InitServer() (client *oss.Client, err error) {

	// 请填写您的AccessKeyId。
	var accessKeyId string = setting.OssSetting.AccessKeyId

	// 请填写您的AccessKeySecret。
	var accessKeySecret string = setting.OssSetting.AccessKeySecret

	var ossEndpoint string = setting.OssSetting.OssEndpoint

	client, err = oss.New(ossEndpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return
	}
	return
}

func DeleteFile(key string) {
	// 创建OSSClient实例。
	client, err := InitServer()
	bucket, err := client.Bucket(setting.OssSetting.OssDefaultBucket)
	if err != nil {
		return
	}
	// 上传文件。
	err = bucket.DeleteObject(key)
	if err != nil {
		logging.Error("Delete file", key, "from oss error", err.Error())
	}
	return
}

// UploadFile 上传本地文件
func UploadFile(localFile string, uploadFile string) (resultFile string, err error) {
	resultFile = ""
	// 创建OSSClient实例。
	client, err := InitServer()
	uploadDir := "image/" + time.Now().Format(constants.DateLayout)

	// <yourObjectName>上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	uploadFile = strings.Trim(uploadFile, "/")
	objectName := fmt.Sprintf("%s/%s", uploadDir, uploadFile) //完整的oss路径
	// <yourLocalFileName>由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt。
	localFileName := localFile
	// 获取存储空间。
	bucket, err := client.Bucket(setting.OssSetting.OssDefaultBucket)
	if err != nil {
		return
	}
	// 上传文件。
	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		return
	}
	resultFile = objectName
	return
}

// UploadFileByte 上传Byte数组文件
func UploadFileByte(key string, uploadByte []byte) (resultfile string, err error) {
	resultfile = ""
	// 创建OSSClient实例。
	client, err := InitServer()
	uploadDir := "image/" + time.Now().Format(constants.DateLayout)

	// <yourObjectName>上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	key = strings.Trim(key, "/")
	objectName := fmt.Sprintf("%s/%s", uploadDir, key) //完整的oss路径
	// <yourLocalFileName>由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt。

	// 获取存储空间。
	bucket, err := client.Bucket(setting.OssSetting.OssDefaultBucket)
	if err != nil {
		return
	}

	// 上传Byte数组。
	err = bucket.PutObject(objectName, bytes.NewReader(uploadByte))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	resultfile = objectName
	return
}

type StsTokenResp struct {
	Region          string `json:"region"`
	AccessKeyId     string `json:"access_key_id"`
	Endpoint        string `json:"endpoint"`
	ResourcesUrl    string `json:"resources_url"`
	AccessKeySecret string `json:"access_key_secret"`
	Bucket          string `json:"bucket"`
	SecurityToken   string `json:"security_token"`
	Path            string `json:"path"`
	ExpiresIn       int    `json:"expires_in"`
}

func GetStsToken() *StsTokenResp {
	stsTokenResp := &StsTokenResp{}

	uploadDir := "image/" + time.Now().Format(constants.DateLayout) + "/"
	nowTime := int(time.Now().Unix())
	stsTokenResp.AccessKeyId = setting.OssSetting.AccessKeyId
	stsTokenResp.Region = setting.OssSetting.OssRegion
	stsTokenResp.Endpoint = setting.OssSetting.OssEndpoint
	stsTokenResp.ResourcesUrl = setting.OssSetting.OssDomain
	stsTokenResp.AccessKeySecret = setting.OssSetting.AccessKeySecret
	stsTokenResp.Bucket = setting.OssSetting.OssDefaultBucket
	stsTokenResp.Path = uploadDir

	var saveStsToken *SaveStsToken
	gredis.Delete(RedisKey)
	if gredis.Exists(RedisKey) {
		saveStsTokenJson, _ := gredis.Get(RedisKey)
		err := json.Unmarshal(saveStsTokenJson, &saveStsToken)
		if err != nil {
			logging.Error("GetStsToken saveStsTokenJson json.Unmarshal", "err", err)
		}
		if saveStsToken.ExpiresIn > 0 && nowTime > saveStsToken.ExpiresIn {
			saveStsToken = GetStsTokenFromAil()
		}
	} else {
		saveStsToken = GetStsTokenFromAil()
	}
	if saveStsToken != nil {
		stsTokenResp.SecurityToken = saveStsToken.SecurityToken
		stsTokenResp.ExpiresIn = saveStsToken.ExpiresIn - nowTime
	}
	return stsTokenResp
}

type SaveStsToken struct {
	SecurityToken   string `json:"security_token"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	ExpiresIn       int    `json:"expires_in"`
}

type Credentials struct {
	AccessKeyId     string
	AccessKeySecret string
	Expiration      time.Time
	SecurityToken   string
}

// AssumedRoleUser the user to AssumedRole
type AssumedRoleUser struct {
	Arn           string
	AssumedRoleId string
}

// Response the response of AssumeRole
type Response struct {
	Credentials     Credentials
	AssumedRoleUser AssumedRoleUser
	RequestId       string
}

// Client sts client
type Client struct {
	AccessKeyId     string
	AccessKeySecret string
	RoleArn         string
	SessionName     string
}

// ServiceError sts service error
type ServiceError struct {
	Code       string
	Message    string
	RequestId  string
	HostId     string
	RawMessage string
	StatusCode int
}

// Error implement interface error
func (e *ServiceError) Error() string {
	return fmt.Sprintf("oss: service returned error: StatusCode=%d, ErrorCode=%s, ErrorMessage=%s, RequestId=%s",
		e.StatusCode, e.Code, e.Message, e.RequestId)
}

// NewClient New STS Client
func NewClient(accessKeyId, accessKeySecret, roleArn, sessionName string) *Client {
	return &Client{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		RoleArn:         roleArn,
		SessionName:     sessionName,
	}
}

// AssumeRole assume role
func (c *Client) AssumeRole(expiredTime uint) (*Response, error) {
	url, err := c.generateSignedURL(expiredTime)
	if err != nil {
		return nil, err
	}

	body, status, err := c.sendRequest(url)
	if err != nil {
		return nil, err
	}

	return c.handleResponse(body, status)
}

// Private function
func (c *Client) generateSignedURL(expiredTime uint) (string, error) {
	signatureNonce := strconv.Itoa(int(time.Now().Unix())) + strconv.Itoa(util.GetRandInt(999, 1))
	queryStr := "SignatureVersion=" + StsSignVersion
	queryStr += "&Format=" + RespBodyFormat
	queryStr += "&Timestamp=" + url.QueryEscape(time.Now().UTC().Format(TimeFormat))
	queryStr += "&RoleArn=" + url.QueryEscape(c.RoleArn)
	queryStr += "&RoleSessionName=" + c.SessionName
	queryStr += "&AccessKeyId=" + c.AccessKeyId
	queryStr += "&SignatureMethod=HMAC-SHA1"
	queryStr += "&Version=" + StsAPIVersion
	queryStr += "&Action=AssumeRole"
	queryStr += "&SignatureNonce=" + signatureNonce
	queryStr += "&DurationSeconds=" + strconv.FormatUint((uint64)(expiredTime), 10)

	// Sort query string
	queryParams, err := url.ParseQuery(queryStr)
	if err != nil {
		return "", err
	}
	result := queryParams.Encode()

	strToSign := HTTPGet + "&" + PercentEncode + "&" + url.QueryEscape(result)

	// Generate signature
	hashSign := hmac.New(sha1.New, []byte(c.AccessKeySecret+"&"))
	hashSign.Write([]byte(strToSign))
	signature := base64.StdEncoding.EncodeToString(hashSign.Sum(nil))

	// Build url
	assumeURL := StsHost + "?" + queryStr + "&Signature=" + url.QueryEscape(signature)

	return assumeURL, nil
}

func (c *Client) sendRequest(url string) ([]byte, int, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

func (c *Client) handleResponse(responseBody []byte, statusCode int) (*Response, error) {
	if statusCode != http.StatusOK {
		se := ServiceError{StatusCode: statusCode, RawMessage: string(responseBody)}
		err := json.Unmarshal(responseBody, &se)
		if err != nil {
			return nil, err
		}
		return nil, &se
	}

	resp := Response{}
	err := json.Unmarshal(responseBody, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func GetStsTokenFromAil() *SaveStsToken {
	nowTime := time.Now().Unix()
	client := NewClient(setting.OssSetting.AccessKeyId, setting.OssSetting.AccessKeySecret, setting.OssSetting.OssRoleARN, setting.OssSetting.OssRoleSessionName)
	result, err := client.AssumeRole(43200)
	logging.Debug("")
	if err != nil {
		logging.Error("GetStsTokenFromAil AssumeRole result", "err", err)
		return nil
	}
	saveStsToken := &SaveStsToken{
		AccessKeyId:     result.Credentials.AccessKeyId,
		AccessKeySecret: result.Credentials.AccessKeySecret,
		SecurityToken:   result.Credentials.SecurityToken,
		ExpiresIn:       int(nowTime) + 7000,
	}
	gredis.Set(RedisKey, saveStsToken, constants.RedisTimeOneHour)
	return saveStsToken
}
