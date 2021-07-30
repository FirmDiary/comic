package datamodels

import "time"

const (
	TypeAdd  = 1
	TypeDecr = 2

	SourceUse   = 1
	SourceShare = 2
)

type QuotaLog struct {
	Id        int64     `json:"id" xorm:"INT(10) not null pk autoincr"`
	UserId    int64     `json:"user_id" xorm:"INT(10) not null index"`
	Type      int       `json:"type" xorm:"TINYINT(1) not null comment('类型 1增加 2减少')"`
	Source    int       `json:"type" xorm:"TINYINT(1) not null default 1 comment('来源 1使用消耗 2完成邀请任务')"`
	Quota     int64     `json:"int" xorm:"INT(10) not null comment('额度')"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}
