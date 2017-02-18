package repo

import (
	"fmt"
	"time"

	"github.com/mkideal/pkg/math/random"

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

const tokenLifeTime = time.Hour * 48

func (repo tokenRepository) NewToken(user *model.User, clientId, scope string) (*model.Token, error) {
	now := time.Now()
	token := &model.Token{
		Uid:           user.Id,
		CreatedAt:     model.FormatTime(now),
		ModifiedAt:    model.FormatTime(now),
		ExpireAt:      model.FormatTime(now.Add(tokenLifeTime)),
		AccessToken:   repo.generateToken(),
		RefreshToken:  repo.generateToken(),
		ResourceOwner: fmt.Sprintf("%d", user.Id),
		Scope:         scope,
		ClientId:      clientId,
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
		token = nil
	}
	return token, err
}

func (repo tokenRepository) RefreshToken(refreshToken, scope string) (*model.Token, error) {
	token := &model.Token{RefreshToken: refreshToken}
	meta := model.TokenMetaVar
	found, err := repo.Get(token, meta.F_refresh_token)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	token.Scope = scope
	token.AccessToken = repo.generateToken()
	token.RefreshToken = repo.generateToken()
	now := time.Now()
	token.CreatedAt = model.FormatTime(now)
	token.ModifiedAt = model.FormatTime(now)
	token.ExpireAt = model.FormatTime(now.Add(tokenLifeTime))
	return token, err
}
