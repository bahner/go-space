package did

import (
	"crypto/rand"
	"errors"
	"strings"
	"time"

	"github.com/mr-tron/base58"
	"github.com/multiformats/go-multibase"
	"golang.org/x/crypto/ed25519"
)

const didPrefix = "did:ipid:"
const didContext = "bafyreict2igxeqinf4uyg26ubff2vginv7hiz3dxiiwb77vwyr3c5jd3rq"

type DIDDocument struct {
	Context        DIDContext        `json:"@context"`
	Authentication Authentication    `json:"authentication"`
	Created        string            `json:"created"`
	ID             string            `json:"id"`
	Previous       map[string]string `json:"previous"`
	Proof          map[string]string `json:"proof"`
	PublicKey      []PublicKey       `json:"publicKey"`
	Signature      Signature         `json:"signature"`
	Updated        string            `json:"updated"`
}

type DIDContext struct {
	Version            int                    `json:"@version"`
	Authentication     IDType                 `json:"authentication"`
	Created            DateTimeType           `json:"created"`
	Cryptosuite        string                 `json:"cryptosuite"`
	DC                 string                 `json:"dc"`
	ID                 string                 `json:"id"`
	Proof              string                 `json:"proof"`
	ProofPurpose       string                 `json:"proofPurpose"`
	ProofValue         string                 `json:"proofValue"`
	PublicKeyMultibase PublicKeyMultibaseType `json:"publicKeyMultibase"`
	SEC                string                 `json:"sec"`
	Type               string                 `json:"type"`
	Updated            DateTimeType           `json:"updated"`
	VerificationMethod IDType                 `json:"verificationMethod"`
	XSD                string                 `json:"xsd"`
}

type IDType struct {
	ID   string `json:"@id"`
	Type string `json:"@type"`
}

type DateTimeType struct {
	ID   string `json:"@id"`
	Type string `json:"@type"`
}

type PublicKeyMultibaseType struct {
	Container string `json:"@container"`
	ID        string `json:"@id"`
	Type      string `json:"@type"`
}

type Authentication struct {
	Type      string   `json:"type"`
	PublicKey []string `json:"publicKey"`
}

type PublicKey struct {
	Curve              string `json:"curve"`
	Expires            string `json:"expires,omitempty"`
	PublicKeyMultibase string `json:"publicKeyMultibase"`
	Type               string `json:"type"`
	Status             string `json:"status,omitempty"`
}

type Signature struct {
	Created        string            `json:"created"`
	Creator        string            `json:"creator"`
	Message        map[string]string `json:"message"`
	SignatureValue string            `json:"signatureValue"`
	Type           string            `json:"type"`
}

func New() (*DIDDocument, ed25519.PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	did := "did:key:" + base58.Encode(publicKey)

	doc := &DIDDocument{
		Context: DIDContext{
			Version: 2,
			Authentication: IDType{
				ID:   "sec:authentication",
				Type: "@id",
			},
			Created: DateTimeType{
				ID:   "dc:created",
				Type: "xsd:dateTime",
			},
			Cryptosuite:  "sec:cryptosuite",
			DC:           "http://purl.org/dc/terms/",
			ID:           "@id",
			Proof:        "sec:proof",
			ProofPurpose: "sec:proofPurpose",
			ProofValue:   "sec:proofValue",
			PublicKeyMultibase: PublicKeyMultibaseType{
				Container: "@set",
				ID:        "sec:publicKeyMultibase",
				Type:      "@id",
			},
			SEC:  "https://w3id.org/security#",
			Type: "@type",
			Updated: DateTimeType{
				ID:   "dc:updated",
				Type: "xsd:dateTime",
			},
			VerificationMethod: IDType{
				ID:   "sec:verificationMethod",
				Type: "@id",
			},
			XSD: "http://www.w3.org/2001/XMLSchema#",
		},
		Authentication: Authentication{
			Type:      "EdDsaSASignatureAuthentication2018",
			PublicKey: []string{did},
		},
		Created: "2018-12-01T03:00:00Z",
		ID:      did,
		Previous: map[string]string{
			"/": "zdpuB1oR3vjYmkDc9ALfY7o6hSt1Hrg2ApXaYAFyiAW5E4NJP",
		},
		Proof: map[string]string{
			"/": "z43AaGF42R2DXsU65bNnHRCypLPr9sg6D7CUws5raiqATVaB1jj",
		},
		PublicKey: []PublicKey{
			{
				Curve:              "ed25519",
				PublicKeyMultibase: base58.Encode(publicKey),
				Type:               "EdDsaPublicKey",
				Status:             "revoked",
			},
			{
				Curve:              "ed25519",
				Expires:            time.Now().Add(365 * 24 * time.Hour).Format(time.RFC3339),
				PublicKeyMultibase: base58.Encode(publicKey),
				Type:               "EdDsaPublicKey",
			},
		},
		Signature: Signature{
			Created:        "2018-12-01T03:00:04Z",
			Creator:        did + "/publicKey/1",
			Message:        map[string]string{"/": "zdpuAyvreXzQHqwv3rL8MaVPjNJjpLLa5Du3HcbpQL41XS35G"},
			SignatureValue: base58.Encode(ed25519.Sign(privateKey, []byte(did))),
			Type:           "ed25519Signature2018",
		},
		Updated: "2018-12-01T03:00:04Z",
	}

	return doc, privateKey, nil
}

// Adds a new public key to the DID Document.
func (doc *DIDDocument) AddPublicKey(pubKey ed25519.PublicKey, isRevoked bool) {
	pk := PublicKey{
		Curve:              "ed25519",
		PublicKeyMultibase: base58.Encode(pubKey),
		Type:               "EdDsaPublicKey",
	}
	if isRevoked {
		pk.Status = "revoked"
	} else {
		pk.Expires = time.Now().Add(365 * 24 * time.Hour).Format(time.RFC3339)
	}
	doc.PublicKey = append(doc.PublicKey, pk)
}

// Signs a message with the private key, returning the signature.
func (doc *DIDDocument) Sign(privateKey ed25519.PrivateKey, message []byte) (string, error) {
	if len(privateKey) == 0 || len(message) == 0 {
		return "", errors.New("empty private key or message")
	}
	signature := ed25519.Sign(privateKey, message)
	return base58.Encode(signature), nil
}

func (doc *DIDDocument) VerifyProof(proof string, message []byte) bool {
	signature, err := base58.Decode(proof)
	if err != nil {
		return false
	}

	for _, pk := range doc.PublicKey {
		keyType := extractKeyType(pk.PublicKeyMultibase)
		if keyType == "ed25519" {
			pubKey, err := base58.Decode(pk.PublicKeyMultibase)
			if err != nil {
				continue
			}
			if ed25519.Verify(pubKey, message, signature) {
				return true
			}
		}
		// Future: Add conditions for other key types
	}
	return false
}
func extractKeyType(did string) string {
	// First, strip off the "did:key:" prefix
	if !strings.HasPrefix(did, didPrefix) {
		return ""
	}
	mbValue := did[len(didPrefix):]

	// Decode the multibase value
	_, decodedBytes, err := multibase.Decode(mbValue)
	if err != nil {
		return ""
	}

	// Check the multicodec prefix
	if len(decodedBytes) > 2 && decodedBytes[0] == 0xed && decodedBytes[1] == 0x01 {
		return "ed25519"
	}

	// Return an empty string for unsupported or unrecognized key types
	return ""
}
