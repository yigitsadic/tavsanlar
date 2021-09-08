package main

import (
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalln("Unable to connect RabbitMQ")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln("Unable to open a channel")
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalln("Unable to declare queue")
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		select {
		case <-signalChannel:
			log.Fatalln("Got quit signal. Quiting...")
		}
	}()

	randomSeed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(randomSeed)

	for {
		time.Sleep(time.Second * time.Duration(random.Intn(4)+1))

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("Hello How are you today?"),
			})
		if err != nil {
			log.Fatalln("unable to publish message")
		}

		log.Println("Published message to queue")
	}
}
