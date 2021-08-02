package services

import (
	"comic/datamodels"
	"comic/repositories"
	"errors"
)

type IUserService interface {
	NewUser(user *datamodels.User) (id int64, err error)
	Get(user *datamodels.User) (has bool)
	DescQuotaByUserId(userId int64) (quota int64, err error)
	AddQuotaByUserId(userId int64, addQuota int64) (quota int64, err error)
}

type UserService struct {
	repository repositories.IUserRepository
}

func NewUserService() IUserService {
	return &UserService{repository: repositories.NewUserRepository()}
}

func (u *UserService) NewUser(user *datamodels.User) (id int64, err error) {
	return u.repository.Create(user)
}

func (u *UserService) Get(user *datamodels.User) (has bool) {
	return u.repository.Get(user)
}

func (u *UserService) DescQuotaByUserId(userId int64) (quota int64, err error) {
	user := &datamodels.User{Id: userId}
	has := u.Get(user)
	if !has {
		return quota, errors.New("用户不存在")
	}
	quota = user.Quota
	if quota < 1 {
		return quota, errors.New("额度不足")
	}
	quota, err = u.repository.DecrQuota(user)
	if err != nil {
		return
	}

	//记录
	quotaRepository := repositories.NewQuotaLogRepository()
	quotaRepository.Add(&datamodels.QuotaLog{
		UserId: userId,
		Type:   datamodels.TypeDecr,
		Source: datamodels.SourceUse,
		Quota:  1,
	})

	return
}

func (u *UserService) AddQuotaByUserId(userId int64, addQuota int64) (quota int64, err error) {
	user := &datamodels.User{Id: userId}
	has := u.Get(user)
	if !has {
		return quota, errors.New("用户不存在")
	}
	quota, err = u.repository.AddQuota(user, addQuota)
	if err != nil {
		return
	}

	//记录
	quotaRepository := repositories.NewQuotaLogRepository()
	quotaRepository.Add(&datamodels.QuotaLog{
		UserId: userId,
		Type:   datamodels.TypeAdd,
		Source: datamodels.SourceShare,
		Quota:  addQuota,
	})

	return
}
