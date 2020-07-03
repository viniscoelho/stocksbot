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

func main() {
	logrus.Infof("I'm alive! :D")
	defer logrus.Infof("I'm dead. :(")

	apiResponse := make(chan []types.Finance, 2)
	thresholds := make(map[string]types.Threshold, 3)

	resetThresholds(thresholds, nil)

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
			if processQuote(q, thresholds) {
				values = append(values, q)
			}
		}

		if len(values) > 0 {
			resetThresholds(thresholds, values)
			apiResponse <- values
			for _, v := range values {
				logrus.Infof("Quotes: %s", v.ToString())
			}
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
		for _, v := range values {
			logrus.Infof("Quotes: %s", v.ToString())
		}
	})
	if err != nil {
		logrus.Fatal(err)
	}

	// At 8:59 on every day-of-week from Monday through Friday
	_, err = c.AddFunc("59 8 * * 1-5", func() {
		// reset thresholds for each day
		resetThresholds(thresholds, nil)
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
