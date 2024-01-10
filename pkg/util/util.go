package util

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"gin-gorm-base/pkg/constants"
	"gin-gorm-base/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"gin-gorm-base/pkg/setting"
	"github.com/google/uuid"
)

func Setup() {
	jwtUserSecret = []byte(setting.AppSetting.JwtUserSecret)
	jwtAdminSecret = []byte(setting.AppSetting.JwtAdminSecret)
}

func UUID() string {
	id := uuid.New()
	return strings.Replace(id.String(), "-", "", -1)
}

//转换时间戳
func StrToTime(toTime string) int {
	if toTime == "" {
		return 0
	}

	//时间 to 时间戳
	loc, _ := time.LoadLocation("Asia/Shanghai") //设置时区
	tt, _ := time.ParseInLocation(constants.TimeLayout, toTime, loc)
	return int(tt.Unix())
}

func DateToTime(toTime string) int {
	if toTime == "" {
		return 0
	}

	//时间 to 时间戳
	loc, _ := time.LoadLocation("Asia/Shanghai") //设置时区
	tt, _ := time.ParseInLocation(constants.DateLayout, toTime, loc)
	return int(tt.Unix())
}

//时间戳转换时间
func TimeToStrTime(toTime int) string {
	if toTime < 1 {
		return ""
	}
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	timeobj := time.Unix(int64(toTime), 0)
	return timeobj.In(cstSh).Format(constants.TimeLayout)
}

func TimeToStrByFormat(toTime int, format string) string {
	if toTime < 1 {
		return ""
	}
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	timeobj := time.Unix(int64(toTime), 0)
	return timeobj.In(cstSh).Format(format)
}

//时间戳转换时间 1608880464 -> 12.1 14:30
func TimeToStrTimeShort(toTime int) string {
	timeLayout := "01.02 15:04"
	if toTime < 1 {
		return ""
	}
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	timeobj := time.Unix(int64(toTime), 0)
	return timeobj.In(cstSh).Format(timeLayout)
}

//时间戳转换时间 1608880464 -> 12月1日 周三
func TimeToWeek(toTime int) string {
	if toTime < 1 {
		return ""
	}
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	timeobj := time.Unix(int64(toTime), 0)
	date := timeobj.In(cstSh).Format("1月2日")

	weekMap := map[string]string{
		"Mon": "周一",
		"Tue": "周二",
		"Wed": "周三",
		"Thu": "周四",
		"Fri": "周五",
		"Sat": "周六",
		"Sun": "周日",
	}
	week := timeobj.In(cstSh).Format("Mon")
	return date + " " + weekMap[week]
}

//时间戳转换日期
func TimeToStrDate(toTime int) string {
	if toTime < 1 {
		return ""
	}
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	timeobj := time.Unix(int64(toTime), 0)
	return timeobj.In(cstSh).Format(constants.DateLayout)
}

//时间戳转换
func TimeToStrShow(toTime int) string {
	ret := ""
	if toTime < 1 {
		return ret
	}
	nowTime := int(time.Now().Unix())
	diff := nowTime - toTime
	if diff < 0 {
		return ret
	}

	if diff < 60 {
		ret = "刚刚"
	} else if 60 <= diff && diff <= 3600 {
		ret = strconv.Itoa(diff/60) + "分钟前"
	} else if 3600 <= diff && diff <= 86400 {
		ret = strconv.Itoa(diff/3600) + "小时前"
	} else if 86400 <= diff && diff <= 86400*30 {
		ret = strconv.Itoa(diff/86400) + "天前"
	} else {
		var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
		timeobj := time.Unix(int64(toTime), 0)
		ret = timeobj.In(cstSh).Format(constants.DateLayout)
	}
	return ret
}

