package userRepoMongo

import (
	"context"
	"log"
	"wcpool/party/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepoMongo struct{}

var ctx = context.Background()

func (repo *UserRepoMongo) GetUsernamesByEmails(mongoDB *mongo.Client, emails []string) map[string]string {
	accountCollection := mongoDB.Database("Account").Collection("AccountCollection")
	opts := options.Find().SetProjection(bson.M{"_id": 1, "username": 1})
	cursor, err := accountCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": emails}}, opts)
	if err != nil {
		log.Default().Panic(err)
		return make(map[string]string)
	}
	defer cursor.Close(ctx)
	results := make(map[string]string)
	for cursor.Next(ctx) {
		emailUsername := models.EmailUsername{}
		cursor.Decode(&emailUsername)
		results[emailUsername.Email] = emailUsername.Username
	}
	return results
}
