package service

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/1260124186-cc/courier-slot-coordinator/internal/domain"
	"github.com/1260124186-cc/courier-slot-coordinator/internal/repository"
)

func TestConcurrentCreationDoesNotExceedZoneCapacity(t *testing.T) {
	dispatch := NewDispatchService(repository.NewMemoryStore())
	arrived := make(chan struct{}, 2)
	release := make(chan struct{})
	dispatch.beforeStoreCreate = func() {
		arrived <- struct{}{}
		<-release
	}

	input := CreateShipmentInput{
		Zone:     "west",
		Window:   "09:00-11:00",
		Packages: []domain.Package{{SKU: "parcel", Units: 10}},
	}
	results := make(chan error, 2)
	var workers sync.WaitGroup
	for _, recipient := range []string{"Ava", "Noah"} {
		workers.Add(1)
		go func(recipient string) {
			defer workers.Done()
			request := input
			request.Recipient = recipient
			_, err := dispatch.CreateShipment(context.Background(), request)
			results <- err
		}(recipient)
	}

	<-arrived
	<-arrived
	close(release)
	workers.Wait()
	close(results)

	successes := 0
	capacityFailures := 0
	for err := range results {
		if err == nil {
			successes++
		}
		if errors.Is(err, domain.ErrZoneCapacity) {
			capacityFailures++
		}
	}
	if successes != 1 || capacityFailures != 1 {
		t.Fatalf("concurrent creates succeeded %d times with %d capacity failures, want 1 and 1", successes, capacityFailures)
	}
}
