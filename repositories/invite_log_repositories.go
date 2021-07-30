package repositories

import (
	"comic/common"
	"comic/datamodels"
	"github.com/go-xorm/xorm"
)

type IInviteLogRepository interface {
	Add(log *datamodels.InviteLog) (err error)
	Get(log *datamodels.InviteLog) (has bool)
}

type InviteLogRepository struct {
	db *xorm.Engine
}

func NewInviteLogRepository() IInviteLogRepository {
	return &InviteLogRepository{common.NewDbEngine()}
}

func (u *InviteLogRepository) Add(log *datamodels.InviteLog) (err error) {
	_, err = u.db.InsertOne(log)
	return
}

func (u *InviteLogRepository) Get(log *datamodels.InviteLog) (has bool) {
	has, _ = u.db.Get(log)
	return
}
