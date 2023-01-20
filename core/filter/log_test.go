package filter

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestAddressLogFilter(t *testing.T) {
	logs := []*types.Log{
		{
			TxHash:  common.HexToHash("0x01"),
			Address: common.HexToAddress("0x02"),
		},
		{
			TxHash:  common.HexToHash("0x03"),
			Address: common.HexToAddress("0x04"),
		},
	}

	target := common.HexToAddress("0x02")
	filter := address{target}

	if filter.Kind() != AddressLogType {
		t.Fatalf("AddressLogFilter invalid type, want: %v got: %v", AddressLogType, filter.Kind())
	}

	p := &Purifier{
		logs: make(map[Kind]map[common.Hash][]*types.Log),
		txs:  make(map[Kind][]*types.Transaction),
	}

	if err := filter.do(p, nil, logs); err != nil {
		t.Fatal(err)
	}

	res := p.Logs()[AddressLogType]

	for addr, logs := range res {
		if !bytes.Equal(target.Hash().Bytes(), addr.Bytes()) {
			t.Fatalf("AddressLogFilter invalid hash, want: %v got: %v", target, addr)
		}

		if !bytes.Equal(common.HexToHash("0x01").Bytes(), logs[0].TxHash.Bytes()) {
			t.Fatalf("AddressLogFilter invalid logs, want: %v got: %v", common.HexToHash("0x01"), logs[0].TxHash)
		}
	}
}

func TestEventLogFilter(t *testing.T) {
	logs := []*types.Log{
		{
			TxHash:  common.HexToHash("0x01"),
			Address: common.HexToAddress("0x02"),
			Topics: []common.Hash{
				common.HexToHash("0x03"),
			},
		},
		{
			TxHash:  common.HexToHash("0x04"),
			Address: common.HexToAddress("0x05"),
			Topics: []common.Hash{
				common.HexToHash("0x06"),
			},
		},
		{
			TxHash:  common.HexToHash("0x07"),
			Address: common.HexToAddress("0x08"),
			Topics:  nil,
		},
	}

	filter := event{}

	filter.id = common.HexToHash("0x02")

	if filter.Kind() != EventLogType {
		t.Fatalf("AddressLogFilter invalid type, want: %v got: %v", EventLogType, filter.Kind())
	}

	p := &Purifier{
		logs: make(map[Kind]map[common.Hash][]*types.Log),
		txs:  make(map[Kind][]*types.Transaction),
	}

	if err := filter.do(p, nil, logs); err != nil {
		t.Fatal(err)
	}

	res := p.Logs()[EventLogType]

	for hash, logs := range res {
		if !bytes.Equal(filter.id.Bytes(), hash.Bytes()) {
			t.Fatalf("AddressLogFilter invalid hash, want: %v got: %v", filter.id, hash)
		}

		if len(logs) != 0 {
			t.Fatalf("EventLogFilter invalid logs, want: [] got: %v", logs)
		}
	}

	// Phase 2
	filter.id = common.HexToHash("0x03")

	p.logs = make(map[Kind]map[common.Hash][]*types.Log)
	p.txs = make(map[Kind][]*types.Transaction)

	if err := filter.do(p, nil, logs); err != nil {
		t.Fatal(err)
	}

	res = p.Logs()[EventLogType]

	for hash, logs := range res {
		if !bytes.Equal(filter.id.Bytes(), hash.Bytes()) {
			t.Fatalf("AddressLogFilter invalid hash, want: %v got: %v", filter.id, hash)
		}

		if !bytes.Equal(common.HexToHash("0x01").Bytes(), logs[0].TxHash.Bytes()) {
			t.Fatalf("EventLogFilter invalid logs, want: %v got: %v", common.HexToHash("0x01"), logs[0].TxHash)
		}
	}
}
