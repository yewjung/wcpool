package controller

import (
	"net/http"
	"wcpool/authorization"
	"wcpool/models"
	partyModels "wcpool/party/models"
	"wcpool/party/service"
	"wcpool/utils"

	"github.com/gorilla/mux"
)

type PartyController struct {
	models.AuthorizableController
}

func (pc *PartyController) GetLeaderboard(storage partyModels.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		partyid := mux.Vars(r)["partyid"]

		// TODO: need to verify that user has access to this party

		// Looking for a valid feature
		ok, _ := pc.CheckAuthorization(w, r, partyid, []authorization.Option{authorization.Option_PARTY_ID})
		if !ok {
			return
		}
		partyService := service.PartyService{Storage: storage}
		leaderboard := partyService.GetLeaderboard(partyid)
		utils.HandleResponse(w, nil, leaderboard)
	}
}

// this method should be called through grpc
func (pc *PartyController) UpdateScore(storage partyModels.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto := utils.GetReqBody(r, partyModels.MemberScore{})
		partyService := service.PartyService{
			Storage: storage,
		}
		err := partyService.UpdateScore(dto.PartyId, dto.Email, dto.Score)
		utils.HandleResponse(w, err, nil)
	}
}

func (pc *PartyController) AddMemberToParty(storage partyModels.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto := utils.GetReqBody(r, partyModels.MemberScore{})
		// TODO: verify permission -> only admin can approve member
		// permission should be a string like :
		// <partyid>$admin
		// use email to find list of permission, use partyid to construct <partyid>$admin
		// user's profile can have a field called permissions = [..., ..., ...]
		ok, _ := pc.CheckAuthorization(w, r, dto.PartyId, []authorization.Option{authorization.Option_PARTY_ID})
		if !ok {
			return
		}
		partyService := service.PartyService{
			Storage: storage,
		}
		err := partyService.AddMemberToParty(dto.PartyId, dto.Email)
		utils.HandleResponse(w, err, nil)
	}
}
func (pc *PartyController) AddParty(storage partyModels.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		party := utils.GetReqBody(r, partyModels.Party{})
		ok, _ := pc.CheckAuthorization(w, r, party.PartyId, nil)
		if !ok {
			return
		}
		partyService := service.PartyService{
			Storage: storage,
		}
		result, err := partyService.AddParty(party.Name)
		utils.HandleResponse(w, err, result)
	}
}
