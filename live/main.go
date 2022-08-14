package main

import (
	"encoding/json"
	"fmt"
	"live/models"
	"wcpool/constants"
	"wcpool/driver"
	"wcpool/utils"

	"github.com/streadway/amqp"
)

func main() {
	ch := driver.ConnectToChannel()
	defer ch.Close()

	// with this channel open, we can then start to interact
	// with the instance and declare Queues that we can publish and
	// subscribe to
	_, err := ch.QueueDeclare(
		constants.SCORE_QUEUE,
		true,
		false,
		false,
		false,
		nil,
	)

	initExchange(ch)

	// Handle any errors if we were unable to create the queue
	if err != nil {
		fmt.Println(err)
	}

	events := make(chan models.UpdateScoreDTO, 5)
	go listenToLiveServer(events)
	go broadcastEvent(events, ch)
	forever := make(chan bool)
	<-forever
}

func listenToLiveServer(events chan<- models.UpdateScoreDTO) {
	dto := models.UpdateScoreDTO{
		PartyId: "",
		Email:   "",
		Score:   0,
	}
	events <- dto
}

func broadcastEvent(events <-chan models.UpdateScoreDTO, channel *amqp.Channel) {
	for event := range events {
		broadcast(event, channel)
	}
}

func broadcast(event interface{}, channel *amqp.Channel) {
	bytes, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Failed to send message")
		return
	}
	channel.Publish(
		"",
		constants.SCORE_QUEUE,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		},
	)

	err = channel.Publish(
		constants.GAME_EVENT_EXCHANGE, // exchange
		"",                            // routing key
		false,                         // mandatory
		false,                         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		})
}

func initExchange(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		constants.GAME_EVENT_EXCHANGE, // name
		constants.FANOUT,              // type
		true,                          // durable
		false,                         // auto-deleted
		false,                         // internal
		false,                         // no-wait
		nil,                           // arguments
	)
	utils.FailOnError(err, "Failed to declare an exchange")
}
