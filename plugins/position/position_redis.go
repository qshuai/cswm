package position

import (
	"erp/models"
	"erp/modules/redis"

	"github.com/astaxie/beego/orm"
)

func AsyncAllPosition() {
	o := orm.NewOrm()
	user := []models.User{}
	o.QueryTable("user").Exclude("username", "").All(&user, "username", "position")

	redis_orm.RedisPool.StorePosition(user)
}

func AsyncOnePosition(user models.User) {
	redis_orm.RedisPool.StoreOnePosition(user)
}

func GetOnePosition(username string) string {
	return redis_orm.RedisPool.GetOnePosition(username)
}
