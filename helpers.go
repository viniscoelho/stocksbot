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

func initializeThresholds() {
	logrus.Infof("Initialising thresholds...")
	values, err := fetchQuotes()
	if err != nil {
		logrus.Fatal(err)
	}

	SAPLowThreshold = values[types.SAPStockCode].RegularMarketPreviousClose - 2
	SAPHighThreshold = values[types.SAPStockCode].RegularMarketPreviousClose + 2

	EURLowThreshold = values[types.EURBRLCode].RegularMarketPreviousClose - 0.05
	EURHighThreshold = values[types.EURBRLCode].RegularMarketPreviousClose + 0.05

	USDLowThreshold = values[types.USDBRLCode].RegularMarketPreviousClose - 0.05
	USDHighThreshold = values[types.USDBRLCode].RegularMarketPreviousClose + 0.05
	logrus.Infof("SAP Threshold: %v %v", SAPLowThreshold, SAPHighThreshold)
	logrus.Infof("EUR Threshold: %v %v", EURLowThreshold, EURHighThreshold)
	logrus.Infof("USD Threshold: %v %v", USDLowThreshold, USDHighThreshold)
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

func processQuote(f types.Finance) bool {
	switch f.Code {
	case types.USDBRLCode:
		if types.FloatCompare(f.Ask, USDLowThreshold) == types.Less {
			USDLowThreshold = f.Ask - 0.03
			return true
		} else if types.FloatCompare(f.Ask, USDHighThreshold) == types.More {
			USDHighThreshold = f.Ask + 0.03
			return true
		}
	case types.EURBRLCode:
		if types.FloatCompare(f.Ask, EURLowThreshold) == types.Less {
			EURLowThreshold = f.Ask - 0.03
			return true
		} else if types.FloatCompare(f.Ask, EURHighThreshold) == types.More {
			EURHighThreshold = f.Ask + 0.03
			return true
		}
	case types.SAPStockCode:
		if types.FloatCompare(f.Ask, SAPLowThreshold) == types.Less {
			SAPLowThreshold = f.Ask - 1.0
			return true
		} else if types.FloatCompare(f.Ask, SAPHighThreshold) == types.More {
			SAPHighThreshold = f.Ask + 1.0
			return true
		}
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
		`, v.Name, v.RegularMarketChangePercent, v.Bid, v.Ask)
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
