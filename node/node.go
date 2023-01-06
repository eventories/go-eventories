package node

import (
	"context"
	"fmt"
	"net/rpc"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/eventories/go-eventories/core"
	"github.com/eventories/go-eventories/core/table"
	"github.com/eventories/go-eventories/p2p"
)

type Node struct {
	engine core.Engine

	tab  *table.Table
	srv  *p2p.Server
	http *rpc.Server
}

func (n *Node) ABI(name string) (*abi.ABI, error) {
	abi := n.tab.ABI(name)
	if abi == nil {
		return nil, fmt.Errorf("non exist abi %s", name)
	}
	return abi, nil
}

func (n *Node) Members() []string {
	return n.srv.Cluster()
}

func (n *Node) AddMember(rawaddr string) error {
	req := &p2p.AddMemberReq{Member: rawaddr}
	b, err := p2p.EncodeRequest(req)
	if err != nil {
		return err
	}
	return n.srv.Commit(context.Background(), nil, b)
}

func (n *Node) DelMember(rawaddr string) error {
	req := &p2p.DelMemberReq{Member: rawaddr}
	b, err := p2p.EncodeRequest(req)
	if err != nil {
		return err
	}
	return n.srv.Commit(context.Background(), nil, b)
}
