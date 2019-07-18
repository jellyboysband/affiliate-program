package parser

import (
	"encoding/json"
	"strings"

	"github.com/gocolly/colly"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Parser struct {
	url       string
	collector *colly.Collector
	log       *logrus.Logger
	documents chan<- json.RawMessage
}

func Build(url, selector string, log *logrus.Logger, ch chan<- json.RawMessage, domains ...string) *Parser {
	collector := colly.NewCollector(colly.AllowedDomains(domains...))

	p := &Parser{collector: collector, log: log, documents: ch, url: url}
	p.collector.OnHTML("a[href]", p.collectLink)
	p.collector.OnHTML(selector, p.parserALI)

	return p
}

func (p *Parser) Start() {
	p.collector.Visit(p.url)
}

func (p *Parser) collectLink(e *colly.HTMLElement) {
	link := e.Attr("href")
	if !strings.Contains(link, "/item/") {
		p.log.Print(link)
		return
	}

	p.log.Println("Next page link found:", link)
	if err := errors.Wrapf(e.Request.Visit(link), "failed to visit"); err != nil {
		p.log.Warn(err)
	}
}

func (p *Parser) parserALI(element *colly.HTMLElement) {
	if !strings.Contains(element.Text, "window.runParams") {
		return
	}

	doc := strings.Trim(element.Text, " \n ")
	doc = strings.TrimPrefix(doc, "window.runParams = {")
	doc = strings.TrimPrefix(doc, "\n ")
	doc = strings.TrimSpace(doc)
	doc = strings.TrimPrefix(doc, "data: ")
	for !strings.HasSuffix(doc, "}}") {
		doc = doc[:len(doc)-1]
	}

	p.documents <- json.RawMessage(doc)

}
