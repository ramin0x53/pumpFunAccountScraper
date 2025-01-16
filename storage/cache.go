package storage

type Cache interface {
	KeyExist(string) (bool, error)
	AddKey(string, string) error
}
