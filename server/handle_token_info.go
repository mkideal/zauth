package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
)

func (svr *Server) handleTokenInfo(w http.ResponseWriter, r *http.Request) {
	argv := new(api.TokenInfoReq)
	err := argv.Parse(r)
	if err != nil {
		log.Warn("TokenInfo parse arguments error: %v, IP=%v", err, httputil.IP(r))
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("TokenInfo request, IP=%v", httputil.IP(r))

	if token := svr.getTokenFromHeader(r); token != "" {
		argv.AccessToken = token
	}

	accessToken, err := svr.tokenRepo.FindToken(argv.AccessToken)
	if err != nil {
		svr.response(w, http.StatusInternalServerError, err)
	}
	if accessToken == nil {
		svr.responseErrorCode(w, api.ErrorCode_TokenNotFound, "token-not-found")
		return
	}
	svr.response(w, http.StatusOK, api.TokenInfoRes{
		Uid:      accessToken.Uid,
		Scope:    accessToken.Scope,
		ExpireAt: accessToken.ExpireAt,
	})
}
