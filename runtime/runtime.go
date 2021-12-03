package runtime

import (
	"fmt"
	"go-gin-example/pkg/setting"
	"log"

	"github.com/streadway/amqp"
)

var mq *amqp.Connection
var channel *amqp.Channel

const BuildingQueuePrefix = "Building"

func init() {
	var (
		username, password, host, port string
	)
	sec, err := setting.Cfg.GetSection("rabbitMQ")
	failOnError(err, "Failed to get section rabbitMQ")

	username = sec.Key("USERNAME").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	port = sec.Key("PORT").String()

	mq, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port))
	failOnError(err, "Failed to connect to RabbitMQ")
	// defer mq.Close()

	channel, err = mq.Channel()
	failOnError(err, "Failed to open a channel")
	fmt.Println("connect to " + fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port))

	// defer channel.Close()
	// ticker := time.NewTicker(5 * time.Second)
	go StartConsume()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func StartConsume() {
	que := DeclareQueue(BuildingQueuePrefix)
	ConsumeMessage(que)
}
