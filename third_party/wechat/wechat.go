package wechat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mkideal/accountd/model"
	"github.com/mkideal/accountd/third_party"
)

const (
	accessTokenURL = "https://api.weixin.qq.com/sns/oauth2/access_token"
	userinfoURL    = "https://api.weixin.qq.com/sns/userinfo"
)

const (
	errorCodeOk = 0
)

func init() {
	third_party.Register(third_party.WeChat, New)
}

type weChat struct {
}

func New() third_party.Client {
	return &weChat{}
}

func (c *weChat) Name() string { return third_party.WeChat }

type response interface {
	ErrorCode() int
	ErrorMsg() string
}

type errorResponse struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (er errorResponse) ErrorCode() int   { return er.Errcode }
func (er errorResponse) ErrorMsg() string { return er.Errmsg }

type accessTokenResponse struct {
	errorResponse

	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionId      string `json:"unionid"`
}

func (c *weChat) request(url string, respObj response) error {
	resp, err := http.Get(url)
	if err != nil {
		return third_party.Error{
			Name:        c.Name(),
			Code:        third_party.NetworkError,
			Description: err.Error(),
		}
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(respObj)
	if err != nil {
		return third_party.Error{
			Name:        c.Name(),
			Code:        third_party.ResponseFormatError,
			Description: err.Error(),
		}
	}
	if respObj.ErrorCode() != errorCodeOk {
		return third_party.Error{
			Name:        c.Name(),
			Code:        strconv.Itoa(respObj.ErrorCode()),
			Description: respObj.ErrorMsg(),
		}
	}
	return nil
}

func (c *weChat) GetAccessToken(clientId, clientSecret, code string) (*third_party.AccessTokenResponse, error) {
	url := accessTokenURL + fmt.Sprintf("?appid=%s&secret=%s&code=%s&grant_type=authorization_code", clientId, clientSecret, code)
	obj := accessTokenResponse{}
	if err := c.request(url, &obj); err != nil {
		return nil, err
	}
	return &third_party.AccessTokenResponse{
		AccessToken:  obj.AccessToken,
		ExpiresIn:    obj.ExpiresIn,
		RefreshToken: obj.RefreshToken,
		Scope:        obj.Scope,
		OpenId:       obj.OpenId,
		Extra: map[string]string{
			"unionid": obj.UnionId,
		},
	}, nil
}

type userInfoResponse struct {
	errorResponse

	OpenId     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionId    string   `json:"unionid"`
}

func (c *weChat) GetUserInfo(accessToken, openId string) (*third_party.UserInfoResponse, error) {
	url := userinfoURL + fmt.Sprintf("?access_token=%s&openid=%s", accessToken, openId)
	obj := userInfoResponse{}
	if err := c.request(url, &obj); err != nil {
		return nil, err
	}
	var sex model.Gender
	switch obj.Sex {
	case 1:
		sex = model.Gender_Male
	case 2:
		sex = model.Gender_Female
	default:
		sex = model.Gender_Secret
	}
	return &third_party.UserInfoResponse{
		OpenId:   obj.OpenId,
		Nickname: obj.Nickname,
		Avatar:   obj.HeadImgURL,
		Sex:      sex,
		Country:  obj.Country,
		Province: obj.Province,
		City:     obj.City,
		Extra: map[string]string{
			"unionid": obj.UnionId,
		},
	}, nil
}
