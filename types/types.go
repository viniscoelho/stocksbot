package types

import (
	"time"

	"github.com/piquette/finance-go"
)

const (
	SAPStockCode = "SAP.F"
	USDBRLCode   = "USDBRL=X"
	EURBRLCode   = "EURBRL=X"

	EPS = 1e-9
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

type Finance struct {
	Code                       string  `json:"code"`
	Name                       string  `json:"name"`
	RegularMarketChangePercent float64 `json:"change_percent"`
	Bid                        float64 `json:"bid"`
	Ask                        float64 `json:"ask"`
	Timestamp                  string  `json:"timestamp"`
}

func NewFinance(q *finance.Quote) Finance {
	return Finance{
		Code:                       q.Symbol,
		Name:                       q.ShortName,
		RegularMarketChangePercent: q.RegularMarketChangePercent,
		Bid:                        q.Bid,
		Ask:                        q.Ask,
		Timestamp:                  time.Now().Format(time.RFC822),
	}
}
