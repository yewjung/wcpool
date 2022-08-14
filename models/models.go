package models

import (
	"context"
	"net/http"
	"wcpool/authorization"
)

type Error struct {
	Message string `json:"errorMessage"`
}

type Key string

func (k Key) String() string {
	return string(k)
}

type AuthorizableController struct {
	AuthClient authorization.AuthorizationClient
}

func (authController *AuthorizableController) CheckAuthorization(w http.ResponseWriter, r *http.Request, partyid string, options []authorization.Option) (bool, string) {
	verRes, err := authController.AuthClient.VerifyPartyID(context.Background(), &authorization.Verification{
		Token:   r.Header.Get("Authorization"),
		Partyid: partyid,
		Options: options,
	})
	if err != nil {
		return false, verRes.Email
	}
	if !verRes.Ok {
		return false, verRes.Email
	}
	return true, verRes.Email
}
