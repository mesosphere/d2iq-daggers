package daggers

import "github.com/caarlos0/env/v6"

// InitConfig initialize new config using env variables and given options
func InitConfig[T any](opts ...Option[T]) (T, error) {
	config := *new(T)

	if err := env.Parse(&config); err != nil {
		return config, err
	}

	for _, o := range opts {
		config = o(config)
	}

	return config, nil
}
