package jsonstream

type config struct {
	startFrom int
	batchSize int
}

type Option interface {
	apply(cfg *config)
}

var _ Option = optionFunc(nil)

type optionFunc func(cfg *config)

func (fn optionFunc) apply(cfg *config) {
	fn(cfg)
}

func WithStartFrom(startFrom int) Option {
	return optionFunc(func(cfg *config) {
		cfg.startFrom = startFrom
	})
}

func WithBatchSize(batchSize int) Option {
	return optionFunc(func(cfg *config) {
		cfg.batchSize = batchSize
	})
}
