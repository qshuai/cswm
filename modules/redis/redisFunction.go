package redis_orm

import (
	"ERP/models"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

//同步mysql数据表permission中的所有数据到redis
func (r *RedisStorage) StorePermission(permission_items []models.Permission) error {
	//获取permission redis存储前缀
	permission_prefix := beego.AppConfig.String("redis::permission_prefix")

	ri := r.pool.Get()
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
		)
		if err != nil {
			logs.Error("存储permission redis错误", err)
		}
	}
	return nil
}

//同步一个用户的permission数据到redis中
func (r *RedisStorage) StoreOnePermission(permission models.Permission) error {
	//获取permission redis存储前缀
	permission_prefix := beego.AppConfig.String("redis::permission_prefix")

	ri := r.pool.Get()
	var err error
	_, err = ri.Do("HMSET", permission_prefix+permission.User.Username,
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
	)
	if err != nil {
		logs.Error("存储permission redis错误", err)
	}

	return nil
}

//获取某个人的permission的一行数据
func (r *RedisStorage) GetOneRowPermission(uid int) map[string]bool {
	//获取permission redis存储前缀
	permission_prefix := beego.AppConfig.String("redis::permission_prefix")

	//获取用户名
	user := models.User{}
	o := orm.NewOrm()
	o.QueryTable("user").Filter("id", uid).One(&user, "username")
	ri := r.pool.Get()
	res, _ := redis.StringMap(ri.Do("HGETALL", permission_prefix+user.Username))
	var maps = make(map[string]bool, len(res))
	for key, value := range res {
		maps[key], _ = strconv.ParseBool(value)
	}
	return maps
}

//获取某个人的permission的某个权限数据
func (r *RedisStorage) GetOneItemPermission(uid int, key string) bool {
	//获取permission redis存储前缀
	permission_prefix := beego.AppConfig.String("redis::permission_prefix")

	//获取用户名
	user := models.User{}
	o := orm.NewOrm()
	o.QueryTable("user").Filter("id", uid).One(&user, "username")

	ri := r.pool.Get()
	b, _ := redis.Bool(ri.Do("HGET", permission_prefix+user.Username, key))
	return b
}

//存储所有user数据表中的position到redis
func (r *RedisStorage) StorePosition(position_items []models.User) error {
	//获取position redis存储前缀
	position_fix := beego.AppConfig.String("redis::position_prefix")

	ri := r.pool.Get()
	var err error
	for _, item := range position_items {
		_, err = ri.Do("SET", position_fix+item.Username, item.Position)
		if err != nil {
			logs.Error("存储position redis错误", err)
		}
	}
	return nil
}

//存储某个user的position到redis
func (r *RedisStorage) StoreOnePosition(position_item models.User) error {
	//获取position redis存储前缀
	position_fix := beego.AppConfig.String("redis::position_prefix")

	ri := r.pool.Get()
	var err error

	_, err = ri.Do("SET", position_fix+position_item.Username, position_item.Position)
	if err != nil {
		logs.Error("存储position redis错误", err)
	}
	return nil
}

//从redis中获取某个人员的position数据
func (r *RedisStorage) GetOnePosition(username string) string {
	//获取position redis存储前缀
	position_fix := beego.AppConfig.String("redis::position_prefix")

	ri := r.pool.Get()
	res, _ := redis.String(ri.Do("GET", position_fix + username))
	return res
}
