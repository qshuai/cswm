package redis_orm

import (
	"ERP/models"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
)

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
