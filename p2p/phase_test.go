package p2p

import (
	"testing"

	"github.com/eventories/go-eventories/database/memdb"
)

func TestSoloPhase(t *testing.T) {
	db, _ := memdb.New()

	key := []byte{1, 2, 3}
	value := []byte{4, 5, 6}

	phase, err := newPhase(nil, doRequest)
	if err != nil {
		t.Fatal(err)
	}

	if err := phase.prepare(key, value); err != nil {
		t.Fatal(err)
	}

	if err := phase.commit(db); err != nil {
		t.Fatal(err)
	}

	phase = nil

	d, err := db.Get(key)
	if err != nil || d == nil {
		t.Fatal(err)
	}
}

func doRequest(req Request, abort bool) error {
	return nil
}
