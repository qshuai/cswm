package controllers

import (
	"github.com/astaxie/beego"
	"html/template"
	"github.com/astaxie/beego/orm"
	"ERP/models"
	"strings"
	"time"
	"strconv"
	"encoding/json"
)

type ProductController struct {
	beego.Controller
}

//获取商品列表
func (c *ProductController) Get() {
	o := orm.NewOrm()
	product := []models.Product{}
	o.QueryTable("product").RelatedSel().OrderBy("-id").All(&product)

	product_byte, _ := json.Marshal(product)
	c.Data["product"] = string(product_byte)
	c.Data["xsrf_token"] = c.XSRFToken()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "product/product_list.html"
}

//删除单条商品信息
func (c *ProductController) Product_item_delete() {
	if c.IsAjax() {
		o := orm.NewOrm()
		product := models.Product{}
		product.Id, _ = c.GetInt("product_id")

		sale_exist := o.QueryTable("sale").Filter("product__id", product.Id).Exist()
		move_exist := o.QueryTable("move").Filter("origin__id", product.Id).Exist()

		if !sale_exist && !move_exist {
			o.Delete(&product)
			c.Data["json"] = ResponseInfo{
				Code:    "success",
				Message: "删除成功",
				Data:    "",
			}
			c.ServeJSON()
		} else {
			switch {
			case sale_exist:
				c.Data["json"] = ResponseInfo{
					Code:    "error",
					Message: "由于此商品存在销售记录，不能删除~",
					Data:    "",
				}
				c.ServeJSON()
			case move_exist:
				c.Data["json"] = ResponseInfo{
					Code:    "error",
					Message: "由于此商品存在移库记录，不能删除~",
					Data:    "",
				}
				c.ServeJSON()
			}
		}

		c.Data["json"] = ResponseInfo{
			Code:    "unknown",
			Message: "未知情况",
			Data:    "",
		}
		c.ServeJSON()
	}
	return
}

//编辑单条商品信息
func (c *ProductController) Product_item_edit() {
	product := models.Product{}
	o := orm.NewOrm()

	product.Id, _ = c.GetInt("product_id")
	//product.Title = c.GetString("title")

	//product.ArtNum = c.GetString("atr_num")
	product.LotNum = c.GetString("lot_num")

	//product.Spec = c.GetString("spec")
	product.Stock, _ = c.GetUint32("stock")
	product.InPrice, _ = c.GetFloat("in_price")
	//product.Unit = c.GetString("unit")

	brand := models.Brand{}
	o.QueryTable("brand").Filter("name", c.GetString("brand")).One(&brand, "id")
	//product.Brand = &brand

	category := models.Category{}
	o.QueryTable("category").Filter("three_stage", c.GetString("three_stage")).One(&category, "id")
	//product.CatNum = &category

	supplier := models.Supplier{}
	o.QueryTable("supplier").Filter("name", c.GetString("supplier")).One(&supplier, "id")
	product.Supplier = &supplier

	store_string := strings.Split(c.GetString("store"), "-")
	store := models.Store{}
	o.QueryTable("store").Filter("pool", store_string[0]).Filter("name", store_string[1]).One(&store, "id")
	product.Store = &store

	product.HasPay, _ = c.GetBool("has_pay_edit")
	product.HasInvioce, _ = c.GetBool("has_invioce_edit")

	product.GetInvioce, _ = time.Parse("2006-1-2", c.GetString("get_invioce_edit"))

	num, err := o.Update(&product, "title", "brand_id", "art_num", "lot_num", "cat_num_id", "spec", "stock", "unit", "store_id", "supplier_id", "in_price", "has_pay", "has_invioce", "get_invioce")
	if num == 1 && err == nil {
		c.Data["url"] = "/product_list"
		c.Data["msg"] = "商品信息修改成功"
		c.TplName = "jump/success.html"
	} else {
		c.Data["url"] = "/product_list"
		c.Data["msg"] = "商品信息修改失败"
		c.TplName = "jump/error.html"
	}

}

//商品添加页面
func (c *ProductController) Add_get() {
	o := orm.NewOrm()
	user := models.User{}
	user.Id, _ = c.GetSession("uid").(int)
	o.QueryTable("user").Filter("id", user.Id).One(&user, "pool_name")

	c.Layout = "common.tpl"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["brand_string"] = GetBrandList()
	c.Data["supplier_string"] = GetSupplierList()
	c.Data["store_string"] = GetStoreList(user.PoolName)
	c.Data["three_stage_string"] = GetThreeStageList()
	c.TplName = "product/product_add.html"
}

