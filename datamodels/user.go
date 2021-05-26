package datamodels

import "time"

type User struct {
    Id        int64     `json:"id" xorm:"INT(10) not null pk autoincr"`
    Openid    string    `json:"openid" xorm:"VARCHAR(150) not null comment('用户openid')"`
    UnionId   string    `json:"union_id" xorm:"VARCHAR(255) not null comment('用户union_id')"`
    CreatedAt time.Time `xorm:"created"`
    UpdatedAt time.Time `xorm:"updated"`
}
