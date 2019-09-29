package parse

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/jellyboysband/eye/cmd/crawler/internal/app"
	"github.com/powerman/structlog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const URL = "https://ru.aliexpress.com/item/%d.html"

type (
	Parser struct {
		collector *colly.Collector
		app       app.App
		log       *structlog.Logger
		chErr     chan<- error
	}
)

func Run(collector *colly.Collector, log *structlog.Logger, id int, application app.App, chErr chan<- error) {
	p := Parser{
		collector: collector,
		app:       application,
		log:       log,
		chErr:     chErr,
	}

	p.collector.OnHTML("script", p.parserALI)

	for {
		url := fmt.Sprintf(URL, id)

		err := p.collector.Visit(url)
		if err != nil && err.Error() != http.StatusText(http.StatusNotFound) {
			chErr <- err
		}
		id++

		// TODO add random sleep.
		time.Sleep(time.Second * 5)
	}
}

const Code = "OK"

// TODO Добавить регулярку
func (p *Parser) parserALI(element *colly.HTMLElement) {
	if !strings.Contains(element.Text, "window.runParams") {
		return
	}

	js := strings.Trim(element.Text, " \n ")
	js = strings.TrimPrefix(js, "window.runParams = {")
	js = strings.TrimPrefix(js, "\n ")
	js = strings.TrimSpace(js)
	js = strings.TrimPrefix(js, "data: ")
	for !strings.HasSuffix(js, "}}") {
		js = js[:len(js)-1]
	}

	page := Page{}
	if err := json.Unmarshal([]byte(js), &page); err != nil {
		p.chErr <- err
	}

	if !valid(&page) {
		return
	}

	p.log.WarnIfFail(func() error { return p.app.Save(convert(&page)) })
}

func valid(page *Page) bool {
	return page.RedirectModule.Code == Code && page.ActionModule.ItemStatus == 0 && page.ActionModule.TotalAvailQuantity > 0
}

func convert(page *Page) app.Document {
	return app.Document{
		Title:         page.TitleModule.Subject,
		Id:            page.PageModule.ProductID,
		URL:           page.PageModule.URL,
		TotalSales:    page.TitleModule.TradeCount,
		RatingProduct: mustFloat(page.TitleModule.Rating.StarSTR),
		Images:        page.ImageModule.ImagePathList,
		TotalComment:  page.TitleModule.Rating.CountFeedback,
		Discount:      page.PriceModule.Discount,
		Max: app.Price{
			Currency: page.PriceModule.Max.Currency,
			Cost:     page.PriceModule.Max.Value,
		},
		Min: app.Price{
			Currency: page.PriceModule.Min.Currency,
			Cost:     page.PriceModule.Min.Value,
		},
		Shop: app.Shop{
			ID:           page.StoreModule.StoreID,
			Name:         page.StoreModule.StoreName,
			Followers:    page.StoreModule.FollowingNumber,
			PositiveRate: mustFloat(page.StoreModule.PositiveRateSTR[:len(page.StoreModule.PositiveRateSTR)-2]), // TODO закомментировать
		},
	}
}

func mustFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(fmt.Errorf("%w:%s", err, "failed to parse seller rate"))
	}
	return f
}
