package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/model"
	"bitbucket.org/mkideal/accountd/oauth2"
)

func (svr *Server) handleToken(w http.ResponseWriter, r *http.Request) {
	ip := httputil.IP(r)
	argv := new(api.TokenReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("Token parse arguments error: %v, IP=%v", err, ip)
		svr.errorResponse(w, r, api.ErrorCode_BadArgument.NewError(err.Error()))
		return
	}
	log.WithJSON(argv).Debug("Token request, IP=%v", ip)

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
		svr.errorResponse(w, r, err)
		return
	}
	if client == nil {
		log.Info("%s: client %s not found", argv.CommandName(), clientId)
		svr.errorResponse(w, r, oauth2.ErrorInvalidClient)
		return
	}

	switch argv.GrantType {
	case oauth2.GrantAuthenticationCode:
		err = svr.grantByAuthorizationCode(w, r, argv, client)
	case oauth2.GrantPassword:
		err = svr.grantByPassword(w, r, argv, client, clientSecret)
	case oauth2.GrantRefreshToken:
		err = svr.grantByRefreshToken(w, r, argv, client, clientSecret)
	default:
		svr.errorResponse(w, r, oauth2.ErrorUnsupportedGrantType)
	}
	if err != nil {
		log.Warn("%s: grant error: %v", argv.CommandName(), err)
		svr.errorResponse(w, r, err)
	}
}

func (svr *Server) grantByAuthorizationCode(w http.ResponseWriter, r *http.Request, argv *api.TokenReq, client *model.Client) error {
	ar, err := svr.authRepo.GetAuthRequest(client.Id, argv.Code)
	if err != nil {
		return err
	}
	if ar == nil {
		return api.NewError(string(oauth2.ErrorInvalidGrant), "code-not-found")
	}
	user, err := svr.userRepo.GetUser(ar.Uid)
	if err != nil {
		return err
	}
	if user == nil {
		return api.NewError(string(oauth2.ErrorInvalidGrant), "user-not-found")
	}
	token, err := svr.tokenRepo.NewToken(user, client.Id, ar.GrantedScopes)
	if err != nil {
		return err
	}
	log.WithJSON(token).Debug("new token")
	svr.authRepo.RemoveAuthRequest(ar.AuthorizationCode)

	svr.response(w, r, api.TokenRes{
		TokenType: oauth2.TokenType,
		Token:     makeTokenInfo(token),
	})
	return nil
}

func (svr *Server) grantByPassword(w http.ResponseWriter, r *http.Request, argv *api.TokenReq, client *model.Client, clientSecret string) error {
	return api.NewError(string(oauth2.ErrorUnsupportedGrantType), "grantByPassword-not-implemented")
}

func (svr *Server) grantByRefreshToken(w http.ResponseWriter, r *http.Request, argv *api.TokenReq, client *model.Client, clientSecret string) error {
	if !model.ValidateClient(client, clientSecret) {
		return api.NewError(string(oauth2.ErrorInvalidGrant), "invalid-client-secret")
	}
	token, err := svr.tokenRepo.RefreshToken(argv.RefreshToken, argv.Scope)
	if err != nil {
		return err
	}

	svr.response(w, r, api.TokenRes{
		TokenType: oauth2.TokenType,
		Token:     makeTokenInfo(token),
	})
	return nil
}
