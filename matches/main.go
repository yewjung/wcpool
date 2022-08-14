package matches

import (
	"log"
	"wcpool/authorization"
	"wcpool/driver"
	"wcpool/matches/controller"
	matchesDriver "wcpool/matches/driver"
	matchesModels "wcpool/matches/models"
	"wcpool/models"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func main() {
	storage := matchesModels.Storage{
		MongoDB:         driver.ConnectMongoDB(),
		MatchRedis:      matchesDriver.ConnectMatchesRedis(),
		PredictionRedis: matchesDriver.ConnectPredictionsRedis(),
	}
	matchController := controller.MatchController{
		AuthorizableController: models.AuthorizableController{
			AuthClient: getSecurityGrpcClient(),
		},
		Storage: storage,
	}
	router := mux.NewRouter()
	router.HandleFunc("/matches", matchController.GetMatchesAndPredictions()).Methods("POST")
	router.HandleFunc("/prediction", matchController.AddPrediction()).Methods("POST")

}

func getSecurityGrpcClient() authorization.AuthorizationClient {
	conn, err := grpc.Dial("security:8085")
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return authorization.NewAuthorizationClient(conn)
}
