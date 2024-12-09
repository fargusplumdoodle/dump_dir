package src

import (
	"github.com/atotto/clipboard"
)

type ClipboardManager interface {
	WriteAll(string) error
}

type SystemClipboard struct{}

func NewSystemClipboard() ClipboardManager {
	return &SystemClipboard{}
}

func (c *SystemClipboard) WriteAll(text string) error {
	return clipboard.WriteAll(text)
}
