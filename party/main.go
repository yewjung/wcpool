package main

import (
	"fmt"
	"log"
	"net/http"
	"party/controller"
	partyDriver "party/driver"
	partyModels "party/models"
	"party/service"
	"wcpool/authorization"
	"wcpool/constants"
	"wcpool/driver"
	"wcpool/models"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func main() {
	partyDB := driver.ConnectPostgresDB()
	mongoDB := driver.ConnectMongoDB()
	redisCache := partyDriver.ConnectRedis()
	storage := partyModels.Storage{
		PostgresPartyDB: partyDB,
		MongoDB:         mongoDB,
		RedisCache:      redisCache,
	}

	router := mux.NewRouter()

	go listenToGameEvents(storage)

	partyController := controller.PartyController{
		AuthorizableController: models.AuthorizableController{
			AuthClient: getSecurityGrpcClient(),
		},
	}
	// /leaderboard/{partyid}/
	router.HandleFunc("/leaderboard/{partyid}", partyController.GetLeaderboard(storage)).Methods("GET")
	// /score data: {partyid, email, score}
	router.HandleFunc("/score", partyController.UpdateScore(storage)).Methods("POST")
	// /member data: {partyid, email}
	router.HandleFunc("/member", partyController.AddMemberToParty(storage)).Methods("POST")
	// /party data: {partyname}
	router.HandleFunc("/party", partyController.AddParty(storage)).Methods("POST")

	fmt.Println("Server is running at port 8090")

	log.Fatal(
		http.ListenAndServe(":8090",
			handlers.CORS(
				handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
				handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
				handlers.AllowedOrigins([]string{"*"}),
			)(router),
		),
	)

}

func getSecurityGrpcClient() authorization.AuthorizationClient {
	conn, err := grpc.Dial("security:8085")
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return authorization.NewAuthorizationClient(conn)
}

func listenToGameEvents(storage partyModels.Storage) {
	ch := driver.ConnectToChannel()
	defer ch.Close()

	msgs, err := ch.Consume(
		constants.SCORE_QUEUE,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully Connected to our RabbitMQ Instance")
	fmt.Println(" [*] - Waiting for messages")

	ps := service.PartyService{
		Storage: storage,
	}
	for d := range msgs {
		go ps.UpdateLiveScore(&d)
	}

}
