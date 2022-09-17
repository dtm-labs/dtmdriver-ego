package manager

import (
	"errors"
	"net/url"
)

var (
	// ErrInvalidScheme is returned when a scheme is invalid
	ErrInvalidScheme = errors.New("invalid scheme")
)

// Parse ...
func Parse(dsn string) (cfg *DSN, err error) {
	u, err := url.Parse(dsn)
	if err != nil || u.Scheme == "" {
		return nil, ErrInvalidScheme
	}

	if b, ok := get(u.Scheme); !ok {
		return nil, ErrInvalidScheme
	} else {
		return b.ParseDSN(dsn)
	}
}
