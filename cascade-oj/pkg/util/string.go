package util

import (
	"regexp"
	"strings"
	"unicode"
)

// BytesCompareIgnoreSpacesAndNewlines 比较两个字节切片，忽略所有空白字符
func BytesCompareIgnoreSpacesAndNewlines(bytes1, bytes2 []byte) bool {
	skipChar := func(c byte) bool {
		return c == ' ' || c == '\n' || c == '\t' || c == '\r'
	}
	index1, index2 := 0, 0
	lenA, lenB := len(bytes1), len(bytes2)
	for index1 < lenA && index2 < lenB {
		for index1 < lenA && skipChar(bytes1[index1]) {
			index1++
		}
		for index2 < lenB && skipChar(bytes2[index2]) {
			index2++
		}
		if index1 < lenA && index2 < lenB {
			if bytes1[index1] != bytes2[index2] {
				return false
			}
			index1++
			index2++
		}
	}
	for index1 < lenA && skipChar(bytes1[index1]) {
		index1++
	}
	for index2 < lenB && skipChar(bytes2[index2]) {
		index2++
	}
	return index1 == lenA && index2 == lenB
}

// StringCompareIgnoreLineEndSpaceAndTextEndEnter 比较字符串，忽略行尾空格和文本末尾换行
func StringCompareIgnoreLineEndSpaceAndTextEndEnter(str1, str2 string) bool {
	str1 = strings.TrimRight(str1, " \n\t\r")
	str2 = strings.TrimRight(str2, " \n\t\r")
	str1 = StringRemoveLineEndSpace(str1)
	str2 = StringRemoveLineEndSpace(str2)
	return str1 == str2
}

// StringRemoveLineEndSpace 移除每行末尾的空格
func StringRemoveLineEndSpace(str string) string {
	strCopy := []byte(str)
	n := len(strCopy)
	res := make([]byte, 0, n)

	for i := 0; i < n; {
		start := i // 记录行有效的起始位置
		for i < n && strCopy[i] != '\n' {
			i++
		}

		// 移除行尾空格
		end := i // 记录行有效的结束位置
		for end > start && strCopy[end-1] == ' ' {
			end--
		}

		// 添加该行
		res = append(res, strCopy[start:end]...)

		// 如果是换行符，添加换行
		if i < n && strCopy[i] == '\n' {
			res = append(res, '\n')
			i++
		}
	}

	return string(res)
}

// StringRemoveLeadingTrailingSpaces 移除每行开头和结尾的空格
func StringRemoveLeadingTrailingSpaces(str string) string {
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n")
}

// CRLF2LF 将Windows换行符转换为Unix换行符
func CRLF2LF(str string) string {
	res := strings.Replace(str, "\r\n", "\n", -1)
	return res
}

// LF2CRLF 将Unix换行符转换为Windows换行符
func LF2CRLF(str string) string {
	replacer := strings.NewReplacer(
		"\r\n", "\r\n", // 确保不会重复处理
		"\n", "\r\n",
	)
	return replacer.Replace(str)
}

// RemoveCR 从字节切片中移除所有\r字符
func RemoveCR(str []byte) []byte {
	var res []byte
	for _, ch := range str {
		if ch != '\r' {
			res = append(res, ch)
		}
	}
	return res
}

// RemoveLF 从字节切片中移除所有\n字符
func RemoveLF(str []byte) []byte {
	var res []byte
	for _, ch := range str {
		if ch != '\n' {
			res = append(res, ch)
		}
	}
	return res
}

// RemoveAllWhitespace 移除字符串中所有空白字符
func RemoveAllWhitespace(str string) string {
	var res strings.Builder
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			res.WriteRune(ch)
		}
	}
	return res.String()
}

// RemoveExtraNewlines 将连续的多个换行压缩为最多2个换行
func RemoveExtraNewlines(s string) string {
	return regexp.MustCompile(`\n{3,}`).ReplaceAllString(s, "\n\n")
}

// NormalizeOutputString 规范化输出字符串（用于比较）
func NormalizeOutputString(output string) string {
	// 统一换行符
	normalized := CRLF2LF(output)

	// 移除首尾空白
	normalized = strings.TrimSpace(normalized)

	// 移除每行首尾的空格
	normalized = StringRemoveLeadingTrailingSpaces(normalized)

	// 压缩多余换行
	normalized = RemoveExtraNewlines(normalized)

	return normalized
}

// LooseNormalizeString 宽松规范化字符串（移除所有空白字符）
func LooseNormalizeString(str string) string {
	return RemoveAllWhitespace(CRLF2LF(str))
}

// CompareStrings 比较两个字符串，支持多种模式
func CompareStrings(str1, str2 string, mode int) bool {
	switch mode {
	case 0: // 完全匹配
		return str1 == str2
	case 1: // 标准化匹配
		return NormalizeOutputString(str1) == NormalizeOutputString(str2)
	case 2: // 宽松匹配（忽略所有空白，仅匹配有效字符）
		return LooseNormalizeString(str1) == LooseNormalizeString(str2)
	case 3: // 忽略行尾空格
		return StringCompareIgnoreLineEndSpaceAndTextEndEnter(str1, str2)
	default:
		return NormalizeOutputString(str1) == NormalizeOutputString(str2)
	}
}

// LimitStringLength 限制字符串长度，超过部分用...表示
func LimitStringLength(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	if maxLen <= 3 {
		return "..."
	}
	return str[:maxLen-3] + "..."
}
