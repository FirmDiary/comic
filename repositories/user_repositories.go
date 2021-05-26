package repositories

import (
	"comic/common"
	"comic/datamodels"
	"github.com/go-xorm/xorm"
)

type IUserRepository interface {
	Create(user *datamodels.User) (int64, error)
	Get(user *datamodels.User) (has bool)
}

type UserRepository struct {
	db *xorm.Engine
}

func NewUserRepository() IUserRepository {
	return &UserRepository{common.NewDbEngine()}
}

func (u *UserRepository) Create(user *datamodels.User) (id int64, err error) {
	return u.db.InsertOne(user)
}

func (u *UserRepository) Get(user *datamodels.User) (has bool) {
	has, _ = u.db.Get(user)
	return
}
