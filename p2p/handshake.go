package p2p

import (
	"errors"
	"time"
)

func (s *Server) doHandshake(peer *peer) (synchronized bool, err error) {
	var (
		seq = s.seq
		bn  = s.bn
	)

	if err := peer.writeMsg(&handshakeMsg{ /*seq, bn*/ }); err != nil {
		return false, err
	}

	msg, err := peer.readMsgWithTimeout(time.Second)
	if err != nil {
		return false, err
	}

	ack, ok := msg.(*h_ackMsg)
	if !ok {
		return false, errors.New("invalid message")
	}

	if seq != ack.Sequence {
		return false, errors.New("invalid sequence")
	}

	return (bn == ack.LatestBN), nil
}

func (s *Server) handshakeHandle(peer *peer, h *handshakeMsg) {
	panic("not impl")
}

func (s *Server) h_ackHandle(peer *peer, ack *h_ackMsg) {
	panic("not impl")
}
