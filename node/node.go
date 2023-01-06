package node

import (
	"net/rpc"

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
