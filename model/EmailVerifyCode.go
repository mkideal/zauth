package model

// 邮箱验证码
type EmailVerifyCode struct {
	Id        int64  // 递增唯一Id
	Email     string // email 地址
	Code      string // 验证码
	ExpiredAt string // 到期时间

}

type EmailVerifyCodeMeta struct {
	F_id        string
	F_email     string
	F_code      string
	F_expiredAt string
}

func (EmailVerifyCodeMeta) Name() string {
	return "EmailVerifyCode"
}

func (EmailVerifyCodeMeta) NumField() int {
	return 4
}

func (EmailVerifyCodeMeta) Field(i int, v EmailVerifyCode) (string, interface{}) {
	switch i {

	case 0:
		return "id", v.Id
	case 1:
		return "email", v.Email
	case 2:
		return "code", v.Code
	case 3:
		return "expiredAt", v.ExpiredAt

	}
	return "", nil
}

func (EmailVerifyCodeMeta) FieldPtr(i int, v *EmailVerifyCode) (string, interface{}) {
	switch i {

	case 0:
		return "id", &v.Id
	case 1:
		return "email", &v.Email
	case 2:
		return "code", &v.Code
	case 3:
		return "expiredAt", &v.ExpiredAt

	}
	return "", nil
}

var EmailVerifyCodeMetaVar = EmailVerifyCodeMeta{

	F_id:        "id",
	F_email:     "email",
	F_code:      "code",
	F_expiredAt: "expiredAt",
}
