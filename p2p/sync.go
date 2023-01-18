package p2p

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

func (s *Server) doSyncronization(peer *peer, start uint64, end uint64) error {
	for ; start <= end; start++ {
		key := []byte(fmt.Sprintf("block-%d", start))
		if err := peer.writeMsg(&syncReqMsg{key}); err != nil {
			return err
		}
		msg, err := peer.readMsgWithTimeout(time.Second)
		if err != nil {
			return err
		}
		res := msg.(*syncResMsg)
		if !bytes.Equal(res.Key, key) {
			return errors.New("invalid key")
		}

		if err := s.db.Put(res.Key, res.Value); err != nil {
			return err
		}

		if err := s.engine.SetBlockNumber(start); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) syncReqHandle(peer *peer, req *syncReqMsg) {
	b, err := s.db.Get(req.Key)
	if err != nil {
		return
	}
	if err := peer.writeMsg(&syncResMsg{Key: req.Key, Value: b}); err != nil {
		return
	}
}
