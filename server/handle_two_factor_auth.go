package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/model"
)

func (svr *Server) handleTwoFactorAuth(w http.ResponseWriter, r *http.Request) {
	ip := httputil.IP(r)
	argv := new(api.TwoFactorAuthReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("TwoFactorAuth parse arguments error: %v, IP=%v", err, ip)
		svr.errorResponse(w, r, api.ErrorCode_BadArgument.NewError(err.Error()))
		return
	}
	log.WithJSON(argv).Debug("TwoFactorAuth request, IP=%v", ip)

	var user *model.User
	switch argv.AuthType {
	case "telno":
		user = svr.telno2faAuth(w, r, argv.AuthId, argv.AuthCode)
	case "email":
		user = svr.email2faAuth(w, r, argv.AuthId, argv.AuthCode)
	default:
		svr.errorResponse(w, r, api.ErrorCode_Unsupported2FaType.NewError(argv.AuthType))
	}
	if user != nil {
		_, err = svr.setSession(w, r, user.Id)
		if err != nil {
			log.Error("%s: set session error: %v", argv.CommandName(), err)
			svr.errorResponse(w, r, err)
			return
		}
		token, err := svr.tokenRepo.NewToken(user, "", "")
		if err != nil {
			log.Error("%s: new token error: %v", argv.CommandName(), err)
			svr.errorResponse(w, r, err)
			return
		}
		svr.response(w, r, api.TwoFactorAuthRes{
			User:  makeUserInfo(user),
			Token: makeTokenInfo(token),
		})
	}
}

func (svr *Server) telno2faAuth(w http.ResponseWriter, r *http.Request, telno, code string) *model.User {
	if !model.IsTelno(telno) {
		log.Info("illegal telno `%s`", telno)
		svr.errorResponse(w, r, api.ErrorCode_IllegalTelno)
		return nil
	}
	vcode, err := svr.telnoVerifyCodeRepo.FindTelnoCode(telno)
	if err != nil {
		log.Error("find telno verify code %s error: %v", telno, err)
		svr.errorResponse(w, r, err)
		return nil
	}
	if vcode == nil {
		log.Info("telno verify code for %s not found", telno)
		svr.errorResponse(w, r, api.ErrorCode_VerifyCodeNotFound)
		return nil
	}
	if model.IsExpired(vcode.ExpireAt) {
		log.Info("telno verify code for %s expired", telno)
		svr.errorResponse(w, r, api.ErrorCode_TelnoVerifyCodeExpired)
		return nil
	}
	user, err := svr.userRepo.GetUserByAccount(telno)
	if err != nil {
		log.Error("get user by account %s error: %v", telno, err)
		svr.errorResponse(w, r, err)
		return nil
	} else if user == nil {
		log.Info("telno %s not found", telno)
		svr.errorResponse(w, r, err)
		return nil
	}
	return user
}

func (svr *Server) email2faAuth(w http.ResponseWriter, r *http.Request, email, code string) *model.User {
	// TODO
	return nil
}
