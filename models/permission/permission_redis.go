package permission

import (
	"github.com/astaxie/beego/orm"
	"ERP/models"
	"ERP/modules/redis"
)

/*
 * 同步mysql中所有用户permission数据到redis（hash）中，但不包括没有设置username的用户
 */
func AsyncMysql2RedisAll() {
	o := orm.NewOrm()
	permission_items := []models.Permission{}
	o.QueryTable("permission").Exclude("user__username", "").RelatedSel().All(&permission_items)
	redis_orm.RedisPool.StorePermission(permission_items)
}

func AsyncMysql2RedisOne(uid int) {
	o := orm.NewOrm()
	permission_item := models.Permission{}
	o.QueryTable("permission").Filter("user__id", uid).RelatedSel().One(&permission_item)
	redis_orm.RedisPool.StoreOnePermission(permission_item)
}

func GetOneRowPermission(uid int) map[string]bool {
	return redis_orm.RedisPool.GetOneRowPermission(uid)
}

func GetOneItemPermission(uid int, key string) bool {
	return redis_orm.RedisPool.GetOneItemPermission(uid, key)
}
