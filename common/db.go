package common

import (
    "comic/config"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/go-xorm/xorm"
)

var engine *xorm.Engine

func NewDbEngine() *xorm.Engine {
    cfg := config.GetConfig().Database

    cntUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", cfg.User, cfg.Password, cfg.Host, cfg.Database)
    engine, err := xorm.NewEngine(cfg.Connection, cntUrl)
    if err != nil {
        panic(err)
    }
    // 显示sql
    engine.ShowSQL(true)
    return engine
}
