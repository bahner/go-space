package did

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/multiformats/go-multicodec"
)

// Ed25519PubKeyMulticodec is a constant representing the multicodec value for the Ed25519 public key
const Ed25519PubKeyMulticodec = multicodec.Ed25519Pub

// DIDDocument defines the structure of a DID Document
type DIDDocument struct {
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

// VerificationMethod defines the structure of a Verification Method
type VerificationMethod struct {
	ID                 string `json:"id"`
	Type               string `json:"type"`
	Controller         string `json:"controller"`
	PublicKeyMultibase string `json:"publicKeyMultibase"`
	Expiry             int64  `json:"expiry,omitempty"`
}

// New initializes a new DIDDocument
func New(identifier string, options map[string]interface{}) (*DIDDocument, error) {
	// Extract components from the identifier
	scheme, method, version, multibaseValue := extractComponents(identifier)

	// Validate the extracted components
	if !isValidComponents(scheme, method, version, multibaseValue) {
		return nil, errors.New("invalidDid")
	}

	// Initialize a new DIDDocument
	doc := initializeDIDDocument(identifier)

	// Create a signature method and add it to the DIDDocument
	signatureMethod, err := createSignatureMethod(identifier, multibaseValue, options)
	if err != nil {
		return nil, err
	}
	updateVerificationMethods(&doc, signatureMethod)

	// Create an encryption method and add it to the DIDDocument
	encryptionMethod, err := createEncryptionMethod(identifier, multibaseValue, options)
	if err != nil {
		return nil, err
	}
	doc.KeyAgreement = append(doc.KeyAgreement, encryptionMethod)

	return &doc, nil
}

// initializeDIDDocument initializes a new DIDDocument with basic fields
func initializeDIDDocument(identifier string) DIDDocument {
	return DIDDocument{
		Context: []string{"https://www.w3.org/ns/did/v1"},
		ID:      identifier,
	}
}

// updateVerificationMethods updates the VerificationMethods in the DIDDocument
func updateVerificationMethods(doc *DIDDocument, vm VerificationMethod) {
	doc.VerificationMethod = append(doc.VerificationMethod, vm)
	doc.Authentication = append(doc.Authentication, vm.ID)
	doc.AssertionMethod = append(doc.AssertionMethod, vm.ID)
	doc.CapabilityDelegation = append(doc.CapabilityDelegation, vm.ID)
	doc.CapabilityInvocation = append(doc.CapabilityInvocation, vm.ID)
}

// extractComponents splits the identifier into its components
func extractComponents(identifier string) (string, string, string, string) {
	components := strings.Split(identifier, ":")
	if len(components) == 3 {
		return components[0], components[1], "1", components[2]
	}
	return components[0], components[1], components[2], components[3]
}

// isValidComponents validates that the components meet certain criteria
func isValidComponents(scheme, method, version, multibaseValue string) bool {
	return scheme == "did" && method == "key" && isValidVersion(version) && multibaseValue[0] == 'z'
}

// isValidVersion validates the version string
func isValidVersion(version string) bool {
	v, err := strconv.Atoi(version)
	return err == nil && v > 0
}

// createSignatureMethod creates a new VerificationMethod for signatures
func createSignatureMethod(identifier, multibaseValue string, options map[string]interface{}) (VerificationMethod, error) {
	return VerificationMethod{
		ID:                 fmt.Sprintf("%s#key1", identifier),
		Type:               "Ed25519VerificationKey2018",
		Controller:         identifier,
		PublicKeyMultibase: multibaseValue,
	}, nil
}

// createEncryptionMethod creates a new VerificationMethod for encryption
func createEncryptionMethod(identifier, multibaseValue string, options map[string]interface{}) (VerificationMethod, error) {
	return VerificationMethod{
		ID:                 fmt.Sprintf("%s#key-encryption", identifier),
		Type:               "X25519KeyAgreementKey2019",
		Controller:         identifier,
		PublicKeyMultibase: multibaseValue,
	}, nil
}

// AddPubKey adds a public key to the DIDDocument
func (doc *DIDDocument) AddPubKey(pubKeyID, pubKeyType string, pubKey ed25519.PublicKey) error {
	// Your original code for adding a public key goes here
	return nil
}

// RemovePubKey removes a public key from the DIDDocument
func (doc *DIDDocument) RemovePubKey(pubKeyID string) error {
	// Your original code for removing a public key goes here
	return nil
}

// UnsignedDID generates the unsigned DIDDocument
func (doc *DIDDocument) UnsignedDID() ([]byte, error) {
	// Remove the signature field before serializing
	doc.Signature = ""
	unsignedDID, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	return unsignedDID, nil
}

// Sign signs the DIDDocument
func (doc *DIDDocument) Sign(privKey ed25519.PrivateKey) error {
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

// ToJSON converts the DIDDocument to JSON format
func (doc *DIDDocument) ToJSON() ([]byte, error) {
	return json.Marshal(doc)
}

// Prune removes expired keys from the DIDDocument.
func (doc *DIDDocument) Prune() {
	var pruned []VerificationMethod
	now := time.Now().Unix() // Current Unix timestamp
	for _, method := range doc.VerificationMethod {
		if method.Expiry == 0 || method.Expiry > now {
			pruned = append(pruned, method)
		}
	}
	doc.VerificationMethod = pruned
}

// Refresh would typically refresh the document from an external source.
// Since this example doesn't integrate external services, let's just return nil.
func (doc *DIDDocument) Refresh() error {
	// Implement your logic to fetch the latest document
	return nil
}

// GetNewestKey retrieves the newest key in the DIDDocument, based on expiry time.
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

// GetAllKeys retrieves all keys from the DIDDocument.
func (doc *DIDDocument) GetAllKeys() []VerificationMethod {
	return doc.VerificationMethod
}

// GetValidKeys retrieves all currently valid keys from the DIDDocument.
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
