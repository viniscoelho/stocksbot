package main

import (
	"fmt"
	"math"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/forex"
	"github.com/piquette/finance-go/quote"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

const (
	EPS = 1e-9

	telegramToken  = ""
	telegramChatID = int64(0)
)

type Finance struct {
	Code                       string  `json:"code"`
	Name                       string  `json:"name"`
	RegularMarketChangePercent float64 `json:"change_percent"`
	Bid                        float64 `json:"bid"`
	Ask                        float64 `json:"ask"`
	Timestamp                  string  `json:"timestamp"`
}

type FloatComp int

const (
	Less FloatComp = iota
	Equal
	More
)

const (
	SAPStockCode = "SAP.F"
	USDBRLCode   = "USDBRL=X"
	EURBRLCode   = "EURBRL=X"
)

var (
	USDLowThreshold  float64
	USDHighThreshold float64
	EURLowThreshold  float64
	EURHighThreshold float64
)

func main() {
	USDLowThreshold = 5.1
	USDHighThreshold = 5.5

	EURLowThreshold = 5.8
	EURHighThreshold = 6.2

	err := sendTelegramMessage("I'm alive!")
	if err != nil {
		logrus.Fatal(err)
	}
	defer sendTelegramMessage("I'm dead. :(")

	apiResponse := make(chan map[string]Finance)

	c := cron.New()
	// At every 5th minute past every hour from 9 through 17
	// on every day-of-week from Monday through Friday
	_, err = c.AddFunc("*/5 9-17 * * 1-5", func() {
		values := make(map[string]Finance)

		// SAP
		sapStock, err := quote.Get(SAPStockCode)
		if err != nil {
			logrus.Fatal(err)
		}

		sapStockToPublish := NewFinance(sapStock)
		publish, err := processCurrency(sapStockToPublish)
		if err != nil {
			logrus.Fatal(err)
		}
		if publish {
			values[SAPStockCode] = sapStockToPublish
		}
		// End SAP

		// Euro
		eurBrl, err := forex.Get(EURBRLCode)
		if err != nil {
			logrus.Fatal(err)
		}

		eurBrlToPublish := NewFinance(&eurBrl.Quote)
		publish, err = processCurrency(eurBrlToPublish)
		if err != nil {
			logrus.Fatal(err)
		}
		if publish {
			values[EURBRLCode] = eurBrlToPublish
		}
		// End Euro

		// USD
		usdBrl, err := forex.Get(USDBRLCode)
		if err != nil {
			logrus.Fatal(err)
		}

		usdBrlToPublish := NewFinance(&usdBrl.Quote)
		publish, err = processCurrency(usdBrlToPublish)
		if err != nil {
			logrus.Fatal(err)
		}
		if publish {
			values[USDBRLCode] = usdBrlToPublish
		}
		// End USD

		if len(values) > 0 {
			apiResponse <- values
			logrus.Infof("%+v", values)
		}
	})
	if err != nil {
		logrus.Fatal(err)
	}

	// At 17:35 on every day-of-week from Monday through Friday
	_, err = c.AddFunc("35 17 * * 1-5", func() {
		values := make(map[string]Finance)
		sapStock, err := quote.Get(SAPStockCode)
		if err != nil {
			logrus.Fatal(err)
		}
		values[SAPStockCode] = NewFinance(sapStock)

		eurBrl, err := forex.Get(EURBRLCode)
		if err != nil {
			logrus.Fatal(err)
		}
		values[EURBRLCode] = NewFinance(&eurBrl.Quote)

		usdBrl, err := forex.Get(USDBRLCode)
		if err != nil {
			logrus.Fatal(err)
		}
		values[USDBRLCode] = NewFinance(&usdBrl.Quote)

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

func processCurrency(f Finance) (bool, error) {
	flag := false
	switch f.Code {
	case USDBRLCode:
		if floatCompare(f.Bid, USDLowThreshold) == Less {
			USDLowThreshold = f.Bid - 0.05
			flag = true
		} else if floatCompare(f.Bid, USDHighThreshold) == More {
			USDHighThreshold = f.Bid + 0.05
			flag = true
		}
	case EURBRLCode:
		if floatCompare(f.Bid, EURLowThreshold) == Less {
			EURLowThreshold = f.Bid - 0.05
			flag = true
		} else if floatCompare(f.Bid, EURHighThreshold) == More {
			EURHighThreshold = f.Bid + 0.05
			flag = true
		}
	case SAPStockCode:
		if floatCompare(f.Bid, 127.0) == More {
			flag = true
		}
	}

	return flag, nil
}

func floatCompare(a, b float64) FloatComp {
	if math.Abs(a-b) < EPS {
		return Less
	} else if math.Abs(a-b) > EPS {
		return More
	}
	return Equal
}

func formatResponse(values map[string]Finance) string {
	fmtRes := "Cotação\n"
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
