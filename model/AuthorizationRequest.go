package model

// 认证请求
type AuthorizationRequest struct {
	Id                int64  // 递增唯一Id
	CreatedAt         string // 创建时间
	AuthorizationCode string // 认证码
	Uid               int64  // 用户Id
	RedirectURI       string // 重定向URI
	ResponseType      string // 返回类型
	State             string // 自定义状态
	ClientId          string // 客户端Id
	GrantedScopes     string // 授权范围
	RequestedScopes   string // 请求范围

}

type AuthorizationRequestMeta struct {
	F_id                string
	F_createdAt         string
	F_authorizationCode string
	F_uid               string
	F_redirectURI       string
	F_responseType      string
	F_state             string
	F_clientId          string
	F_grantedScopes     string
	F_requestedScopes   string
}

func (AuthorizationRequestMeta) Name() string {
	return "AuthorizationRequest"
}

func (AuthorizationRequestMeta) NumField() int {
	return 10
}

func (AuthorizationRequestMeta) Field(i int, v AuthorizationRequest) (string, interface{}) {
	switch i {

	case 0:
		return "id", v.Id
	case 1:
		return "createdAt", v.CreatedAt
	case 2:
		return "authorizationCode", v.AuthorizationCode
	case 3:
		return "uid", v.Uid
	case 4:
		return "redirectURI", v.RedirectURI
	case 5:
		return "responseType", v.ResponseType
	case 6:
		return "state", v.State
	case 7:
		return "clientId", v.ClientId
	case 8:
		return "grantedScopes", v.GrantedScopes
	case 9:
		return "requestedScopes", v.RequestedScopes

	}
	return "", nil
}

func (AuthorizationRequestMeta) FieldPtr(i int, v *AuthorizationRequest) (string, interface{}) {
	switch i {

	case 0:
		return "id", &v.Id
	case 1:
		return "createdAt", &v.CreatedAt
	case 2:
		return "authorizationCode", &v.AuthorizationCode
	case 3:
		return "uid", &v.Uid
	case 4:
		return "redirectURI", &v.RedirectURI
	case 5:
		return "responseType", &v.ResponseType
	case 6:
		return "state", &v.State
	case 7:
		return "clientId", &v.ClientId
	case 8:
		return "grantedScopes", &v.GrantedScopes
	case 9:
		return "requestedScopes", &v.RequestedScopes

	}
	return "", nil
}

var AuthorizationRequestMetaVar = AuthorizationRequestMeta{

	F_id:                "id",
	F_createdAt:         "createdAt",
	F_authorizationCode: "authorizationCode",
	F_uid:               "uid",
	F_redirectURI:       "redirectURI",
	F_responseType:      "responseType",
	F_state:             "state",
	F_clientId:          "clientId",
	F_grantedScopes:     "grantedScopes",
	F_requestedScopes:   "requestedScopes",
}
