package editor

type MockFileOpener struct {
	CalledOpenEditor bool
}

func (m *MockFileOpener) OpenEditor(fileName string) error {
	m.CalledOpenEditor = true
	return nil
}
