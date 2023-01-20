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
	// logs stores items by Kind. Items are mapped Logs to
	// requested values ​​(e.g. address, topic).
	logs map[Kind]map[common.Hash][]*types.Log

	// txs maps the result of the requested filter to Kind.
	txs map[Kind][]*types.Transaction

	filters map[Kind]Filter
}

func New(filters ...Filter) *Purifier {
	p := &Purifier{
		logs:    make(map[Kind]map[common.Hash][]*types.Log),
		txs:     make(map[Kind][]*types.Transaction),
		filters: make(map[Kind]Filter, 0),
	}

	for _, df := range defaultFilters {
		p.filters[df.Kind()] = df
	}

	for _, filter := range filters {
		// If it already exists, it will not be overwritten. 'defaultFilters'
		// needs to be run first, but overwriting it changes the order.
		if _, ok := p.filters[filter.Kind()]; ok {
			continue
		}

		p.filters[filter.Kind()] = filter
	}

	return p
}

// Returns a list of Kinds of processed/upcoming filters.
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
		if err := p.filtering(filter, eth, txs); err != nil {
			return err
		}
	}

	return nil
}

func (p *Purifier) Logs() map[Kind]map[common.Hash][]*types.Log {
	return p.logs
}

func (p *Purifier) Txs() map[Kind][]*types.Transaction {
	return p.txs
}

func (p *Purifier) filtering(filter Filter, eth *interaction.Interactor, txs []*types.Transaction) error {
	var err error

	switch f := filter.(type) {
	case transactionFilter:
		err = f.do(p, eth, txs)

	case logFilter:
		r := make([]*types.Log, 0)

		logs := p.logs[AllLogsType]
		if logs == nil {
			return errors.New("transaction has no logs, but log filter requested")
		}

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
