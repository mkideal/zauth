package server

import (
	"bitbucket.org/mkideal/accountd/model"
)

type UserRepository interface {
	Add(user *model.User) error
	Update(user *model.User) (int, error)
	Remove(uid int64) (int, error)
	Find(uid int64) (*model.User, error)
	FindByAccount(account string) (*model.User, error)
}
