package server

import (
	"fmt"
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/model"
)

func (svr *Server) handleUser(w http.ResponseWriter, r *http.Request) {
	argv := new(api.UserReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("User parse arguments error: %v, IP=%v", err, httputil.IP(r))
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("User request, IP=%v", httputil.IP(r))
	var user *model.User
	if argv.Uid > 0 {
		user, err = svr.userRepo.GetUser(argv.Uid)
	} else if argv.Account != "" {
		user, err = svr.userRepo.GetUserByAccount(argv.Account)
	} else {
		log.Info("%s: missing arguments uid and account", argv.CommandName())
		svr.responseErrorCode(w, api.ErrorCode_MissingArgument, "missing arguments uid and account")
		return
	}
	if err != nil {
		log.Error("%s: get user error: %v", argv.CommandName(), err)
		svr.response(w, http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		var desc string
		if argv.Uid > 0 {
			desc = fmt.Sprintf("uid %d not found", argv.Uid)
		} else {
			desc = fmt.Sprintf("account %s not found", argv.Account)
		}
		log.Info("%s: %s", argv.CommandName(), desc)
		svr.responseErrorCode(w, api.ErrorCode_UserNotFound, desc)
		return
	}
	svr.response(w, http.StatusOK, makeUserInfo(user))
}
