// formatter/formatter.go
package formatter

import (
	"unicode"
)

// FormatText formats text by adding spaces between:
// 1. CJK and Latin characters
// 2. CJK and numbers
func FormatText(text string) string {
	if text == "" {
		return ""
	}

	// 先处理标点替换
	text = replacePunctuation(text)

	var result []rune
	runes := []rune(text)
	
	// Process all characters
	for i := range runes {
		// Add current character
		result = append(result, runes[i])
		
		// If this is not the last character
		if i < len(runes)-1 {
			curr := runes[i]
			next := runes[i+1]
			
			// Check if we need to add a space between current and next character
			if needSpace(curr, next) && next != ' ' && curr != ' ' {
				result = append(result, ' ')
			}
		}
	}
	
	return string(result)
}

// needSpace determines if a space is needed between two characters
func needSpace(curr, next rune) bool {
	// Rule 1: Add space between CJK and Latin characters
	cjkLatinSpace := (isCJK(curr) && isLatin(next)) || (isLatin(curr) && isCJK(next))
	
	// Rule 2: Add space between CJK and numbers
	cjkNumberSpace := (isCJK(curr) && unicode.IsNumber(next)) || (unicode.IsNumber(curr) && isCJK(next))
	
	return cjkLatinSpace || cjkNumberSpace
}

// isCJK checks if a character is a CJK character
func isCJK(r rune) bool {
	// Unicode ranges for Chinese, Japanese, and Korean characters
	return unicode.Is(unicode.Han, r) || // Chinese
		unicode.Is(unicode.Hiragana, r) || // Japanese Hiragana
		unicode.Is(unicode.Katakana, r) || // Japanese Katakana
		unicode.Is(unicode.Hangul, r) // Korean
}

// isLatin checks if a character is a Latin character
func isLatin(r rune) bool {
	return unicode.IsLetter(r) && !isCJK(r)
}

// replacePunctuation 替换中文后面的英文标点为对应的中文标点，并删除多余空格
func replacePunctuation(text string) string {
	if text == "" {
		return ""
	}

	var result []rune
	runes := []rune(text)
	length := len(runes)

	for i := 0; i < length; i++ {
		curr := runes[i]
		result = append(result, curr)

		// 如果当前字符是 CJK，且下一个字符是需要替换的英文标点
		if i < length-1 && isCJK(curr) {
			switch runes[i+1] {
			case ',':
				result = append(result, '，')
				i++
				// 跳过逗号后面的空格
				for i+1 < length && runes[i+1] == ' ' {
					i++
				}
			case '.':
				result = append(result, '。')
				i++
				// 跳过句号后面的空格
				for i+1 < length && runes[i+1] == ' ' {
					i++
				}
			}
		}
	}

	return string(result)
}
