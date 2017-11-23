package redis_orm

import (
	"strconv"

	"erp/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
)

//同步mysql数据表permission中的所有数据到redis
func (r *RedisStorage) StorePermission(permission_items []models.Permission) error {
	//获取permission redis存储前缀
	permission_prefix := beego.AppConfig.String("redis::permission_prefix")

	ri := r.pool.Get()
	defer ri.Close()
	var err error
	for _, item := range permission_items {
		_, err = ri.Do("HMSET", permission_prefix+item.User.Username,
			"AddMember", item.AddMember,
			"EditMember", item.EditMember,
			"ActiveMember", item.ActiveMember,
			"AddConsumer", item.AddConsumer,
			"EditConsumer", item.EditConsumer,
			"ViewConsumer", item.ViewConsumer,
			"AddBrand", item.AddBrand,
			"AddDealer", item.AddDealer,
			"ViewDealer", item.ViewDealer,
			"AddSupplier", item.AddSupplier,
			"ViewSupplier", item.ViewSupplier,
			"AddProduct", item.AddProduct,
			"InputInPrice", item.InputInPrice,
			"ViewProductStore", item.ViewProductStore,
			"ViewStock", item.ViewStock,
			"ViewInPrice", item.ViewInPrice,
			"EditProduct", item.EditProduct,
			"DeleteProduct", item.DeleteProduct,
			"OutputProduct", item.OutputProduct,
			"ViewSale", item.ViewSale,
			"ViewSaleConsumer", item.ViewSaleConsumer,
			"ViewSaleInPrice", item.ViewSaleInPrice,
			"EditSale", item.EditSale,
			"OperateCategory", item.OperateCategory,
			"RequestMove", item.RequestMove,
			"ResponseMove", item.ResponseMove,
			"ViewMove", item.ViewMove,
			"AddStore", item.AddStore,
			"ViewStore", item.ViewStore,
			"OperateOtherStore", item.OperateOtherStore,
		)
		if err != nil {
			logs.Error("存储permission redis错误", err)
		}
	}
	return nil
}

//同步一个用户的permission数据到redis中
func (r *RedisStorage) StoreOnePermission(username string, permission models.Permission) error {
	//获取permission redis存储前缀
	permission_prefix := beego.AppConfig.String("redis::permission_prefix")

	ri := r.pool.Get()
	defer ri.Close()
	var err error
	_, err = ri.Do("HMSET", permission_prefix+username,
		"AddMember", permission.AddMember,
		"EditMember", permission.EditMember,
		"ActiveMember", permission.ActiveMember,
		"AddConsumer", permission.AddConsumer,
		"EditConsumer", permission.EditConsumer,
		"ViewConsumer", permission.ViewConsumer,
		"AddBrand", permission.AddBrand,
		"AddDealer", permission.AddDealer,
		"ViewDealer", permission.ViewDealer,
		"AddSupplier", permission.AddSupplier,
		"ViewSupplier", permission.ViewSupplier,
		"AddProduct", permission.AddProduct,
		"InputInPrice", permission.InputInPrice,
		"ViewProductStore", permission.ViewProductStore,
		"ViewStock", permission.ViewStock,
		"ViewInPrice", permission.ViewInPrice,
		"EditProduct", permission.EditProduct,
		"DeleteProduct", permission.DeleteProduct,
		"OutputProduct", permission.OutputProduct,
		"ViewSale", permission.ViewSale,
		"ViewSaleConsumer", permission.ViewSaleConsumer,
		"ViewSaleInPrice", permission.ViewSaleInPrice,
		"EditSale", permission.EditSale,
		"OperateCategory", permission.OperateCategory,
		"RequestMove", permission.RequestMove,
		"ResponseMove", permission.ResponseMove,
		"ViewMove", permission.ViewMove,
		"AddStore", permission.AddStore,
		"ViewStore", permission.ViewStore,
		"OperateOtherStore", permission.OperateOtherStore,
	)
	if err != nil {
		logs.Error("存储permission redis错误", err)
	}

	return nil
}

