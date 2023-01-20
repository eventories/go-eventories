package filter

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/eventories/go-eventories/core/interaction"
)

func TestPurifier(t *testing.T) {
	client, err := ethclient.Dial("https://rpc.ankr.com/eth")
	if err != nil {
		t.Fatal(err)
	}

	bn, _ := client.BlockNumber(context.Background())

	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(bn)-5))
	if err != nil {
		t.Fatal(err)
	}

	txs := block.Transactions()

	filter := New(DeployFilter(), CoinTransferFilter())

	if err := filter.Run(interaction.New(client), txs); err != nil {
		t.Fatal(err)
	}

	fmt.Println(filter.Txs())
	fmt.Println(filter.Logs())
}
