package datamodels

type App struct {
	Id           int64  `json:"id" xorm:"INT(10) not null pk autoincr"`
	Name         string `json:"name" xorm:"VARCHAR(64) not null comment('名称')"`
	AppId        string `json:"app_id" xorm:"CHAR(18) not null comment('appid')"`
	DefaultQuota int64  `json:"int" xorm:"INT(10) not null comment('初始额度')"`
	AppSecret    string `json:"app_secret" xorm:"CHAR(32) not null comment('appsecret')"`
}