//获取某个人的permission的一行数据
func (r *RedisStorage) GetOneRowPermission(username string) map[string]bool {
	//获取permission redis存储前缀
	permission_prefix := beego.AppConfig.String("redis::permission_prefix")

	ri := r.pool.Get()
	defer ri.Close()
	res, _ := redis.StringMap(ri.Do("HGETALL", permission_prefix+username))
	var maps = make(map[string]bool, len(res))
	for key, value := range res {
		maps[key], _ = strconv.ParseBool(value)
	}
	return maps
}

//获取某个人的permission的某个权限数据
func (r *RedisStorage) GetOneItemPermission(username string, key string) bool {
	//获取permission redis存储前缀
	permission_prefix := beego.AppConfig.String("redis::permission_prefix")

	ri := r.pool.Get()
	defer ri.Close()
	b, _ := redis.Bool(ri.Do("HGET", permission_prefix+username, key))
	return b
}

//存储所有user数据表中的position到redis
func (r *RedisStorage) StorePosition(position_items []models.User) error {
	//获取position redis存储前缀
	userdata_prefix := beego.AppConfig.String("redis::userdata_prefix")

	ri := r.pool.Get()
	defer ri.Close()
	var err error
	for _, item := range position_items {
		_, err = ri.Do("HSET", userdata_prefix+item.Username, "position", item.Position)
		if err != nil {
			logs.Error("存储position redis错误", err)
		}
	}
	return nil
}

//存储某个user的position到redis
func (r *RedisStorage) StoreOnePosition(position_item models.User) error {
	//获取position redis存储前缀
	userdata_prefix := beego.AppConfig.String("redis::userdata_prefix")

	ri := r.pool.Get()
	defer ri.Close()
	var err error

	_, err = ri.Do("HSET", userdata_prefix+position_item.Username, "position", position_item.Position)
	if err != nil {
		logs.Error("存储position redis错误", err)
	}
	return nil
}

//从redis中获取某个人员的position数据
func (r *RedisStorage) GetOnePosition(username string) string {
	//获取position redis存储前缀
	userdata_prefix := beego.AppConfig.String("redis::userdata_prefix")

	ri := r.pool.Get()
	defer ri.Close()
	res, _ := redis.String(ri.Do("HGET", userdata_prefix+username, "position"))
	return res
}

//增加某人的消息数量
func (r *RedisStorage) IncrOneMessage(username string) error {
	//获取position redis存储前缀
	userdata_prefix := beego.AppConfig.String("redis::userdata_prefix")

	ri := r.pool.Get()
	defer ri.Close()
	_, err := ri.Do("HINCRBY", userdata_prefix+username, "message", 1)
	return err
}

//减少某人的消息数量
func (r *RedisStorage) DecrOneMessage(username string) error {
	//获取position redis存储前缀
	userdata_prefix := beego.AppConfig.String("redis::userdata_prefix")

	ri := r.pool.Get()
	defer ri.Close()
	_, err := ri.Do("HINCRBY", userdata_prefix+username, "message", -1)
	return err
}

//获取某人的消息数量
func (r *RedisStorage) GetOneMessageNum(username string) int {
	//获取position redis存储前缀
	userdata_prefix := beego.AppConfig.String("redis::userdata_prefix")

	ri := r.pool.Get()
	defer ri.Close()
	res, _ := redis.Int(ri.Do("HGET", userdata_prefix+username, "message"))
	return res
}

//存储所有的未读message到redis
type M struct {
	Username string
	Num      string
}

func (r *RedisStorage) StoreAllMessage2Redis(message []M) error {
	//获取position redis存储前缀
	userdata_prefix := beego.AppConfig.String("redis::userdata_prefix")

	ri := r.pool.Get()
	defer ri.Close()

	for _, item := range message {
		ri.Do("HSET", userdata_prefix+item.Username, "message", item.Num)
	}
	return nil
}

//修改key
func (r *RedisStorage) RenameKey(old, new string) error {
	userdata_prefix := beego.AppConfig.String("redis::userdata_prefix")
	permission_prefix := beego.AppConfig.String("redis::permission_prefix")

	ri := r.pool.Get()
	defer ri.Close()
	ri.Do("RENAME", userdata_prefix+old, userdata_prefix+new)
	ri.Do("RENAME", permission_prefix+old, permission_prefix+new)
	return nil
}
