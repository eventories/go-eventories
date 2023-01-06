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
)

type Msg interface {
	PhaseID() [8]byte
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

	default:
		err = errors.New("invalid msg type")
	}

	return nil, err
}

func (a *ackMsg) PhaseID() [8]byte { return a.ID }
func (*ackMsg) Kind() byte         { return ackMsgType }

func (p *prepareMsg) PhaseID() [8]byte { return p.ID }
func (*prepareMsg) Kind() byte         { return prepareMsgType }

func (c *commitMsg) PhaseID() [8]byte { return c.ID }
func (*commitMsg) Kind() byte         { return commitMsgType }

func (a *abortMsg) PhaseID() [8]byte { return a.ID }
func (*abortMsg) Kind() byte         { return abortMsgType }
