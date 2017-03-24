package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/model"
)

func (svr *Server) handleSMSCode(w http.ResponseWriter, r *http.Request) {
	ip := httputil.IP(r)
	argv := new(api.SMSCodeReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("SMSCode parse arguments error: %v, IP=%v", err, ip)
		svr.errorResponse(w, r, api.ErrorCode_BadArgument.NewError(err.Error()))
		return
	}
	log.WithJSON(argv).Debug("SMSCode request, IP=%v", ip)

	if !model.IsTelno(argv.Telno) {
		return
	}
	var (
		maxInterval = svr.config.TelnoVerifyCodeMaxInterval
		expiration  = svr.config.TelnoVerifyCodeExpiration
	)
	vcode, err := svr.telnoVerifyCodeRepo.NewTelnoCode(svr.config.TelnoVerifyCodeLength, argv.Telno, maxInterval, expiration)
	if err != nil {
		log.Error("%s: new telno verify code for %s error: %v", argv.CommandName(), argv.Telno, err)
		svr.errorResponse(w, r, err)
		return
	}
	if vcode == nil {
		log.Info("%s: new telno verify code for %s too often", argv.CommandName(), argv.Telno)
		svr.errorResponse(w, r, api.ErrorCode_TelnoVerifyCodeTooOften)
		return
	}
	sms := svr.config.SMS
	if err := svr.telnoVerifyCodeRepo.SendTelnoCode(vcode, sms.SMSURL, sms.SMSUsername, sms.SMSPassword, sms.SMSMsgFormat); err != nil {
		log.Warn("send SMS code to telno %s error: %v", argv.Telno, err)
		svr.errorResponse(w, r, api.ErrorCode_FailedToSendSMSCode.NewError(err.Error()))
		svr.telnoVerifyCodeRepo.RemoveTelnoCode(argv.Telno)
		return
	}
	svr.response(w, r, api.SMSCodeRes{})
}
