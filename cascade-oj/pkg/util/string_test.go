package util

import (
	"bytes"
	"strings"
	"testing"
)

// TestBytesCompareIgnoreSpacesAndNewlines 测试字节比较函数
func TestBytesCompareIgnoreSpacesAndNewlines(t *testing.T) {
	tests := []struct {
		name     string
		bytes1   []byte
		bytes2   []byte
		expected bool
	}{
		// 基本测试
		{"完全相同", []byte("hello"), []byte("hello"), true},
		{"完全不同", []byte("hello"), []byte("world"), false},

		// 空格测试
		{"忽略空格", []byte("hello world"), []byte("hello  world"), true},
		{"忽略多个空格", []byte("a b c"), []byte("a  b   c"), true},

		// 换行测试
		{"忽略换行", []byte("hello\nworld"), []byte("hello world"), true},
		{"不同换行符", []byte("hello\r\nworld"), []byte("hello\nworld"), true},

		// 制表符测试
		{"忽略制表符", []byte("hello\tworld"), []byte("hello world"), true},

		// 混合测试
		{"混合空白", []byte("a\nb\tc d"), []byte("a b c d"), true},
		{"混合空白2", []byte("a\n\nb\t\tc d"), []byte("a b c d"), true},

		// 边界测试
		{"空字节", []byte(""), []byte(""), true},
		{"空字节与非空", []byte(""), []byte(" "), true}, // 忽略空格后都为空
		{"只有空白", []byte(" \n\t\r"), []byte(""), true},

		// 字符测试
		{"中间字符不同", []byte("abc def"), []byte("abc xyz"), false},
		{"长度不同", []byte("hello"), []byte("hello world"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BytesCompareIgnoreSpacesAndNewlines(tt.bytes1, tt.bytes2)
			if result != tt.expected {
				t.Errorf("BytesCompareIgnoreSpacesAndNewlines(%q, %q) = %v, want %v",
					string(tt.bytes1), string(tt.bytes2), result, tt.expected)
			}

			// 测试对称性
			result2 := BytesCompareIgnoreSpacesAndNewlines(tt.bytes2, tt.bytes1)
			if result != result2 {
				t.Errorf("函数不对称: BytesCompareIgnoreSpacesAndNewlines(%q, %q) = %v, 反向 = %v",
					string(tt.bytes1), string(tt.bytes2), result, result2)
			}
		})
	}
}

// TestStringRemoveLineEndSpace 测试移除行尾空格函数
func TestStringRemoveLineEndSpace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// 基本测试
		{"无行尾空格", "hello world", "hello world"},
		{"单行尾空格", "hello world   ", "hello world"},
		{"多行尾空格", "hello world      ", "hello world"},

		// 多行测试
		{"两行都有行尾空格", "line1  \nline2  ", "line1\nline2"},
		{"第一行有行尾空格", "line1  \nline2", "line1\nline2"},
		{"第二行有行尾空格", "line1\nline2  ", "line1\nline2"},

		// 混合测试
		{"行中有空格", "hello world  \nanother line", "hello world\nanother line"},
		{"行首空格保留", "  hello world  ", "  hello world"},
		{"空行", "\n\n", "\n\n"},
		{"空行带空格", "  \n  \n   ", "\n\n"}, // 行尾空格以及文本末尾空格被移除，留下空行

		// 边界测试
		{"空字符串", "", ""},
		{"只有空格", "   ", ""},
		{"只有换行", "\n", "\n"},
		{"Windows换行符", "line1  \r\nline2", "line1  \r\nline2"}, // 注意：需要先转换换行符

		// 复杂测试
		{
			"多行混合",
			"  line1 with spaces   \nline2\t\n   line3   \n\nline4",
			"  line1 with spaces\nline2\t\n   line3\n\nline4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringRemoveLineEndSpace(tt.input)
			if result != tt.expected {
				t.Errorf("StringRemoveLineEndSpace(%q) = %q, want %q",
					tt.input, result, tt.expected)
			}

			// 测试幂等性：再次应用应该得到相同结果
			result2 := StringRemoveLineEndSpace(result)
			if result != result2 {
				t.Errorf("函数不是幂等的: StringRemoveLineEndSpace(%q) = %q, 再次应用 = %q",
					result, result, result2)
			}
		})
	}
}

