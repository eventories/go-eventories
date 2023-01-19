package p2p

import (
	"errors"
	"time"
)

var (
	ackTimeout = time.Second

	errInvalidID = errors.New("invalid ID")
)

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			continue
		}

		go func(peer *peer) {
			defer func() {
				time.Sleep(time.Second)
				peer.conn.Close()
			}()

			s.readLoop(peer)
		}(&peer{conn: conn})
	}
}

// is need to seperate by protocol?
func (s *Server) readLoop(peer *peer) {
	for {
		msg, err := peer.readMsg()
		if err != nil {
			return
		}

		s.handle(peer, msg)
	}
}

func (s *Server) handle(peer *peer, msg Msg) {
	switch msg.Kind() {
	case handshakeMsgType:
		s.handshakeHandle(peer, msg.(*handshakeMsg))

	case syncReqMsgType:
		s.syncReqHandle(peer, msg.(*syncReqMsg))

	case prepareMsgType:
		s.prepareHandle(peer, msg.(*prepareMsg))

	case commitMsgType:
		s.commitHandle(peer, msg.(*commitMsg))

	case abortMsgType:
		s.abortHandle(peer, msg.(*abortMsg))
	}
}
