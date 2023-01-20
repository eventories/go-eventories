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

/*
	CoinTransferFilter = Kind(0) + iota
	DeployFilter
	SpectificDeployFilter
	AllTransactionsFilter
	AllLogsFilter

	// Log
	AddressLogsFilter
	EventLogsFilter
*/
func NewAllTransactionsFilter() Filter {
	return &allTransactions{}
}

func NewAllLogsFilter() Filter {
	return &allLogs{}
}

func NewCoinTransferFilter() Filter {
	return &coinTransfer{}
}

func NewDeployFilter() Filter {
	return &deploy{}
}

func NewSpectificDeployFilter(abi *abi.ABI) Filter {
	return &spectificDeploy{abi}
}

func NewAddressLogFilter(target common.Address) Filter {
	return &address{target}
}

func NewEventLogFilter(eventID common.Hash) Filter {
	return &event{eventID}
}
