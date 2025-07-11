package storage

import "io"

type Storage interface {
	GetBuffer(key string) ([]byte, error)
	GetStream(key string) (io.ReadCloser, error)
	Put(key string, data io.Reader) error
	Delete(key string) error
	List(path string) ([]string, error)
	GetPresignedURL(key string) (string, error)
}
