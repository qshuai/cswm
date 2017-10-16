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

//同步一个用户的permission数据到redis中
func AsyncMysql2RedisOne(username string) {
	o := orm.NewOrm()
	permission_item := models.Permission{}
	o.QueryTable("permission").Filter("user__username", username).One(&permission_item)
	redis_orm.RedisPool.StoreOnePermission(username, permission_item)
}

//获取某个人的permission的一行数据
func GetOneRowPermission(username string) map[string]bool {
	return redis_orm.RedisPool.GetOneRowPermission(username)
}

//获取某个人的permission的某个权限数据
func GetOneItemPermission(username, key string) bool {
	return redis_orm.RedisPool.GetOneItemPermission(username, key)
}
