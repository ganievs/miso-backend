package storage

import "io"

// MockStorage is a mock implementation of the Storage interface.
type MockStorage struct {
	ListFunc            func(prefix string) ([]string, error)
	GetPresignedURLFunc   func(key string) (string, error)
	GetStreamFunc         func(key string) (io.ReadCloser, error)
}

func (m *MockStorage) GetBuffer(key string) ([]byte, error) {
	return nil, nil
}

func (m *MockStorage) GetStream(key string) (io.ReadCloser, error) {
	if m.GetStreamFunc != nil {
		return m.GetStreamFunc(key)
	}
	return nil, nil
}

func (m *MockStorage) Put(key string, data io.Reader) error {
	return nil
}

func (m *MockStorage) Delete(key string) error {
	return nil
}

func (m *MockStorage) List(prefix string) ([]string, error) {
	if m.ListFunc != nil {
		return m.ListFunc(prefix)
	}
	return nil, nil
}

func (m *MockStorage) GetPresignedURL(key string) (string, error) {
	if m.GetPresignedURLFunc != nil {
		return m.GetPresignedURLFunc(key)
	}
	return "", nil
}
