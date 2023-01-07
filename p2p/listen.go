package p2p

import (
	"bytes"
	"errors"
	"time"
)

const (
	ackTimeout = time.Second
)

var (
	errInvalidID = errors.New("invalid ID")
)

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			continue
		}

		go s.handle(&peer{conn})
	}
}

func (s *Server) handle(peer *peer) {
	defer func() {
		time.Sleep(time.Second)
		peer.conn.Close()
	}()

	for {
		msg, err := peer.readMsg()
		if err != nil {
			return
		}

		switch msg.Kind() {
		case prepareMsgType:
			prepare, ok := msg.(*prepareMsg)
			if !ok {
				return
			}

			s.prepareHandle(peer, prepare)

		case commitMsgType:
			commit, ok := msg.(*commitMsg)
			if !ok {
				return
			}

			s.commitHandle(peer, commit)

		case abortMsgType:
			abort, ok := msg.(*abortMsg)
			if !ok {
				return
			}

			s.abortHandle(peer, abort)
		}
	}
}

func (s *Server) prepareHandle(peer *peer, prepare *prepareMsg) {
	s.local.id = prepare.ID
	s.local.key = prepare.Key
	s.local.value = prepare.Value

	var (
		err error
		msg Msg
	)

	defer func() {
		if err != nil {
			peer.writeMsg(&abortMsg{s.local.id})
			s.local.id = [8]byte{}
			s.local.key = nil
			s.local.value = nil
		}
	}()

	if err = peer.writeMsg(&ackMsg{s.local.id}); err != nil {
		return
	}

	msg, err = peer.readMsgWithTimeout(ackTimeout)
	if err != nil {
		return
	}

	if prepare.ID != msg.(*ackMsg).ID {
		err = errInvalidID
		return
	}
}

func (s *Server) commitHandle(peer *peer, commit *commitMsg) {
	if !bytes.Equal(s.local.id[:], commit.ID[:]) {
		return
	}

	var (
		err error
		msg Msg
		req Request
	)

	defer func() {
		if err != nil {
			// Revert
			if req != nil {
				s.doRequest(req, true)
			}
			if req == nil {
				s.db.Delete(s.local.key)
			}
		}
	}()

	// Request
	if bytes.Contains(s.local.value, requestPrefix) {
		req, err = decodeRequest(s.local.value)
		if err != nil {
			return
		}
		if err := s.doRequest(req, false); err != nil {
			return
		}
	} else {
		// Data
		if err := s.db.Put(s.local.key, s.local.value); err != nil {
			return
		}
	}

	if err = peer.writeMsg(&ackMsg{s.local.id}); err != nil {
		return
	}

	msg, err = peer.readMsgWithTimeout(ackTimeout)
	if err != nil {
		return
	}

	if s.local.id != msg.(*ackMsg).ID {
		err = errInvalidID
		return
	}
}

func (s *Server) abortHandle(peer *peer, abort *abortMsg) {
	if !bytes.Equal(s.local.id[:], abort.ID[:]) {
		return
	}

	// Revert
	if bytes.Contains(s.local.value, requestPrefix) {
		// Request
		req, err := decodeRequest(s.local.value)
		if err != nil {
			return
		}
		if err := s.doRequest(req, true); err != nil {
			return
		}
	} else {
		// Data
		if err := s.db.Delete(s.local.key); err != nil {
			return
		}
	}
}

// Request should also contain logic for revert that take that action.
func (s *Server) doRequest(req Request, abort bool) error {
	switch req.Kind() {
	case addMemberType:
		add := req.(*AddMemberReq)

		if abort {
			return s.election.DelMember(add.Member)
		}

		return s.election.AddMember(add.Member)

	case delMemberType:
		del := req.(*DelMemberReq)

		if abort {
			return s.election.AddMember(del.Member)
		}

		return s.election.DelMember(del.Member)
	}

	return errors.New("invalid request type")
}
