package main

import (
	"fmt"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/piquette/finance-go/forex"
	"github.com/piquette/finance-go/quote"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/stocksbot/types"
)

const (
	telegramToken  = ""
	telegramChatID = int64(0)
)

var (
	USDLowThreshold  float64
	USDHighThreshold float64
	EURLowThreshold  float64
	EURHighThreshold float64
)

func main() {
	USDLowThreshold = 5.1
	USDHighThreshold = 5.45

	EURLowThreshold = 5.8
	EURHighThreshold = 6.15

	err := sendTelegramMessage("I'm alive!")
	if err != nil {
		logrus.Fatal(err)
	}
	defer sendTelegramMessage("I'm dead. :(")

	apiResponse := make(chan map[string]types.Finance)

	c := cron.New()
	// At every 5th minute past every hour from 9 through 17
	// on every day-of-week from Monday through Friday
	_, err = c.AddFunc("*/5 9-17 * * 1-5", func() {
		values, err := fetchQuotes()
		if err != nil {
			logrus.Fatal(err)
		}

		for k, q := range values {
			if !processQuote(q) {
				delete(values, k)
			}
		}

		if len(values) > 0 {
			apiResponse <- values
			logrus.Infof("%+v", values)
		}
		logrus.Infof("USD: %v %v, EUR: %v %v", USDLowThreshold, USDHighThreshold, EURLowThreshold, EURHighThreshold)
	})
	if err != nil {
		logrus.Fatal(err)
	}

	// At every hour from 9 through 17
	// on every day-of-week from Monday through Friday
	_, err = c.AddFunc("*/60 9-17 * * 1-5", func() {
		values, err := fetchQuotes()
		if err != nil {
			logrus.Fatal(err)
		}

		if len(values) > 0 {
			apiResponse <- values
			logrus.Infof("%+v", values)
		}
		logrus.Infof("USD: %v %v, EUR: %v %v", USDLowThreshold, USDHighThreshold, EURLowThreshold, EURHighThreshold)
	})
	if err != nil {
		logrus.Fatal(err)
	}

	// At 17:35 on every day-of-week from Monday through Friday
	_, err = c.AddFunc("35 17 * * 1-5", func() {
		values, err := fetchQuotes()
		if err != nil {
			logrus.Fatal(err)
		}

		apiResponse <- values
		logrus.Info(values)
	})
	if err != nil {
		logrus.Fatal(err)
	}
	c.Start()

	for {
		select {
		case resp := <-apiResponse:
			err := sendTelegramMessage(formatResponse(resp))
			if err != nil {
				logrus.Fatal(err)
			}
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

func processQuote(f types.Finance) bool {
	switch f.Code {
	case types.USDBRLCode:
		if types.FloatCompare(f.Bid, USDLowThreshold) == types.Less {
			USDLowThreshold = f.Bid - 0.05
			return true
		} else if types.FloatCompare(f.Bid, USDHighThreshold) == types.More {
			USDHighThreshold = f.Bid + 0.05
			return true
		}
	case types.EURBRLCode:
		if types.FloatCompare(f.Bid, EURLowThreshold) == types.Less {
			EURLowThreshold = f.Bid - 0.05
			return true
		} else if types.FloatCompare(f.Bid, EURHighThreshold) == types.More {
			EURHighThreshold = f.Bid + 0.05
			return true
		}
	case types.SAPStockCode:
		if types.FloatCompare(f.Bid, 127.0) == types.More {
			return true
		}
	}

	return false
}

func formatResponse(values map[string]types.Finance) string {
	fmtRes := "Cotações\n"
	for _, v := range values {
		fmtRes += fmt.Sprintf(`
		Moeda: %s
		Variação Cambial: %v%%
		Venda: %v
		Compra: %v
		Horário: %s
		`, v.Name, v.RegularMarketChangePercent, v.Bid, v.Ask, v.Timestamp)
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
