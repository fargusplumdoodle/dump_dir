package src

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type TokenEstimator struct {
}

func NewTokenEstimator() *TokenEstimator {
	return &TokenEstimator{}
}

func (te *TokenEstimator) EstimateTokens(content string) int {
	if content == "" {
		return 0
	}

	words := strings.Fields(content)
	count := 0

	for _, word := range words {
		count += te.estimateWordTokens(word)
	}

	// Account for whitespace and formatting
	count += te.estimateWhitespaceTokens(content)

	return count
}

// estimateWordTokens analyzes a single word and estimates its token count
func (te *TokenEstimator) estimateWordTokens(word string) int {
	switch {
	case te.containsNumber(word):
		return te.estimateNumberTokens(word)
	case te.isCompoundWord(word):
		return te.estimateCompoundTokens(word)
	case te.containsSpecialChars(word):
		return te.estimateSpecialCharTokens(word)
	case te.isURLOrPath(word):
		return te.estimateURLPathTokens(word)
	case te.containsUnicode(word):
		return te.estimateUnicodeTokens(word)
	default:
		return te.estimateBasicWordTokens(word)
	}
}

// containsNumber checks if the string contains any numerical digits
func (te *TokenEstimator) containsNumber(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// estimateNumberTokens estimates tokens for strings containing numbers
func (te *TokenEstimator) estimateNumberTokens(s string) int {
	numCount := 0
	for _, r := range s {
		if unicode.IsDigit(r) {
			numCount++
		}
	}
	// Numbers are often split into 2-3 digit chunks
	return (numCount+2)/3 + 1 // Add 1 for any surrounding context
}

// isCompoundWord checks if the word is a compound word (camelCase or snake_case)
func (te *TokenEstimator) isCompoundWord(s string) bool {
	return te.isCamelCase(s) || strings.Contains(s, "_")
}

// isCamelCase checks if the string contains camelCase patterns
func (te *TokenEstimator) isCamelCase(s string) bool {
	var (
		hasLower     bool
		hasUpper     bool
		lastWasLower bool
	)

	for _, r := range s {
		if unicode.IsUpper(r) {
			if lastWasLower {
				return true
			}
			hasUpper = true
		}
		if unicode.IsLower(r) {
			hasLower = true
			lastWasLower = true
		} else {
			lastWasLower = false
		}
	}

	return hasUpper && hasLower
}

// estimateCompoundTokens estimates tokens for compound words
func (te *TokenEstimator) estimateCompoundTokens(s string) int {
	parts := strings.Split(s, "_")
	count := 0

	for _, part := range parts {
		if te.isCamelCase(part) {
			// Count transitions from lower to upper case
			var lastWasLower bool
			upperCount := 0
			for _, r := range part {
				if unicode.IsUpper(r) && lastWasLower {
					upperCount++
				}
				lastWasLower = unicode.IsLower(r)
			}
			count += upperCount + 1
		} else {
			count++
		}
	}

	return max(count, 1)
}

// containsSpecialChars checks if the string contains special characters
func (te *TokenEstimator) containsSpecialChars(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '_' {
			return true
		}
	}
	return false
}

// estimateSpecialCharTokens estimates tokens for strings with special characters
func (te *TokenEstimator) estimateSpecialCharTokens(s string) int {
	specialCount := 0
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '_' {
			specialCount++
		}
	}
	return specialCount + 1 // Add 1 for the main word content
}

// isURLOrPath checks if the string appears to be a URL or file path
func (te *TokenEstimator) isURLOrPath(s string) bool {
	return strings.Contains(s, "://") ||
		strings.HasPrefix(s, "www.") ||
		strings.Contains(s, ".com") ||
		strings.Contains(s, ".org") ||
		strings.Contains(s, ".net") ||
		strings.Contains(s, "/") ||
		strings.Contains(s, "\\")
}

// estimateURLPathTokens estimates tokens for URLs and file paths
func (te *TokenEstimator) estimateURLPathTokens(s string) int {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '/' || r == '\\' || r == '.' || r == ':' || r == '-' || r == '_'
	})
	return len(parts) + 1 // Add 1 for the separators themselves
}

// containsUnicode checks if the string contains non-ASCII characters
func (te *TokenEstimator) containsUnicode(s string) bool {
	return utf8.RuneCountInString(s) != len(s)
}

// estimateUnicodeTokens estimates tokens for strings with Unicode characters
func (te *TokenEstimator) estimateUnicodeTokens(s string) int {
	return utf8.RuneCountInString(s)
}

// estimateBasicWordTokens estimates tokens for basic words
func (te *TokenEstimator) estimateBasicWordTokens(s string) int {
	if len(s) > 12 {
		return 2
	}
	return 1
}

// estimateWhitespaceTokens estimates tokens for whitespace and formatting
func (te *TokenEstimator) estimateWhitespaceTokens(content string) int {
	return strings.Count(content, "\n") + 1
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
