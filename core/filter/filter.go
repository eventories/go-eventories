package filter

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/eventories/go-eventories/core/interaction"
)

type Filter interface {
	Kind() Kind
}

type transactionFilter interface {
	Kind() Kind
	do(*Purifier, *interaction.Interactor, []*types.Transaction) error
}

type logFilter interface {
	Kind() Kind
	do(*Purifier, *interaction.Interactor, []*types.Log) error
}

func AllTransactionsFilter() Filter {
	return &allTransactions{}
}

func AllLogsFilter() Filter {
	return &allLogs{}
}

func CoinTransferFilter() Filter {
	return &coinTransfer{}
}

func DeployFilter() Filter {
	return &deploy{}
}

func SpectificDeployFilter(abi *abi.ABI) Filter {
	return &spectificDeploy{abi}
}

func AddressLogFilter(target common.Address) Filter {
	return &address{target}
}

func EventLogFilter(eventID common.Hash) Filter {
	return &event{eventID}
}
