package main

import (
	"flag"
	"fmt"

	"github.com/gocolly/colly"
	"github.com/jellyboysband/eye/cmd/crawler/internal/app"
	"github.com/jellyboysband/eye/cmd/crawler/internal/filter"
	"github.com/jellyboysband/eye/cmd/crawler/internal/parse"
	"github.com/pkg/errors"
	"github.com/powerman/structlog"
	"github.com/streadway/amqp"
)

var (
	log = structlog.New()

	cfg struct {
		domain    string
		startID   int
		queueName string
		rateMin   float64

		appID      string
		rabbitPass string
		rabbitUser string
		rabbitHost string
		rabbitPort int
	}
)

func init() {
	flag.StringVar(&cfg.domain, "domain", "ru.aliexpress.com", "Main domain for parse.")
	flag.IntVar(&cfg.startID, "startID", 33035919711+380000, "First parsing product.")

	flag.StringVar(&cfg.appID, "app.id", "parser", "Application ID")

	flag.StringVar(&cfg.queueName, "queue.name", "documents_ali", "Queue name.")
	flag.StringVar(&cfg.rabbitPass, "rabbit.pass", "rabbitmq", "RabbitMQ password.")
	flag.StringVar(&cfg.rabbitUser, "rabbit.user", "rabbitmq", "RabbitMQ user.")
	flag.StringVar(&cfg.rabbitHost, "rabbit.host", "localhost", "RabbitMQ host.")
	flag.IntVar(&cfg.rabbitPort, "rabbit.port", 5672, "RabbitMQ port.")

	flag.Float64Var(&cfg.rateMin, "rate-min", 0.05, "Min our rate for send in queue.")

	flag.Parse()
}

func main() {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d",
		cfg.rabbitUser,
		cfg.rabbitPass,
		cfg.rabbitHost,
		cfg.rabbitPort,
	)
	log.Println(url)
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	defer log.WarnIfFail(conn.Close)
	log.Info("Parser started")
	log.Fatal(run(conn))
}

func run(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return errors.Wrap(err, "failed to get channel rabbit")
	}
	defer log.WarnIfFail(ch.Close)

	queue, err := ch.QueueDeclare(cfg.queueName, false, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "failed to declare queue")
	}

	f := filter.New(ch, queue, cfg.appID)
	collector := colly.NewCollector(colly.AllowedDomains(cfg.domain))
	application := app.New(f, cfg.rateMin)

	chErr := make(chan error)
	go parse.Run(collector, log, cfg.startID, application, chErr)
	for err := range chErr {
		log.Warn(err)
	}

	return nil
}
