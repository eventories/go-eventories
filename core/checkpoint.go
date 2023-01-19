package core

import (
	"encoding/binary"
	"io/fs"
	"io/ioutil"
	"sync/atomic"
)

const defaultPath = "/checkpoint"

type Checkpoint struct {
	kind string
	n    uint64
}

func NewCheckpoint(kind string) *Checkpoint {
	_, err := ioutil.ReadDir(defaultPath)
	if err != nil {
		return &Checkpoint{defaultPath, 0}
	}

	n, err := ioutil.ReadFile(defaultPath + "/" + kind)
	if err != nil {
		return &Checkpoint{defaultPath, 0}
	}

	return &Checkpoint{defaultPath, binary.BigEndian.Uint64(n)}
}

func (c *Checkpoint) Checkpoint() uint64 {
	return atomic.LoadUint64(&c.n)
}

func (c *Checkpoint) SetCheckpoint(n uint64) error {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)

	if err := ioutil.WriteFile(defaultPath+"/"+c.kind, b, fs.FileMode(1)); err != nil {
		return err
	}

	atomic.StoreUint64(&c.n, n)
	return nil
}

func (c *Checkpoint) Increase() error {
	atomic.AddUint64(&c.n, 1)
	// write files
	return nil
}
