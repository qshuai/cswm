package redis_orm

import "ERP/models"

type RedisInterface interface {
	StorePermission(permission_items []models.Permission) error
}
