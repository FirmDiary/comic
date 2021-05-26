package datamodels

import (
    "time"
)

const (
    TypeFace = 1
    TypeAll  = 2
)

type Upload struct {
    Id        int64     `json:"id" xorm:"INT(10) not null pk autoincr"`
    UserId    int64     `json:"user_id" xorm:"INT(10) not null index comment('用户ID')"`
    File      string    `json:"file" xorm:"VARCHAR(50) not null default '' comment('文件名称')"`
    Type      int       `json:"type" xorm:"TINYINT(1) not null index default 1 comment('转换类型 1人脸 2其他')"`
    CreatedAt time.Time `xorm:"created"`
    UpdatedAt time.Time `xorm:"updated"`
}
