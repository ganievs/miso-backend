package storage

import "io"

type Storage interface {
	Get(key string) ([]byte, error)
	Put(key string, data io.Reader) error
	Delete(key string) error
	List(path string) ([]string, error)
}
