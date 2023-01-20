package filter

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/eventories/go-eventories/core/interaction"
)

type (
	Address struct {
		Target common.Address
	}

	Event struct {
		ID common.Hash
	}
)

//
func (a *Address) Kind() Kind { return AddressLogsFilter }

func (a *Address) Do(p *Purifier, eth *interaction.Interactor, logs []*types.Log) error {
	rlogs := make([]*types.Log, 0)
	for _, log := range logs {
		if bytes.Equal(log.Address.Bytes(), a.Target.Bytes()) {
			rlogs = append(rlogs, log)
		}
	}

	r := make(map[common.Hash][]*types.Log)
	r[a.Target.Hash()] = rlogs

	p.logs[a.Kind()] = r

	return nil
}

//
func (e *Event) Kind() Kind { return EventLogsFilter }

func (e *Event) Do(p *Purifier, eth *interaction.Interactor, logs []*types.Log) error {
	rlogs := make([]*types.Log, 0)
	for _, log := range logs {
		if len(log.Topics) != 0 {
			if bytes.Equal(log.Topics[0].Bytes(), e.ID[:]) {
				rlogs = append(rlogs, log)
			}
		}
	}

	r := make(map[common.Hash][]*types.Log)
	r[e.ID] = rlogs

	p.logs[e.Kind()] = r

	return nil
}
