package stringx

import (
	"strings"
	"unicode"
	"unsafe"
)

// StringToBytes 将字符串转换为字节数组的函数
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		string
		Cap int
	}{s, len(s)}))
}

// BytesToString 将字节切片转换为字符串的函数
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Copy函数接受一个字符串s作为输入，并返回该字符串的副本。
func Copy(s string) string {
	// 使用unsafe包直接操作内存以创建输入字符串s的副本。
	// 这种方法可能存在风险，应谨慎使用，因为它绕过了Go的内存安全性。
	return *(*string)(unsafe.Pointer(&struct {
		string
		Cap int
	}{s, len(s)}))
}

// PadStart 在字符串开始位置填充指定的字符，直到字符串达到指定的长度。
func PadStart(str string, size int, defaultValue ...string) string {
	// 如果字符串长度已经大于等于指定长度，则直接返回原字符串。
	if len(str) >= size {
		return str
	}

	// 如果提供了默认值，则使用默认值进行填充，否则使用空格进行填充。
	if len(defaultValue) > 0 {
		str = strings.Repeat(defaultValue[0], size-len(str)) + str
	} else {
		str = strings.Repeat(" ", size-len(str)) + str
	}

	return str[:size]
}

// PadEnd 函数接受一个字符串和一个大小参数，以及一个可选的默认值参数。
// 如果输入字符串长度小于指定大小，则使用默认值或空格填充到指定大小；如果输入字符串长度大于等于指定大小，则返回原始字符串。
// 如果提供了默认值参数，则使用默认值填充到指定大小；否则使用空格填充。
func PadEnd(str string, size int, defaultValue ...string) string {
	if len(str) >= size {
		return str
	}

	if len(defaultValue) > 0 {
		str += strings.Repeat(defaultValue[0], size-len(str))
	} else {
		str += strings.Repeat(" ", size-len(str))
	}

	return str[:size]
}

// IsNumeric 判断传入的字符串是否为纯数字
func IsNumeric(str string) bool {
	if str == "" {
		return false
	}

	for _, s := range str {
		if !unicode.IsDigit(s) {
			return false
		}
	}

	return true
}
