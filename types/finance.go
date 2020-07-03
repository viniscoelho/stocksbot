package types

import (
	"fmt"

	"github.com/piquette/finance-go"
)

type simpleFinance struct {
	code                       string
	name                       string
	regularMarketChangePercent float64
	regularMarketPreviousClose float64
	bid                        float64
	ask                        float64
}

func NewFinance(q *finance.Quote) *simpleFinance {
	return &simpleFinance{
		code:                       q.Symbol,
		name:                       q.ShortName,
		regularMarketChangePercent: q.RegularMarketChangePercent,
		regularMarketPreviousClose: q.RegularMarketPreviousClose,
		bid:                        q.Bid,
		ask:                        q.Ask,
	}
}

func (sf simpleFinance) Code() string {
	return sf.code
}

func (sf simpleFinance) Name() string {
	return sf.name
}

func (sf simpleFinance) RegularMarketChangePercent() float64 {
	return sf.regularMarketChangePercent
}

func (sf simpleFinance) RegularMarketPreviousClose() float64 {
	return sf.regularMarketPreviousClose
}

func (sf simpleFinance) Bid() float64 {
	return sf.bid
}

func (sf simpleFinance) Ask() float64 {
	return sf.ask
}

func (sf simpleFinance) ToString() string {
	return fmt.Sprintf("%+v", sf)
}
