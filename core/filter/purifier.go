package filter

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/eventories/go-eventories/core/interaction"
)

var defaultFilters = []Filter{&allLogs{}}

// Block Purifier
type Purifier struct {
	logs map[Kind]map[common.Hash][]*types.Log
	txs  map[Kind][]*types.Transaction

	filters map[Kind]Filter
}

func New() *Purifier {
	filters := make(map[Kind]Filter)

	for _, df := range defaultFilters {
		filters[df.Kind()] = df
	}

	return &Purifier{
		logs:    make(map[Kind]map[common.Hash][]*types.Log),
		txs:     make(map[Kind][]*types.Transaction),
		filters: filters,
	}
}

func (p *Purifier) RegisterFilter(fs ...Filter) {
	for _, f := range fs {
		if _, ok := p.filters[f.Kind()]; ok {
			continue
		}

		p.filters[f.Kind()] = f
	}
}

func (p *Purifier) Filters() []Kind {
	kinds := make([]Kind, 0, len(p.filters))
	for kind := range p.filters {
		kinds = append(kinds, kind)
	}
	return kinds
}

func (p *Purifier) Filtering(eth *interaction.Interactor, txs []*types.Transaction) error {
	// defaultFilters are performed first.
	for _, filter := range p.filters {
		if err := p.route(filter, eth, txs); err != nil {
			return err
		}
	}

	return nil
}

func (p *Purifier) Log(kind Kind) map[common.Hash][]*types.Log {
	if _, ok := p.filters[kind]; !ok {
		return nil
	}
	return p.logs[kind]
}

func (p *Purifier) Tx(kind Kind) []*types.Transaction {
	if _, ok := p.filters[kind]; !ok {
		return nil
	}
	return p.txs[kind]
}

func (p *Purifier) route(filter Filter, eth *interaction.Interactor, txs []*types.Transaction) error {
	var err error

	switch f := filter.(type) {
	case transactionFilter:
		err = f.do(p, eth, txs)

	case logFilter:
		r := make([]*types.Log, 0)

		logs := p.logs[AllLogsFilter]
		for _, log := range logs {
			r = append(r, log...)
		}

		err = f.do(p, eth, r)

	default:
		err = errors.New("invalid filter type")
	}

	if err != nil {
		p.logs = make(map[Kind]map[common.Hash][]*types.Log)
		p.txs = make(map[Kind][]*types.Transaction)
		return err
	}

	return nil
}
