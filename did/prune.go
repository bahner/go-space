package did

import "time"

// Prune removes expired keys from the DID.
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
