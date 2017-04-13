package server

import (
	"fmt"
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"github.com/mkideal/accountd/api"
	"github.com/mkideal/accountd/model"
)

func (svr *Server) handleUser(w http.ResponseWriter, r *http.Request) {
	ip := httputil.IP(r)
	argv := new(api.UserReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("User parse arguments error: %v, IP=%v", err, ip)
		svr.errorResponse(w, r, api.ErrorCode_BadArgument.NewError(err.Error()))
		return
	}
	log.WithJSON(argv).Debug("User request, IP=%v", ip)
	var user *model.User
	if argv.Uid > 0 {
		user, err = svr.userRepo.GetUser(argv.Uid)
	} else if argv.Account != "" {
		user, err = svr.userRepo.GetUserByAccount(argv.Account)
	} else {
		log.Info("%s: missing arguments uid and account", argv.CommandName())
		svr.errorResponse(w, r, api.ErrorCode_MissingArgument.NewError("missing arguments uid and account"))
		return
	}
	if err != nil {
		log.Error("%s: get user error: %v", argv.CommandName(), err)
		svr.errorResponse(w, r, err)
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
		svr.errorResponse(w, r, api.ErrorCode_UserNotFound.NewError(desc))
		return
	}
	svr.response(w, r, makeUserInfo(user))
}
