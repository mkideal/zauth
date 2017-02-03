package model

// 访问令牌
type AccessToken struct {
	Id            int64  // 递增唯一Id
	Uid           int64  // 用户Id
	CreatedAt     string // 创建时间
	ModifiedAt    string // 修改时间
	ExpireAt      string // 到期时间
	Token         string // 令牌
	RefreshToken  string // 刷新用令牌
	ResourceOwner string // 资源所有者
	ClientId      string // 客户Id
	Scopes        string // 可访问权限范围

}

type AccessTokenMeta struct {
	F_id            string
	F_uid           string
	F_createdAt     string
	F_modifiedAt    string
	F_expireAt      string
	F_token         string
	F_refreshToken  string
	F_resourceOwner string
	F_clientId      string
	F_scopes        string
}

func (AccessTokenMeta) Name() string {
	return "AccessToken"
}

func (AccessTokenMeta) NumField() int {
	return 10
}

func (AccessTokenMeta) Field(i int, v AccessToken) (string, interface{}) {
	switch i {

	case 0:
		return "id", v.Id
	case 1:
		return "uid", v.Uid
	case 2:
		return "createdAt", v.CreatedAt
	case 3:
		return "modifiedAt", v.ModifiedAt
	case 4:
		return "expireAt", v.ExpireAt
	case 5:
		return "token", v.Token
	case 6:
		return "refreshToken", v.RefreshToken
	case 7:
		return "resourceOwner", v.ResourceOwner
	case 8:
		return "clientId", v.ClientId
	case 9:
		return "scopes", v.Scopes

	}
	return "", nil
}

func (AccessTokenMeta) FieldPtr(i int, v *AccessToken) (string, interface{}) {
	switch i {

	case 0:
		return "id", &v.Id
	case 1:
		return "uid", &v.Uid
	case 2:
		return "createdAt", &v.CreatedAt
	case 3:
		return "modifiedAt", &v.ModifiedAt
	case 4:
		return "expireAt", &v.ExpireAt
	case 5:
		return "token", &v.Token
	case 6:
		return "refreshToken", &v.RefreshToken
	case 7:
		return "resourceOwner", &v.ResourceOwner
	case 8:
		return "clientId", &v.ClientId
	case 9:
		return "scopes", &v.Scopes

	}
	return "", nil
}

var AccessTokenMetaVar = AccessTokenMeta{

	F_id:            "id",
	F_uid:           "uid",
	F_createdAt:     "createdAt",
	F_modifiedAt:    "modifiedAt",
	F_expireAt:      "expireAt",
	F_token:         "token",
	F_refreshToken:  "refreshToken",
	F_resourceOwner: "resourceOwner",
	F_clientId:      "clientId",
	F_scopes:        "scopes",
}
