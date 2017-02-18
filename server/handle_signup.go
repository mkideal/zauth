package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/model"
)

func (svr *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
	ip := httputil.IP(r)
	argv := new(api.SignupReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("Signup parse arguments error: %v, IP=%v", err, ip)
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("Signup request, IP=%v", ip)
	// 检查 accountType,account,password
	var (
		isNormalUsername    = model.IsNormalUsername(argv.Account)
		isEmail             = model.IsEmail(argv.Account)
		isTelno             = model.IsTelno(argv.Account)
		isCustomAccountType = isNormalUsername || isEmail || isTelno
	)
	if !isNormalUsername && argv.AccountType == int(model.AccountType_Normal) {
		log.Info("%s: invalid username `%s`", argv.CommandName(), argv.Account)
		svr.responseErrorCode(w, api.ErrorCode_IllegalUsername, "illegal-username")
		return
	}
	if !isEmail && argv.AccountType == int(model.AccountType_Email) {
		log.Info("%s: invalid email `%s`", argv.CommandName(), argv.Account)
		svr.responseErrorCode(w, api.ErrorCode_IllegalEmail, "illegal-email")
		return
	}
	if !isTelno && argv.AccountType == int(model.AccountType_Telno) {
		log.Info("%s: invalid telno `%s`", argv.CommandName(), argv.Account)
		svr.responseErrorCode(w, api.ErrorCode_IllegalTelno, "illegal-telno")
		return
	}
	if isCustomAccountType {
		if !model.IsPassword(argv.Password) {
			log.Info("%s: invalid password", argv.CommandName())
			svr.responseErrorCode(w, api.ErrorCode_IllegalPassword, "illegal-password")
			return
		}
	} else {
		// FIXME: 暂时不支持其它账号类型注册
		log.Info("%s: unsupported accountType %d", argv.CommandName(), argv.AccountType)
		svr.responseErrorCode(w, api.ErrorCode_IllegalAccountType, "illegal-account-type")
		return
	}

	if found, _ := svr.userRepo.AccountExist(argv.Account); found {
		log.Info("%s: account %s duplicated", argv.CommandName(), argv.Account)
		svr.responseErrorCode(w, api.ErrorCode_AccountDuplicated, "account-duplicated")
		return
	}

	user := new(model.User)
	user.AccountType = model.AccountType(argv.AccountType)
	user.CreatedIp = httputil.IP(r)
	user.Account = argv.Account
	user.Nickname = argv.Nickname
	if err := svr.userRepo.AddUser(user, argv.Password); err != nil {
		if found, _ := svr.userRepo.AccountExist(argv.Account); found {
			log.Info("%s: account %s duplicated", argv.CommandName(), argv.Account)
			return
		}
		log.Error("%s: add user %s error: %v", argv.CommandName(), argv.Account, err)
		svr.response(w, http.StatusInternalServerError, err)
		return
	}
	res := api.SignupRes{
		Uid:      user.Id,
		Account:  user.Account,
		Nickname: user.Nickname,
	}
	svr.response(w, http.StatusOK, res)
}
