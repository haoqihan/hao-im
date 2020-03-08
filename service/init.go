package service

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"hao-im/model"
	"log"
)

var DbEngin *xorm.Engine

func init() {
	var err error
	driverName := "mysql"
	dataSourceName := "root:12345678@(127.0.0.1:3306)/chat?charset=utf8"
	DbEngin, err = xorm.NewEngine(driverName, dataSourceName)
	if err != nil {
		log.Fatal(err.Error())
	}
	// 是否显示sql语句
	DbEngin.ShowSQL(true)
	DbEngin.ShowSQL(true)
	// 设置数据库最大打开连接数
	DbEngin.SetMaxOpenConns(2)
	// 创建user表 Contact Community
	DbEngin.Sync2(new(model.User), new(model.Contact), new(model.Community))

	fmt.Println("init data base ok")

}
