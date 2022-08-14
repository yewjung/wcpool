package controller

import (
	"net/http"
	"wcpool/authorization"
	matchesModels "wcpool/matches/models"
	"wcpool/matches/service"
	"wcpool/models"
	"wcpool/utils"
)

type MatchController struct {
	models.AuthorizableController
	Storage                   matchesModels.Storage
	matchAndPredictionService *service.MatchAndPredictionService
}

func (mc *MatchController) GetMatchesAndPredictions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto := utils.GetReqBody(r, matchesModels.MatchRequestDTO{})
		ok, email := mc.CheckAuthorization(w, r, dto.Partyid, []authorization.Option{authorization.Option_PARTY_ID})
		if !ok {
			return
		}
		mps := mc.getMatchAndPredictionService()
		result, err := mps.GetMatchesAndPredictions(dto.Matchday, email, dto.Partyid)
		utils.HandleResponse(w, err, result)
	}
}

func (mc *MatchController) AddPrediction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto := utils.GetReqBody(r, matchesModels.PredictionDTO{})
		ok, email := mc.CheckAuthorization(w, r, dto.PartyID, []authorization.Option{authorization.Option_PARTY_ID})
		if !ok {
			utils.SendError(w, http.StatusUnauthorized, models.Error{
				Message: "Unauthorized action",
			})
			return
		}
		mps := mc.getMatchAndPredictionService()
		key := matchesModels.MatchEmailPartyKey{
			MatchID: dto.MatchID,
			Email:   email,
			PartyID: dto.PartyID,
		}
		err := mps.AddPrediction(key, matchesModels.Prediction{
			GoalA: dto.GoalA,
			GoalB: dto.GoalB,
			Score: dto.Score,
		})
		utils.HandleResponse(w, err, nil)
	}
}

func (mc *MatchController) getMatchAndPredictionService() *service.MatchAndPredictionService {
	if mc.matchAndPredictionService == nil {
		mc.matchAndPredictionService = &service.MatchAndPredictionService{
			Storage: mc.Storage,
		}
	}
	return mc.matchAndPredictionService
}
