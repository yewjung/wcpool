package events

import (
	"fmt"
	"log"
	"net/http"
	"wcpool/constants"
	"wcpool/driver"
	"wcpool/utils"

	"github.com/streadway/amqp"
)

var eventChannel <-chan amqp.Delivery

func main() {
	router := http.NewServeMux()

	eventChannel = subscribeToLiveServer()

	router.HandleFunc("/event", sseHandler)

	log.Fatal(http.ListenAndServe(":3500", router))
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Client connected")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		fmt.Println("Could not init http.Flusher")
	}

	for {
		select {
		case message := <-eventChannel:
			fmt.Println("case message... sending message")
			fmt.Println(message.Body)
			fmt.Fprintf(w, "data: %s\n\n", message.Body)
			flusher.Flush()
		case <-r.Context().Done():
			fmt.Println("Client closed connection")
			return
		}
	}

}

func subscribeToLiveServer() <-chan amqp.Delivery {
	ch := driver.ConnectToChannel()
	defer ch.Close()

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

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,                        // queue name
		"",                            // routing key
		constants.GAME_EVENT_EXCHANGE, // exchange
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utils.FailOnError(err, "Failed to register a consumer")

	return msgs
}
