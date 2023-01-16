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

		go func(peer *peer) {
			defer func() {
				time.Sleep(time.Second)
				peer.conn.Close()
			}()

			peer.registerProtocol("HANDSHAKE")
			peer.registerProtocol("SYNC")

			synchronized, err := s.doHandshake(peer)
			if err != nil {
				panic("noop")
			}

			if !synchronized {
				if err := s.doSyncronization(peer); err != nil {
					panic("noop")
				}
			}

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

	case h_ackMsgType:
		s.h_ackHandle(peer, msg.(*h_ackMsg))

	case syncReqMsgType:
		s.syncReqHandle(peer, msg.(*syncReqMsg))

	case syncResMsgType:
		s.syncResHandle(peer, msg.(*syncResMsg))

	case prepareMsgType:
		s.prepareHandle(peer, msg.(*prepareMsg))

	case commitMsgType:
		s.commitHandle(peer, msg.(*commitMsg))

	case abortMsgType:
		s.abortHandle(peer, msg.(*abortMsg))
	}
}
