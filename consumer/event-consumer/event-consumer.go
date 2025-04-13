package event_consumer

import (
	"SPBHistoryBot/events"
	"log"
	"sync"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func NewConsumer(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	wg := sync.WaitGroup{}
	for _, event := range events {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("got new event: %s", event.Text)

			if err := c.processor.Process(event); err != nil {
				log.Printf("can't handle event: %s", err.Error())
			}
		}()

	}
	wg.Wait()
	return nil
}
