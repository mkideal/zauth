package repo

import (
	"fmt"
	"time"

	"github.com/mkideal/pkg/math/random"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/model"
)

type tokenRepository struct {
	*SqlRepository
}

func NewTokenRepository(sqlRepo *SqlRepository) TokenRepository {
	return tokenRepository{SqlRepository: sqlRepo}
}

func (repo tokenRepository) generateToken() string {
	return random.String(64, random.CryptoSource)
}

const (
	tokenLifeTime        = time.Hour * 48      // 2 days
	refreshTokenLifeTime = time.Hour * 24 * 30 // 30 days
)

func (repo tokenRepository) NewToken(user *model.User, clientId, scope string) (*model.Token, error) {
	now := time.Now()
	token := &model.Token{
		Uid:                  user.Id,
		CreatedAt:            model.FormatTime(now),
		AccessTokenExpireAt:  model.FormatTime(now.Add(tokenLifeTime)),
		RefreshTokenExpireAt: model.FormatTime(now.Add(refreshTokenLifeTime)),
		AccessToken:          repo.generateToken(),
		RefreshToken:         repo.generateToken(),
		ResourceOwner:        fmt.Sprintf("%d", user.Id),
		Scope:                scope,
		ClientId:             clientId,
	}
	err := repo.Insert(token)
	if err != nil {
		token = nil
	}
	return token, err
}

func (repo tokenRepository) GetToken(accessToken string) (*model.Token, error) {
	token := &model.Token{AccessToken: accessToken}
	found, err := repo.Get(token)
	if !found || err != nil {
		return nil, err
	}
	if model.IsExpired(token.AccessTokenExpireAt) {
		return nil, api.ErrorCode_TokenExpired
	}
	return token, err
}

func (repo tokenRepository) RefreshToken(refreshToken, scope string) (*model.Token, error) {
	token := &model.Token{RefreshToken: refreshToken}
	meta := model.TokenMetaVar
	// get token by refreshToken
	found, err := repo.Get(token, meta.F_refresh_token)
	if !found || err != nil {
		return nil, err
	}
	// remove got token
	if err := repo.Remove(token); err != nil {
		return nil, err
	}
	if model.IsExpired(token.RefreshTokenExpireAt) {
		return nil, api.ErrorCode_TokenExpired
	}
	// generate new token
	now := time.Now()
	token.Scope = scope
	token.AccessToken = repo.generateToken()
	token.RefreshToken = repo.generateToken()
	token.CreatedAt = model.FormatTime(now)
	token.AccessTokenExpireAt = model.FormatTime(now.Add(tokenLifeTime))
	token.RefreshTokenExpireAt = model.FormatTime(now.Add(refreshTokenLifeTime))
	err = repo.Insert(token)
	return token, err
}
