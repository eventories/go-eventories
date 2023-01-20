package filter

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	// Transfer
	tx1 = types.NewTransaction(0, common.HexToAddress("0x01"), big.NewInt(100), 0, big.NewInt(0), nil)
	// Call
	tx2 = types.NewTransaction(0, common.HexToAddress("0x02"), big.NewInt(0), 0, big.NewInt(0), []byte{1, 2, 3, 4})
	// Deploy
	tx3 = types.NewContractCreation(0, nil, 0, nil, []byte{4, 8, 7, 4})
)

func TestAllTransactionsFilter(t *testing.T) {
	filter := allTransactions{}

	if filter.Kind() != AllTransactionsType {
		t.Fatalf("AllTransactionsFilter invalid type, want: %v got: %v", EventLogType, filter.Kind())
	}

	p := &Purifier{
		logs: make(map[Kind]map[common.Hash][]*types.Log),
		txs:  make(map[Kind][]*types.Transaction),
	}

	if err := filter.do(p, nil, []*types.Transaction{tx1, tx2, tx3}); err != nil {
		t.Fatal(err)
	}

	res := p.Txs()[AllTransactionsType]

	if !bytes.Equal(tx1.Hash().Bytes(), res[0].Hash().Bytes()) {
		t.Fatalf("AllTransactionsFilter invalid txHash, want: %v got: %v", tx1.Hash().Bytes(), res[0].Hash().Bytes())
	}

	if !bytes.Equal(tx2.Hash().Bytes(), res[1].Hash().Bytes()) {
		t.Fatalf("AllTransactionsFilter invalid txHash, want: %v got: %v", tx2.Hash().Bytes(), res[1].Hash().Bytes())
	}

	if !bytes.Equal(tx3.Hash().Bytes(), res[2].Hash().Bytes()) {
		t.Fatalf("AllTransactionsFilter invalid txHash, want: %v got: %v", tx2.Hash().Bytes(), res[2].Hash().Bytes())
	}
}

func TestCoinTransferFilter(t *testing.T) {
	filter := coinTransfer{}

	if filter.Kind() != CoinTransferType {
		t.Fatalf("CoinTransferFilter invalid type, want: %v got: %v", EventLogType, filter.Kind())
	}

	p := &Purifier{
		logs: make(map[Kind]map[common.Hash][]*types.Log),
		txs:  make(map[Kind][]*types.Transaction),
	}

	if err := filter.do(p, nil, []*types.Transaction{tx1, tx2, tx3}); err != nil {
		t.Fatal(err)
	}

	res := p.Txs()[CoinTransferType]

	if !bytes.Equal(tx1.Hash().Bytes(), res[0].Hash().Bytes()) {
		t.Fatalf("TestCoinTransferFilter invalid txHash, want: %v got: %v", tx1.Hash(), res[0].Hash())
	}
}

func TestDeployFilter(t *testing.T) {
	filter := deploy{}

	if filter.Kind() != DeployType {
		t.Fatalf("DeployFilter invalid type, want: %v got: %v", EventLogType, filter.Kind())
	}

	p := &Purifier{
		logs: make(map[Kind]map[common.Hash][]*types.Log),
		txs:  make(map[Kind][]*types.Transaction),
	}

	if err := filter.do(p, nil, []*types.Transaction{tx1, tx2, tx3}); err != nil {
		t.Fatal(err)
	}

	res := p.Txs()[DeployType]
	if !bytes.Equal(tx3.Hash().Bytes(), res[0].Hash().Bytes()) {
		t.Fatalf("TestDeployFilter invalid txHash, want: %v got: %v", tx3.Hash(), res[0].Hash())
	}

	if !bytes.Equal(tx3.Data(), res[0].Data()) {
		t.Fatalf("TestDeployFilter invalid data, want: %v got: %v", tx3.Data(), res[0].Data())
	}
}

func TestSpectificDeploy(t *testing.T) {
	methods := make(map[string]abi.Method)
	methods["method1"] = abi.Method{ID: []byte{4}}
	methods["method2"] = abi.Method{ID: []byte{8}}
	methods["method3"] = abi.Method{ID: []byte{7}}

	target := &abi.ABI{
		Methods: methods,
	}

	filter := spectificDeploy{}

	if filter.Kind() != SpectificDeployType {
		t.Fatalf("SpectificDeployFilter invalid type, want: %v got: %v", EventLogType, filter.Kind())
	}

	filter.abi = target

	p := &Purifier{
		logs: make(map[Kind]map[common.Hash][]*types.Log),
		txs:  make(map[Kind][]*types.Transaction),
	}

	if err := filter.do(p, nil, []*types.Transaction{tx1, tx2, tx3}); err != nil {
		t.Fatal(err)
	}

	res := p.Txs()[SpectificDeployType]
	if !bytes.Equal(tx3.Hash().Bytes(), res[0].Hash().Bytes()) {
		t.Fatalf("TestSpectificDeploy invalid txHash, want: %v got: %v", tx3.Hash(), res[0].Hash())
	}

	// Phase 2
	methods = make(map[string]abi.Method)
	methods["method1"] = abi.Method{ID: []byte{14}}
	methods["method2"] = abi.Method{ID: []byte{18}}
	methods["method3"] = abi.Method{ID: []byte{17}}

	target = &abi.ABI{
		Methods: methods,
	}

	filter.abi = target

	p.logs = make(map[Kind]map[common.Hash][]*types.Log)
	p.txs = make(map[Kind][]*types.Transaction)

	if err := filter.do(p, nil, []*types.Transaction{tx1, tx2, tx3}); err != nil {
		t.Fatal(err)
	}

	res = p.Txs()[SpectificDeployType]

	if len(res) != 0 {
		t.Fatalf("TestSpectificDeploy invalid, want: [] got: %v", res[0])
	}
}
