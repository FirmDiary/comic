package services

import (
    "comic/datamodels"
    "comic/repositories"
)

type IUserService interface {
    NewUser(user *datamodels.User) (id int64, err error)
    Get(user *datamodels.User) (has bool)
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
