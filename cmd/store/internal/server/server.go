package server

import (
	"context"
	"encoding/json"
	"github.com/jellyboysband/eye/cmd/store/internal/app"
	"github.com/pkg/errors"
	"github.com/powerman/structlog"
	"github.com/streadway/amqp"
	"time"
)

type (
	Server struct {
		appID         string
		app           app.App
		expectedAppID string
	}

	document struct {
		Images        []string `json:"images"`
		OurRating     float64  `json:"our_rating"`
		StoreID       string   `json:"store_id"`
		Title         string   `json:"title"`
		ID            int      `json:"id"`
		URL           string   `json:"url"`
		TotalSales    int      `json:"total_sales"`
		RatingProduct float64  `json:"rating_product"`
		TotalComment  int      `json:"total_comment"`
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
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Followers    int    `json:"followers"`
		PositiveRate float64 `json:"positive_rate"`
	}
)

func New(appID string, app app.App, expectedAppID string) *Server {
	return &Server{appID, app, expectedAppID}
}

const timeoutSaveProduct = time.Second * 100

func (s *Server) Listen(insertProducts <-chan amqp.Delivery, sendProducts *amqp.Channel, queueName string, log *structlog.Logger) error {

	// TODO refactor
	for val := range insertProducts {
		// TODO Добавить дополнительные проверки
		log.Println(string(val.Body))
		switch {
		case val.AppId != s.expectedAppID:
			log.Warn("unknown source", "appID:", val.AppId)
			continue
		case val.ContentType != contentType:
			log.Warn("unknown content type", "content type:", val.ContentType)
			continue
		}

		// TODO move
		d := &document{}
		err := json.Unmarshal(val.Body, &d)
		if err != nil {
			log.Warn(errors.Wrap(err, "failed to unmarshal document"))
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeoutSaveProduct)
		product, err := s.app.Save(ctx, documentToProduct(d))
		if err != nil {
			log.Warn(errors.Wrap(err, "failed to save product"))
			continue
		}

		// TODO move
		js, err := json.Marshal(productToDocument(product))
		if err != nil {
			log.Warn(errors.Wrap(err, "failed to convert json"))
			continue
		}

		// TODO дописать publishing
		publishing := amqp.Publishing{AppId: s.appID, ContentType: contentType, Body: js}
		err = sendProducts.Publish("", queueName, false, false, publishing)
		if err != nil {
			log.Warn(errors.Wrap(err, "failed to publish document"))
			continue
		}

		cancel()
	}

	return nil
}

const (
	contentType = "application/json"
)

func productToDocument(product *app.Product) document {
	return document{
		StoreID:       product.Id,
		Images:        product.Images,
		OurRating:     product.OurRating,
		Title:         product.Title,
		ID:            product.AliID,
		URL:           product.URL,
		TotalSales:    product.TotalSales,
		RatingProduct: product.RatingProduct,
		TotalComment:  product.TotalComment,
		Discount:      product.Discount,
		Max: price{
			Currency: product.Max.Currency,
			Cost:     product.Max.Cost,
		},
		Min: price{
			Currency: product.Min.Currency,
			Cost:     product.Min.Cost,
		},
		Shop: shop{
			ID:           product.Shop.ID,
			Name:         product.Shop.Name,
			Followers:    product.Shop.Followers,
			PositiveRate: product.Shop.PositiveRate,
		},
	}
}

func documentToProduct(doc *document) app.ArgSaveProduct {
	return app.ArgSaveProduct{
		AliID:         doc.ID,
		OurRating:     doc.OurRating,
		Rating:        doc.RatingProduct,
		Images:        doc.Images,
		URL:           doc.URL,
		Title:         doc.Title,
		TotalSales:    doc.TotalSales,
		RatingProduct: doc.RatingProduct,
		TotalComment:  doc.TotalComment,
		Discount:      doc.Discount,
		Max: app.Price{
			Currency: doc.Max.Currency,
			Cost:     doc.Max.Cost,
		},
		Min: app.Price{
			Currency: doc.Min.Currency,
			Cost:     doc.Min.Cost,
		},
		Shop: app.Shop{
			ID:           doc.Shop.ID,
			Name:         doc.Shop.Name,
			Followers:    doc.Shop.Followers,
			PositiveRate: doc.Shop.PositiveRate,
		},
	}
}
