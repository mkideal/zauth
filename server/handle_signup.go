package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"github.com/mkideal/accountd/api"
	"github.com/mkideal/accountd/model"
	"github.com/mkideal/accountd/repo"
	"github.com/mkideal/accountd/third_party"
)

func (svr *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
	ip := httputil.IP(r)
	argv := new(api.SignupReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("Signup parse arguments error: %v, IP=%v", err, ip)
		svr.errorResponse(w, r, api.ErrorCode_BadArgument.NewError(err.Error()))
		return
	}
	log.WithJSON(argv).Debug("Signup request, IP=%v", ip)
	// 检查 accountType,account,password
	var (
		isNormalUsername    = model.IsNormalUsername(argv.Account)
		isEmail             = model.IsEmail(argv.Account)
		isTelno             = model.IsTelno(argv.Account)
		isCustomAccountType = isNormalUsername || isEmail || isTelno
		isThirdParty        bool
		opts                []repo.UserAddOption
	)
	if !isNormalUsername && argv.AccountType == int(model.AccountType_Normal) {
		log.Info("%s: invalid username `%s`", argv.CommandName(), argv.Account)
		svr.errorResponse(w, r, api.ErrorCode_IllegalUsername)
		return
	}
	if !isEmail && argv.AccountType == int(model.AccountType_Email) {
		log.Info("%s: invalid email `%s`", argv.CommandName(), argv.Account)
		svr.errorResponse(w, r, api.ErrorCode_IllegalEmail)
		return
	}
	if !isTelno && argv.AccountType == int(model.AccountType_Telno) {
		log.Info("%s: invalid telno `%s`", argv.CommandName(), argv.Account)
		svr.errorResponse(w, r, api.ErrorCode_IllegalTelno)
		return
	}
	if isCustomAccountType {
		if !model.IsPassword(argv.Password) {
			log.Info("%s: invalid password", argv.CommandName())
			svr.errorResponse(w, r, api.ErrorCode_IllegalPassword)
			return
		}
	} else {
		// 第三方账号注册
		name := third_party.GetNameByType(model.AccountType(argv.AccountType))
		third, ok := svr.third_parties[name]
		if !ok {
			log.Info("%s: unsupported accountType %d", argv.CommandName(), argv.AccountType)
			svr.errorResponse(w, r, api.ErrorCode_IllegalAccountType)
			return
		}
		var (
			accessToken = argv.ThirdAccessToken
			openId      = argv.ThirdOpenId
		)
		if accessToken == "" {
			resp, err := third.GetAccessToken(argv.ThirdClientId, argv.ThirdClientSecret, argv.Account)
			if err != nil {
				log.Info("%s: third_party error: %v", err)
				svr.errorResponse(w, r, api.ErrorCode_ThirdPartyError.NewError(err.Error()))
				return
			}
			accessToken = resp.AccessToken
			openId = resp.OpenId
		}
		resp2, err := third.GetUserInfo(accessToken, openId)
		if err != nil {
			log.Info("%s: third_party %s error: %v", argv.CommandMethod(), name, err)
			svr.errorResponse(w, r, api.ErrorCode_ThirdPartyError.NewError(err.Error()))
			return
		}
		log.WithJSON(resp2).Info("third_party %s UserInfoResponse", name)
		argv.Account = model.JoinAccount(model.AccountType(argv.AccountType), resp2.OpenId)
		argv.Nickname = resp2.Nickname
		opts = append(opts,
			repo.WithGender(resp2.Sex),
			repo.WithCountry(resp2.Country),
			repo.WithProvince(resp2.Province),
			repo.WithCity(resp2.City),
		)
		isThirdParty = true
	}

	if oldUser, _ := svr.userRepo.GetUserByAccount(argv.Account); oldUser != nil {
		if isThirdParty {
			// 如果第三方账号注册时，账号已经存在,则直接取得token
			token, err := svr.createToken(argv.CommandMethod(), oldUser, w, r)
			if err != nil {
				svr.errorResponse(w, r, err)
				return
			} else {
				svr.response(w, r, api.SignupRes{
					User:  makeUserInfo(oldUser),
					Token: makeTokenInfo(token),
				})
			}
		} else {
			log.Info("%s: account %s duplicated", argv.CommandName(), argv.Account)
			svr.errorResponse(w, r, api.ErrorCode_AccountDuplicated)
			return
		}
	}

	user := new(model.User)
	user.AccountType = model.AccountType(argv.AccountType)
	user.CreatedIp = httputil.IP(r)
	user.Account = argv.Account
	user.Nickname = argv.Nickname
	if err := svr.userRepo.AddUser(user, argv.Password, opts...); err != nil {
		if found, _ := svr.userRepo.AccountExist(argv.Account); found {
			log.Info("%s: account %s duplicated", argv.CommandName(), argv.Account)
			return
		}
		log.Error("%s: add user %s error: %v", argv.CommandName(), argv.Account, err)
		svr.errorResponse(w, r, err)
		return
	}
	token, err := svr.createToken(argv.CommandMethod(), user, w, r)
	if err != nil {
		svr.errorResponse(w, r, err)
		return
	}
	svr.response(w, r, api.SignupRes{
		User:  makeUserInfo(user),
		Token: makeTokenInfo(token),
	})
}
