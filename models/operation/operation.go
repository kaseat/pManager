package operation

// Type is market operation type
type Type string

const (
	// Buy operation
	Buy Type = "buy"
	// Sell operation
	Sell Type = "sell"
	// BrokerageFee operation
	BrokerageFee Type = "brokerageFee"
	// ExchangeFee operation
	ExchangeFee Type = "exchangeFee"
	// PayIn operation
	PayIn Type = "payIn"
	// PayOut operation
	PayOut Type = "payOut"
	// Coupon operation
	Coupon Type = "coupon"
	// AccInterestBuy operation
	AccInterestBuy Type = "accruedInterestBuy"
	// AccInterestSell operation
	AccInterestSell Type = "accruedInterestSell"
	// Unknown operation
	Unknown Type = "unknown"
)
