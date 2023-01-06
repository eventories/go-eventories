package api

import "github.com/ethereum/go-ethereum/accounts/abi"

type Backend interface {
	ABI(string) (*abi.ABI, error)
}
