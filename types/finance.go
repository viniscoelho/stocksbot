package types

import (
	"fmt"

	"github.com/piquette/finance-go"
)

const (
	EquityType   = "EQUITY"
	CurrencyType = "CURRENCY"
)

type simpleFinance struct {
	code                       string
	quoteType                  string
	name                       string
	regularMarketChangePercent float64
	regularMarketPreviousClose float64
	regularMarketPrice         float64
	bid                        float64
	ask                        float64
}

func NewFinance(q *finance.Equity) *simpleFinance {
	return &simpleFinance{
		code:                       q.Symbol,
		quoteType:                  string(q.QuoteType),
		name:                       q.ShortName,
		regularMarketChangePercent: q.RegularMarketChangePercent,
		regularMarketPreviousClose: q.RegularMarketPreviousClose,
		regularMarketPrice:         q.RegularMarketPrice,
		bid:                        q.Bid,
		ask:                        q.Ask,
	}
}

func (sf simpleFinance) Code() string {
	return sf.code
}

func (sf simpleFinance) QuoteType() string {
	return sf.quoteType
}

func (sf simpleFinance) Name() string {
	return sf.name
}

func (sf simpleFinance) RegularMarketPrice() float64 {
	return sf.regularMarketPrice
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

func (sf simpleFinance) FormatResponse() string {
	switch sf.quoteType {
	case EquityType:
		return fmt.Sprintf(`
		Nome: %s
		Variação: %v%%
		Valor: %v
		`, sf.Name(), sf.RegularMarketChangePercent(), sf.RegularMarketPrice())
	case CurrencyType:
		return fmt.Sprintf(`
		Nome: %s
		Variação: %v%%
		Compra: %v
		Venda: %v
		`, sf.Name(), sf.RegularMarketChangePercent(), sf.Bid(), sf.Ask())
	}

	return "type not supported"
}

func (sf simpleFinance) ToString() string {
	return fmt.Sprintf("%+v", sf)
}
