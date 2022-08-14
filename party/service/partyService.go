package service

import (
	"context"
	"encoding/json"
	"log"
	partyModels "party/models"
	partynamerepomongo "party/repository/partyNameRepoMongo"
	"party/repository/partyRepoSql"
	partyscorerepomongo "party/repository/partyScoreRepoMongo"
	"party/repository/userRepoMongo"
	"sync"
	"wcpool/models"
	"wcpool/utils"

	"github.com/go-redis/redis/v9"
	"github.com/streadway/amqp"
)

type PartyService struct {
	Storage partyModels.Storage
}

var ctx = context.Background()

func (ps *PartyService) GetLeaderboard(partyid string) partyModels.Leaderboard {
	return utils.GetFromCacheOrFunc(context.Background(), ps.Storage.RedisCache, models.Key(partyid), ps.GetLeaderboardFromDB)
}

func (ps *PartyService) GetLeaderboardFromDB(partyidKey models.Key) partyModels.Leaderboard {
	postgresDB := ps.Storage.PostgresPartyDB
	partyRepo := partyRepoSql.PartyRepo{
		DB: postgresDB,
	}
	emails, err := partyRepo.GetPartyMemberIDs(partyidKey.String())
	if err != nil {
		log.Default().Panic(err)
		return partyModels.Leaderboard{}
	}
	// find all usernames based on emails from user db (mongoDB)
	var wg sync.WaitGroup
	wg.Add(3)
	var emailUsernames map[string]string
	go func() {
		userRepo := userRepoMongo.UserRepoMongo{}
		emailUsernames = userRepo.GetUsernamesByEmails(ps.Storage.MongoDB, emails)
		wg.Done()
	}()

	// find party name
	var partyName string
	go func() {
		partyNameRepo := partynamerepomongo.PartyNameRepoMongo{}
		partyName, _ = partyNameRepo.GetPartyName(ps.Storage.MongoDB, partyidKey.String())
		wg.Done()
	}()

	// find score for each email
	var emailScore map[string]int
	go func() {
		scoreRepo := partyscorerepomongo.PartyScoreRepoMongo{}
		emailScore = scoreRepo.GetScoresByIDs(ps.Storage.MongoDB, partyidKey.String(), emails)
		wg.Done()
	}()

	wg.Wait()

	// construct leaderboard
	members := make([]partyModels.Member, len(emailScore))
	for _, email := range emails {
		member := partyModels.Member{
			Email:    email,
			Username: emailUsernames[email],
			Score:    emailScore[email],
		}
		members = append(members, member)
	}
	leaderboard := partyModels.Leaderboard{
		Name:    partyName,
		Members: members,
	}

	// return leaderboard
	return leaderboard
}

func (ps *PartyService) UpdateScore(partyid string, email string, score int) error {
	// policy: write around cache
	// remove leaderboard entry from cache
	ps.removePartyFromCache(ps.Storage.RedisCache, partyid)

	// update score db with new score
	scoreRepo := partyscorerepomongo.PartyScoreRepoMongo{}
	err := scoreRepo.UpdateScore(ps.Storage.MongoDB, partyid, email, score)
	if err != nil {
		log.Default().Panic(err)
	}
	return err
}

func (ps *PartyService) AddMemberToParty(partyid string, email string) error {
	// remove from cache
	ps.removePartyFromCache(ps.Storage.RedisCache, partyid)

	// add record to party (member) db
	memberRepo := partyRepoSql.PartyRepo{DB: ps.Storage.PostgresPartyDB}
	err := memberRepo.AddMemberToParty(partyid, email)
	if err != nil {
		return err
	}

	// add record to score db
	scoreRepo := partyscorerepomongo.PartyScoreRepoMongo{}
	return scoreRepo.AddScore(ps.Storage.MongoDB, partyid, email, 0)
}

func (ps *PartyService) AddParty(name string) (interface{}, error) {
	// add new party to party name db
	partyNameRepo := partynamerepomongo.PartyNameRepoMongo{}
	return partyNameRepo.AddParty(ps.Storage.MongoDB, name)
}

func (ps *PartyService) removePartyFromCache(cache *redis.Client, partyid string) {
	err := cache.Del(ctx, partyid).Err()
	if err != nil {
		log.Default().Panic(err)
	}
}

func (ps *PartyService) UpdateLiveScore(d *amqp.Delivery) {
	dto := partyModels.UpdateScoreDTO{}
	err := json.Unmarshal(d.Body, &dto)
	if err != nil {
		d.Reject(true)
		return
	}
	err = ps.UpdateScore(dto.PartyId, dto.Email, dto.Score)
	if err != nil {
		d.Reject(true)
		return
	}
	d.Ack(false)
}