//商品添加post
func (c *ProductController) Add_post() {
	product := models.Product{}

	user := models.User{}
	user.Id = c.GetSession("uid").(int)
	product.User = &user
	//product.Title = c.GetString("title")
	//product.ArtNum = c.GetString("atr_num")
	product.LotNum = c.GetString("lot_num")
	//product.Unit = c.GetString("unit")
	product.HasPay, _ = c.GetBool("has_pay")
	product.HasInvioce, _ = c.GetBool("has_invioce")

	o := orm.NewOrm()

	brand := models.Brand{}
	o.QueryTable("brand").Filter("name", c.GetString("brand")).One(&brand, "id")
	//product.Brand = &brand

	category := models.Category{}
	o.QueryTable("category").Filter("three_stage", c.GetString("three_stage")).One(&category, "id")
	//product.CatNum = &category

	supplier := models.Supplier{}
	o.QueryTable("supplier").Filter("name", c.GetString("supplier")).One(&supplier, "id")
	product.Supplier = &supplier

	store_string := strings.Split(c.GetString("store"), "-")
	store := models.Store{}
	o.QueryTable("store").Filter("pool", store_string[0]).Filter("name", store_string[1]).One(&store, "id")
	product.Store = &store

	product.GetInvioce, _ = time.Parse("2006-1-2", c.GetString("get_invioce"))

	//spec_slice := c.GetStrings("spec")
	//stock_slice := c.GetStrings("stock")
	//inprice_slice := c.GetStrings("in_price")

	//for index, item := range spec_slice {
	//	product.Spec = item
	//	stock_temp, _ := strconv.ParseUint(stock_slice[index], 10, 0)
	//	product.Stock = uint32(stock_temp)
	//	product.InPrice, _ = strconv.ParseFloat(inprice_slice[index], 64)
	//	_, err := o.Insert(&product)
	//
	//	//防止出现重复主键值
	//	product.Id++
	//
	//	if err != nil {
	//		c.Data["msg"] = "添加商品失败~"
	//		c.Data["url"] = "/product_add"
	//		c.TplName = "jump/error.html"
	//		return
	//	}
	//}
	c.Data["msg"] = "添加商品成功~"
	c.Data["url"] = "/product_list"
	c.TplName = "jump/error.html"
}

//ajax通过货号一键填充表单
func (c *ProductController) SearchByCatnum() {
	if !c.IsAjax() {
		return
	}

	art_num := c.GetString("art_num")

	o := orm.NewOrm()
	product := []models.Product{}
	product_spec := []models.Product{}
	product_temp := models.Product{}

	//获取不同spec类别
	o.QueryTable("product").Distinct().Filter("art_num", art_num).All(&product_spec, "spec")

	//通过spec和art_num限定，循环查询
	//for _, item := range product_spec {
		//o.QueryTable("product").Filter("art_num", art_num).Filter("spec", item.Spec).RelatedSel().One(&product_temp)
		product = append(product, product_temp)
	//}

	c.Data["json"] = product
	c.ServeJSON()
}

func (c *ProductController) Product_track() {

}

//管理员添加商品模板
func (c *ProductController) ProductTemplateList() {
	c.Layout = "common.tpl"
	c.TplName = "product/product_template_list.html"
}

//商品模板添加页面
func (c *ProductController) ProductTemplateAdd() {
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["brand_string"] = GetBrandList()
	c.Data["supplier_string"] = GetSupplierList()
	c.Data["three_stage_string"] = GetThreeStageList()
	c.Layout = "common.tpl"
	c.TplName = "product/product_template_add.html"
}

//商品模板添加提交
func (c *ProductController) ProductTemplateAddPost() {
	product := models.ProductTemplate{}

	product.Title = c.GetString("title")
	product.ArtNum = c.GetString("atr_num")
	product.Unit = c.GetString("unit")

	o := orm.NewOrm()

	brand := models.Brand{}
	o.QueryTable("brand").Filter("name", c.GetString("brand")).One(&brand, "id")
	product.Brand = &brand

	category := models.Category{}
	o.QueryTable("category").Filter("three_stage", c.GetString("three_stage")).One(&category, "id")
	product.CatNum = &category

	supplier := models.Supplier{}
	o.QueryTable("supplier").Filter("name", c.GetString("supplier")).One(&supplier, "id")
	//product.Supplier = &supplier

	spec_slice := c.GetStrings("spec")
	inprice_slice := c.GetStrings("in_price")

	for index, item := range spec_slice {
		product.Spec = item
		if inprice_slice[index] != "" {
			product.InPrice, _ = strconv.ParseFloat(inprice_slice[index], 64)
		}
		_, err := o.Insert(&product)

		//防止出现重复主键值
		product.Id++

		if err != nil {
			c.Data["msg"] = "添加商品模板失败~"
			c.Data["url"] = "/product_template_add"
			c.TplName = "jump/error.html"
			return
		}
	}
	c.Data["msg"] = "添加商品模板成功~"
	c.Data["url"] = "/product_template_list"
	c.TplName = "jump/error.html"
}

//商品模板编辑提交
func (c *ProductController) ProductTemplateEditPost() {

}

//公共函数
//获取品牌列表
func GetBrandList() string {
	o := orm.NewOrm()
	brand := []models.Brand{}
	o.QueryTable("brand").All(&brand, "name")
	var brand_string string
	for _, item := range brand {
		brand_string += item.Name + ", "
	}
	return brand_string
}

//获取供应商列表
func GetSupplierList() string {
	o := orm.NewOrm()
	supplier := []models.Supplier{}
	o.QueryTable("supplier").All(&supplier, "name")
	var supplier_string string
	for _, item := range supplier {
		supplier_string += item.Name + ", "
	}
	return supplier_string
}

//获取库房列表
func GetStoreList(pool_name string) string {

	var store_string string

	if strings.Contains(pool_name, "-") {
		store_string = pool_name + ","
	} else {
		o := orm.NewOrm()
		store := []models.Store{}
		o.QueryTable("store").Filter("pool", pool_name).All(&store, "name")

		for _, item := range store {
			store_string += pool_name + "-" + item.Name + ", "
		}
	}

	return store_string
}

//获取三级分类列表
func GetThreeStageList() string {
	o := orm.NewOrm()
	three_stage := []models.Category{}
	o.QueryTable("category").Exclude("three_stage", "-").Filter("is_hidden", false).All(&three_stage, "three_stage")
	var three_stage_string string
	for _, item := range three_stage {
		three_stage_string += item.ThreeStage + ", "
	}
	return three_stage_string
}
