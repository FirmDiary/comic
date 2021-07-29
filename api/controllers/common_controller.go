package controllers

import (
	"comic/common"
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
	//[]{}
	etc := []map[string]string{
		{
			"origin": "https://comic-img.zwww.cool/images/etco1.png",
			"res":    "https://comic-img.zwww.cool/images/etcr1.png",
			"desc":   "《生活多美好》（1946年）曾经被美国电影学术权威的美国电影协会评价为“百年百部励志片榜首”，“影史最有力的影片TOP1”",
		},
		{
			"origin": "https://comic-img.zwww.cool/images/etco2.png",
			"res":    "https://comic-img.zwww.cool/images/etcr2.png",
			"desc":   "《城市之光》（1931年）查理卓别林的作品，本片位列影史百大影片第11位，是美国国家电影保护局指定典藏珍品",
		},
		{
			"origin": "https://comic-img.zwww.cool/images/etco3.png",
			"res":    "https://comic-img.zwww.cool/images/etcr3.png",
			"desc":   "《鬼子来了》（2000年）姜文导演指导的影片，获得戛纳国际电影节评委会大奖",
		},
		{
			"origin": "https://comic-img.zwww.cool/images/etco4.png",
			"res":    "https://comic-img.zwww.cool/images/etcr4.png",
			"desc":   "《罗马假日》（1953年）奥黛丽赫本凭借此片获得奥斯卡最佳女主角。",
		},
		{
			"origin": "https://comic-img.zwww.cool/images/etco5.png",
			"res":    "https://comic-img.zwww.cool/images/etcr5.png",
			"desc":   "《辛德勒的名单》（1993年）该片改编自澳大利亚小说家托马斯·肯尼利的同名小说，讲述了一名身在波兰的的德国人辛德勒，在二战时雇用了1100多名犹太人在他的工厂工作，帮助他们逃过被屠杀的劫数。获得了第66届奥斯卡金像奖最佳影片等7个奖项",
		},
		{
			"origin": "https://comic-img.zwww.cool/images/etco6.png",
			"res":    "https://comic-img.zwww.cool/images/etcr6.png",
			"desc":   "《影》凭借着布景、服装等一系列手段，将拍摄物全部变为黑白，并附上了中国特有的特色，那就是水墨画。得益于整体统一的水墨画美学体系，影片画面干净写意，拍摄起动作场面来，飘逸、灵动、清隽。",
		},
	}
	return common.ReSuccessData(etc)
}
