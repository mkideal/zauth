package qq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"bitbucket.org/mkideal/accountd/model"
	"bitbucket.org/mkideal/accountd/third_party"
)

const (
	accessTokenURL = "https://graph.qq.com/oauth2.0/token"
	openIdURL      = "https://graph.qq.com/oauth2.0/me"
	userinfoURL    = "https://graph.qq.com/user/get_user_info"
)

const (
	errorCodeOk = 0
)

func init() {
	third_party.Register(third_party.QQ, New)
}

type qq struct {
}

func New() third_party.Client {
	return &qq{}
}

func (c *qq) Name() string { return third_party.QQ }

type response interface {
	ErrorCode() int
	ErrorMsg() string
}

type errorResponse struct {
	Errcode int    `json:"ret"`
	Errmsg  string `json:"msg"`
}

func (er errorResponse) ErrorCode() int   { return er.Errcode }
func (er errorResponse) ErrorMsg() string { return er.Errmsg }

func (c *qq) request(url string, respObj response) error {
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

func (c *qq) GetAccessToken(clientId, clientSecret, code string) (*third_party.AccessTokenResponse, error) {
	return nil, third_party.Error{
		Name:        c.Name(),
		Code:        third_party.UnsupportedAPI,
		Description: "GetAccessToken not implemented",
	}
}

type openIdResponse struct {
	errorResponse
	ClientId string `json:"client_id"`
	OpenId   string `json:"openid"`
}

type userInfoResponse struct {
	errorResponse
	Nickname       string `json:"nickname"`
	FigureURL      string `json:"figureurl"`
	FigureURL_1    string `json:"figureurl_1"`
	FigureURL_2    string `json:"figureurl_2"`
	FigureURL_qq_1 string `json:"figureurl_qq_1"`
	FigureURL_qq_2 string `json:"figureurl_qq_2"`
	Gender         string `json:"gender"`
}

func (c *qq) GetUserInfo(accessToken, openId string) (*third_party.UserInfoResponse, error) {
	url := openIdURL + fmt.Sprintf("?access_token=%s", accessToken)
	openIdResp := openIdResponse{}
	if err := c.request(url, &openIdResp); err != nil {
		return nil, err
	}
	url = userinfoURL + fmt.Sprintf("?access_token=%s&oauth_consumer_key=%s&openid=%s", accessToken, openIdResp.ClientId, openIdResp.OpenId)
	respObj := userInfoResponse{}
	if err := c.request(url, &respObj); err != nil {
		return nil, err
	}
	var sex model.Gender
	if respObj.Gender == "男" {
		sex = model.Gender_Male
	} else {
		sex = model.Gender_Female
	}
	var avatar string
	switch {
	case respObj.FigureURL_qq_2 != "":
		avatar = respObj.FigureURL_qq_2
	case respObj.FigureURL_qq_1 != "":
		avatar = respObj.FigureURL_qq_1
	case respObj.FigureURL_2 != "":
		avatar = respObj.FigureURL_2
	case respObj.FigureURL_1 != "":
		avatar = respObj.FigureURL_1
	case respObj.FigureURL != "":
		avatar = respObj.FigureURL
	}
	return &third_party.UserInfoResponse{
		OpenId:   openIdResp.OpenId,
		Nickname: respObj.Nickname,
		Sex:      sex,
		Avatar:   avatar,
	}, nil
}
