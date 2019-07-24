package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

/*
Consumer struct
*/
type Consumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      string
	exchange   string
	name       string
}

/*
MessageHandler struct
*/
type MessageHandler struct {
}

func (*MessageHandler) handle(msg map[string]*string) {
	// persist
}

/*
Start starts consumer
*/
func (c *Consumer) Start(ctx context.Context, handler MessageHandler) {

	var values map[string]*string

	errorChan := make(chan error, 1)

	c.queue = "Metrics"
	c.exchange = "proxy.new-request"
	c.name = "Metrics"

	c.connect()

	for {
		select {
		case <-ctx.Done():
			defer c.connection.Close()
			defer c.channel.Close()
			ctx.Err()
			return
		case err := <-errorChan:
			failOnError(err, "", "")
			defer c.connection.Close()
			defer c.channel.Close()
			return
		default:
			mc := c.messageChannel()

			for d := range mc {

				err := json.Unmarshal(d.Body, &values)

				failOnError(err, "Unmarshal error", "Error")
				b, err := json.MarshalIndent(values, "", "  ")
				if err != nil {
					fmt.Println("error:", err)
				}
				fmt.Print(string(b))

				if err := d.Ack(false); err != nil {
					log.Printf("Error acknowledging message : %s", err)
				} else {
					log.Printf("Acknowledged message")
				}
			}
		}
	}

}

func failOnError(err error, msgerr string, msgsuc string) {
	if err != nil {
		log.Fatalf("%s: %s", msgerr, err)

	} else {
		fmt.Printf("%s\n", msgsuc)
	}

}

func (c *Consumer) connect() {
	fmt.Println("Connecting to RabbitMQ ...")
	var err error
	c.connection, err = amqp.Dial("") //Insert the  connection string
	failOnError(err, "RabbitMQ connection failure", "RabbitMQ Connection Established")

	c.channel, err = c.connection.Channel()
	failOnError(err, "Failed to open a channel", "Opened the channel")

	_, err = c.channel.QueueDeclare(
		c.queue, // queue name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // args
	)

	failOnError(err, "Failed to declare the queue", "Declared the queue")

	err = c.channel.QueueBind(
		c.name,     // name
		"*",        // routing key
		c.exchange, // exchange
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to publish a message ", "Published the message")

}

/*
readMetrics allows to send metric information from a HTTP Request to a AMQP instance.

*/
func (c *Consumer) messageChannel() <-chan amqp.Delivery {

	consumeChannel, err := c.channel.Consume(
		c.queue, //routing key
		c.name,  // consumer name
		false,   //auto ack
		false,   //exclusive
		false,   // no local
		false,   // no wait
		nil,     // args
	)

	failOnError(err, "Failed to publish a message ", "Published the message")

	return consumeChannel

}
