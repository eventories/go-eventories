package p2p

import (
	"errors"
	"time"
)

func doHandshake(peer *peer, backend *Server) (synchronized bool, err error) {
	var (
		seq = backend.seq
		bn  = backend.bn
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
