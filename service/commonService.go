package service

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"sort"
	"strings"
	"tiktoklive/utils"
)

var (
	ConnManger       *utils.ConnectionManager
	LogObj           *LogService
	TiktokSignerUrl  = ""
	TiktokBaseUrl    = "https://www.tiktok.com/"
	TiktokWebcastUrl = "https://webcast.us.tiktok.com/webcast/"
	secret           = "lUKG128F7po2jK54R5V23N76d6308WHDITNO20LSDJ-6NMV"
	defaultGetParams = map[string]string{
		"aid":                 "1988",
		"app_language":        "en-US",
		"app_name":            "tiktok_web",
		"browser_language":    "en",
		"browser_name":        "Mozilla",
		"browser_online":      "true",
		"browser_platform":    "Win32",
		"browser_version":     "5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.63 Safari/537.36",
		"cookie_enabled":      "true",
		"cursor":              "",
		"internal_ext":        "",
		"device_platform":     "web",
		"focus_state":         "true",
		"from_page":           "user",
		"history_len":         "4",
		"is_fullscreen":       "false",
		"is_page_visible":     "true",
		"did_rule":            "3",
		"fetch_rule":          "1",
		"identity":            "audience",
		"last_rtt":            "0",
		"live_id":             "12",
		"resp_content_type":   "protobuf",
		"screen_height":       "1152",
		"screen_width":        "2048",
		"tz_name":             "Europe/Berlin",
		"referer":             "https://www.tiktok.com/",
		"root_referer":        "https://www.tiktok.com/",
		"version_code":        "180800",
		"webcast_sdk_version": "1.3.0",
		"update_version_code": "1.3.0",
	}
	commHeader = map[string]string{
		"User-Agent":      "5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.63 Safari/537.36",
		"Accept":          "text/html,application/json,application/protobuf",
		"Origin":          "https://www.tiktok.com/",
		"Accept-Language": "en-US,en;q=0.9",
		"Connection":      "close",
		"Cache-Control":   "max-age=0",
		"referer":         "https://www.tiktok.com/",
	}
	wssClose = errors.New("web socket close")
)

// 参数与值映射
func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string)
	for key, value := range m {
		out[key] = value
	}
	return out
}

// 拼接key=value&
func assembleParams(m map[string]string) string {
	keys := sortMapKeys(m)
	params := ""
	for _, key := range keys {
		value := m[key]
		params += key + "=" + value + "&"
	}
	// 去掉最后一个多余的 '&'
	params = params[:len(params)-1]
	return params
}

// sortMapKeys升级排序
func sortMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func BuildSign(reqBody map[string]string) string {
	verifyStr := assembleParams(reqBody)
	//拼接secret
	verifyStr = verifyStr + "&secret=" + secret

	signMd5 := md5.Sum([]byte(verifyStr))

	return strings.ToUpper(hex.EncodeToString(signMd5[:]))
}

// GetCommonHeader 获取公共头
func GetCommonHeader(h map[string]string) map[string]string {
	header := make(map[string]string)
	for key, value := range h {
		header[key] = value
	}
	return header
}

// ResponseJson 响应json数据格式
func ResponseJson(w http.ResponseWriter, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

// GenerateUniqueString 生成随机数
func GenerateUniqueString(length int) string {
	// 计算生成字节数
	byteLength := length * 3 / 4

	// 生成随机字节
	randomBytes := make([]byte, byteLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}

	// 使用 base64 编码生成唯一字符串
	uniqueString := base64.RawURLEncoding.EncodeToString(randomBytes)

	// 截取指定长度的字符串
	if len(uniqueString) > length {
		uniqueString = uniqueString[:length]
	}

	return uniqueString
}
