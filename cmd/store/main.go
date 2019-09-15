package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jellyboysband/eye/cmd/store/internal/app"
	"github.com/jellyboysband/eye/cmd/store/internal/server"
	"github.com/jellyboysband/eye/cmd/store/internal/store"
	"github.com/pkg/errors"
	"github.com/powerman/structlog"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net"
	"strconv"
)

var (
	log = structlog.New()

	cfg struct {
		dbPass     string
		dbUsername string
		dbHost     string
		dbPort     int

		InsertQueueName string
		SendQueueName   string
		appID           string
		expectedAppID   string
		rabbitPass      string
		rabbitUser      string
		rabbitHost      string
		rabbitPort      int
	}
)

func init() {
	flag.StringVar(&cfg.dbPass, "db.pass", "mongo", "DB password.")
	flag.StringVar(&cfg.dbUsername, "db.username", "mongo", "DB username.")
	flag.StringVar(&cfg.dbHost, "db.host", "localhost", "DB host.")
	flag.IntVar(&cfg.dbPort, "db.port", 27017, "DB port.")

	flag.StringVar(&cfg.appID, "app.id", "parser", "Application ID")
	flag.StringVar(&cfg.expectedAppID, "expected.app.id", "parser", "Expected application ID")

	flag.StringVar(&cfg.InsertQueueName, "insert.queue.name", "filtered_products", "Insert queue name.")
	flag.StringVar(&cfg.SendQueueName, "send.queue.name", "ready_products", "Send queue name.")
	flag.StringVar(&cfg.rabbitPass, "rabbit.pass", "rabbitmq", "RabbitMQ password.")
	flag.StringVar(&cfg.rabbitUser, "rabbit.user", "rabbitmq", "RabbitMQ user.")
	flag.StringVar(&cfg.rabbitHost, "rabbit.host", "localhost", "RabbitMQ host.")
	flag.IntVar(&cfg.rabbitPort, "rabbit.port", 5672, "RabbitMQ port.")

	flag.Parse()
}

func main() {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d",
		cfg.rabbitUser,
		cfg.rabbitPass,
		cfg.rabbitHost,
		cfg.rabbitPort,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	defer log.WarnIfFail(conn.Close)

	opts := options.Client().ApplyURI("mongodb://" + net.JoinHostPort(cfg.dbHost, strconv.Itoa(cfg.dbPort))).
		SetAuth(options.Credential{
			Username: cfg.dbUsername,
			Password: cfg.dbPass,
		})

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Store started")
	log.Fatal(run(conn, client))
}

func run(conn *amqp.Connection, client *mongo.Client) error {
	ch, err := conn.Channel()
	if err != nil {
		return errors.Wrap(err, "failed to get channel rabbit")
	}
	defer log.WarnIfFail(ch.Close)

	queue, err := ch.QueueDeclare(cfg.SendQueueName, false, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "failed to declare queue")
	}

	msgs, err := ch.Consume(
		cfg.InsertQueueName, // queue
		"",                  // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		return errors.Wrapf(err, "failed to consume")
	}

	storage := store.New(client)
	application := app.New(storage)
	s := server.New(cfg.appID, application, cfg.expectedAppID)
	return s.Listen(msgs, ch, queue.Name, log)
}
