package did

import (
	"crypto/ed25519"
	"errors"
	"time"
)

// AddPubKey adds a public key to the DID
func (doc *DIDDocument) AddPubKey(pubKeyID, pubKeyType string, pubKey ed25519.PublicKey) error {
	// Your original code for adding a public key goes here
	return nil
}

// RemovePubKey removes a public key from the DID
func (doc *DIDDocument) RemovePubKey(pubKeyID string) error {
	// Your original code for removing a public key goes here
	return nil
}

// GetNewestKey retrieves the newest key in the DID, based on expiry time.
func (doc *DIDDocument) GetNewestKey() (*VerificationMethod, error) {
	var newest *VerificationMethod
	for _, method := range doc.VerificationMethod {
		if newest == nil || method.Expiry > newest.Expiry {
			newest = &method
		}
	}
	if newest == nil {
		return nil, errors.New("no valid keys found")
	}
	return newest, nil
}

// GetAllKeys retrieves all keys from the DID.
func (doc *DIDDocument) GetAllKeys() []VerificationMethod {
	return doc.VerificationMethod
}

// GetValidKeys retrieves all currently valid keys from the DID.
func (doc *DIDDocument) GetValidKeys() ([]VerificationMethod, error) {
	var validKeys []VerificationMethod
	now := time.Now().Unix() // Current Unix timestamp
	for _, method := range doc.VerificationMethod {
		if method.Expiry == 0 || method.Expiry > now {
			validKeys = append(validKeys, method)
		}
	}
	if len(validKeys) == 0 {
		return nil, errors.New("no valid keys found")
	}
	return validKeys, nil
}
