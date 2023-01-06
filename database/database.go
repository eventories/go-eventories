package database

type Database interface {
	Get(key []byte) ([]byte, error)
	Has(key []byte) bool
	Put(key []byte, value []byte) error
}
