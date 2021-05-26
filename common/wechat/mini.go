package wechat

import (
	"comic/common"
	"comic/config"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
)

type Mini struct {
	wc  *wechat.Wechat
	app *miniprogram.MiniProgram
}

func NewMini() *Mini {
	globalCfg := config.GetConfig()
	offCfg := &miniConfig.Config{
		AppID:     globalCfg.Mini.AppID,
		AppSecret: globalCfg.Mini.AppSecret,
		Cache:     common.NewRedis(),
	}

	wc := wechat.NewWechat()
	return &Mini{
		wc:  wc,
		app: wc.GetMiniProgram(offCfg),
	}
}
