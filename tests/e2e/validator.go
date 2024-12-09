package e2e

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// Validator provides methods to validate both console and Clipboard Output
type Validator struct {
	t         *testing.T
	output    string
	clipboard string
}

// NewOutputValidator creates a new Validator instance
func NewOutputValidator(t *testing.T, result *Result) *Validator {
	return &Validator{
		t:         t,
		output:    result.Output,
		clipboard: result.Clipboard,
	}
}

// AssertFileInOutput checks if a specific file is present in both console Output and Clipboard
func (v *Validator) AssertFileInOutput(filePath string) *Validator {
	// Check console Output
	if !strings.Contains(v.output, filePath) {
		v.t.Errorf("Expected console Output to contain file path %q", filePath)
	}

	// Check Clipboard content for START and END file markers
	startMarker := fmt.Sprintf("START FILE: %s", filePath)
	endMarker := fmt.Sprintf("END FILE: %s", filePath)

	if !strings.Contains(v.clipboard, startMarker) {
		v.t.Errorf("Expected Clipboard to contain start marker for file %q", filePath)
	}
	if !strings.Contains(v.clipboard, endMarker) {
		v.t.Errorf("Expected Clipboard to contain end marker for file %q", filePath)
	}

	return v
}

// AssertSuccessfulRun validates that the execution was successful
func (v *Validator) AssertSuccessfulRun() *Validator {
	requiredMessages := []string{
		"Total files found:",
		"Total lines across all parsed files:",
		"Estimated tokens:",
		"✅ File contents have been copied to clipboard.",
	}

	for _, msg := range requiredMessages {
		if !strings.Contains(v.output, msg) {
			v.t.Errorf("Expected Output to contain %q", msg)
		}
	}

	return v
}

// AssertError checks if the Output contains a specific error message
func (v *Validator) AssertError(expectedError string) *Validator {
	if !strings.Contains(v.output, expectedError) {
		v.t.Errorf("Expected error message %q in Output", expectedError)
	}
	return v
}

// AssertFileCount validates the total number of files found
func (v *Validator) AssertFileCount(expectedCount int) *Validator {
	re := regexp.MustCompile(`Total files found: (\d+)`)
	matches := re.FindStringSubmatch(v.output)

	if len(matches) < 2 {
		v.t.Error("Could not find total files count in Output")
		return v
	}

	count, _ := strconv.Atoi(matches[1])
	if count != expectedCount {
		v.t.Errorf("Expected %d files, got %d", expectedCount, count)
	}

	return v
}

// AssertLineCount validates the total number of lines
func (v *Validator) AssertLineCount(expectedCount int) *Validator {
	re := regexp.MustCompile(`Total lines across all parsed files: (\d+)`)
	matches := re.FindStringSubmatch(v.output)

	if len(matches) < 2 {
		v.t.Error("Could not find total lines count in Output")
		return v
	}

	count, _ := strconv.Atoi(matches[1])
	if count != expectedCount {
		v.t.Errorf("Expected %d lines, got %d", expectedCount, count)
	}

	return v
}

// AssertTokenCount validates the estimated token count
func (v *Validator) AssertTokenCount(expectedCount int) *Validator {
	// Extract token count from Output, handling both raw numbers and k-formatted numbers
	re := regexp.MustCompile(`Estimated tokens: (\d+(?:\.\d+)?k?)`)
	matches := re.FindStringSubmatch(v.output)

	if len(matches) < 2 {
		v.t.Error("Could not find token count in Output")
		return v
	}

	tokenStr := matches[1]
	var actualCount float64

	if strings.HasSuffix(tokenStr, "k") {
		// Handle k-formatted numbers (e.g., "18k" or "18.5k")
		numStr := strings.TrimSuffix(tokenStr, "k")
		baseCount, _ := strconv.ParseFloat(numStr, 64)
		actualCount = baseCount * 1000
	} else {
		actualCount, _ = strconv.ParseFloat(tokenStr, 64)
	}

	// Allow for some flexibility in token estimation (±5%)
	tolerance := float64(expectedCount) * 0.05
	if actualCount < float64(expectedCount)-tolerance || actualCount > float64(expectedCount)+tolerance {
		v.t.Errorf("Expected approximately %d tokens, got %.0f (tolerance: ±%.0f)",
			expectedCount, actualCount, tolerance)
	}

	return v
}

// AssertFileTooLarge checks if a file is properly marked as too large in the clipboard
func (v *Validator) AssertFileTooLarge(filePath string, expectedSize int) *Validator {
	startMarker := fmt.Sprintf("START FILE: %s", filePath)
	sizeMarker := fmt.Sprintf("<FILE TOO LARGE: %d bytes>", expectedSize)
	endMarker := fmt.Sprintf("END FILE: %s", filePath)

	// Check for start marker
	if !strings.Contains(v.clipboard, startMarker) {
		v.t.Errorf("Expected clipboard to contain start marker for file %q", filePath)
	}

	// Check for size warning
	if !strings.Contains(v.clipboard, sizeMarker) {
		v.t.Errorf("Expected clipboard to contain size warning %q for file %q", sizeMarker, filePath)
	}

	// Check for end marker
	if !strings.Contains(v.clipboard, endMarker) {
		v.t.Errorf("Expected clipboard to contain end marker for file %q", filePath)
	}

	// Check markers appear in correct order
	clipboardContent := v.clipboard
	startIndex := strings.Index(clipboardContent, startMarker)
	sizeIndex := strings.Index(clipboardContent, sizeMarker)
	endIndex := strings.Index(clipboardContent, endMarker)

	if !(startIndex < sizeIndex && sizeIndex < endIndex) {
		v.t.Errorf("File markers and size warning are not in correct order for file %q", filePath)
	}

	return v
}

// AssertEmptyFile checks if a file is properly marked as empty in the clipboard
func (v *Validator) AssertEmptyFile(filePath string) *Validator {
	startMarker := fmt.Sprintf("START FILE: %s", filePath)
	emptyMarker := "<EMPTY FILE>"
	endMarker := fmt.Sprintf("END FILE: %s", filePath)

	// Check for start marker
	if !strings.Contains(v.clipboard, startMarker) {
		v.t.Errorf("Expected clipboard to contain start marker for empty file %q", filePath)
	}

	// Check for empty file marker
	if !strings.Contains(v.clipboard, emptyMarker) {
		v.t.Errorf("Expected clipboard to contain empty file marker for file %q", filePath)
	}

	// Check for end marker
	if !strings.Contains(v.clipboard, endMarker) {
		v.t.Errorf("Expected clipboard to contain end marker for file %q", filePath)
	}

	// Check markers appear in correct order
	clipboardContent := v.clipboard
	startIndex := strings.Index(clipboardContent, startMarker)
	emptyIndex := strings.Index(clipboardContent, emptyMarker)
	endIndex := strings.Index(clipboardContent, endMarker)

	if !(startIndex < emptyIndex && emptyIndex < endIndex) {
		v.t.Errorf("File markers and empty marker are not in correct order for file %q", filePath)
	}

	return v
}
