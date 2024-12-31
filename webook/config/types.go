package config

type config struct {
	DB          DBConfig
	RedisConfig RedisConfig
}

type DBConfig struct {
	NSD string
}
type RedisConfig struct {
	Addr string
}
