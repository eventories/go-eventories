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

type (
	AllLogs struct{}

	AllTransactions struct{}

	CoinTransfer struct{}

	Deploy struct{}

	SpectificDeploy struct {
		ABI *abi.ABI
	}
)

//
func (a *AllLogs) Kind() Kind { return AllLogsFilter }

func (a *AllLogs) Do(p *Purifier, eth *interaction.Interactor, txs []*types.Transaction) error {
	r := make(map[common.Hash][]*types.Log)
	for _, tx := range txs {
		logs, err := eth.GetTransactionLogs(context.Background(), tx.Hash())
		if err != nil {
			return err
		}
		r[tx.Hash()] = logs
	}

	p.logs[a.Kind()] = r

	return nil
}

//
func (a *AllTransactions) Kind() Kind { return AllTransactionsFilter }

func (a *AllTransactions) Do(p *Purifier, eth *interaction.Interactor, txs []*types.Transaction) error {
	p.txs[a.Kind()] = txs
	return nil
}

//
func (c *CoinTransfer) Kind() Kind { return CoinTransferFilter }

func (c *CoinTransfer) Do(p *Purifier, eth *interaction.Interactor, txs []*types.Transaction) error {
	r := make([]*types.Transaction, 0)

	for _, tx := range txs {
		if tx.Value().Cmp(common.Big0) != 0 {
			r = append(r, tx)
		}
	}

	p.txs[c.Kind()] = r

	return nil
}

//
func (d *Deploy) Kind() Kind { return DeployFilter }

func (d *Deploy) Do(p *Purifier, eth *interaction.Interactor, txs []*types.Transaction) error {
	r := make([]*types.Transaction, 0)

	for _, tx := range txs {
		if tx.To() == nil && tx.Data() != nil {
			r = append(r, tx)
		}
	}

	p.txs[d.Kind()] = r

	return nil
}

//
func (s *SpectificDeploy) Kind() Kind { return SpectificDeployFilter }

func (s *SpectificDeploy) Do(p *Purifier, eth *interaction.Interactor, txs []*types.Transaction) error {
	if s.ABI == nil {
		return errors.New("must be set ABI")
	}

	r := make([]*types.Transaction, 0)

	for _, tx := range txs {
		if tx.To() == nil && tx.Data() != nil {
			data := tx.Data()
			methods := s.ABI.Methods

			satisfy := true
			for _, method := range methods {
				if !bytes.Contains(data, method.ID) {
					satisfy = false
				}
			}

			if satisfy {
				r = append(r, tx)
			}
		}
	}

	p.txs[s.Kind()] = r

	return nil
}
