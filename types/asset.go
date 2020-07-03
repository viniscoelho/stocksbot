package types

import (
	"fmt"

	"github.com/piquette/finance-go"
)

const (
	EquityType   = "EQUITY"
	CurrencyType = "CURRENCY"
)

type simpleAsset struct {
	code                       string
	quoteType                  string
	name                       string
	regularMarketChangePercent float64
	regularMarketPreviousClose float64
	regularMarketPrice         float64
	bid                        float64
	ask                        float64
}

func NewAsset(q *finance.Equity) *simpleAsset {
	return &simpleAsset{
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

func (sf simpleAsset) Code() string {
	return sf.code
}

func (sf simpleAsset) QuoteType() string {
	return sf.quoteType
}

func (sf simpleAsset) Name() string {
	return sf.name
}

func (sf simpleAsset) RegularMarketPrice() float64 {
	return sf.regularMarketPrice
}

func (sf simpleAsset) RegularMarketChangePercent() float64 {
	return sf.regularMarketChangePercent
}

func (sf simpleAsset) RegularMarketPreviousClose() float64 {
	return sf.regularMarketPreviousClose
}

func (sf simpleAsset) Bid() float64 {
	return sf.bid
}

func (sf simpleAsset) Ask() float64 {
	return sf.ask
}

func (sf simpleAsset) FormatResponse() string {
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

func (sf simpleAsset) ToString() string {
	return fmt.Sprintf("%+v", sf)
}
