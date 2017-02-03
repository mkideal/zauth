package model

// 手机验证码
type TelnoVerifyCode struct {
	Id        int64  // 递增唯一Id
	Telno     string // 手机号码
	Code      string // 验证码
	ExpiredAt string // 到期时间

}

type TelnoVerifyCodeMeta struct {
	F_id        string
	F_telno     string
	F_code      string
	F_expiredAt string
}

func (TelnoVerifyCodeMeta) Name() string {
	return "TelnoVerifyCode"
}

func (TelnoVerifyCodeMeta) NumField() int {
	return 4
}

func (TelnoVerifyCodeMeta) Field(i int, v TelnoVerifyCode) (string, interface{}) {
	switch i {

	case 0:
		return "id", v.Id
	case 1:
		return "telno", v.Telno
	case 2:
		return "code", v.Code
	case 3:
		return "expiredAt", v.ExpiredAt

	}
	return "", nil
}

func (TelnoVerifyCodeMeta) FieldPtr(i int, v *TelnoVerifyCode) (string, interface{}) {
	switch i {

	case 0:
		return "id", &v.Id
	case 1:
		return "telno", &v.Telno
	case 2:
		return "code", &v.Code
	case 3:
		return "expiredAt", &v.ExpiredAt

	}
	return "", nil
}

var TelnoVerifyCodeMetaVar = TelnoVerifyCodeMeta{

	F_id:        "id",
	F_telno:     "telno",
	F_code:      "code",
	F_expiredAt: "expiredAt",
}
