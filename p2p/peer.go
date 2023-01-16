package p2p

import (
	"errors"
	"net"
	"sync"
	"time"
)

type peer struct {
	conn *net.TCPConn

	mu        sync.Mutex
	protocols map[string]struct{}
}

func (p *peer) registerProtocol(protocol string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.protocols[protocol] = struct{}{}
}

func (p *peer) protocol(protocol string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, ok := p.protocols[protocol]
	return ok
}

func (p *peer) writeMsg(msg Msg) error {
	b, err := encodeMsg(msg)
	if err != nil {
		return err
	}
	_, err = p.conn.Write(b)
	return err
}

func (p *peer) readMsg() (Msg, error) {
	b := make([]byte, 2048)
	n, err := p.conn.Read(b)
	if err != nil {
		return nil, err
	}
	return decodeMsg(b[:n])
}

func (p *peer) readMsgWithTimeout(timeout time.Duration) (Msg, error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	recv := make(chan Msg, 1)
	go func() {
		msg, err := p.readMsg()
		if err != nil {
			return
		}
		recv <- msg
	}()

	for {
		select {
		case msg := <-recv:
			return msg, nil
		case <-timer.C:
			return nil, errors.New("timeout")
		}
	}
}
