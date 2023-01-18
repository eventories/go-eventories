package core

import (
	"github.com/eventories/go-eventories/core/interaction"
	"github.com/eventories/go-eventories/database"
)

type Fetcher struct {
	eth *interaction.Interactor

	cp *checkpoint
	db database.Database
}

func (f *Fetcher) BlockNumber() uint64 {
	return f.cp.blockNumber()
}

func (f *Fetcher) SetBlockNumber(bn uint64) error {
	return f.cp.setBlockNumber(bn)
}
