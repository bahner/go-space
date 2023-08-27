package did

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
)

// UnsignedDID generates the unsigned DID
func (doc *DID) UnsignedDID() ([]byte, error) {
	// Remove the signature field before serializing
	doc.Signature = ""
	unsignedDID, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	return unsignedDID, nil
}

// Sign signs the DID
func (doc *DID) Sign(privKey ed25519.PrivateKey) error {
	unsignedDID, err := doc.UnsignedDID()
	if err != nil {
		return err
	}
	// Sign the unsigned DID Document
	signature := ed25519.Sign(privKey, unsignedDID)
	// Encode the signature in base64 and add it to the document
	doc.Signature = base64.StdEncoding.EncodeToString(signature)
	return nil
}
