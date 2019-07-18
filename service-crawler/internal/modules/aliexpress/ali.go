package aliexpress

import (
	"crawler/internal/lib/parser"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

const (
	url    = "https://ru.aliexpress.com/item/32358502254.html?spm=a2g01.12597576.p99adbb.16.4c45124ax0CMgO&gps-id=6272461&scm=1007.16233.119941.0&scm_id=1007.16233.119941.0&scm-url=1007.16233.119941.0&pvid=a7a42fc8-1d59-4a72-9a0f-2edde4324f1f"
	domain = "ru.aliexpress.com"
	Name   = "Aliexpress"
)

type (
	Page struct {
		TitleModule TitleModule `json:"titleModule"`
		PriceModule PriceModule `json:"priceModule"`
		StoreModule StoreModule `json:"storeModule"`
	}

	TitleModule struct {
		Rating     FeedbackRating `json:"feedbackRating"`
		TradeCount int            `json:"tradeCount"`
		Subject    string         `json:"subject"`
	}

	FeedbackRating struct {
		StarSTR       string `json:"averageStar"`
		CountFeedback int    `json:"totalValidNum"`
	}

	StoreModule struct {
		FollowingNumber int    `json:"followingNumber"`
		PositiveRateSTR string `json:"positiveRate"`
	}

	PriceModule struct {
		Discount float64 `json:"discount"`
		Max      Amount  `json:"maxAmount"`
		Min      Amount  `json:"minActivityAmount"`
	}

	Amount struct {
		Currency string  `json:"currency"`
		Value    float64 `json:"value"`
	}
)

func RunCollector(log *logrus.Logger) error {
	errs := make(chan error)
	documents := make(chan json.RawMessage)

	p := parser.Build(url, "script", log, documents, domain)

	go Listen(log, documents, errs)
	go p.Start()

	return <-errs
}

func Listen(log *logrus.Logger, documents <-chan json.RawMessage, errs chan<- error) {
	for js := range documents {
		var page Page
		if err := json.Unmarshal(js, &page); err != nil {
			errs <- err
		}
		log.Print(page)
	}
}
