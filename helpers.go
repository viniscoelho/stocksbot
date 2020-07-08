package main

import (
	"fmt"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/piquette/finance-go/equity"
	"github.com/sirupsen/logrus"
	"github.com/stocksbot/types"
)

func resetThresholds(thresholds map[string]types.Threshold, filter []types.Asset) {
	logrus.Infof("Resetting thresholds...")
	values, err := fetchQuotes()
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		for n, t := range thresholds {
			logrus.Infof("Threshold %s: %v", n, t.ToString())
		}
	}()

	// update all
	if filter == nil {
		for _, c := range types.DefaultAssets {
			switch c {
			case types.SAPStockCode:
				thresholds[c] = types.NewThreshold(
					values[c].RegularMarketPreviousClose()-1.0,
					values[c].RegularMarketPreviousClose()+1.0)
			case types.EURBRLCode, types.USDBRLCode:
				thresholds[c] = types.NewThreshold(
					values[c].RegularMarketPreviousClose()-0.025,
					values[c].RegularMarketPreviousClose()+0.025)
			}
		}
		return
	}

	// update only for those which have changed
	for _, f := range filter {
		code := f.Code()
		switch code {
		case types.SAPStockCode:
			thresholds[code] = types.NewThreshold(
				values[code].RegularMarketPrice()-1.0,
				values[code].RegularMarketPrice()+1.0)
		case types.EURBRLCode, types.USDBRLCode:
			thresholds[code] = types.NewThreshold(
				values[code].Ask()-0.025,
				values[code].Ask()+0.025)
		}
	}
}

func fetchQuotes() (map[string]types.Asset, error) {
	symbols := []string{types.SAPStockCode, types.EURBRLCode, types.USDBRLCode}
	values := make(map[string]types.Asset)

	assets := equity.List(symbols)
	for assets.Next() {
		q := assets.Equity()
		values[q.Symbol] = types.NewAsset(q)
		logrus.Infof("Fetched %s: %+v", q.Symbol, q)
	}
	if assets.Err() != nil {
		return values, assets.Err()
	}

	// fetchQuotes should not return an empty list of values
	if len(values) == 0 {
		return values, fmt.Errorf("error: query returned an empty list")
	}

	return values, nil
}

// processQuote returns true if a quote had significant changes in a period
// of time; returns false otherwise
func processQuote(f types.Asset, thresholds map[string]types.Threshold) bool {
	var price float64
	code := f.Code()

	switch f.QuoteType() {
	case types.EquityType:
		price = f.RegularMarketPrice()
	default:
		price = f.Ask()
	}

	if _, ok := thresholds[code]; !ok {
		logrus.Warn("Unregistered asset -- this should not happen")
		return false
	}

	if types.FloatCompare(price, thresholds[code].LowerBound()) == types.Less {
		return true
	} else if types.FloatCompare(price, thresholds[code].UpperBound()) == types.More {
		return true
	}

	return false
}

func formatResponse(values []types.Asset) string {
	fmtRes := fmt.Sprintf("Cotações %s\n", time.Now().Format(time.RFC822))
	for _, v := range values {
		fmtRes += v.FormatResponse()
	}
	return fmtRes
}

func sendTelegramMessage(content string) error {
	botAPI, err := telegram.NewBotAPI(telegramToken)
	if err != nil {
		return err
	}

	msgConfig := telegram.NewMessage(telegramChatID, content)
	_, err = botAPI.Send(msgConfig)
	return err
}
