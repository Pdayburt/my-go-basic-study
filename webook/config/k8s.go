//go:build k8s

// 使用k8s这个编译标签
package config

var Config = config{
	DB:          DBConfig{NSD: "root:root@tcp(webook-mysql:3309)/webook"},
	RedisConfig: RedisConfig{Addr: "webook-redis:16379"},
}
