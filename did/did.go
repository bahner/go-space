package did

import (
	"errors"

	"github.com/multiformats/go-multicodec"
)

// Ed25519PubKeyMulticodec is a constant representing the multicodec value for the Ed25519 public key
const Ed25519PubKeyMulticodec = multicodec.Ed25519Pub

// DID defines the structure of a DID Document
type DID struct {
	Context              []string             `json:"@context"`
	ID                   string               `json:"id"`
	Signature            string               `json:"signature,omitempty"`
	VerificationMethod   []VerificationMethod `json:"verificationMethod"`
	Authentication       []string             `json:"authentication"`
	AssertionMethod      []string             `json:"assertionMethod"`
	CapabilityDelegation []string             `json:"capabilityDelegation"`
	CapabilityInvocation []string             `json:"capabilityInvocation"`
	KeyAgreement         []VerificationMethod `json:"keyAgreement"`
}

// New initializes a new DID
func New(identifier string, options map[string]interface{}) (*DID, error) {
	// Extract components from the identifier
	scheme, method, version, multibaseValue := extractComponents(identifier)

	// Validate the extracted components
	if !isValidComponents(scheme, method, version, multibaseValue) {
		return nil, errors.New("invalidDid")
	}

	// Initialize a new DID
	doc := initializeDID(identifier)

	// Create a signature method and add it to the DID
	signatureMethod, err := createSignatureMethod(identifier, multibaseValue, options)
	if err != nil {
		return nil, err
	}
	updateVerificationMethods(&doc, signatureMethod)

	// Create an encryption method and add it to the DID
	encryptionMethod, err := createEncryptionMethod(identifier, multibaseValue, options)
	if err != nil {
		return nil, err
	}
	doc.KeyAgreement = append(doc.KeyAgreement, encryptionMethod)

	return &doc, nil
}

// initializeDID initializes a new DID with basic fields
func initializeDID(identifier string) DID {
	return DID{
		Context: []string{"https://www.w3.org/ns/did/v1"},
		ID:      identifier,
	}
}
