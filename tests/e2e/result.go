package e2e

import (
	"fmt"
	"strings"
)

// Result contains the results of running the command
type Result struct {
	env       *Environment
	Output    string
	Clipboard string
	Err       error
}

// AssertOutputContains checks if the command Output contains expected content
func (r *Result) AssertOutputContains(expected string) *Result {
	if !strings.Contains(r.Output, expected) {
		r.env.t.Errorf("Expected output to contain %q, got: %q", expected, r.Output)
	}
	return r
}

// AssertClipboardContains checks if the Clipboard contains expected content
func (r *Result) AssertClipboardContains(expected string) *Result {
	if !strings.Contains(r.Clipboard, expected) {
		r.env.t.Errorf("Expected clipboard to contain %q, got: %q", expected, r.Clipboard)
	}
	return r
}

// AssertNoError checks if the command completed without error
func (r *Result) AssertNoError() *Result {
	if r.Err != nil {
		r.env.t.Errorf("Expected no error, got: %v", r.Err)
	}
	return r
}

// AssertError checks if the command completed with an error
func (r *Result) AssertError() *Result {
	if r.Err == nil {
		r.env.t.Error("Expected an error, got none")
	}
	return r
}

func (r *Result) PrintOutput() {
	fmt.Println(r.Output)
}

func (r *Result) PrintClipboard() {
	fmt.Println(r.Clipboard)
}
