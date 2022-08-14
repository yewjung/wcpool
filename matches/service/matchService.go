package service

import (
	"context"
	matchModels "wcpool/matches/models"
	"wcpool/models"
	"wcpool/utils"

	"wcpool/matches/repository"
)

type MatchService struct {
	Storage   matchModels.Storage
	matchRepo *repository.MatchMongo
}

func (ms *MatchService) GetMatchesByMatchday(matchday string) []matchModels.Match {
	return utils.GetFromCacheOrFunc(context.Background(), ms.Storage.MatchRedis, models.Key(matchday), ms.getMatchRepo().GetMatchesByMatchday)
}

func (ms *MatchService) getMatchRepo() *repository.MatchMongo {
	if ms.matchRepo == nil {
		ms.matchRepo = &repository.MatchMongo{
			Client: ms.Storage.MongoDB,
		}
	}
	return ms.matchRepo
}