// TestStringCompareIgnoreLineEndSpaceAndTextEndEnter 测试忽略行尾空格的字符串比较
func TestStringCompareIgnoreLineEndSpaceAndTextEndEnter(t *testing.T) {
	tests := []struct {
		name     string
		str1     string
		str2     string
		expected bool
	}{
		// 完全相同
		{"完全相同", "hello", "hello", true},

		// 行尾空格不同
		{"行尾空格", "hello", "hello ", true},
		{"多行尾空格", "hello", "hello   ", true},
		{"多行行尾空格", "line1\nline2", "line1 \nline2 ", true},

		// 文本末尾换行
		{"末尾换行", "hello", "hello\n", true},
		{"末尾多个换行", "hello", "hello\n\n", true},
		{"末尾混合空白", "hello", "hello \n\t\r", true},

		// 实际内容不同
		{"内容不同", "hello", "world", false},
		{"行中不同", "hello world", "hello there", false},

		// 边界测试
		{"空字符串", "", "", true},
		{"空白字符串", " ", "", true},
		{"多个空白", " \n\t\r", "", true},

		// 复杂测试
		{
			"复杂相同",
			"line1\nline2  \nline3\n",
			"line1 \nline2\nline3 \n\n",
			true,
		},
		{
			"复杂不同",
			"line1\nline2\nline3",
			"line1\nlineX\nline3",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringCompareIgnoreLineEndSpaceAndTextEndEnter(tt.str1, tt.str2)
			if result != tt.expected {
				t.Errorf("StringCompareIgnoreLineEndSpaceAndTextEndEnter(%q, %q) = %v, want %v",
					tt.str1, tt.str2, result, tt.expected)
			}

			// 测试对称性
			result2 := StringCompareIgnoreLineEndSpaceAndTextEndEnter(tt.str2, tt.str1)
			if result != result2 {
				t.Errorf("函数不对称: StringCompareIgnoreLineEndSpaceAndTextEndEnter(%q, %q) = %v, 反向 = %v",
					tt.str1, tt.str2, result, result2)
			}
		})
	}
}

// TestCRLF2LFAndF2CRLF 测试换行符转换函数
func TestCRLF2LFAndF2CRLF(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedRLF2LF  string
		expectedF2CRLF string
	}{
		// 基本测试
		{"纯LF", "line1\nline2\n", "line1\nline2\n", "line1\r\nline2\r\n"},
		{"纯CRLF", "line1\r\nline2\r\n", "line1\nline2\n", "line1\r\nline2\r\n"},
		{"混合", "line1\nline2\r\nline3", "line1\nline2\nline3", "line1\r\nline2\r\nline3"},
		{"单个字符", "a\nb\r\nc", "a\nb\nc", "a\r\nb\r\nc"},

		// 边界测试
		{"空字符串", "", "", ""},
		{"只有LF", "\n", "\n", "\r\n"},
		{"只有CRLF", "\r\n", "\n", "\r\n"},
		{"多个LF", "\n\n\n", "\n\n\n", "\r\n\r\n\r\n"},
		{"多个CRLF", "\r\n\r\n\r\n", "\n\n\n", "\r\n\r\n\r\n"},

		// 修复的bug测试：不会出现 \r\r\n
		{"避免重复", "line1\r\nline2\nline3\r\n", "line1\nline2\nline3\n", "line1\r\nline2\r\nline3\r\n"},
		{"混合复杂", "\r\n\n\r\n", "\n\n\n", "\r\n\r\n\r\n"},

		// 多次转换测试
		{"多行", "a\nb\nc\nd", "a\nb\nc\nd", "a\r\nb\r\nc\r\nd"},
		{"带内容的CRLF", "hello\r\nworld\r\ntest", "hello\nworld\ntest", "hello\r\nworld\r\ntest"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试 CRLF2LF
			result := CRLF2LF(tt.input)
			if result != tt.expectedRLF2LF {
				t.Errorf("CRLF2LF(%q) = %q, want %q",
					tt.input, result, tt.expectedRLF2LF)
			}

			// 测试 LF2CRLF
			result2 := LF2CRLF(tt.input)
			if result2 != tt.expectedF2CRLF {
				t.Errorf("LF2CRLF(%q) = %q, want %q",
					tt.input, result2, tt.expectedF2CRLF)
			}

			// 测试 LF2CRLF 不会产生 \r\r\n
			if strings.Contains(result2, "\r\r\n") {
				t.Errorf("LF2CRLF(%q) = %q, contains \\r\\r\\n which is wrong",
					tt.input, result2)
			}

			// 测试 CRLF2LF 的幂等性
			result3 := CRLF2LF(result)
			if result != result3 {
				t.Errorf("CRLF2LF 不是幂等的: CRLF2LF(%q) = %q, 再次应用 = %q",
					result, result, result3)
			}

			// 测试 LF2CRLF 在特定条件下的幂等性
			// 如果输入已经是CRLF格式，再次转换应该不变
			if tt.input == tt.expectedF2CRLF {
				result4 := LF2CRLF(result2)
				if result2 != result4 {
					t.Errorf("LF2CRLF 不是幂等的: LF2CRLF(%q) = %q, 再次应用 = %q",
						result2, result2, result4)
				}
			}
		})
	}

	// 额外测试：验证不会出现 \r\r\n
	t.Run("验证不产生重复CR", func(t *testing.T) {
		// 测试各种可能产生重复的情况
		testCases := []string{
			"hello\r\nworld",
			"\r\n\r\n\r\n",
			"a\r\nb\nc\r\n",
			"\n\r\n\n",
		}

		for _, input := range testCases {
			result := LF2CRLF(input)
			if strings.Contains(result, "\r\r\n") {
				t.Errorf("LF2CRLF(%q) = %q, contains \\r\\r\\n",
					input, result)
			}
		}
	})
}

