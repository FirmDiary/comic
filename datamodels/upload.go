package datamodels

import (
	"time"
)

const (
	PlateDeepAi = 1
	PlateU2net  = 2
	PlateMTAi   = 3

	TypeFace = 1 //u2net 人脸
	TypeAll  = 2 //u2net 全部
)

type Upload struct {
	Id        int64     `json:"id" xorm:"INT(10) not null pk autoincr"`
	UserId    int64     `json:"user_id" xorm:"INT(10) not null index comment('用户ID')"`
	File      string    `json:"file" xorm:"VARCHAR(50) not null default '' comment('文件名称')"`
	Type      int       `json:"type" xorm:"TINYINT(1) not null index default 0 comment('转换类型')"`
	Plate     int       `json:"type" xorm:"TINYINT(1) not null index default 0 comment('转换平台 0无 1.DeepAi 2.U2net 3.百度Ai')"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}
