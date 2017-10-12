package redis_orm

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func (r *RedisStorage) Set_user(key string, value int) error {
	con := r.pool.Get()
	_, err := con.Do("SET", key, value)
	if err != nil {
		errors.New("redis operation error")
	}
	fmt.Println(redis.String(con.Do("GET", key)))
	return nil
}
