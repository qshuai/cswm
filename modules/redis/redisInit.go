package redis_orm

import "github.com/astaxie/beego"

var RedisPool RedisInterface

//获取redis配置信息
func init() {
	port, _ := beego.AppConfig.Int("redis::port")
	DB, _ := beego.AppConfig.Int("redis::db")
	IdleTimeout, _ := beego.AppConfig.Int("redis::idletimeout")
	MaxIdleConns, _ := beego.AppConfig.Int("redis::maxidleconns")
	MaxOpenConns, _ := beego.AppConfig.Int("redis::maxopenconns")
	InitialOpenConns, _ := beego.AppConfig.Int("redis::initialopenconns")

	redis_config := Redis{
		Host:             beego.AppConfig.String("redis::host"),
		Port:             port,
		Password:         beego.AppConfig.String("redis::password"),
		DB:               DB,
		IdleTimeout:      IdleTimeout,
		MaxIdleConns:     MaxIdleConns,
		MaxOpenConns:     MaxOpenConns,
		InitialOpenConns: InitialOpenConns,
	}

	var err error
	RedisPool, err = NewRedis(redis_config)
	if err != nil {

	}
}
