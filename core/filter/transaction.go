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
	_ = transactionFilter(&allLogs{})
	_ = transactionFilter(&allTransactions{})
	_ = transactionFilter(&coinTransfer{})
	_ = transactionFilter(&deploy{})
	_ = transactionFilter(&spectificDeploy{})
)

type allLogs struct{}

func (a *allLogs) Kind() Kind { return AllLogsFilter }

func (a *allLogs) do(p *Purifier, eth *interaction.Interactor, txs []*types.Transaction) error {
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

type allTransactions struct{}

func (a *allTransactions) Kind() Kind { return AllTransactionsFilter }

func (a *allTransactions) do(p *Purifier, eth *interaction.Interactor, txs []*types.Transaction) error {
	p.txs[a.Kind()] = txs
	return nil
}

type coinTransfer struct{}

func (c *coinTransfer) Kind() Kind { return CoinTransferFilter }

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

func (d *deploy) Kind() Kind { return DeployFilter }

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
	ABI *abi.ABI
}

func (s *spectificDeploy) Kind() Kind { return SpectificDeployFilter }

func (s *spectificDeploy) do(p *Purifier, eth *interaction.Interactor, txs []*types.Transaction) error {
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
