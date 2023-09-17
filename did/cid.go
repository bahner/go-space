package did

import (
	"encoding/json"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

// ToJSON converts the DID to JSON format
func (d *DIDDocument) MarshalToJSON() ([]byte, error) {
	return json.Marshal(d)
}

// CIDify creates a CID for the DID
// The cid is created from the entire DID document
// and is what will be pubslished to IPNS.
func (d *DIDDocument) CID() (string, error) {

	// Convert DID to bytes
	DID, err := d.MarshalToJSON()
	if err != nil {
		return "", err
	}

	// Create a multihash
	hashed, err := multihash.Sum(DID, multihash.SHA2_256, -1)
	if err != nil {
		return "", err
	}

	// Create a CID from the multihash
	c := cid.NewCidV1(cid.Raw, hashed)

	return c.String(), nil
}
