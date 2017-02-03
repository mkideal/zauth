package model

// 用户信息
type User struct {
	Id                int64       // 随机唯一Id
	AccountType       AccountType // 账号类型
	Account           string      // 账号
	Nickname          string      // 昵称
	Avatar            string      // 头像
	QRCode            string      // 二维码
	Gender            Gender      // 性别
	Birthday          string      // 生日
	IdCardType        IdCardType  // 身份证件类型
	IdCard            string      // 证件唯一标识
	EncryptedPassword string      // 加密后密码
	PasswordSalt      string      // 加密密码的盐
	CreatedAt         string      // 账号创建时间
	CreatedIP         string      // 账号创建时IP
	LastLoginAt       string      // 最后登陆时间
	LastLoginIP       string      // 最后登陆时IP

}

type UserMeta struct {
	F_id                string
	F_accountType       string
	F_account           string
	F_nickname          string
	F_avatar            string
	F_QRCode            string
	F_gender            string
	F_birthday          string
	F_idCardType        string
	F_idCard            string
	F_encryptedPassword string
	F_passwordSalt      string
	F_createdAt         string
	F_createdIP         string
	F_lastLoginAt       string
	F_lastLoginIP       string
}

func (UserMeta) Name() string {
	return "User"
}

func (UserMeta) NumField() int {
	return 16
}

func (UserMeta) Field(i int, v User) (string, interface{}) {
	switch i {

	case 0:
		return "id", v.Id
	case 1:
		return "accountType", v.AccountType
	case 2:
		return "account", v.Account
	case 3:
		return "nickname", v.Nickname
	case 4:
		return "avatar", v.Avatar
	case 5:
		return "QRCode", v.QRCode
	case 6:
		return "gender", v.Gender
	case 7:
		return "birthday", v.Birthday
	case 8:
		return "idCardType", v.IdCardType
	case 9:
		return "idCard", v.IdCard
	case 10:
		return "encryptedPassword", v.EncryptedPassword
	case 11:
		return "passwordSalt", v.PasswordSalt
	case 12:
		return "createdAt", v.CreatedAt
	case 13:
		return "createdIP", v.CreatedIP
	case 14:
		return "lastLoginAt", v.LastLoginAt
	case 15:
		return "lastLoginIP", v.LastLoginIP

	}
	return "", nil
}

func (UserMeta) FieldPtr(i int, v *User) (string, interface{}) {
	switch i {

	case 0:
		return "id", &v.Id
	case 1:
		return "accountType", &v.AccountType
	case 2:
		return "account", &v.Account
	case 3:
		return "nickname", &v.Nickname
	case 4:
		return "avatar", &v.Avatar
	case 5:
		return "QRCode", &v.QRCode
	case 6:
		return "gender", &v.Gender
	case 7:
		return "birthday", &v.Birthday
	case 8:
		return "idCardType", &v.IdCardType
	case 9:
		return "idCard", &v.IdCard
	case 10:
		return "encryptedPassword", &v.EncryptedPassword
	case 11:
		return "passwordSalt", &v.PasswordSalt
	case 12:
		return "createdAt", &v.CreatedAt
	case 13:
		return "createdIP", &v.CreatedIP
	case 14:
		return "lastLoginAt", &v.LastLoginAt
	case 15:
		return "lastLoginIP", &v.LastLoginIP

	}
	return "", nil
}

var UserMetaVar = UserMeta{

	F_id:                "id",
	F_accountType:       "accountType",
	F_account:           "account",
	F_nickname:          "nickname",
	F_avatar:            "avatar",
	F_QRCode:            "QRCode",
	F_gender:            "gender",
	F_birthday:          "birthday",
	F_idCardType:        "idCardType",
	F_idCard:            "idCard",
	F_encryptedPassword: "encryptedPassword",
	F_passwordSalt:      "passwordSalt",
	F_createdAt:         "createdAt",
	F_createdIP:         "createdIP",
	F_lastLoginAt:       "lastLoginAt",
	F_lastLoginIP:       "lastLoginIP",
}
