package wechat

import (
	"comic/common"
	"comic/repositories"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
)

type Mini struct {
	wc  *wechat.Wechat
	app *miniprogram.MiniProgram
}

func NewMini(appId int64) *Mini {
	appRepository := repositories.NewAppRepository()
	has, app := appRepository.Get(appId)
	if !has {
		panic("app不存在" + string(appId))
	}
	offCfg := &miniConfig.Config{
		AppID:     app.AppId,
		AppSecret: app.AppSecret,
		Cache:     common.NewRedis(),
	}

	wc := wechat.NewWechat()
	return &Mini{
		wc:  wc,
		app: wc.GetMiniProgram(offCfg),
	}
}
