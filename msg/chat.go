package msg

import (
	"crypto/ed25519"
	"mime"

	"github.com/bahner/go-ma/msg"
)

const (
	CHAT_SERIALIZATION = "plaintext"
	CHAT_MESSAGE_TYPE  = "chat"
)

var CHAT_CONTENT_TYPE_PARAMS = map[string]string{
	"type": CHAT_MESSAGE_TYPE,
}

// New creates a new Message instance
func Chat(
	from string,
	to string,
	content []byte,
	priv_key ed25519.PrivateKey) (*msg.Message, error) {

	mimeType := msg.CONTENT_TYPE + "+" + CHAT_SERIALIZATION
	contentType := mime.FormatMediaType(mimeType, CHAT_CONTENT_TYPE_PARAMS)

	return msg.New(from, to, contentType, content, priv_key)
}
