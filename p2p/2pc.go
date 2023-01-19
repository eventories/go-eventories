package p2p

import (
	"bytes"
	"errors"
)

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

	// Checks that is received abortMsg.
	ack, ok := msg.(*ackMsg)
	if !ok {
		// Occur panic if it is not an abortMsg.
		s.handle(peer, msg.(*abortMsg))
		return
	}

	if prepare.ID != ack.ID {
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
			peer.writeMsg(&abortMsg{commit.ID})
			// Revert
			if req != nil {
				s.doRequest(req, true)
			}

			s.db.Delete(s.local.key)
		}
	}()

	req, err = decodeRequest(s.local.value)
	if err == nil {
		// Request
		if err := s.doRequest(req, false); err != nil {
			return
		}
	}

	if err := s.db.Put(s.local.key, s.local.value); err != nil {
		return
	}

	if err = peer.writeMsg(&ackMsg{s.local.id}); err != nil {
		return
	}

	msg, err = peer.readMsgWithTimeout(ackTimeout)
	if err != nil {
		return
	}

	// Checks that is received abortMsg.
	ack, ok := msg.(*ackMsg)
	if !ok {
		// Occur panic if it is not an abortMsg.
		s.handle(peer, msg.(*abortMsg))
		return
	}

	if s.local.id != ack.ID {
		err = errInvalidID
		return
	}
}

func (s *Server) abortHandle(peer *peer, abort *abortMsg) {
	if !bytes.Equal(s.local.id[:], abort.ID[:]) {
		return
	}

	// Do revert.
	req, err := decodeRequest(s.local.value)
	if err == nil {
		// Request
		s.doRequest(req, true)
	} else {
		// Data
		s.db.Delete(s.local.key)
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
