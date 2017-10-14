package redis_orm

import "ERP/models"

type RedisInterface interface {
	//权限操作
	StorePermission(permission_items []models.Permission) error
	StoreOnePermission(permission models.Permission) error
	GetOneRowPermission(uid int) map[string]bool
	GetOneItemPermission(uid int, key string) bool

	//职位操作
	StorePosition(position_items []models.User) error
	StoreOnePosition(position_item models.User) error
	GetOnePosition(username string) string
}
