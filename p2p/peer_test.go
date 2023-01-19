package p2p

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"testing"
)

func TestPeerRW(t *testing.T) {
	test := make(map[string]Msg)
	test["127.0.0.1:51131"] = &prepareMsg{[8]byte{1, 1, 1, 1, 1, 1, 1, 1}, []byte{1, 2, 3}, []byte{4, 5, 6}}
	test["127.0.0.1:51132"] = &commitMsg{[8]byte{2, 2, 2, 2, 2, 2, 2, 2}}
	test["127.0.0.1:51133"] = &abortMsg{[8]byte{3, 3, 3, 3, 3, 3, 3, 3}}

	server := "127.0.0.1:51134"
	mock := makeListener(server)
	defer mock.Close()

	for sender, msg := range test {
		go func(s string, m Msg) {
			addr, err := net.ResolveTCPAddr("tcp", server)
			if err != nil {
				panic("invalid addr")
			}

			local, err := net.ResolveTCPAddr("tcp", s)
			if err != nil {
				panic("invalid addr")
			}

			conn, err := net.DialTCP("tcp", local, addr)
			if err != nil {
				panic(err)
			}

			peer := &peer{conn: conn}
			peer.writeMsg(m)
			peer.conn.Close()
		}(sender, msg)
	}

	recv := 0

	for {
		conn, err := mock.AcceptTCP()
		if err != nil {
			panic(err)
		}

		peer := &peer{conn: conn}
		defer peer.conn.Close()

		msg, err := peer.readMsg()
		if err != nil {
			return
		}

		want, ok := test[peer.conn.RemoteAddr().String()]
		if !ok {
			panic("invalid sender")
		}

		if !compareMsg(msg, want) {
			panic(fmt.Sprintf("want: %v, got: %v", want, msg))
		}

		recv++

		if recv == len(test) {
			return
		}
	}
}

func compareMsg(a Msg, b Msg) bool {
	b1, _ := json.Marshal(a)
	b2, _ := json.Marshal(b)
	return bytes.Equal(b1, b2)
}

func makeListener(raw string) *net.TCPListener {
	addr, err := net.ResolveTCPAddr("tcp", raw)
	if err != nil {
		panic("mockServer generation failed: invalid raw address")
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic("mockServer generation failed")
	}
	return listener
}
