package redis_orm

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Redis struct {
	Host             string
	Port             int
	Password         string
	DB               int
	IdleTimeout      int
	MaxIdleConns     int
	MaxOpenConns     int
	InitialOpenConns int
}

type RedisStorage struct {
	pool   *redis.Pool
	config Redis
}

func NewRedis(config Redis) (RedisInterface, error) {
	r := &RedisStorage{
		config: config,
		pool: &redis.Pool{
			MaxIdle:     config.MaxIdleConns,
			MaxActive:   config.MaxOpenConns,
			IdleTimeout: time.Duration(config.IdleTimeout) * time.Second,
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}

	r.pool.Dial = r.createRedisConnect

	//初始化连接数量
	if config.InitialOpenConns > config.MaxIdleConns {
		config.InitialOpenConns = config.MaxIdleConns
	} else if config.InitialOpenConns == 0 {
		config.InitialOpenConns = 1
	}

	return r, r.initRedisConnect()
}

func (r *RedisStorage) createRedisConnect() (redis.Conn, error) {
	return redis.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", r.config.Host, r.config.Port),
		redis.DialDatabase(r.config.DB),
		redis.DialPassword(r.config.Password),
	)
}

func (r *RedisStorage) initRedisConnect() error {
	cons := make([]redis.Conn, r.config.InitialOpenConns)
	defer func() {
		for _, c := range cons {
			if c != nil {
				c.Close()
			}
		}
	}()

	for i := 0; i < r.config.InitialOpenConns; i++ {
		cons[i] = r.pool.Get()
		if _, err := cons[i].Do("PING"); err != nil {
			return err
		}
	}

	return nil
}
