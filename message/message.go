package message

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

const Version = "1"

type Message struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Version   string    `json:"version"`
	Data      []byte    `json:"data"`
	Created   time.Time `json:"created"`
	Received  time.Time `json:"received"`
	Signature []byte    `json:"signature"`
}

func New(from string, to string, data []byte) *Message {
	return &Message{
		From:    from,
		To:      to,
		Version: Version,
		Created: time.Now(),
		Data:    data,
	}
}

func Marshal(m *Message) ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) UnsignedMessage() *Message {

	// This returns
	c := &Message{
		To:      m.To,
		From:    m.From,
		Data:    m.Data,
		Created: m.Created,
		Version: m.Version,
	}

	return c
}

func Unmarshal(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (m *Message) Sign(privateKey crypto.PrivKey) error {
	data, err := Marshal(m.UnsignedMessage())
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	sig, err := privateKey.Sign(data)
	if err != nil {
		return fmt.Errorf("failed to sign message: %v", err)
	}
	m.Signature = sig
	return nil
}

func (m *Message) Verify() (bool, error) {
	publicKey, err := extractPublicKeyFromIPNS(m.From)
	if err != nil {
		return false, err
	}

	data, err := Marshal(m.UnsignedMessage())
	if err != nil {
		return false, fmt.Errorf("failed to marshal message: %v", err)
	}

	isValid, err := publicKey.Verify(data, m.Signature)
	if err != nil {
		return false, fmt.Errorf("error verifying message: %v", err)
	}

	return isValid, nil
}

func extractPublicKeyFromIPNS(ipns string) (crypto.PubKey, error) {
	pid, err := peer.Decode(ipns)
	if err != nil {
		return nil, err
	}

	return pid.ExtractPublicKey()
}
