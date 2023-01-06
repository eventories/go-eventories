package p2p

import (
	"errors"
	"math/rand"
	"net"
	"time"
)

const (
	storage = 0 + iota
	wait
	prepare
	commit
)

type phase struct {
	id         [8]byte
	state      uint32
	key, value []byte

	cohorts []*peer
	msgCh   chan Msg

	do func(Request, bool) error
}

func newPhase(members []string, do func(Request, bool) error) (*phase, error) {
	if members == nil {
		return nil, errors.New("empty cluster")
	}
	peers := make([]*peer, 0, len(members))

	resCh := make(chan *peer, len(members))
	for _, member := range members {
		go func(rawaddr string) {
			addr, err := net.ResolveTCPAddr("tcp", rawaddr)
			if err != nil {
				resCh <- nil
				return
			}

			conn, err := net.DialTCP("tcp", nil, addr)
			if err != nil {
				resCh <- nil
				return
			}

			resCh <- &peer{conn}
		}(member)
	}

	var err error

	for i := 0; i < cap(resCh); i++ {
		peer := <-resCh
		if peer == nil {
			err = errors.New("invalid member address")
		}
		peers = append(peers, peer)
	}

	if err != nil {
		for _, peer := range peers {
			if peer != nil {
				peer.conn.Close()
			}
		}
		return nil, err
	}

	phase := &phase{
		id:      randomIDGenerator(),
		state:   wait,
		key:     make([]byte, 0),
		value:   make([]byte, 0),
		cohorts: peers,
		msgCh:   nil,
		do:      do,
	}

	phase.msgCh = phase.readCohorts()

	return phase, nil
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

func randomIDGenerator() [8]byte {
	var buf = [8]byte{}
	rng := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	rng.Read(buf[:])
	return buf
}
