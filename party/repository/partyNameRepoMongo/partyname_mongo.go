package partynamerepomongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PartyNameRepoMongo struct{}

var ctx = context.Background()

func (repo *PartyNameRepoMongo) GetPartyName(mongoDB *mongo.Client, partyid string) (string, error) {
	partyNameCollection := mongoDB.Database("Party").Collection("PartyName")
	opts := options.FindOne().SetProjection(bson.M{"_id": 0})
	objId, _ := primitive.ObjectIDFromHex(partyid)
	result := partyNameCollection.FindOne(ctx, bson.M{"_id": objId}, opts)
	var partyname string
	err := result.Decode(&partyname) // im not sure if this will work
	if err != nil {
		log.Default().Panic(err)
	}
	return partyname, err
}

func (repo *PartyNameRepoMongo) AddParty(mongoDB *mongo.Client, partyname string) (interface{}, error) {
	partyNameCollection := mongoDB.Database("Party").Collection("PartyName")
	result, err := partyNameCollection.InsertOne(ctx, bson.M{"name": partyname})
	if err != nil {
		log.Default().Panic(err)
		return "", err
	}
	oid, _ := result.InsertedID.(primitive.ObjectID)
	return oid.Hex(), nil
}
