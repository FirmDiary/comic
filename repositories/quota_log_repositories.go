package repositories

import (
	"comic/common"
	"comic/datamodels"
	"github.com/go-xorm/xorm"
)

type IQuotaLogRepository interface {
	Add(log *datamodels.QuotaLog) (err error)
}

type QuotaLogRepository struct {
	db *xorm.Engine
}

func NewQuotaLogRepository() IQuotaLogRepository {
	return &QuotaLogRepository{common.NewDbEngine()}
}

func (u *QuotaLogRepository) Add(log *datamodels.QuotaLog) (err error) {
	_, err = u.db.InsertOne(log)
	return
}
