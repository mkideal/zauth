package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/model"
	"bitbucket.org/mkideal/accountd/oauth2"
)

func (svr *Server) handleAccessToken(w http.ResponseWriter, r *http.Request) {
	argv := new(api.AccessTokenReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("AccessToken parse arguments error: %v, IP=%v", err, httputil.IP(r))
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("AccessToken request, IP=%v", httputil.IP(r))

	clientId, clientSecret, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		log.Warn("%s: BasicAuth failure", argv.CommandName())
		return
	}
	log.Debug("%s: clientId=%s, clientSecret=%s", argv.CommandName(), clientId, clientSecret)

	client, err := svr.clientRepo.GetClient(clientId)
	if err != nil {
		log.Error("%s: get client %s error: %v", argv.CommandName(), clientId, err)
		svr.errorResponse(argv.CommandName(), w, err)
		return
	}
	if client == nil {
		log.Info("%s: client %s not found", argv.CommandName(), clientId)
		svr.oauthErrorResponse(argv.CommandName(), w, oauth2.ErrorInvalidClient)
		return
	}

	switch argv.GrantType {
	case oauth2.GrantAuthenticationCode:
		err = svr.grantByAuthorizationCode(w, argv, client)
	case oauth2.GrantPassword:
		err = svr.grantByPassword(w, argv, client, clientSecret)
	case oauth2.GrantRefreshToken:
		err = svr.grantByRefreshToken(w, argv, client, clientSecret)
	default:
		svr.oauthErrorResponse(argv.CommandName(), w, oauth2.ErrorUnsupportedGrantType)
	}
	if err != nil {
		log.Error("%s: grant error: %v", argv.CommandName(), err)
		svr.errorResponse(argv.CommandName(), w, err)
	}
}

func (svr *Server) grantByAuthorizationCode(w http.ResponseWriter, argv *api.AccessTokenReq, client *model.Client) error {
	ar, err := svr.authRepo.GetAuthRequest(client.Id, argv.Code)
	if err != nil {
		return err
	}
	if ar == nil {
		return oauth2.NewError(oauth2.ErrorInvalidGrant, "code-not-found")
	}
	user, err := svr.userRepo.GetUser(ar.Uid)
	if err != nil {
		return err
	}
	if user == nil {
		return oauth2.NewError(oauth2.ErrorInvalidGrant, "user-not-found")
	}
	accessToken, err := svr.tokenRepo.NewToken(client, user, ar.GrantedScopes)
	if err != nil {
		return err
	}
	log.WithJSON(accessToken).Debug("new token")
	svr.authRepo.RemoveAuthRequest(ar.Id)

	svr.response(w, http.StatusOK, api.AccessTokenRes{
		TokenType:    oauth2.TokenType,
		Scope:        accessToken.Scope,
		AccessToken:  accessToken.Token,
		RefreshToken: accessToken.RefreshToken,
		ExpireAt:     accessToken.ExpireAt,
	})
	return nil
}

func (svr *Server) grantByPassword(w http.ResponseWriter, argv *api.AccessTokenReq, client *model.Client, clientSecret string) error {
	return oauth2.NewError(oauth2.ErrorUnsupportedGrantType, "grantByPassword-not-implemented")
}

func (svr *Server) grantByRefreshToken(w http.ResponseWriter, argv *api.AccessTokenReq, client *model.Client, clientSecret string) error {
	if !model.ValidateClint(client, clientSecret) {
		return oauth2.NewError(oauth2.ErrorInvalidGrant, "invalid-client-secret")
	}
	accessToken, err := svr.tokenRepo.RefreshToken(client, argv.RefreshToken, argv.Scope)
	if err != nil {
		return err
	}

	svr.response(w, http.StatusOK, api.AccessTokenRes{
		TokenType:    oauth2.TokenType,
		Scope:        accessToken.Scope,
		AccessToken:  accessToken.Token,
		RefreshToken: accessToken.RefreshToken,
		ExpireAt:     accessToken.ExpireAt,
	})
	return nil
}
