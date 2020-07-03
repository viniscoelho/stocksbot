package types

const (
	// default quotes
	SAPStockCode = "SAP.F"
	USDBRLCode   = "USDBRL=X"
	EURBRLCode   = "EURBRL=X"

	EPS = 1e-9
)

var (
	DefaultQuotes = []string{SAPStockCode, USDBRLCode, EURBRLCode}
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

type Finance interface {
	Code() string
	Name() string
	RegularMarketChangePercent() float64
	RegularMarketPreviousClose() float64
	Bid() float64
	Ask() float64
}

type Threshold interface {
	LowerBound() float64
	UpperBound() float64
	UpdateBounds(low, high float64)
}
