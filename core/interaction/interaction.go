package interaction

import (
	"bytes"
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Interactor struct {
	client *ethclient.Client
}

func (i *Interactor) ChainID(ctx context.Context) (*big.Int, error) {
	return i.client.ChainID(ctx)
}

func (i *Interactor) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	return i.client.BlockNumber(ctx)
}

func (i *Interactor) GetTransactionsByNumber(ctx context.Context, blockNumber uint64) ([]*types.Transaction, error) {
	block, err := i.client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, err
	}
	return block.Transactions(), nil
}

func (i *Interactor) GetTransactionHashesByNumber(ctx context.Context, blockNumber uint64) ([]common.Hash, error) {
	hashes := make([]common.Hash, 0)

	block, err := i.client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, err
	}

	for _, hash := range block.Transactions() {
		hashes = append(hashes, hash.Hash())
	}

	return hashes, nil
}

func (i *Interactor) GetTransactionLogs(ctx context.Context, txHash common.Hash) ([]*types.Log, error) {
	receipt, err := i.client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}
	return receipt.Logs, nil
}

func (i *Interactor) GetCode(ctx context.Context, address *common.Address) ([]byte, error) {
	var (
		code []byte
		err  error
	)

	code, err = i.client.CodeAt(ctx, *address, nil)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(code, make([]byte, 0)) {
		code, err = i.eip1967(ctx, address)
	}

	return code, err
}

func (i *Interactor) eip1967(ctx context.Context, address *common.Address) ([]byte, error) {
	return i.client.StorageAt(ctx, *address, common.HexToHash("0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"), nil)
}
