package msg

import (
	"erp/modules/redis"
	"github.com/astaxie/beego/orm"
)

func IncrOneMessage(username string) error {
	return redis_orm.RedisPool.IncrOneMessage(username)
}

func DecrOneMessage(username string) error {
	return redis_orm.RedisPool.DecrOneMessage(username)
}

func GetOneMessageNum(username string) int {
	return redis_orm.RedisPool.GetOneMessageNum(username)
}

func AsyncAllMessage2Redis() error {
	mm := []redis_orm.M{}

	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("user.username, count(to_id) as num").
		From("message").
		InnerJoin("user").
		On("user.id = message.to_id").
		Where("is_read = 0").
		GroupBy("to_id")
	sql := qb.String()
	o := orm.NewOrm()
	o.Raw(sql).QueryRows(&mm)
	return redis_orm.RedisPool.StoreAllMessage2Redis(mm)
}
