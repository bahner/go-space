package global

import (
	"context"

	"github.com/bahner/go-myspace/p2p/pubsub"
)

func startPubSubService(ctx context.Context) (*pubsub.Service, error) {

	defer globalWg.Done()

	log.Info("Starting Global PubSubService")
	ps := pubsub.New()

	globalWg.Add(1)
	go PubSubService.Start(ctx, globalWg)

	return ps, nil
}
