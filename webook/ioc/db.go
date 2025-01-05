package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"my-go-basic-study/webook/config"
	"my-go-basic-study/webook/internal/repository/dao"
)

func InitDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.NSD))
	//db, err := gorm.Open(mysql.Open("root:root@tcp(webook-mysql:3309)/webook"))
	//db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		//panic 整个go routine 结束
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
