package p2p

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/eventories/election"
	"github.com/eventories/go-eventories/core"
	"github.com/eventories/go-eventories/database"
)

type recoverReq struct {
	peer  *peer
	start uint64
	end   uint64
}

type Server struct {
	seq *core.Checkpoint
	// seq uint64

	logger *log.Logger

	listener *net.TCPListener
	election *election.Election
	db       database.Database

	local *phase // storage mode

	recover chan recoverReq
}

func NewServer(listener *net.TCPListener, election *election.Election, db database.Database) *Server {
	s := &Server{
		seq:      core.NewCheckpoint("seq"),
		listener: listener,
		election: election,
		db:       db,
		recover:  make(chan recoverReq),
	}

	s.local = &phase{
		id:      [8]byte{},
		state:   storage,
		cohorts: nil,
		key:     nil,
		value:   nil,
	}

	return s
}

func (s *Server) Run() error {
	if err := s.election.Run(); err != nil {
		return err
	}

	go s.loop()
	go s.acceptLoop()

	return nil
}

// temp
func (s *Server) loop() {
	for {
		select {
		case r := <-s.recover:
			go s.doSyncronization(r.peer, r.start, r.end)
		}
	}
}

func (s *Server) LocalAddr() *net.TCPAddr {
	return s.listener.Addr().(*net.TCPAddr)
}

func (s *Server) LeaderIP() net.IP {
	panic("not")
}

func (s *Server) Cluster() []string {
	return nil
	// return s.election.Cluster()
}

func (s *Server) Role() election.Role {
	return s.election.Role()
}

func (s *Server) Commit(ctx context.Context, key []byte, value []byte) error {
	if s.election.Role() != election.Leader {
		return errors.New("not leader")
	}

	var (
		cluster = s.election.Cluster()
		cohorts = make([]*peer, 0, len(cluster))
	)

	defer func() {
		for _, cohort := range cohorts {
			cohort.conn.Close()
		}
	}()

	for _, addr := range cluster {
		peer, err := DialTCP(addr, s)
		if err != nil {
			return err
		}
		cohorts = append(cohorts, peer)
	}

	phase, err := newPhase(s.seq.Checkpoint(), cohorts, s.doRequest)
	if err != nil {
		return err
	}

	if err := phase.prepare(key, value); err != nil {
		return err
	}

	if err := phase.commit(s.db); err != nil {
		return err
	}

	s.seq.Increase()

	return nil
}

// 2-phase-commit
// cluster manage
