package filter

import (
	"bytes"
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/eventories/go-eventories/core/interaction"
)

var (
	batchSize = 20

	_ = transactionFilter(&allLogs{})
	_ = transactionFilter(&allTransactions{})
	_ = transactionFilter(&coinTransfer{})
	_ = transactionFilter(&deploy{})
	_ = transactionFilter(&spectificDeploy{})
)

type allLogs struct{}

func (a *allLogs) Kind() Kind { return AllLogsType }

func (a *allLogs) do(p *Purifier, eth *interaction.Interactor, txs []*types.Transaction) error {
	ch := make(chan map[common.Hash][]*types.Log, (len(txs)/batchSize)+1)

	for i := 0; i < len(txs); i += batchSize {
		end := i + batchSize
		if end >= len(txs) {
			end = len(txs) - 1
		}

		go func(batch []*types.Transaction) {
			temp := make(map[common.Hash][]*types.Log)

			for _, tx := range txs {
				logs, err := eth.GetTransactionLogs(context.Background(), tx.Hash())
				if err != nil {
					ch <- nil
					return
				}

				// if len(logs) != 0 {
				// 	temp[tx.Hash()] = logs
				// }
				// Includes 'nil'
				temp[tx.Hash()] = logs
			}
			ch <- temp
		}(txs[i:end])
	}

	r := make(map[common.Hash][]*types.Log)
	for i := 0; i < cap(ch); i++ {
		res := <-ch
		if res == nil {
			return errors.New("allLogs failure")
		}

		for k, v := range res {
			r[k] = v
		}
	}

	p.logs[a.Kind()] = r

	return nil
}

type allTransactions struct{}

func (a *allTransactions) Kind() Kind { return AllTransactionsType }

func (a *allTransactions) do(p *Purifier, eth *interaction.Interactor, txs []*types.Transaction) error {
	p.txs[a.Kind()] = txs
	return nil
}

type coinTransfer struct{}

func (c *coinTransfer) Kind() Kind { return CoinTransferType }

func (c *coinTransfer) do(p *Purifier, eth *interaction.Interactor, txs []*types.Transaction) error {
	r := make([]*types.Transaction, 0)

	for _, tx := range txs {
		if tx.Value().Cmp(common.Big0) != 0 {
			r = append(r, tx)
		}
	}

	p.txs[c.Kind()] = r

	return nil
}

type deploy struct{}

func (d *deploy) Kind() Kind { return DeployType }

func (d *deploy) do(p *Purifier, eth *interaction.Interactor, txs []*types.Transaction) error {
	r := make([]*types.Transaction, 0)

	for _, tx := range txs {
		if tx.To() == nil && tx.Data() != nil {
			r = append(r, tx)
		}
	}

	p.txs[d.Kind()] = r

	return nil
}

type spectificDeploy struct {
	abi *abi.ABI
}

func (s *spectificDeploy) Kind() Kind { return SpectificDeployType }

func (s *spectificDeploy) do(p *Purifier, eth *interaction.Interactor, txs []*types.Transaction) error {
	if s.abi == nil {
		return errors.New("must be set ABI")
	}

	r := make([]*types.Transaction, 0)

	for _, tx := range txs {
		if tx.To() == nil && tx.Data() != nil {
			data := tx.Data()
			methods := s.abi.Methods

			satisfied := true
			for _, method := range methods {
				if !bytes.Contains(data, method.ID) {
					satisfied = false
					break
				}
			}

			if satisfied {
				r = append(r, tx)
			}
		}
	}

	p.txs[s.Kind()] = r

	return nil
}
