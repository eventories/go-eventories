package table

import (
	"bytes"
	"encoding/json"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/eventories/go-eventories/database"
)

const (
	prefix = "table-abi-"
)

type Table struct {
	db database.Database

	mu    sync.Mutex
	cache map[string]*abi.ABI
}

func (t *Table) ABI(name string) *abi.ABI {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.abi(name)
}

func (t *Table) Kinds() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	res := make([]string, 0, len(t.cache))
	for kind := range t.cache {
		res = append(res, kind)
	}
	return res
}

func (t *Table) Register(name string, abi *abi.ABI) error {
	b, err := json.Marshal(abi)
	if err != nil {
		return err
	}

	if err := t.db.Put([]byte(prefix+name), b); err != nil {
		return err
	}

	t.mu.Lock()
	t.cache[name] = abi
	t.mu.Unlock()

	return nil
}

func (t *Table) Deregister(name string) error {
	if err := t.db.Put([]byte(prefix+name), nil); err != nil {
		return nil
	}

	t.mu.Lock()
	delete(t.cache, name)
	t.mu.Unlock()
	return nil
}

func (t *Table) Satisfy(name string, code []byte) bool {
	t.mu.Lock()
	methodIDs := t.methodIDs(t.abi(name))
	t.mu.Unlock()

	if methodIDs == nil {
		return false
	}

	for _, id := range methodIDs {
		if !bytes.Contains(code, id) {
			return false
		}
	}

	return true
}

// The caller must hold t.mu.
func (t *Table) abi(name string) *abi.ABI {
	if abi, ok := t.cache[name]; ok {
		return abi
	}

	b, err := t.db.Get([]byte(prefix + name))
	if err != nil {
		return nil
	}

	var abi abi.ABI
	if err := json.Unmarshal(b, &abi); err != nil {
		panic("noop")
	}

	t.cache[name] = &abi

	return &abi
}

func (t *Table) methodIDs(abi *abi.ABI) (res [][]byte) {
	events := abi.Methods
	res = make([][]byte, 0, len(events))

	for _, method := range events {
		res = append(res, method.ID)
	}
	return res
}
