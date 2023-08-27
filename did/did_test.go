package did

import (
	"crypto/ed25519"
	"encoding/json"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	identifier := "did:key:z6MkuHcZ9fNjD3F6Nz1DNvG4b8c4FcP3s3tVfBsriXcHAeWZ"
	options := make(map[string]interface{})
	doc, err := New(identifier, options)
	if err != nil {
		t.Errorf("Failed to create a new DID document: %v", err)
		return
	}
	if doc.ID != identifier {
		t.Errorf("DID identifier mismatch: got %v, want %v", doc.ID, identifier)
	}
}

func TestUnsignedDID(t *testing.T) {
	identifier := "did:key:z6MkuHcZ9fNjD3F6Nz1DNvG4b8c4FcP3s3tVfBsriXcHAeWZ"
	options := make(map[string]interface{})
	doc, _ := New(identifier, options)

	unsignedDID, err := doc.UnsignedDID()
	if err != nil {
		t.Errorf("Failed to generate unsigned DID: %v", err)
		return
	}
	var parsed map[string]interface{}
	json.Unmarshal(unsignedDID, &parsed)
	if _, exists := parsed["signature"]; exists {
		t.Errorf("Signature field should not exist in unsigned DID")
	}
}

func TestSign(t *testing.T) {
	_, privKey, _ := ed25519.GenerateKey(nil)
	identifier := "did:key:z6MkuHcZ9fNjD3F6Nz1DNvG4b8c4FcP3s3tVfBsriXcHAeWZ"
	options := make(map[string]interface{})
	doc, _ := New(identifier, options)

	err := doc.Sign(privKey)
	if err != nil {
		t.Errorf("Failed to sign the DID document: %v", err)
		return
	}
	if doc.Signature == "" {
		t.Errorf("Signature should not be empty after signing")
	}
}

func TestPrune(t *testing.T) {
	identifier := "did:key:z6MkuHcZ9fNjD3F6Nz1DNvG4b8c4FcP3s3tVfBsriXcHAeWZ"
	options := make(map[string]interface{})
	doc, _ := New(identifier, options)

	// Add an expired key
	doc.VerificationMethod = append(doc.VerificationMethod, VerificationMethod{
		ID:     "expired-key",
		Expiry: time.Now().Unix() - 1,
	})

	doc.Prune()

	for _, method := range doc.VerificationMethod {
		if method.ID == "expired-key" {
			t.Errorf("Prune failed to remove the expired key")
			return
		}
	}
}
