package auth

// Provides a Builder for MutualTLSStrategy authentication strategy
func MutualTLSBuilder(certificate string) Builder {
	return func(config Config) (Strategy, error) {
		panic("Not implemented")
	}
}
