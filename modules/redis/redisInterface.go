package redis_orm

import "ERP/models"

type RedisInterface interface {
	StorePermission(permission_items []models.Permission) error
	StoreOnePermission(permission models.Permission) error
	GetOneRowPermission(uid int) map[string]bool
	GetOneItemPermission(uid int, key string) bool
}
