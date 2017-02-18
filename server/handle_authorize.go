package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/model"
	"bitbucket.org/mkideal/accountd/oauth2"
)

func (svr *Server) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	ip := httputil.IP(r)
	argv := new(api.AuthorizeReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("Authorize parse arguments error: %v, IP=%v", err, ip)
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("Authorize request, IP=%v", ip)

	var (
		session *model.Session
		user    *model.User
	)
	session = svr.getSession(r)
	if session != nil {
		user, err = svr.userRepo.GetUser(session.Uid)
	}
	if user == nil || err != nil {
		r.ParseForm()
		params := url.Values{
			"return_to": {fmt.Sprintf("%s?%s", svr.config.Pages.Authorize, r.Form.Encode())},
		}.Encode()
		log.Debug("params: %s", params)
		uri := fmt.Sprintf("%s?%s", svr.config.Pages.Login, params)
		http.Redirect(w, r, uri, http.StatusFound)
		return
	}

	client, err := svr.clientRepo.GetClient(argv.ClientId)
	if err != nil {
		log.Error("%s: get client %s error: %v", argv.CommandName(), argv.ClientId, err)
		svr.errorResponse(argv.CommandName(), w, err)
		return
	}
	if client == nil {
		log.Info("%s: client %s not found", argv.CommandName(), argv.ClientId)
		svr.oauthErrorResponse(argv.CommandName(), w, oauth2.ErrorInvalidClient)
		return
	}

	values := map[string]interface{}{"state": argv.State}

	if argv.ResponseType != oauth2.ResponseCode {
		authErr := oauth2.NewError(oauth2.ErrorUnsupportedResponseType, "must-be-code")
		log.Info("%s: response_type must be `code`, but got `%s`", argv.CommandName(), argv.ResponseType)
		params := authErr.EncodeWith(values)
		uri := fmt.Sprintf("%s?%s", client.CallbackURL, params)
		http.Redirect(w, r, uri, http.StatusFound)
		return
	}
	uid := argv.Uid
	if session := svr.getSession(r); session != nil {
		uid = session.Uid
	}

	ar, err := svr.authRepo.NewAuthRequest(client, uid, argv.State, argv.Scope, oauth2.ResponseCode)
	if err != nil {
		log.Error("%s: NewAuthRequest error: %v", argv.CommandName(), err)
		authErr := oauth2.WrapError(err)
		params := authErr.EncodeWith(values)
		uri := fmt.Sprintf("%s?%s", client.CallbackURL, params)
		http.Redirect(w, r, uri, http.StatusFound)
		return
	}

	params := url.Values{
		"code":  {ar.AuthorizationCode},
		"state": {argv.State},
	}.Encode()
	uri := fmt.Sprintf("%s?%s", client.CallbackURL, params)
	http.Redirect(w, r, uri, http.StatusFound)
}
