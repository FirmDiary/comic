package controllers

import (
	"comic/common"
	"comic/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"net/http"
)

type CommonController struct {
	Ctx iris.Context
}

func (c *CommonController) BeforeActivation(b mvc.BeforeActivation) {
	b.HandleMany(http.MethodGet, "/old/etc", "OldEtc")
}

func (c *CommonController) OldEtc() common.Response {
	etc := []map[string]string{
		{
			"origin":    "https://comic-img.zwww.cool/images/etco1.png",
			"res":       "https://comic-img.zwww.cool/images/etcr1.png",
			"direction": services.DirectionColumn,
			"desc":      "《辛德勒的名单》（1993年）该片改编自澳大利亚小说家托马斯·肯尼利的同名小说，荣获第66届奥斯卡金像奖最佳影片等7个奖项",
		},
		{
			"origin":    "https://comic-img.zwww.cool/images/etco2.png",
			"res":       "https://comic-img.zwww.cool/images/etcr2.png",
			"direction": services.DirectionColumn,
			"desc":      "摄影：Rene Maltete",
		},
		{
			"origin":    "https://comic-img.zwww.cool/images/etco3.png",
			"res":       "https://comic-img.zwww.cool/images/etcr3.png",
			"direction": services.DirectionColumn,
			"desc":      "《鬼子来了》（2000年）姜文导演指导的影片，获得戛纳国际电影节评委会大奖",
		},
		{
			"origin":    "https://comic-img.zwww.cool/images/etco4.png",
			"res":       "https://comic-img.zwww.cool/images/etcr4.png",
			"direction": services.DirectionColumn,
			"desc":      "摄影：Elliott Erwitt",
		},
		{
			"origin":    "https://comic-img.zwww.cool/images/etco5.png",
			"res":       "https://comic-img.zwww.cool/images/etcr5.png",
			"direction": services.DirectionColumn,
			"desc":      "《城市之光》（1931年）查理卓别林的作品，本片位列影史百大影片第11位，是美国国家电影保护局指定典藏珍品",
		},
		{
			"origin":    "https://comic-img.zwww.cool/images/etco6.png",
			"res":       "https://comic-img.zwww.cool/images/etcr6.png",
			"direction": services.DirectionColumn,
			"desc":      "《影》凭借着布景、服装等一系列手段，将拍摄物全部变为黑白，并附上了中国特有的特色，那就是水墨画。",
		},
	}
	return common.ReSuccessData(etc)
}
