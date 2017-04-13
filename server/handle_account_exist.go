package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"github.com/mkideal/accountd/api"
	"github.com/mkideal/accountd/model"
)

func (svr *Server) handleAccountExist(w http.ResponseWriter, r *http.Request) {
	ip := httputil.IP(r)
	argv := new(api.AccountExistReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("AccountExist parse arguments error: %v, IP=%v", err, ip)
		svr.errorResponse(w, r, api.ErrorCode_BadArgument.NewError(err.Error()))
		return
	}
	log.WithJSON(argv).Debug("AccountExist request, IP=%v", ip)
	if !model.IsNormalUsername(argv.Username) {
		log.Info("%s: illegal username: `%s`", argv.CommandName(), argv.Username)
		svr.errorResponse(w, r, api.ErrorCode_IllegalUsername)
		return
	}
	found, err := svr.userRepo.AccountExist(argv.Username)
	if err != nil {
		log.Error("%s: find account %s error: %v", argv.CommandName(), argv.Username, err)
		svr.errorResponse(w, r, err)
		return
	}
	svr.response(w, r, api.AccountExistRes{Existed: found})
}
