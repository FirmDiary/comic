package repositories

import (
	"comic/common"
	"comic/datamodels"
	"github.com/go-xorm/xorm"
)

type IUploadRepository interface {
	Create(upload *datamodels.Upload) (int64, error)
}

type UploadRepository struct {
	db *xorm.Engine
}

func NewUploadRepository() IUploadRepository {
	return &UploadRepository{common.NewDbEngine()}
}

func (u UploadRepository) Create(upload *datamodels.Upload) (int64, error) {
	_, err := u.db.InsertOne(upload)
	if err != nil {
		return 0, err
	}
	return upload.UserId, nil
}
