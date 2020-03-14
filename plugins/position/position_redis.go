package position

import (
	"github.com/astaxie/beego/orm"
	"github.com/qshuai/cswm/models"
	"github.com/qshuai/cswm/modules/redis"
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
