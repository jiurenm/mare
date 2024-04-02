package randx

import (
	"crypto/rand"
	"unsafe"
)

const alpha = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// RandString 生成指定长度的随机字符串
func RandString(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}

	for i, b := range bytes {
		bytes[i] = alpha[b%byte(len(alpha))]
	}

	return bytesToString(bytes)
}

// RandNumber 生成指定长度的随机数字字符串
func RandNumber(n int) string {
	const alphaNum = "1234567890"

	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}

	for i, b := range bytes {
		bytes[i] = alphaNum[b%byte(len(alphaNum))]
	}

	return bytesToString(bytes)
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
