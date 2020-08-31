package instrument

// Type represents instrument type
type Type string

const (
	// Stock represents stocks
	Stock Type = "Stock"
	// Bond represents bonds
	Bond Type = "Bond"
	// EtfStock represents ETF with stocks
	EtfStock Type = "EtfStock"
	// EtfBond represents ETF with bonds
	EtfBond Type = "EtfBond"
	// EtfMixed represents ETF with mixed securities
	EtfMixed Type = "EtfMixed"
	// EtfGold represents ETF with gold
	EtfGold Type = "EtfGold"
	// EtfCurrency represents ETF with currency
	EtfCurrency Type = "EtfCurrency"
	// Currency represents currency
	Currency Type = "Currency"
)
