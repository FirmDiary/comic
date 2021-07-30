package datamodels

import "time"

type InviteLog struct {
	Id        int64     `json:"id" xorm:"INT(10) not null pk autoincr"`
	UserId    int64     `json:"user_id" xorm:"INT(10) not null index comment('受邀请者id')"`
	InviteId  int64     `json:"invite_id" xorm:"INT(10) not null index comment('邀请者id')"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}
