package cribbage

import (
	"bytes"
	"encoding/gob"
)

const (
	TO_CRIB = 0
)

type Message struct {
	Type int
	Val  uint
}

// Encode converts a message into a byte slice
func Encode(m *Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(m)
	if err != nil {
		return nil, err
	}
	encoded := buf.Bytes()
	return encoded, nil
}

// Decode converts a byte slice into a Message
func Decode(raw []byte) (*Message, error) {
	m := &Message{}
	buf := bytes.NewBuffer(raw)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
