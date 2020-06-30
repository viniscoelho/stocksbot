package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

const (
	EPS = 1e-9

	telegramToken  = ""
	telegramChatID = int64(0)

	coinsAPI = "https://economia.awesomeapi.com.br/all/USD-BRL,EUR-BRL"
)

type CurrencyDTO struct {
	Code      string `json:"code"`
	CodeIn    string `json:"codein"`
	Name      string `json:"name"`
	PctChange string `json:"pctChange"`
	Bid       string `json:"bid"`
	Ask       string `json:"ask"`
	Timestamp string `json:"create_date"`
}

type FloatComp int

const (
	Less FloatComp = iota
	Equal
	More
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

	apiResponse := make(chan map[string]CurrencyDTO)

	c := cron.New()
	// At every 5th minute past every hour from 9 through 17
	// on every day-of-week from Monday through Friday
	c.AddFunc("*/5 9-17 * * 1-5", func() {
		values, err := getCurrencyValues()
		if err != nil {
			logrus.Fatal(err)
		}

		publish, err := processCurrency(values)
		if err != nil {
			logrus.Fatal(err)
		}

		if publish {
			apiResponse <- values
			logrus.Info(values)
		}
	})
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

func processCurrency(values map[string]CurrencyDTO) (bool, error) {
	flag := false
	for _, v := range values {
		cur, err := strconv.ParseFloat(v.Bid, 64)
		if err != nil {
			return flag, err
		}

		switch v.Code {
		case "USD":
			if floatCompare(cur, USDLowThreshold) == Less {
				USDLowThreshold = cur - 0.05
				flag = true
			} else if floatCompare(cur, USDHighThreshold) == More {
				USDHighThreshold = cur + 0.05
				flag = true
			}
		case "EUR":
			if floatCompare(cur, EURLowThreshold) == Less {
				EURLowThreshold = cur - 0.05
				flag = true
			} else if floatCompare(cur, EURHighThreshold) == More {
				EURHighThreshold = cur + 0.05
				flag = true
			}
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

func formatResponse(values map[string]CurrencyDTO) string {
	fmtRes := "Cotações\n"
	for _, v := range values {
		fmtRes += fmt.Sprintf(`
		Moeda: %s
		Variação Cambial: %s
		Venda: %s
		Compra: %s
		Horário: %s
		`, v.Name, v.PctChange, v.Bid, v.Ask, v.Timestamp)
	}
	return fmtRes
}

func getCurrencyValues() (map[string]CurrencyDTO, error) {
	resp, err := http.Get(coinsAPI)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch currencies: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body content: %s", err)
	}
	logrus.Infof("%s", string(body))

	dto := make(map[string]CurrencyDTO)
	err = json.Unmarshal(body, &dto)
	if err != nil {
		return nil, fmt.Errorf("failed to parse body content: %s", err)
	}

	return dto, nil
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
