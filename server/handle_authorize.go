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
	argv := new(api.AuthorizeReq)
	err := argv.Parse(r)
	if err != nil {
		log.Warn("Authorize parse arguments error: %v, IP=%v", err, httputil.IP(r))
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("Authorize request, IP=%v", httputil.IP(r))

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
		svr.errorResponse(argv.CommandName(), w, err)
		return
	}
	if client == nil {
		svr.oauthErrorResponse(argv.CommandName(), w, oauth2.ErrorInvalidClient)
		return
	}

	values := map[string]interface{}{"state": argv.State}

	if argv.ResponseType != oauth2.ResponseCode {
		authErr := oauth2.NewError(oauth2.ErrorUnsupportedResponseType, "must-be-code")
		params := authErr.EncodeWith(values)
		uri := fmt.Sprintf("%s?%s", client.CallbackURL, params)
		http.Redirect(w, r, uri, http.StatusFound)
		return
	}

	ar, err := svr.authRepo.NewAuthRequest(client, argv.Uid, argv.State, argv.Scope, oauth2.ResponseCode)
	if err != nil {
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