// TestRemoveCRAndRemoveLF 测试移除换行符函数
func TestRemoveCrAndRemoveLf(t *testing.T) {
	// 测试 RemoveCR
	t.Run("RemoveCR", func(t *testing.T) {
		tests := []struct {
			input    []byte
			expected []byte
		}{
			{[]byte("hello\rworld"), []byte("helloworld")},
			{[]byte("\rhello\r\rworld\r"), []byte("helloworld")},
			{[]byte(""), []byte("")},
			{[]byte("\r"), []byte("")},
			{[]byte("hello\nworld\r"), []byte("hello\nworld")},
		}

		for _, tt := range tests {
			result := RemoveCR(tt.input)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("RemoveCR(%q) = %q, want %q",
					string(tt.input), string(result), string(tt.expected))
			}
		}
	})

	// 测试 RemoveLf
	t.Run("RemoveLf", func(t *testing.T) {
		tests := []struct {
			input    []byte
			expected []byte
		}{
			{[]byte("hello\nworld"), []byte("helloworld")},
			{[]byte("\nhello\n\nworld\n"), []byte("helloworld")},
			{[]byte(""), []byte("")},
			{[]byte("\n"), []byte("")},
			{[]byte("hello\r\nworld"), []byte("hello\rworld")},
		}

		for _, tt := range tests {
			result := RemoveLF(tt.input)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("RemoveLF(%q) = %q, want %q",
					string(tt.input), string(result), string(tt.expected))
			}
		}
	})
}

// TestRemoveAllWhitespace 测试移除所有空白字符
func TestRemoveAllWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"无空白", "hello", "hello"},
		{"有空格", "hello world", "helloworld"},
		{"有换行", "hello\nworld", "helloworld"},
		{"有制表符", "hello\tworld", "helloworld"},
		{"混合空白", "hello \n\t\r world", "helloworld"},
		{"全空白", " \n\t\r ", ""},
		{"空字符串", "", ""},
		{"Unicode空白", "hello\u2003world", "helloworld"}, // 全角空格
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveAllWhitespace(tt.input)
			if result != tt.expected {
				t.Errorf("RemoveAllWhitespace(%q) = %q, want %q",
					tt.input, result, tt.expected)
			}

			// 测试幂等性
			result2 := RemoveAllWhitespace(result)
			if result != result2 {
				t.Errorf("函数不是幂等的: RemoveAllWhitespace(%q) = %q, 再次应用 = %q",
					result, result, result2)
			}
		})
	}
}

