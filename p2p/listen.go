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

		p := &peer{conn: conn, protocols: make(map[string]struct{})}
		go s.readLoop(p)

		p.registerProtocol("HANDSHAKE")
		p.registerProtocol("SYNC")

		bn, err := s.doHandshake(p)
		if err != nil {
			p.conn.Close()
			return
		}

		if s.engine.BlockNumber() < bn {
			if err := s.doSyncronization(p, s.engine.BlockNumber(), bn); err != nil {
				p.conn.Close()
				return
			}
		}

		go func(peer *peer) {
			defer func() {
				time.Sleep(time.Second)
				peer.conn.Close()
			}()

			// Start 2pc protocol.
			peer.registerProtocol("2PC")
		}(p)
	}
}

// is need to seperate by protocol?
func (s *Server) readLoop(peer *peer) {
	for {
		msg, err := peer.readMsg()
		if err != nil {
			return
		}

		if !peer.protocol(msg.Protocol()) {
			// Ignore not support protocol messages.
			continue
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
