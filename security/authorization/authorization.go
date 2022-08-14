package authorization

import (
	context "context"
	"security/controller"
	"security/models"
	auth "wcpool/authorization"
)

type AuthorizationServerImpl struct {
	auth.UnimplementedAuthorizationServer
	Storage models.Storage
}

func (authServer AuthorizationServerImpl) VerifyPartyID(ctx context.Context, verification *auth.Verification) (*auth.VerificationResult, error) {
	token := verification.GetToken()
	authService := controller.AuthUserService{DB: authServer.Storage.PostgresDB}
	ok, email := authService.IsTokenStillValid(token)
	if !ok {
		return &auth.VerificationResult{
			Ok:    false,
			Email: email,
		}, nil
	}
	accountService := controller.AccountService{MongoDB: authServer.Storage.MongoDB, Cache: authServer.Storage.RedisCache}
	account := accountService.FindByEmail(email)

	verificationMethods := map[auth.Option]models.VerificationMethod{
		auth.Option_PARTY_ID: accountService.IsUserFromParty,
		auth.Option_IS_ADMIN: accountService.IsUserAdminOfParty,
	}
	for _, option := range verification.GetOptions() {
		method, ok := verificationMethods[option]
		if !ok || !method(account, verification.Partyid) {
			return &auth.VerificationResult{
				Ok:    false,
				Email: email,
			}, nil
		}
	}

	return &auth.VerificationResult{
		Ok:    true,
		Email: email,
	}, nil

}