// TestRemoveExtraNewlines 测试移除多余换行
func TestRemoveExtraNewlines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"无多余换行", "line1\nline2", "line1\nline2"},
		{"三个换行", "line1\n\n\nline2", "line1\n\nline2"},
		{"多个换行", "line1\n\n\n\n\nline2", "line1\n\nline2"},
		{"开头多余换行", "\n\n\nline1", "\n\nline1"},
		{"结尾多余换行", "line1\n\n\n", "line1\n\n"},
		{"空行", "\n\n\n", "\n\n"},
		{"空字符串", "", ""},
		{"混合", "a\n\n\nb\nc\n\n\nd", "a\n\nb\nc\n\nd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveExtraNewlines(tt.input)
			if result != tt.expected {
				t.Errorf("RemoveExtraNewlines(%q) = %q, want %q",
					tt.input, result, tt.expected)
			}

			// 测试幂等性
			result2 := RemoveExtraNewlines(result)
			if result != result2 {
				t.Errorf("函数不是幂等的: RemoveExtraNewlines(%q) = %q, 再次应用 = %q",
					result, result, result2)
			}
		})
	}
}

// TestNormalizeOutputString 测试输出规范化
func TestNormalizeOutputString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// 基本测试
		{
			"Windows换行+行尾空格",
			"line1  \r\nline2  \r\n",
			"line1\nline2",
		},
		{
			"混合空白",
			"  line1  \n\n  line2  \n  line3  \n\n\n",
			"line1\n\nline2\nline3",
		},
		{
			"只有空白",
			"  \n\n  \t\n  ",
			"",
		},
		{
			"空字符串",
			"",
			"",
		},
		{
			"复杂案例",
			"  Result: 42  \r\n\r\n  Time: 1.23s  \r\n\r\n\r\n",
			"Result: 42\n\nTime: 1.23s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeOutputString(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeOutputString(%q) = %q, want %q",
					tt.input, result, tt.expected)
			}

			// 测试幂等性
			result2 := NormalizeOutputString(result)
			if result != result2 {
				t.Errorf("函数不是幂等的: NormalizeOutputString(%q) = %q, 再次应用 = %q",
					result, result, result2)
			}
		})
	}
}

// LooseNormalizeString 测试宽松规范化
func TestLooseNormalizeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"普通字符串", "hello world", "helloworld"},
		{"换行符", "hello\nworld", "helloworld"},
		{"Windows换行", "hello\r\nworld", "helloworld"},
		{"制表符", "hello\tworld", "helloworld"},
		{"混合", "a\nb\tc d\re", "abcde"},
		{"空字符串", "", ""},
		{"全空白", " \n\t\r ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LooseNormalizeString(tt.input)
			if result != tt.expected {
				t.Errorf("LooseNormalizeString(%q) = %q, want %q",
					tt.input, result, tt.expected)
			}

			// 测试幂等性
			result2 := LooseNormalizeString(result)
			if result != result2 {
				t.Errorf("函数不是幂等的: LooseNormalizeString(%q) = %q, 再次应用 = %q",
					result, result, result2)
			}
		})
	}
}

// TestCompareStrings 测试字符串比较函数
func TestCompareStrings(t *testing.T) {
	tests := []struct {
		name     string
		str1     string
		str2     string
		mode     int
		expected bool
	}{
		// 模式0：完全匹配
		{"完全匹配-相同", "hello", "hello", 0, true},
		{"完全匹配-不同", "hello", "hello ", 0, false},
		{"完全匹配-换行不同", "hello\n", "hello", 0, false},

		// 模式1：标准化匹配
		{"标准化-相同", "hello", "hello", 1, true},
		{"标准化-行尾空格", "hello", "hello ", 1, true},
		{"标准化-换行", "hello", "hello\n", 1, true},
		{"标准化-Windows换行", "hello\nworld", "hello\r\nworld", 1, true},
		{"标准化-多行", "a\nb", "a \nb ", 1, true},
		{"标准化-不同内容", "hello", "world", 1, false},

		// 模式2：宽松匹配（忽略所有空白）
		{"宽松-相同", "hello", "hello", 2, true},
		{"宽松-有空格", "hello world", "helloworld", 2, true},
		{"宽松-有换行", "hello\nworld", "helloworld", 2, true},
		{"宽松-不同", "hello", "helloo", 2, false},
		// 模式3：忽略行尾空格
		{"忽略行尾-相同", "hello", "hello", 3, true},
		{"忽略行尾-行尾空格", "hello", "hello ", 3, true},
		{"忽略行尾-多行", "a\nb", "a \nb ", 3, true},
		{"忽略行尾-行中空格", "hel lo", "hello", 3, false}, // 行中空格不能忽略

		// 默认模式
		{"默认模式-相同", "hello", "hello", 999, true},
		{"默认模式-行尾空格", "hello", "hello ", 999, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareStrings(tt.str1, tt.str2, tt.mode)
			if result != tt.expected {
				t.Errorf("CompareStrings(%q, %q, %d) = %v, want %v",
					tt.str1, tt.str2, tt.mode, result, tt.expected)
			}

			// 测试对称性（除了模式0完全匹配外）
			if tt.mode != 0 {
				result2 := CompareStrings(tt.str2, tt.str1, tt.mode)
				if result != result2 {
					t.Errorf("函数不对称: CompareStrings(%q, %q, %d) = %v, 反向 = %v",
						tt.str1, tt.str2, tt.mode, result, result2)
				}
			}
		})
	}
}

