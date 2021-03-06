package types

const (
	// default quotes
	SAPStockCode = "SAP.F"
	USDBRLCode   = "USDBRL=X"
	EURBRLCode   = "EURBRL=X"

	EPS = 1e-9
)

var (
	DefaultAssets = []string{SAPStockCode, USDBRLCode, EURBRLCode}
)

type FloatComp int

const (
	Less FloatComp = iota
	Equal
	More
)

func FloatCompare(a, b float64) FloatComp {
	if a+EPS < b {
		return Less
	} else if a-EPS > b {
		return More
	}
	return Equal
}

type Asset interface {
	Code() string
	QuoteType() string
	Name() string
	RegularMarketPrice() float64
	RegularMarketChangePercent() float64
	RegularMarketPreviousClose() float64
	Bid() float64
	Ask() float64
	FormatResponse() string
	ToString() string
}

type Threshold interface {
	LowerBound() float64
	UpperBound() float64
	UpdateBounds(low, high float64)
	ToString() string
}
