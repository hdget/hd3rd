package smsforwarder

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	convertUtils "github.com/hdget/hdutils/convert"
	"github.com/pkg/errors"
	"net/url"
)

// GenerateSignature 签名算法
func generateSignature(timestamp, secret string) (string, error) {
	if timestamp == "" || secret == "" {
		return "", errors.New("timestamp or secret is empty")
	}

	// 构建待签名字符串
	stringToSign := fmt.Sprintf("%s\n%s", timestamp, secret)
	stringToSignEnc := []byte(stringToSign)

	// 生成HMAC SHA256签名
	hmacCode := hmac.New(sha256.New, convertUtils.StringToBytes(secret))
	hmacCode.Write(stringToSignEnc)
	signature := base64.StdEncoding.EncodeToString(hmacCode.Sum(nil))

	// 对签名进行URL编码
	encodedSign := url.QueryEscape(signature)

	return encodedSign, nil
}
