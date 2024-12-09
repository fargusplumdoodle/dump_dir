package tests

type MockClipboard struct {
	Content string
}

func NewMockClipboard() *MockClipboard {
	return &MockClipboard{}
}

func (m *MockClipboard) WriteAll(text string) error {
	m.Content = text
	return nil
}
