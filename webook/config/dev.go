//go:build !k8s

// Package config 不使用k8s这个编译标签
package config

var Config = config{
	DB:          DBConfig{NSD: "root:root@tcp(localhost:13316)/webook"},
	RedisConfig: RedisConfig{Addr: "localhost:6379"},
}
