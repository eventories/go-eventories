package core

import (
	"encoding/binary"
	"io/fs"
	"io/ioutil"
	"sync/atomic"
)

type checkpoint struct {
	path  string
	block uint64
}

func newCheckpoint(path string) *checkpoint {
	_, err := ioutil.ReadDir(path)
	if err != nil {
		return &checkpoint{path, 0}
	}

	bn, err := ioutil.ReadFile(path + "checkpoint")
	if err != nil {
		return &checkpoint{path, 0}
	}

	return &checkpoint{path, binary.BigEndian.Uint64(bn)}
}

func (c *checkpoint) blockNumber() uint64 {
	return atomic.LoadUint64(&c.block)
}

func (c *checkpoint) setBlockNumber(bn uint64) error {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, bn)

	if err := ioutil.WriteFile(c.path+"checkpoint", b, fs.FileMode(1)); err != nil {
		return err
	}

	atomic.StoreUint64(&c.block, bn)
	return nil
}
