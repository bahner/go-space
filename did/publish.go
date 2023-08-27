package did

import (
	"fmt"
	"time"

	"github.com/ipfs/boxo/ipns"
	"github.com/ipfs/boxo/path"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/libp2p/go-libp2p/core/crypto"
	log "github.com/sirupsen/logrus"
)

// PublishOption holds options for publishing a DID document
type PublishOption struct {
	TTL *time.Duration
	EOL *time.Time
}

// PublishOptions initializes default options for publishing
func PublishOptions() *PublishOption {
	ipnsRecordTTL := 24 * time.Hour
	ipnsRecordEOL := time.Now().Add(24 * time.Hour)
	return &PublishOption{
		TTL: &ipnsRecordTTL,
		EOL: &ipnsRecordEOL,
	}
}

// Publish publishes the DID document
func (d *DID) Publish(privateKey crypto.PrivKey, opts *PublishOption) error {
	if err := validatePublishOptions(opts); err != nil {
		return err
	}

	cid, err := d.CID()
	if err != nil {
		log.Errorf("Failed to get CID of DID document: %s", err)
		return fmt.Errorf("could not get CID of DID document: %w", err)
	}

	// Convert CID to a path
	path := path.FromString(cid)

	ipnsRecord, err := createIPNSRecord(privateKey, path, opts)
	if err != nil {
		return err
	}

	serializedIpnsRecord, err := serializeIPNSRecord(ipnsRecord)
	if err != nil {
		return err
	}

	return publishToIPFS(d.ID, serializedIpnsRecord)
}

func validatePublishOptions(opts *PublishOption) error {
	if opts.TTL == nil || opts.EOL == nil {
		log.Errorf("TTL or EOL in PublishOption is nil")
		return fmt.Errorf("TTL or EOL in PublishOption cannot be nil")
	}
	return nil
}

func createIPNSRecord(privateKey crypto.PrivKey, path path.Path, opts *PublishOption) (*ipns.Record, error) {
	ipnsRecord, err := ipns.NewRecord(privateKey, path, 1, *opts.EOL, *opts.TTL)
	if err != nil {
		log.Errorf("Failed to create new IPNS record: %s", err)
		return nil, fmt.Errorf("could not create new IPNS record: %w", err)
	}
	return ipnsRecord, nil
}

func serializeIPNSRecord(ipnsRecord *ipns.Record) (uint64, error) {
	serializedIpnsRecord, err := ipnsRecord.Sequence()
	if err != nil {
		log.Errorf("Failed to serialize IPNS record: %s", err)
		return 0, fmt.Errorf("could not serialize IPNS record: %w", err)
	}
	return serializedIpnsRecord, nil
}

func publishToIPFS(didID string, serializedIpnsRecord uint64) error {
	sh := shell.NewLocalShell()
	if sh == nil {
		log.Errorf("Could not connect to the local IPFS node")
		return fmt.Errorf("could not connect to the local IPFS node")
	}

	apiLifetime := 5 * time.Minute
	apiTTL := 1 * time.Minute
	ipnsPath := "/ipns/" + didID
	resolve := true

	response, err := sh.PublishWithDetails(
		ipnsPath,
		fmt.Sprint(serializedIpnsRecord),
		apiLifetime,
		apiTTL,
		resolve,
	)
	if err != nil {
		log.Errorf("Failed to publish IPNS record: %s", err)
		return fmt.Errorf("could not publish IPNS record: %w", err)
	}
	log.Debugf("Successfully published IPNS record. Response: %v", response)

	return nil
}
