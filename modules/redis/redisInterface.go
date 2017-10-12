package redis_orm

type redisInterface interface {
	Set_user(key string, value int) error
}