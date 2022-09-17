package manager

var (
	// m is a map from scheme to dsn builder.
	m = make(map[string]DSNParser)
)

// Register ...
func Register(b DSNParser) {
	m[b.Scheme()] = b
}

// Get returns the dsn builder registered with the given scheme.
//
// If no builder is register with the scheme, nil will be returned.
func get(scheme string) (b DSNParser, ok bool) {
	b, ok = m[scheme]

	return
}
