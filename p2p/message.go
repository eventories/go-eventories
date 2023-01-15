package p2p

import (
	"encoding/json"
	"errors"
)

const (
	ackMsgType = byte(1) + iota
	prepareMsgType
	commitMsgType
	abortMsgType

	handshakeMsgType
	h_ackMsgType
	h_nackMsgType

	syncReqMsgType
	syncResMsgType
)

type Msg interface {
	Kind() byte
}

// 2-phase-commit
type (
	ackMsg struct {
		ID [8]byte
	}

	prepareMsg struct {
		ID    [8]byte
		Key   []byte
		Value []byte
	}

	commitMsg struct {
		ID [8]byte
	}

	abortMsg struct {
		ID [8]byte
	}
)

// handkshake
type (
	handshakeMsg struct {
		// Sequence uint64
		// LatestBN uint64
	}

	h_ackMsg struct {
		Sequence uint64
		LatestBN uint64
	}
)

// sync
type (
	syncReqMsg struct{}

	syncResMsg struct{}
)

func encodeMsg(msg Msg) ([]byte, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	rb := make([]byte, len(b)+1)
	rb[0] = msg.Kind()
	copy(rb[1:], b)

	return rb, nil
}

func decodeMsg(b []byte) (Msg, error) {
	if len(b) < 1 {
		return nil, errors.New("too short")
	}

	var err error

	switch b[0] {
	case ackMsgType:
		var m ackMsg
		if err = json.Unmarshal(b[1:], &m); err == nil {
			return &m, nil
		}

	case prepareMsgType:
		var m prepareMsg
		if err = json.Unmarshal(b[1:], &m); err == nil {
			return &m, nil
		}

	case commitMsgType:
		var m commitMsg
		if err = json.Unmarshal(b[1:], &m); err == nil {
			return &m, nil
		}

	case abortMsgType:
		var m abortMsg
		if err = json.Unmarshal(b[1:], &m); err == nil {
			return &m, nil
		}

	case handshakeMsgType:
		var m handshakeMsg
		if err = json.Unmarshal(b[1:], &m); err == nil {
			return &m, nil
		}

	case h_ackMsgType:
		var m h_ackMsg
		if err = json.Unmarshal(b[1:], &m); err == nil {
			return &m, nil
		}

	case syncReqMsgType:
		var m syncReqMsg
		if err = json.Unmarshal(b[1:], &m); err == nil {
			return &m, nil
		}

	case syncResMsgType:
		var m syncResMsg
		if err = json.Unmarshal(b[1:], &m); err == nil {
			return &m, nil
		}

	default:
		err = errors.New("invalid msg type")
	}

	return nil, err
}

func (*ackMsg) Kind() byte       { return ackMsgType }
func (*prepareMsg) Kind() byte   { return prepareMsgType }
func (*commitMsg) Kind() byte    { return commitMsgType }
func (*abortMsg) Kind() byte     { return abortMsgType }
func (*handshakeMsg) Kind() byte { return handshakeMsgType }
func (*h_ackMsg) Kind() byte     { return h_ackMsgType }
func (*syncReqMsg) Kind() byte   { return syncReqMsgType }
func (*syncResMsg) Kind() byte   { return syncResMsgType }
