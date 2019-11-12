package filter

import (
	"encoding/json"
	"github.com/jellyboysband/eye/cmd/crawler/internal/app"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"log"
)

const (
	contentType = "application/json"
)

type (
	store struct {
		Name  string
		ch    *amqp.Channel
		appID string
	}

	document struct {
		Title         string   `json:"title"`
		OurRating     float64  `json:"our_rating"`
		ID            int      `json:"id"`
		URL           string   `json:"url"`
		TotalSales    int      `json:"total_sales"`
		RatingProduct float64  `json:"rating_product"`
		TotalComment  int      `json:"total_comment"`
		Images        []string `json:"images"`
		Discount      float64  `json:"discount"`
		Max           price    `json:"max"`
		Min           price    `json:"min"`
		Shop          shop     `json:"shop"`
	}

	price struct {
		Currency string  `json:"currency"`
		Cost     float64 `json:"cost"`
	}

	shop struct {
		ID           int     `json:"id"`
		Name         string  `json:"name"`
		Followers    int     `json:"followers"`
		PositiveRate float64 `json:"positive_rate"`
	}
)

func New(ch *amqp.Channel, queue amqp.Queue, appID string) app.Store {
	return &store{Name: queue.Name, ch: ch, appID: appID}
}

func (s *store) Send(d app.Document, rate float64) error {
	js, err := json.Marshal(convertDocument(d, rate))
	if err != nil {
		return errors.Wrap(err, "failed to convert json")
	}

	log.Println(string(js))
	// TODO дописать publishing
	publishing := amqp.Publishing{
		AppId:       s.appID,
		ContentType: contentType,
		Body:        js,
		DeliveryMode:2,
	}

	err = s.ch.Publish("", s.Name, false, false, publishing)
	if err != nil {
		return errors.Wrap(err, "failed to publish document")
	}

	return nil
}

func convertDocument(d app.Document, ourRating float64) document {
	return document{
		Title:         d.Title,
		OurRating:     ourRating,
		ID:            d.Id,
		URL:           d.URL,
		TotalSales:    d.TotalSales,
		RatingProduct: d.RatingProduct,
		TotalComment:  d.TotalComment,
		Images:        d.Images,
		Discount:      d.Discount,
		Max: price{
			Currency: d.Max.Currency,
			Cost:     d.Max.Cost,
		},
		Min: price{
			Currency: d.Min.Currency,
			Cost:     d.Min.Cost,
		},
		Shop: shop{
			ID:           d.Shop.ID,
			Name:         d.Shop.Name,
			Followers:    d.Shop.Followers,
			PositiveRate: d.Shop.PositiveRate,
		},
	}
}
