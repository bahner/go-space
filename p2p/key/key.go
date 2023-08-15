package key

import (
	"crypto/rand"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/mr-tron/base58"
)

func PrintEd25519KeyAndExit() {

	privKey, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}
	marshalledPrivKey, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		log.Fatalf("Failed to marshal private key: %v", err)
	}
	encodedPrivKey := base58.Encode(marshalledPrivKey)
	fmt.Println(encodedPrivKey)

	os.Exit(0)

}

func CreateIdentity(privKey string) crypto.PrivKey {

	// Decode the secret key
	decoded, err := base58.Decode(privKey)
	if err != nil {
		log.Fatalf("Failed to decode base58 secret key: %v", err)
	}
	p, err := crypto.UnmarshalPrivateKey(decoded)
	if err != nil {
		log.Fatalf("Failed to unmarshal private key: %v", err)
	}

	return p

}
