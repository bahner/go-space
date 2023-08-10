package global

import (
	"context"
	"sync"

	"github.com/bahner/go-myspace/p2p/pubsub"
)

func startPubSubService(ctx context.Context, wg *sync.WaitGroup) (*pubsub.Service, error) {

	defer wg.Done()

	log.Info("Starting Global PubSubService")
	ps := pubsub.New()

	wg.Add(1)
	go PubSubService.Start(ctx, wg)

	return ps, nil
}
