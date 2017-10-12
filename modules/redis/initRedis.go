package redis_orm

import "github.com/astaxie/beego"

func init() {
	//获取配置信息
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

	obj, err := NewRedis(redis_config)
	if err != nil {

	}
	obj.Set_user("name", 45)
}
