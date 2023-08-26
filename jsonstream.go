package jsonstream

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

const (
	DefaultBatchSize = 100
	DefaultStartFrom = 1
)

type Entry[T any] struct {
	Value T
	Err   error
}

func Unmarshal[T any](ctx context.Context, r io.Reader, opts ...Option) (<-chan Entry[T], error) {
	cfg := &config{
		batchSize: DefaultBatchSize,
		startFrom: DefaultStartFrom,
	}

	for _, opt := range opts {
		opt.apply(cfg)
	}

	dec := json.NewDecoder(r)

	if err := skipTokensUntil(dec, cfg.startFrom); err != nil {
		return nil, fmt.Errorf("skip tokens: %w", err)
	}

	out := make(chan Entry[T], cfg.batchSize)

	go func() {
		defer close(out)

		for dec.More() {
			entry := Entry[T]{}

			if err := dec.Decode(&entry.Value); err != nil {
				entry.Err = fmt.Errorf("decode: %w", err)
			}

			select {
			case out <- entry:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, nil
}

func skipTokensUntil(dec *json.Decoder, startFrom int) error {
	for i := 0; i < startFrom; i++ {
		if _, err := dec.Token(); err != nil {
			return err
		}
	}

	return nil
}
