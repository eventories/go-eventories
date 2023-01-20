package filter

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/eventories/go-eventories/core/interaction"
)

var (
	_ = logFilter(&address{})
	_ = logFilter(&event{})
)

type address struct {
	target common.Address
}

func (a *address) Kind() Kind { return AddressLogFilter }

func (a *address) do(p *Purifier, eth *interaction.Interactor, logs []*types.Log) error {
	rlogs := make([]*types.Log, 0)
	for _, log := range logs {
		if bytes.Equal(a.target.Bytes(), log.Address.Bytes()) {
			rlogs = append(rlogs, log)
		}
	}

	r := make(map[common.Hash][]*types.Log)
	r[a.target.Hash()] = rlogs

	p.logs[a.Kind()] = r

	return nil
}

type event struct {
	id common.Hash
}

func (e *event) Kind() Kind { return EventLogFilter }

func (e *event) do(p *Purifier, eth *interaction.Interactor, logs []*types.Log) error {
	rlogs := make([]*types.Log, 0)
	for _, log := range logs {
		if len(log.Topics) != 0 {
			if bytes.Equal(e.id[:], log.Topics[0].Bytes()) {
				rlogs = append(rlogs, log)
			}
		}
	}

	r := make(map[common.Hash][]*types.Log)
	r[e.id] = rlogs

	p.logs[e.Kind()] = r

	return nil
}
