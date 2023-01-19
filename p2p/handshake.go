package p2p

import (
	"errors"
	"net"
	"time"
)

func DialTCP(rawaddr string, backend *Server) (*peer, error) {
	addr, err := net.ResolveTCPAddr("tcp", rawaddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	peer := &peer{conn: conn}

	if err := peer.writeMsg(&handshakeMsg{}); err != nil {
		peer.conn.Close()
		return nil, err
	}

	msg, err := peer.readMsgWithTimeout(time.Second)
	if err != nil {
		return nil, err
	}

	ack, ok := msg.(*h_ackMsg)
	if !ok {
		return nil, errors.New("invalid message")
	}

	if backend.engine.BlockNumber() == ack.LatestBN {
		return peer, nil
	}

	if backend.engine.BlockNumber() < ack.LatestBN {
		if err := backend.doSyncronization(peer, backend.engine.BlockNumber(), ack.LatestBN); err != nil {
			peer.conn.Close()
			return nil, err
		}
	}

	panic("invalid protocl")
}

func (s *Server) handshakeHandle(peer *peer, h *handshakeMsg) {
	peer.writeMsg(&h_ackMsg{LatestBN: s.engine.BlockNumber()})
}
