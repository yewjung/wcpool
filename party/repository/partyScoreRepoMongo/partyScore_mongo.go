package partyscorerepomongo

import (
	"context"
	"fmt"
	"log"
	"party/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PartyScoreRepoMongo struct{}

var ctx = context.Background()

func (repo *PartyScoreRepoMongo) GetScoresByIDs(mongoDB *mongo.Client, partyid string, emails []string) map[string]int {
	scoreCollection := mongoDB.Database("Party").Collection("ScoreCollection")
	ids := make([]string, len(emails))
	for _, email := range emails {
		ids = append(ids, repo.constructID(partyid, email))
	}
	cursor, err := scoreCollection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		log.Default().Panic(err)
		return make(map[string]int)
	}
	defer cursor.Close(context.TODO())
	results := make(map[string]int)
	for cursor.Next(context.TODO()) {
		score := models.EmailScore{}
		cursor.Decode(&score)
		results[score.Email] = score.Score
	}
	return results
}

func (repo *PartyScoreRepoMongo) UpdateScore(mongoDB *mongo.Client, partyid string, email string, score int) error {
	scoreCollection := mongoDB.Database("Party").Collection("ScoreCollection")
	_, err := scoreCollection.UpdateByID(ctx, repo.constructID(partyid, email), bson.M{"$set": bson.M{"score": score}})
	if err != nil {
		log.Default().Panic(err)
		return err
	}
	return nil
}

func (repo *PartyScoreRepoMongo) AddScore(mongoDB *mongo.Client, partyid string, email string, score int) error {
	scoreCollection := mongoDB.Database("Party").Collection("ScoreCollection")
	doc := bson.M{"_id": repo.constructID(partyid, email), "score": score}
	_, err := scoreCollection.InsertOne(ctx, doc)
	if err != nil {
		log.Default().Panic(err)
	}
	return err

}

func (repo *PartyScoreRepoMongo) constructID(partyid string, email string) string {
	return fmt.Sprintf(partyid, email)
}