//md5 加密
func Md5String(str string) string {
	if str == "" {
		return ""
	}

	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

//去重切片
func RemoveDuplicateElement(languages []string) []string {
	if len(languages) < 1 {
		return nil
	}
	result := make([]string, 0, len(languages))
	temp := map[string]struct{}{}
	for _, item := range languages {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func RemoveDuplicateInt(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func DomainFullUrl(key string) string {
	return setting.AppSetting.PrefixUrl + "/" + key
}

// RandomString returns a random string with a fixed length
func RandomString(n int, allowedChars ...[]rune) string {
	var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	var letters []rune

	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func FilterEmoji(content string) string {
	new_content := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			new_content += string(value)
		}
	}
	return new_content
}

func ExcelExportData(header []string, data []map[string]interface{}, sheetName string, fileName string) string {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	style := &xlsx.Style{}
	style.Fill = *xlsx.NewFill("solid", "EFEFDE", "EFEFDE")
	style.Border = xlsx.Border{RightColor: "FF"}
	file = xlsx.NewFile()
	sheet, _ = file.AddSheet(sheetName)
	row = sheet.AddRow()
	for i := 0; i < len(header); i++ { //looping from 0 to the length of the array
		cell = row.AddCell()
		cell.Value = header[i]
		cell.SetStyle(style)
	}
	for _, obj := range data {
		row = sheet.AddRow()
		for i := 0; i < len(header); i++ {
			switch obj[header[i]].(type) {
			case string:
				obj[header[i]] = obj[header[i]]
				break
			case int:
				obj[header[i]] = strconv.Itoa(obj[header[i]].(int))
				break
			case float64:
				obj[header[i]] = strconv.FormatFloat(obj[header[i]].(float64), 'E', -1, 64)
				break
			}
			cell = row.AddCell()
			cell.Value = obj[header[i]].(string)
		}
	}
	url := fileName + ".xlsx"
	file.Save(url)
	return url
}

// 10进制转成26进制字母 A:0,B:1
func DecToChar(num int64) string {
	s26 := strconv.FormatInt(num, 26)

	res := ""
	for _, c := range s26 {
		if c <= 57 {
			res += string(c + 17)
		} else {
			res += string(c - 22)
		}
	}
	return res
}

func ShowSubstr(s string, l int) string {
	if len(s) <= l {
		return s
	}
	ss, sl, rl, rs := "", 0, 0, []rune(s)
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			rl = 1
		} else {
			rl = 2
		}

		if sl+rl > l {
			break
		}
		sl += rl
		ss += string(r)
	}
	return ss
}

func IsInArray(item string, arr []string) bool {
	setMap := make(map[string]struct{})
	for _, val := range arr {
		setMap[val] = struct{}{}
	}
	if _, ok := setMap[item]; ok {
		return true
	}
	return false
}

func IntIsInArray(item int, arr []int) bool {
	setMap := make(map[int]struct{})
	for _, val := range arr {
		setMap[val] = struct{}{}
	}
	if _, ok := setMap[item]; ok {
		return true
	} else {
		return false
	}
}

func GetRandInt(maxInt, minInt int) int {
	max := int64(maxInt)
	min := int64(minInt)
	rand.Seed(time.Now().UnixNano())
	num := rand.Int63n(max-min) + min
	return int(num)
}

func GetUidFromHeader(c *gin.Context) (uid int, roleId int, isAdmin bool, err error) {
	const BEARER_SCHEMA_USER = "Bearer "

	authHeader := c.GetHeader("Authorization")
	logging.Debug("GetUidFromHeader get header", "Authorization", authHeader)

	if len(authHeader) < len(BEARER_SCHEMA_USER) {
		return 0, 0, false, fmt.Errorf("Token is empty")
	}

	token := authHeader[len(BEARER_SCHEMA_USER):]
	if token == "" {
		return 0, 0, false, nil
	}

	tokenType := authHeader[0:1] //第一个字节用来标识是用户的 token 还是 admin 的
	if tokenType == "A" {
		res, err := ParseAdminToken(token)
		if err != nil {
			return 0, 0, false, err
		}
		return res.Id, res.RoleId, true, nil
	} else {
		res, err := ParseUserToken(token)
		if err != nil {
			return 0, 0, false, err
		}
		return res.Id, 0, false, nil
	}
}

func FloatRound(f float64, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}

func GetInterfaceToString(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case time.Time:
		t, _ := value.(time.Time)
		key = t.String()
		// 2022-11-23 11:29:07 +0800 CST  这类格式把尾巴去掉
		key = strings.Replace(key, " +0800 CST", "", 1)
		key = strings.Replace(key, " +0000 UTC", "", 1)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}
