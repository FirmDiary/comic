package repositories

import (
	"comic/common"
	"comic/datamodels"
	"github.com/go-xorm/xorm"
)

type IUserRepository interface {
	Create(user *datamodels.User) (int64, error)
	Get(user *datamodels.User) (has bool)
	DecrQuota(user *datamodels.User) (int64, error)
	AddQuota(user *datamodels.User, addQuota int64) (quota int64, err error)
}

type UserRepository struct {
	db *xorm.Engine
}

func NewUserRepository() IUserRepository {
	return &UserRepository{common.NewDbEngine()}
}

func (u *UserRepository) DecrQuota(user *datamodels.User) (quota int64, err error) {
	quota = user.Quota
	user.Quota = quota - 1
	u.db.Id(user.Id).Update(user)

	return
}

func (u *UserRepository) AddQuota(user *datamodels.User, addQuota int64) (quota int64, err error) {
	quota = user.Quota
	user.Quota = quota + addQuota
	u.db.Id(user.Id).Update(user)

	return
}

func (u *UserRepository) Create(user *datamodels.User) (id int64, err error) {
	return u.db.InsertOne(user)
}

func (u *UserRepository) Get(user *datamodels.User) (has bool) {
	has, _ = u.db.Get(user)
	return
}
