package manager

import (
	"errors"
	"net/url"
)

var (
	errInvalidScheme = errors.New("invalid scheme")
)

func Parse(dsn string) (cfg *DSN, err error) {
	u, err := url.Parse(dsn)
	if err != nil || u.Scheme == "" {
		return
	}
	return Get(u.Scheme).ParseDSN(dsn)
}
