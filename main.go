package main

import (
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

	SAPLowThreshold  float64
	SAPHighThreshold float64
)

func main() {
	logrus.Infof("I'm alive! :D")
	defer logrus.Infof("I'm dead. :(")

	initializeThresholds()

	apiResponse := make(chan []types.Finance, 2)

	c := cron.New()
	// At every 5th minute past every hour from 9 through 17
	// on every day-of-week from Monday through Friday
	_, err := c.AddFunc("*/5 9-17 * * 1-5", func() {
		quotes, err := fetchQuotes()
		if err != nil {
			logrus.Fatal(err)
		}

		values := make([]types.Finance, 0)
		for _, q := range quotes {
			if processQuote(q) {
				values = append(values, q)
			}
		}

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
		quotes, err := fetchQuotes()
		if err != nil {
			logrus.Fatal(err)
		}

		values := make([]types.Finance, 0)
		for _, q := range quotes {
			values = append(values, q)
		}

		apiResponse <- values
		logrus.Infof("%+v", values)
	})
	if err != nil {
		logrus.Fatal(err)
	}

	// At 8:59 on every day-of-week from Monday through Friday
	_, err = c.AddFunc("59 8 * * 1-5", func() {
		// reset thresholds for each day
		initializeThresholds()
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
