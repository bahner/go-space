package did

import "fmt"

// VerificationMethod defines the structure of a Verification Method
type VerificationMethod struct {
	ID                 string `json:"id"`
	Type               string `json:"type"`
	Controller         string `json:"controller"`
	PublicKeyMultibase string `json:"publicKeyMultibase"`
	Expiry             int64  `json:"expiry,omitempty"`
}

// updateVerificationMethods updates the VerificationMethods in the DID
func updateVerificationMethods(doc *DIDDocument, vm VerificationMethod) {
	doc.VerificationMethod = append(doc.VerificationMethod, vm)
	doc.Authentication = append(doc.Authentication, vm.ID)
	doc.AssertionMethod = append(doc.AssertionMethod, vm.ID)
	doc.CapabilityDelegation = append(doc.CapabilityDelegation, vm.ID)
	doc.CapabilityInvocation = append(doc.CapabilityInvocation, vm.ID)
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
