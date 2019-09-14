package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/powerman/structlog"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type (
	document struct {
		Va            primitive.ObjectID `json:"_id"`
		Title         string             `json:"title"`
		ID            int                `json:"id"`
		URL           string             `json:"url"`
		TotalSales    int                `json:"total_sales"`
		RatingProduct string             `json:"rating_product"`
		TotalComment  int                `json:"total_comment"`
		Discount      float64            `json:"discount"`
		Max           price              `json:"max"`
		Min           price              `json:"min"`
		Shop          shop               `json:"shop"`
	}

	price struct {
		Currency string  `json:"currency"`
		Cost     float64 `json:"cost"`
	}

	shop struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Followers    int    `json:"followers"`
		PositiveRate string `json:"positive_rate"`
	}
)

func main() {
	log := structlog.New()

	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer log.WarnIfFail(conn.Close)

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer log.WarnIfFail(ch.Close)

	q, err := ch.QueueDeclare(
		"documents_ali", // name
		false,           // durable
		false,           // delete when usused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	cread := options.Credential{
		Username: "mongo",
		Password: "mongo",
	}
	opts := options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(cread)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	collection := client.Database("test").Collection("test")

	for d := range msgs {
		var val document
		err := json.Unmarshal(d.Body, &val)
		if err != nil {
			log.Fatal(err)
		}
		val.Va = primitive.NewObjectIDFromTimestamp(time.Now())
		insertResult, err := collection.InsertOne(context.TODO(), val)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(insertResult)

		var val2 document
		f := bson.D{
			{"id", bson.D{{"$gte", val.ID}}},
		}
		err = collection.FindOne(context.TODO(), f).Decode(&val2)
		log.Println(val2)
		log.Println(err)

		var arr []document
		filter := bson.M{}
		c, err := collection.Find(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		err = c.All(context.TODO(), &arr)
		log.Println(arr)
	}
}
