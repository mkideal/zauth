package repo

import (
	"fmt"
	"time"

	"github.com/mkideal/pkg/math/random"

	"bitbucket.org/mkideal/accountd/model"
)

type tokenRepository struct {
	SqlRepository
}

func NewTokenRepository(sqlRepo SqlRepository) TokenRepository {
	return tokenRepository{SqlRepository: sqlRepo}
}

func (repo tokenRepository) generateToken() string {
	return random.String(64, random.CryptoSource)
}

const tokenLifeTime = time.Hour * 48

func (repo tokenRepository) NewToken(client *model.Client, user *model.User, scope string) (*model.AccessToken, error) {
	now := time.Now()
	token := &model.AccessToken{
		Uid:           user.Id,
		CreatedAt:     model.FormatTime(now),
		ModifiedAt:    model.FormatTime(now),
		ExpireAt:      model.FormatTime(now.Add(tokenLifeTime)),
		Token:         repo.generateToken(),
		RefreshToken:  repo.generateToken(),
		ResourceOwner: fmt.Sprintf("%d", user.Id),
		ClientId:      client.Id,
		Scope:         client.Scope,
	}
	err := repo.insert(token)
	if err != nil {
		token = nil
	}
	return token, err
}

func (repo tokenRepository) GetToken(token string) (*model.AccessToken, error) {
	accessToken := &model.AccessToken{Token: token}
	found, err := repo.getByFields(accessToken, model.AccessTokenMetaVar.F_token)
	if !found || err != nil {
		accessToken = nil
	}
	return accessToken, err
}

func (repo tokenRepository) RefreshToken(client *model.Client, refreshToken, scope string) (*model.AccessToken, error) {
	accessToken := &model.AccessToken{ClientId: client.Id, RefreshToken: refreshToken}
	meta := model.AccessTokenMetaVar
	found, err := repo.getByFields(accessToken, meta.F_client_id, meta.F_refresh_token)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	accessToken.Scope = scope
	accessToken.Token = repo.generateToken()
	accessToken.RefreshToken = repo.generateToken()
	now := time.Now()
	accessToken.CreatedAt = model.FormatTime(now)
	accessToken.ModifiedAt = model.FormatTime(now)
	accessToken.ExpireAt = model.FormatTime(now.Add(tokenLifeTime))
	return accessToken, err
}
