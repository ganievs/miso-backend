package storage

type Storage interface {
	Get(key string) ([]byte, error)
	Put(key string, data []byte) error
	Delete(key string) error
	List(path string) ([]string, error)
}
