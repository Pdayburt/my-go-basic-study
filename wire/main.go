package wire

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"my-go-basic-study/wire/repository"
	"my-go-basic-study/wire/repository/dao"
)

func main() {
	db, err := gorm.Open(mysql.Open("dsn"))
	if err != nil {
		panic("failed to connect database")
	}
	userDao := dao.NewUser(db)
	userRepository := repository.NewUserRepository(userDao)
	fmt.Println(userRepository)

}
