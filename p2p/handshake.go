package p2p

import (
	"errors"
	"time"
)

func (s *Server) doHandshake(peer *peer) (uint64, error) {
	if err := peer.writeMsg(&handshakeMsg{ /*seq, bn*/ }); err != nil {
		return 0, err
	}

	msg, err := peer.readMsgWithTimeout(time.Second)
	if err != nil {
		return 0, err
	}

	ack, ok := msg.(*h_ackMsg)
	if !ok {
		return 0, errors.New("invalid message")
	}

	return ack.LatestBN, nil
}

func (s *Server) handshakeHandle(peer *peer, h *handshakeMsg) {
	peer.writeMsg(&h_ackMsg{LatestBN: s.engine.BlockNumber()})
}
