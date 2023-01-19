package p2p

import (
	"bytes"
	"errors"
	"math/rand"
	"time"

	"github.com/eventories/go-eventories/database"
)

const (
	storage = 0 + iota
	wait
	prepare
	commit
)

type phase struct {
	seq        uint64
	id         [8]byte
	state      uint32
	key, value []byte

	cohorts []*peer
	msgCh   chan Msg

	do func(req Request, abort bool) error
}

func newPhase(seq uint64, cohorts []*peer, do func(Request, bool) error) (*phase, error) {
	phase := &phase{
		seq:     seq,
		id:      randomIDGenerator(),
		state:   wait,
		key:     make([]byte, 0),
		value:   make([]byte, 0),
		cohorts: cohorts,
		msgCh:   nil,
		do:      do,
	}

	if cohorts == nil {
		phase.msgCh = make(chan Msg, 1)
		return phase, nil
	}

	phase.msgCh = phase.readCohorts()

	return phase, nil
}

func (p *phase) prepare(key []byte, value []byte) error {
	if p.state != wait {
		return errors.New("invalid phase state")
	}

	var (
		id   = randomIDGenerator()
		want = len(p.cohorts) + 1 // Include node selfs.
		got  = 0
	)

	p.id = id
	p.key = key
	p.value = value

	p.msgCh <- &ackMsg{id}
	p.broadcast(&prepareMsg{p.seq, id, key, value})

	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	for {
		select {
		case msg := <-p.msgCh:
			switch msg.Kind() {
			case ackMsgType:
				ack := msg.(*ackMsg)
				if bytes.Equal(ack.ID[:], id[:]) {
					got++
				}

			case abortMsgType:
				abort := msg.(*abortMsg)
				if bytes.Equal(abort.ID[:], id[:]) {
					return errors.New("got abortMsg")
				}
			}

			if got == want {
				p.broadcast(&ackMsg{id})
				p.state = prepare
				return nil
			}

		case <-timer.C:
			p.broadcast(&abortMsg{id})
			return errors.New("timeout")
		}
	}
}

func (p *phase) commit(db database.Database) (err error) {
	if p.state != prepare {
		return errors.New("invalid phase state")
	}

	var (
		req  Request = nil
		want         = len(p.cohorts) + 1 // Include node selfs.
		got          = 0
	)

	defer func() {
		// Do revert.
		if err != nil {
			p.broadcast(&abortMsg{p.id})

			if req != nil {
				p.do(req, true)
			}

			db.Delete(p.key)
		}
	}()

	// Committing.
	req, err = decodeRequest(p.value)
	if err == nil {
		// Request
		if err = p.do(req, false); err != nil {
			return
		}
	}

	if err = db.Put(p.key, p.value); err != nil {
		return
	}

	p.msgCh <- &ackMsg{p.id}
	p.broadcast(&commitMsg{p.id})

	// Aggregate ack.
	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	for {
		select {
		case msg := <-p.msgCh:
			switch msg.Kind() {
			case ackMsgType:
				ack := msg.(*ackMsg)
				if bytes.Equal(ack.ID[:], p.id[:]) {
					got++
				}

			case abortMsgType:
				abort := msg.(*abortMsg)
				if bytes.Equal(abort.ID[:], p.id[:]) {
					return errors.New("got abortMsg")
				}
			}

			if got == want {
				p.broadcast(&ackMsg{p.id})
				return nil
			}

		case <-timer.C:
			return errors.New("timeout")
		}
	}
}

func (p *phase) readCohorts() chan Msg {
	resCh := make(chan Msg, len(p.cohorts))
	for _, p := range p.cohorts {
		go func(target *peer) {
			for {
				msg, err := target.readMsg()
				if err != nil {
					return
				}
				resCh <- msg
			}
		}(p)
	}
	return resCh
}

func (p *phase) broadcast(msg Msg) {
	for _, cohort := range p.cohorts {
		cohort.writeMsg(msg)
	}
	time.Sleep(25 * time.Millisecond)
}

func randomIDGenerator() [8]byte {
	var buf = [8]byte{}
	rng := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	rng.Read(buf[:])
	return buf
}