// TestLimitStringLength 测试字符串长度限制
func TestLimitStringLength(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		// 短于限制
		{"短字符串", "hello", 10, "hello"},
		{"等于限制", "hello", 5, "hello"},

		// 长于限制
		{"需要截断", "hello world", 8, "hello..."},
		{"需要截断2", "hello world", 6, "hel..."},
		{"长字符串", "abcdefghijklmnopqrstuvwxyz", 10, "abcdefg..."},

		// 边界情况
		{"空字符串", "", 5, ""},
		{"最大长度3", "hello", 3, "..."},
		{"最大长度小于3", "hello", 2, "..."},
		{"最大长度0", "hello", 0, "..."},
		{"最大长度负数", "hello", -1, "..."},

		// Unicode测试
		{"Unicode字符", "你好世界", 3, "..."}, // 中文字符每个占3字节
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LimitStringLength(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("LimitStringLength(%q, %d) = %q, want %q",
					tt.input, tt.maxLen, result, tt.expected)
			}
			if len(result) > tt.maxLen && result != "..." {
				t.Errorf("LimitStringLength returned a string longer than maxLen: got %q, want %q",
					result, tt.expected)
			}
		})
	}
}

// TestAllFunctionsWithUnicode 测试Unicode支持
func TestAllFunctionsWithUnicode(t *testing.T) {
	t.Run("Unicode字符串处理", func(t *testing.T) {
		// 测试各种Unicode字符
		unicodeStr := "Hello 世界\nこんにちは\n안녕하세요  "

		// 测试 StringRemoveLineEndSpace
		result := StringRemoveLineEndSpace(unicodeStr)
		expected := "Hello 世界\nこんにちは\n안녕하세요"
		if result != expected {
			t.Errorf("StringRemoveLineEndSpace with Unicode failed: got %q, want %q",
				result, expected)
		}

		// 测试 NormalizeOutputString
		result2 := NormalizeOutputString(unicodeStr)
		expected2 := "Hello 世界\nこんにちは\n안녕하세요"
		if result2 != expected2 {
			t.Errorf("NormalizeOutputString with Unicode failed: got %q, want %q",
				result2, expected2)
		}

		// 测试 RemoveAllWhitespace
		result3 := RemoveAllWhitespace(unicodeStr)
		expected3 := "Hello世界こんにちは안녕하세요"
		if result3 != expected3 {
			t.Errorf("RemoveAllWhitespace with Unicode failed: got %q, want %q",
				result3, expected3)
		}
	})
}

// BenchmarkStringRemoveLineEndSpace 性能测试
func BenchmarkStringRemoveLineEndSpace(b *testing.B) {
	// 准备测试数据
	testString := `Line 1 with some spaces at the end    
Line 2 without spaces
Line 3 with tabs	at the end		
Line 4 with multiple   spaces   in   the   middle    
Line 5 normal
Line 6 with spaces    

`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringRemoveLineEndSpace(testString)
	}
}

// BenchmarkNormalizeOutputString 性能测试
func BenchmarkNormalizeOutputString(b *testing.B) {
	// 准备测试数据
	testString := `  Result: 42  
	
  Time: 1.23 seconds  
	
  Memory: 128 MB  
	
	
  Status: OK    
	
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NormalizeOutputString(testString)
	}
}
