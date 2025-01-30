package main

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	/*s := pflag.String("name", "Jack", "姓名")
	pflag.Parse()
	fmt.Println("pflag.String: ", *s)*/
	initLogger()
	initViper()
	serve := InitWebService()
	serve.Run(":8080")

}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	type Person struct {
		Name string
	}

	zap.ReplaceGlobals(logger)
	zap.L().Info("绑定参数测试",
		zap.Error(errors.New("这是绑定参数测试错误")),
		zap.String("name", "jack"),
		zap.Int8("age", 18),
		zap.Any("ip", Person{Name: "Rose"}))
}

func initViper() {
	viper.SetConfigFile("config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
