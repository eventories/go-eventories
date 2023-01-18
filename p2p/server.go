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

type Server struct {
	engine *core.Fetcher

	logger *log.Logger

	listener *net.TCPListener
	election *election.Election
	db       database.Database

	local *phase // storage mode
}

func NewServer(listener *net.TCPListener, election *election.Election, db database.Database) *Server {
	s := &Server{
		engine:   nil,
		listener: listener,
		election: election,
		db:       db,
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

	go s.acceptLoop()

	return nil
}

func (s *Server) LocalAddr() *net.TCPAddr {
	return s.listener.Addr().(*net.TCPAddr)
}

func (s *Server) LeaderIP() net.IP {
	panic("not")
}

func (s *Server) Cluster() []string {
	return s.election.Cluster()
}

func (s *Server) Role() election.Role {
	return s.election.Role()
}

func (s *Server) Commit(ctx context.Context, key []byte, value []byte) error {
	if s.election.Role() != election.Leader {
		return errors.New("not leader")
	}

	phase, err := newPhase(s.election.Cluster(), s.doRequest)
	if err != nil {
		return err
	}

	if err := phase.prepare(key, value); err != nil {
		return err
	}

	if err := phase.commit(s.db); err != nil {
		return err
	}

	return nil
}

// 2-phase-commit
// cluster manage
