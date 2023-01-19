package core

import (
	"github.com/eventories/go-eventories/core/interaction"
	"github.com/eventories/go-eventories/database"
)

type Fetcher struct {
	eth *interaction.Interactor

	cp *Checkpoint
	db database.Database
}

func NewFetcher(interact *interaction.Interactor, db database.Database) (*Fetcher, error) {
	return &Fetcher{
		eth: interact,
		cp:  NewCheckpoint("blockNumber"),
	}, nil
}

func (f *Fetcher) BlockNumber() uint64 {
	return f.cp.Checkpoint()
}

func (f *Fetcher) SetBlockNumber(bn uint64) error {
	return f.cp.SetCheckpoint(bn)
}
