package redis_orm

import "github.com/qshuai/cswm/models"

type RedisInterface interface {
	//权限操作
	StorePermission(permission_items []models.Permission) error
	StoreOnePermission(username string, permission models.Permission) error
	GetOneRowPermission(username string) map[string]bool
	GetOneItemPermission(username, key string) bool

	//职位操作
	StorePosition(position_items []models.User) error
	StoreOnePosition(position_item models.User) error
	GetOnePosition(username string) string

	//消息操作
	IncrOneMessage(username string) error
	DecrOneMessage(username string) error
	GetOneMessageNum(username string) int
	StoreAllMessage2Redis(message []M) error

	//修改key
	RenameKey(old, new string) error
}
