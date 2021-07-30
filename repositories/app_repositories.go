package repositories

import (
	"comic/common"
	"comic/datamodels"
	"github.com/go-xorm/xorm"
)

type IAppRepository interface {
	Get(id int64) (has bool, app *datamodels.App)
}

type AppRepository struct {
	db *xorm.Engine
}

func NewAppRepository() IAppRepository {
	return &AppRepository{common.NewDbEngine()}
}

func (u *AppRepository) Get(id int64) (has bool, app *datamodels.App) {
	app = &datamodels.App{Id: id}
	has, _ = u.db.Get(app)
	return
}
