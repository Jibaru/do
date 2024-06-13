package reader

type MockFileReader struct {
	ReadFn func(filename string) (string, error)
}

func (m *MockFileReader) Read(filename string) (string, error) {
	return m.ReadFn(filename)
}
