package main

import (
	"fmt"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/piquette/finance-go/forex"
	"github.com/piquette/finance-go/quote"
	"github.com/sirupsen/logrus"
	"github.com/stocksbot/types"
)

func resetThresholds(thresholds map[string]types.Threshold, filter []types.Finance) {
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
		for _, c := range types.DefaultQuotes {
			switch c {
			case types.SAPStockCode:
				thresholds[c] = types.NewThreshold(
					values[c].RegularMarketPreviousClose()-2.0,
					values[c].RegularMarketPreviousClose()+2.0)
			case types.EURBRLCode, types.USDBRLCode:
				thresholds[c] = types.NewThreshold(
					values[c].RegularMarketPreviousClose()-0.05,
					values[c].RegularMarketPreviousClose()+0.05)
			}
		}
		return
	}

	// update only those which have changed
	for _, f := range filter {
		code := f.Code()
		switch code {
		case types.SAPStockCode:
			thresholds[code] = types.NewThreshold(
				values[code].RegularMarketPreviousClose()-2.0,
				values[code].RegularMarketPreviousClose()+2.0)
		case types.EURBRLCode, types.USDBRLCode:
			thresholds[code] = types.NewThreshold(
				values[code].RegularMarketPreviousClose()-0.05,
				values[code].RegularMarketPreviousClose()+0.05)
		}
	}
}

func fetchQuotes() (map[string]types.Finance, error) {
	values := make(map[string]types.Finance)

	// SAP
	sapStock, err := quote.Get(types.SAPStockCode)
	if err != nil {
		return nil, err
	}
	values[types.SAPStockCode] = types.NewFinance(sapStock)

	// EUR
	eurBrl, err := forex.Get(types.EURBRLCode)
	if err != nil {
		return nil, err
	}
	values[types.EURBRLCode] = types.NewFinance(&eurBrl.Quote)

	// USD
	usdBrl, err := forex.Get(types.USDBRLCode)
	if err != nil {
		return nil, err
	}
	values[types.USDBRLCode] = types.NewFinance(&usdBrl.Quote)

	return values, nil
}

func processQuote(f types.Finance, thresholds map[string]types.Threshold) bool {
	code := f.Code()
	if types.FloatCompare(f.Ask(), thresholds[code].LowerBound()) == types.Less {
		return true
	} else if types.FloatCompare(f.Ask(), thresholds[code].LowerBound()) == types.More {
		return true
	}

	return false
}

func formatResponse(values []types.Finance) string {
	fmtRes := fmt.Sprintf("Cotações %s\n", time.Now().Format(time.RFC822))
	for _, v := range values {
		fmtRes += fmt.Sprintf(`
		Nome: %s
		Variação: %v%%
		Compra: %v
		Venda: %v
		`, v.Name(), v.RegularMarketChangePercent(), v.Bid(), v.Ask())
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
