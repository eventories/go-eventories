package p2p

import (
	"encoding/json"
	"errors"
)

const (
	addMemberType = byte(1) + iota
	delMemberType
)

type Request interface {
	Kind() byte
}

var requestPrefix = []byte("r-e-q-")

// m
type (
	addMember struct {
		Member string
	}

	delMember struct {
		Member string
	}
)

func encodeRequest(req Request) ([]byte, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	rb := make([]byte, 1+len(requestPrefix)+len(b))
	rb[0] = req.Kind()

	idx := 1
	idx += copy(rb[idx:idx+len(requestPrefix)], requestPrefix)
	idx += copy(rb[idx:], b)

	return rb, nil
}

func decodeRequest(b []byte) (Request, error) {
	if len(b) < 1 {
		return nil, errors.New("too short")
	}

	var err error

	prefixIdx := 1 + len(requestPrefix)
	switch b[0] {
	case addMemberType:
		var add addMember
		if err = json.Unmarshal(b[prefixIdx:], &add); err == nil {
			return &add, nil
		}

	case delMemberType:
		var del delMember
		if err = json.Unmarshal(b[prefixIdx:], &del); err == nil {
			return &del, nil
		}

	default:
		err = errors.New("invalid request type")
	}

	return nil, err
}

func (*addMember) Kind() byte { return addMemberType }
func (*delMember) Kind() byte { return delMemberType }
