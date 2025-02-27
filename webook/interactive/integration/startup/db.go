package startup

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"my-go-basic-study/webook/interactive/repository/dao"
)

func InitDb() *gorm.DB {

	type MysqlConfig struct {
		DSN string `yaml:"dsn"`
	}
	var mysqlConfig MysqlConfig
	err := viper.UnmarshalKey("db.mysql", &mysqlConfig)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(mysqlConfig.DSN))
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
