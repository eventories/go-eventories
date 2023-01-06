package api

import "github.com/ethereum/go-ethereum/accounts/abi"

type ClientAPI struct {
	b Backend
}

func (c *ClientAPI) CurrentABI(name string) (*abi.ABI, error) {
	return c.b.ABI(name)
}
