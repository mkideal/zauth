package model

// 访问客户端
type Client struct {
	Id          string
	Secret      string
	Name        string
	Scopes      string
	Description string
	CallbackURL string
}

type ClientMeta struct {
	F_id          string
	F_secret      string
	F_name        string
	F_scopes      string
	F_description string
	F_callbackURL string
}

func (ClientMeta) Name() string {
	return "Client"
}

func (ClientMeta) NumField() int {
	return 6
}

func (ClientMeta) Field(i int, v Client) (string, interface{}) {
	switch i {

	case 0:
		return "id", v.Id
	case 1:
		return "secret", v.Secret
	case 2:
		return "name", v.Name
	case 3:
		return "scopes", v.Scopes
	case 4:
		return "description", v.Description
	case 5:
		return "callbackURL", v.CallbackURL

	}
	return "", nil
}

func (ClientMeta) FieldPtr(i int, v *Client) (string, interface{}) {
	switch i {

	case 0:
		return "id", &v.Id
	case 1:
		return "secret", &v.Secret
	case 2:
		return "name", &v.Name
	case 3:
		return "scopes", &v.Scopes
	case 4:
		return "description", &v.Description
	case 5:
		return "callbackURL", &v.CallbackURL

	}
	return "", nil
}

var ClientMetaVar = ClientMeta{

	F_id:          "id",
	F_secret:      "secret",
	F_name:        "name",
	F_scopes:      "scopes",
	F_description: "description",
	F_callbackURL: "callbackURL",
}
