package instrument

// Type represents instrument type
type Type string

const (
	// Stock represents stocks
	Stock Type = "Stock"
	// Bond represents bonds
	Bond Type = "Bond"
	// Etf represents ETF
	Etf Type = "Etf"
	// Currency represents currency
	Currency Type = "Currency"
)
