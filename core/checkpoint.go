package core

import (
	"encoding/binary"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
)

var (
	extension = ".checkpoint"
)

type Checkpoint struct {
	path string
	kind string
	n    uint64
}

func NewCheckpoint(basePath string, kind string) *Checkpoint {
	path, err := defaultPath(basePath, runtime.GOOS)
	if err != nil {
		panic(err)
	}

	if _, err := ioutil.ReadDir(path); err != nil {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			panic(err)
		}
		return &Checkpoint{path, kind, 0}
	}

	n, err := ioutil.ReadFile(filepath.Join(path, filepath.Base(kind+extension)))
	if err != nil {
		return &Checkpoint{path, kind, 0}
	}

	return &Checkpoint{path, kind, binary.BigEndian.Uint64(n)}
}

func (c *Checkpoint) Checkpoint() uint64 {
	return atomic.LoadUint64(&c.n)
}

func (c *Checkpoint) SetCheckpoint(n uint64) error {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)

	if err := ioutil.WriteFile(filepath.Join(c.path, filepath.Base(c.kind+extension)), b, fs.FileMode(0644)); err != nil {
		return err
	}

	atomic.StoreUint64(&c.n, n)
	return nil
}

func (c *Checkpoint) Increase() error {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, atomic.LoadUint64(&c.n)+1)

	if err := ioutil.WriteFile(filepath.Join(c.path, filepath.Base(c.kind+extension)), b, fs.FileMode(0644)); err != nil {
		return err
	}

	atomic.AddUint64(&c.n, 1)
	return nil
}

// There is little possibility of adding a separate logic for each
// OS other than the path.
func defaultPath(base string, os string) (string, error) {
	switch os {
	case "darwin":
		fallthrough
	case "freebsd":
		fallthrough
	case "linux":
		return "./" + base, nil
	case "windows":
		return ".\\" + base, nil
	default:
		return "", errors.New("not supported OS: " + os)
	}
}
