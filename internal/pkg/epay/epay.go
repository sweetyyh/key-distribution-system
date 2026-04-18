package epay

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// PayParams 发起支付的参数
type PayParams struct {
	Money      string // 金额，保留两位小数
	Name       string // 商品名称
	NotifyURL  string // 服务器异步通知地址
	OutTradeNo string // 商户订单号
	PID        string // 商户ID
	ReturnURL  string // 页面跳转通知地址
	SiteName   string // 网站名称
	Type       string // 支付类型: alipay / wxpay / qqpay
	Key        string // 商户密钥
	SubmitURL  string // 支付平台 submit.php 地址
}

// BuildPayURL 构造带签名的支付跳转 URL（供前端直接跳转）
func BuildPayURL(p PayParams) string {
	params := map[string]string{
		"money":       p.Money,
		"name":        p.Name,
		"notify_url":  p.NotifyURL,
		"out_trade_no": p.OutTradeNo,
		"pid":         p.PID,
		"return_url":  p.ReturnURL,
		"sitename":    p.SiteName,
		"type":        p.Type,
	}
	sign := calcSign(params, p.Key)

	q := url.Values{}
	for k, v := range params {
		q.Set(k, v)
	}
	q.Set("sign", sign)
	q.Set("sign_type", "MD5")

	return p.SubmitURL + "?" + q.Encode()
}

// VerifyCallback 验证支付平台异步通知的签名
// params 是从回调请求中收到的所有参数（包括 sign 和 sign_type）
// key 是商户密钥
func VerifyCallback(params map[string]string, key string) bool {
	sign, ok := params["sign"]
	if !ok {
		return false
	}
	// 剔除 sign 和 sign_type 后重新计算
	filtered := make(map[string]string, len(params))
	for k, v := range params {
		if k == "sign" || k == "sign_type" {
			continue
		}
		filtered[k] = v
	}
	expected := calcSign(filtered, key)
	return strings.EqualFold(sign, expected)
}

// calcSign 按键名 ASCII 升序排列后拼接并 MD5
func calcSign(params map[string]string, key string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if params[k] == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+params[k])
	}
	raw := strings.Join(parts, "&") + key

	h := md5.New()
	io.WriteString(h, raw)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// QueryOrder 查询单个订单状态（可选用）
func QueryOrder(apiURL, pid, key, outTradeNo string) (string, error) {
	u := fmt.Sprintf("%s?act=order&pid=%s&key=%s&out_trade_no=%s", apiURL, pid, key, outTradeNo)
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	return string(b), err
}
