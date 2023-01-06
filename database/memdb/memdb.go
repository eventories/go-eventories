package memdb

import "sync"

type memoryDB struct {
	mu sync.Mutex
	db map[string][]byte
}

func New() (*memoryDB, error) {
	return &memoryDB{
		db: make(map[string][]byte),
	}, nil
}

func (m *memoryDB) Get(key []byte) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.db[string(key)], nil
}

func (m *memoryDB) Has(key []byte) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.db[string(key)]
	return ok
}

func (m *memoryDB) Put(key []byte, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.db[string(key)] = value
	return nil
}

func (m *memoryDB) Delete(key []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.db, string(key))
	return nil
}
