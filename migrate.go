package main

import (
	"comic/common"
	"comic/datamodels"
	"log"
)
//执行数据库迁移
func main() {
	db := common.NewDbEngine()
	db.SetMaxOpenConns(2)

	//https://www.kancloud.cn/kancloud/xorm-manual-zh-cn/56006
	db.Sync2(
		new(datamodels.Upload),
		new(datamodels.User),
	)
	log.Println("init database success")
}
