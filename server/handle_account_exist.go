package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/model"
)

func (svr *Server) handleAccountExist(w http.ResponseWriter, r *http.Request) {
	argv := new(api.AccountExistReq)
	err := argv.Parse(r)
	if err != nil {
		log.Warn("AccountExist parse arguments error: %v, IP=%v", err, httputil.IP(r))
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("AccountExist request, IP=%v", httputil.IP(r))
	if !model.IsNormalUsername(argv.Username) {
		svr.responseErrorCode(w, api.ErrorCode_IllegalUsername, "illegal-username-format")
		return
	}
	found, err := svr.userRepo.AccountExist(argv.Username)
	if err != nil {
		svr.response(w, http.StatusInternalServerError, err)
		return
	}
	svr.response(w, http.StatusOK, api.AccountExistRes{Existed: found})
}
