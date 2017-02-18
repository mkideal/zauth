package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/model"
)

func (svr *Server) handleAccountExist(w http.ResponseWriter, r *http.Request) {
	ip := httputil.IP(r)
	argv := new(api.AccountExistReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("AccountExist parse arguments error: %v, IP=%v", err, ip)
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("AccountExist request, IP=%v", ip)
	if !model.IsNormalUsername(argv.Username) {
		log.Info("%s: illegal username: `%s`", argv.CommandName(), argv.Username)
		svr.responseErrorCode(w, api.ErrorCode_IllegalUsername, "illegal-username-format")
		return
	}
	found, err := svr.userRepo.AccountExist(argv.Username)
	if err != nil {
		log.Error("%s: find account %s error: %v", argv.CommandName(), argv.Username, err)
		svr.response(w, http.StatusInternalServerError, err)
		return
	}
	svr.response(w, http.StatusOK, api.AccountExistRes{Existed: found})
}
