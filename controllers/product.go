package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"erp/models"
	"erp/plugins/permission"
	"erp/plugins/position"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type ProductController struct {
	beego.Controller
}

type product struct {
	Id           int
	UserName     string
	Title        string
	BrandName    string
	ArtNum       string
	LotNum       string
	ThreeStage   string
	Spec         string
	Stock        string
	Unit         string
	Pool         string
	StoreName    string
	InTime       time.Time
	SupplierName string
	InPrice      string
	HasPay       bool
	HasInvoice   bool
	GetInvoice   time.Time
}

//定义querybuiler查询结果的接受结构体
type product_template struct {
	Id         int
	ThreeStage string
	Title      string
	ArtNum     string
	Spec       string
	Unit       string
	Suppliers  string
	InPrice    float64
	BrandName  string
	DealerName string
}

//每次获取的product数据条数
const ProductLimit = 100
const TemplateLimit = 100

//获取商品列表
func (c *ProductController) Get() {
	username := c.GetSession("username").(string)
	user := models.User{}
	o := orm.NewOrm()
	o.QueryTable("user").Filter("username", username).One(&user, "position", "pool_name")
	operate_other_store := !permission.GetOneItemPermission(username, "OperateOtherStore")

	qb, _ := orm.NewQueryBuilder("mysql")
	p := []product{}
	if user.Position != "超级管理员" || operate_other_store {
		if user.PoolName != "" {
			qb.Select("product.id", "product.title", "product.art_num", "product.lot_num", "product.spec", "product.stock", "product.unit",
				"product.in_time", "product.in_price", "product.has_pay", "product.has_invoice", "product.get_invoice",
				"brand.name as brand_name", "supplier.name as supplier_name", "category.three_stage", "user.name as user_name",
				"store.pool", "store.name as store_name").
				From("product").
				LeftJoin("brand").
				On("brand.id = product.brand_id").
				LeftJoin("supplier").
				On("supplier.id = product.supplier_id").
				LeftJoin("category").
				On("category.id = product.cat_num_id").
				LeftJoin("user").
				On("user.id = product.user_id").
				InnerJoin("store").
				On("store.id = product.store_id AND store.name = ?").
				OrderBy("in_time").Desc().
				Limit(ProductLimit)
			sql := qb.String()
			o.Raw(sql, user.PoolName).QueryRows(&p)
		}
	} else {
		qb.Select("product.id", "product.title", "product.art_num", "product.lot_num", "product.spec", "product.stock", "product.unit",
			"product.in_time", "product.in_price", "product.has_pay", "product.has_invoice", "product.get_invoice",
			"brand.name as brand_name", "supplier.name as supplier_name", "category.three_stage", "user.name as user_name",
			"store.pool", "store.name as store_name").
			From("product").
			LeftJoin("brand").
			On("brand.id = product.brand_id").
			LeftJoin("supplier").
			On("supplier.id = product.supplier_id").
			LeftJoin("category").
			On("category.id = product.cat_num_id").
			LeftJoin("user").
			On("user.id = product.user_id").
			LeftJoin("store").
			On("store.id = product.store_id").
			OrderBy("in_time").Desc().
			Limit(ProductLimit)
		sql := qb.String()
		o.Raw(sql).QueryRows(&p)
	}

	view_product_store := !permission.GetOneItemPermission(username, "ViewProductStore")
	view_stock := !permission.GetOneItemPermission(username, "ViewStock")
	view_in_price := !permission.GetOneItemPermission(username, "ViewInPrice")

	if view_product_store || view_stock || view_in_price {
		length := len(p)
		for index := 0; index < length; index++ {
			if view_product_store {
				p[index].Pool = "**"
				p[index].StoreName = "**"
			}
			if view_stock {
				p[index].Stock = "***"
			}
			if view_in_price {
				p[index].InPrice = "***"
			}
		}

	}

	product_byte, _ := json.Marshal(p)
	c.Data["product"] = string(product_byte)
	c.Data["xsrf_token"] = c.XSRFToken()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "product/product_list.html"
}

//ajax加载更多商品
func (c *ProductController) ProductLoadMore() {
	if !c.IsAjax() {
		return
	}

	username := c.GetSession("username").(string)
	user := models.User{}
	o := orm.NewOrm()
	o.QueryTable("user").Filter("username", username).One(&user, "position", "pool_name")
	operate_other_store := !permission.GetOneItemPermission(username, "OperateOtherStore")

	offset, _ := c.GetInt("offset")
	qb, _ := orm.NewQueryBuilder("mysql")
	p := []product{}
	if user.Position != "超级管理员" || operate_other_store {
		if user.PoolName != "" {
			qb.Select("product.id", "product.title", "product.art_num", "product.lot_num", "product.spec", "product.stock", "product.unit",
				"product.in_time", "product.in_price", "product.has_pay", "product.has_invoice", "product.get_invoice",
				"brand.name as brand_name", "supplier.name as supplier_name", "category.three_stage", "user.name as user_name",
				"store.pool", "store.name as store_name").
				From("product").
				LeftJoin("brand").
				On("brand.id = product.brand_id").
				LeftJoin("supplier").
				On("supplier.id = product.supplier_id").
				LeftJoin("category").
				On("category.id = product.cat_num_id").
				LeftJoin("user").
				On("user.id = product.user_id").
				InnerJoin("store").
				On("store.id = product.store_id AND store.name = ?").
				OrderBy("in_time").Desc().
				Limit(ProductLimit).
				Offset(offset * ProductLimit)
			sql := qb.String()
			o.Raw(sql, user.PoolName).QueryRows(&p)
		}
	} else {
		qb.Select("product.id", "product.title", "product.art_num", "product.lot_num", "product.spec", "product.stock", "product.unit",
			"product.in_time", "product.in_price", "product.has_pay", "product.has_invoice", "product.get_invoice",
			"brand.name as brand_name", "supplier.name as supplier_name", "category.three_stage", "user.name as user_name",
			"store.pool", "store.name as store_name").
			From("product").
			LeftJoin("brand").
			On("brand.id = product.brand_id").
			LeftJoin("supplier").
			On("supplier.id = product.supplier_id").
			LeftJoin("category").
			On("category.id = product.cat_num_id").
			LeftJoin("user").
			On("user.id = product.user_id").
			LeftJoin("store").
			On("store.id = product.store_id").
			OrderBy("in_time").Desc().
			Limit(ProductLimit).
			Offset(offset * ProductLimit)
		sql := qb.String()
		o.Raw(sql).QueryRows(&p)
	}

	view_product_store := !permission.GetOneItemPermission(username, "ViewProductStore")
	view_stock := !permission.GetOneItemPermission(username, "ViewStock")
	view_in_price := !permission.GetOneItemPermission(username, "ViewInPrice")
	if view_product_store || view_stock || view_in_price {
		length := len(p)
		for index := 0; index < length; index++ {
			if view_product_store {
				p[index].Pool = "**"
				p[index].StoreName = "**"
			}
			if view_stock {
				p[index].Stock = "***"
			}
			if view_in_price {
				p[index].InPrice = "***"
			}
		}
	}

	product_byte, _ := json.Marshal(p)
	c.Data["json"] = string(product_byte)
	c.ServeJSON()
}

//删除单条商品信息
func (c *ProductController) Product_item_delete() {
	if c.IsAjax() {
		//判断当前用户是否会有删除商品的权限
		if !permission.GetOneItemPermission(c.GetSession("username").(string), "DeleteProduct") {
			c.Abort("401")
		}
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
	//判断当前用户时候有权限
	un := c.GetSession("username").(string)
	if !permission.GetOneItemPermission(un, "EditProduct") {
		c.Abort("401")
	}

	if position.GetOnePosition(un) != "超级管理员" {
		if !permission.GetOneItemPermission(un, "OperateOtherStore") {
			user := models.User{}
			o := orm.NewOrm()
			o.QueryTable("user").Filter("username", un).One(&user, "pool_name")
			if user.PoolName != c.GetString("store") {
				c.Abort("401")
			}
		}
	}

	product := models.Product{}
	o := orm.NewOrm()

	product.Id, _ = c.GetInt("product_id")
	product.Title = c.GetString("title")

	product.ArtNum = c.GetString("atr_num")
	product.LotNum = c.GetString("lot_num")

	product.Spec = c.GetString("spec")
	product.Stock, _ = c.GetUint32("stock")
	product.InPrice, _ = c.GetFloat("in_price")
	product.Unit = c.GetString("unit")

	product.Brand = GetBrand(c.GetString("brand"))

	product.CatNum = GetCategory(c.GetString("three_stage"))

	product.Supplier = GetSupplier(c.GetString("supplier"))

	store_string := strings.Split(c.GetString("store"), "-")
	store := models.Store{}
	o.QueryTable("store").Filter("pool", store_string[0]).Filter("name", store_string[1]).One(&store, "id")
	product.Store = &store

	product.HasPay, _ = c.GetBool("has_pay_edit")
	product.HasInvoice, _ = c.GetBool("has_invioce_edit")

	product.GetInvoice, _ = time.Parse("2006-1-2", c.GetString("get_invioce_edit"))

	num, err := o.Update(&product, "title", "brand_id", "art_num", "lot_num", "cat_num_id", "spec", "stock", "unit", "store_id", "supplier_id", "in_price", "has_pay", "has_invoice", "get_invoice")
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
	//判断当前用户时候有权限
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "AddProduct") {
		c.Abort("401")
	}

	o := orm.NewOrm()
	user := models.User{}
	user.Id, _ = c.GetSession("uid").(int)
	o.QueryTable("user").Filter("id", user.Id).One(&user, "pool_name")

	c.Data["art_num_string"] = GetArtNumList()
	c.Data["store_string"] = user.PoolName + ","
	c.Layout = "common.tpl"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "product/product_add.html"
}

//商品添加post
func (c *ProductController) Add_post() {
	//判断当前用户时候有权限
	un := c.GetSession("username").(string)
	if !permission.GetOneItemPermission(un, "AddProduct") {
		c.Abort("401")
	}

	if position.GetOnePosition(un) != "超级管理员" {
		if !permission.GetOneItemPermission(un, "OperateOtherStore") {
			user := models.User{}
			o := orm.NewOrm()
			o.QueryTable("user").Filter("username", un).One(&user, "pool_name")
			if user.PoolName != c.GetString("store") {
				c.Abort("401")
			}
		}
	}

	product := models.Product{}
	user := models.User{}
	user.Id = c.GetSession("uid").(int)
	product.User = &user
	product.Title = c.GetString("title")
	product.ArtNum = c.GetString("atr_num")
	product.LotNum = c.GetString("lot_num")
	product.Unit = c.GetString("unit")
	product.HasPay, _ = c.GetBool("has_pay")

	o := orm.NewOrm()

	product.Brand = GetBrand(c.GetString("brand"))

	product.CatNum = GetCategory(c.GetString("three_stage"))

	product.Supplier = GetSupplier(c.GetString("supplier"))

	store_string := c.GetString("store")
	store := models.Store{}
	o.QueryTable("store").Filter("pool", beego.AppConfig.String("defaultstore")).Filter("name", store_string).One(&store, "id")
	product.Store = &store

	spec_slice := c.GetStrings("spec")
	stock_slice := c.GetStrings("stock")
	inprice_slice := c.GetStrings("in_price")

	for index, item := range spec_slice {
		product.Spec = item
		stock_temp, _ := strconv.ParseUint(stock_slice[index], 10, 0)
		product.Stock = uint32(stock_temp)
		product.InPrice, _ = strconv.ParseFloat(inprice_slice[index], 64)
		_, err := o.Insert(&product)

		//防止出现重复主键值
		product.Id++

		if err != nil {
			c.Data["msg"] = "添加商品失败~"
			c.Data["url"] = "/product_add"
			c.TplName = "jump/error.html"
			return
		}
	}
	c.Data["msg"] = "添加商品成功~"
	c.Data["url"] = "/product_list"
	c.TplName = "jump/error.html"
}

//ajax通过货号一键填充表单
func (c *ProductController) SearchByCatnum() {
	if !c.IsAjax() {
		return
	}

	//判断当前用户时候有权限
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "AddProduct") {
		c.Abort("401")
	}

	type ajax_template struct {
		Title      string
		BrandName  string
		ThreeStage string
		Spec       string
		Unit       string
		Suppliers  string
		InPrice    float64
	}

	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("product_template.title", "product_template.spec", "product_template.unit", "product_template.suppliers",
		"product_template.in_price", "brand.name as brand_name", "category.three_stage").
		From("product_template").
		InnerJoin("brand").
		On("brand.id = product_template.brand_id").
		InnerJoin("category").
		On("category.id = product_template.cat_num_id").
		Where("art_num = ?")
	sql := qb.String()
	o := orm.NewOrm()
	p := []ajax_template{}
	o.Raw(sql, c.GetString("art_num")).QueryRows(&p)

	c.Data["json"] = p
	c.ServeJSON()
}

func (c *ProductController) Product_track() {
	username := c.GetSession("username").(string)
	if !permission.GetOneItemPermission(username, "ViewSale") {
		c.Abort("401")
	}

	user := models.User{}
	o := orm.NewOrm()
	o.QueryTable("user").Filter("username", username).One(&user, "position", "pool_name")
	operate_other_store := !permission.GetOneItemPermission(username, "OperateOtherStore")

	pid, _ := c.GetInt(":pid")
	pro := models.Product{}
	o.QueryTable("product").Filter("id", pid).One(&pro, "art_num")

	sale := []salelist{}
	qb, _ := orm.NewQueryBuilder("mysql")

	if user.Position != "超级管理员" || operate_other_store {
		if user.PoolName != "" {
			qb.Select("sale.id", "sale.out_price", "sale.num", "sale.send", "sale.has_invoice", "sale.invoice_num", "sale.no", "sale.comment",
				"sale.send_invoice", "sale.get_invoice", "sale.get_money", "sale.get_date", "sale.created", "store.pool", "store.name as store_name",
				"product.title", "product.art_num", "product.in_price", "brand.name as brand", "product.unit", "product.spec", "user.name as salesman_name", "consumer.name as consumer_name").
				From("sale").
				LeftJoin("product").
				On("product.id = sale.product_id AND product.art_num = ?").
				InnerJoin("brand").
				On("product.brand_id = brand.id").
				LeftJoin("user").
				On("user.id = sale.salesman_id").
				LeftJoin("consumer").
				On("consumer.id = sale.consumer_id").
				InnerJoin("store").
				On("store.id = sale.store_id AND store.name = ?").
				OrderBy("created").Desc()

			sql := qb.String()
			o.Raw(sql, pro.ArtNum, user.PoolName, pid).QueryRows(&sale)
		}
	} else {
		qb.Select("sale.id", "sale.out_price", "sale.num", "sale.send", "sale.has_invoice", "sale.invoice_num", "sale.no", "sale.comment",
			"sale.send_invoice", "sale.get_invoice", "sale.get_money", "sale.get_date", "sale.created", "store.pool", "store.name as store_name",
			"product.title", "product.art_num", "product.in_price", "brand.name as brand", "product.unit", "product.spec", "user.name as salesman_name", "consumer.name as consumer_name").
			From("sale").
			LeftJoin("product").
			On("product.id = sale.product_id AND product.art_num = ?").
			InnerJoin("brand").
			On("product.brand_id = brand.id").
			LeftJoin("user").
			On("user.id = sale.salesman_id").
			LeftJoin("consumer").
			On("consumer.id = sale.consumer_id").
			InnerJoin("store").
			On("store.id = sale.store_id").
			OrderBy("created").Desc()
		sql := qb.String()
		o.Raw(sql, pro.ArtNum).QueryRows(&sale)
	}

	view_consumer := !permission.GetOneItemPermission(username, "ViewSaleConsumer")
	view_in_price := !permission.GetOneItemPermission(username, "ViewSaleInPrice")
	if view_consumer || view_in_price {
		length := len(sale)
		for index := 0; index < length; index++ {
			if view_consumer {
				sale[index].ConsumerName = "***"
			}
			if view_in_price {
				sale[index].InPrice = "***"
			}
		}
	}
	c.Data["sale"] = sale
	c.Layout = "common.tpl"
	c.TplName = "product/sale_info.html"
}

//管理员添加商品模板
func (c *ProductController) ProductTemplateList() {
	pos := position.GetOnePosition(c.GetSession("username").(string))
	if pos != "超级管理员" && pos != "总库管理员" {
		c.Abort("401")
	}

	pt := []product_template{}

	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("product_template.id", "product_template.title", "product_template.art_num",
		"product_template.spec", "product_template.unit", "product_template.suppliers",
		"product_template.in_price", "brand.name as brand_name", "category.three_stage",
		"dealer.name as dealer_name").
		From("product_template").
		InnerJoin("brand").
		On("brand.id = product_template.brand_id").
		InnerJoin("category").
		On("category.id = product_template.cat_num_id").
		LeftJoin("dealer").
		On("dealer.id = product_template.dealer_id").
		OrderBy("id").Desc().
		Limit(TemplateLimit)
	sql := qb.String()
	o := orm.NewOrm()
	o.Raw(sql).QueryRows(&pt)

	template_byte, _ := json.Marshal(pt)
	c.Data["template"] = string(template_byte)

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["brand_string"] = GetBrandList()
	c.Data["supplier_string"] = GetSupplierList()
	c.Data["three_stage_string"] = GetThreeStageList()
	c.Layout = "common.tpl"
	c.TplName = "product/product_template_list.html"
}

//加载更多模板数据（ajax）
func (c *ProductController) TemplateLoadMore() {
	if !c.IsAjax() {
		return
	}
	pos := position.GetOnePosition(c.GetSession("username").(string))
	if pos != "超级管理员" && pos != "总库管理员" {
		c.Abort("401")
	}

	offset, _ := c.GetInt("offset")
	pt := []product_template{}
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("product_template.id", "product_template.title", "product_template.art_num",
		"product_template.spec", "product_template.unit", "product_template.suppliers",
		"product_template.in_price", "brand.name as brand_name", "category.three_stage",
		"dealer.name as dealer_name").
		From("product_template").
		InnerJoin("brand").
		On("brand.id = product_template.brand_id").
		InnerJoin("category").
		On("category.id = product_template.cat_num_id").
		LeftJoin("dealer").
		On("dealer.id = product_template.dealer_id").
		OrderBy("id").Desc().
		Limit(TemplateLimit).
		Offset(TemplateLimit * offset)
	sql := qb.String()
	o := orm.NewOrm()
	o.Raw(sql).QueryRows(&pt)

	template_byte, _ := json.Marshal(pt)
	c.Data["json"] = string(template_byte)
	c.ServeJSON()
}

//商品模板添加页面
func (c *ProductController) ProductTemplateAdd() {
	pos := position.GetOnePosition(c.GetSession("username").(string))
	if pos != "超级管理员" && pos != "总库管理员" {
		c.Abort("401")
	}

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["brand_string"] = GetBrandList()
	c.Data["supplier_string"] = GetSupplierList()
	c.Data["three_stage_string"] = GetThreeStageList()
	c.Layout = "common.tpl"
	c.TplName = "product/product_template_add.html"
}

//商品模板添加提交
func (c *ProductController) ProductTemplateAddPost() {
	pos := position.GetOnePosition(c.GetSession("username").(string))
	if pos != "超级管理员" && pos != "总库管理员" {
		c.Abort("401")
	}

	product_template := models.ProductTemplate{}

	product_template.Title = c.GetString("title")
	product_template.Unit = c.GetString("unit")

	o := orm.NewOrm()

	//通过货号判断，数据库中时候已经存在该类商品
	exist := o.QueryTable("product_template").Filter("art_num", product_template.ArtNum).Exist()
	if exist {
		c.Data["msg"] = "该货号的商品模板已经存在~"
		c.Data["url"] = "/product_template_add"
		c.TplName = "jump/error.html"
		return
	}

	product_template.Brand = GetBrand(c.GetString("brand"))

	product_template.CatNum = GetCategory(c.GetString("three_stage"))

	product_template.Suppliers = c.GetString("supplier")

	spec_slice := c.GetStrings("spec")
	inprice_slice := c.GetStrings("in_price")
	art_num_slice := c.GetStrings("atr_num")

	for index, item := range spec_slice {
		product_template.Spec = item
		product_template.ArtNum = art_num_slice[index]
		if inprice_slice[index] != "" {
			product_template.InPrice, _ = strconv.ParseFloat(inprice_slice[index], 64)
		} else {
			product_template.InPrice = 0
		}
		_, err := o.Insert(&product_template)

		//防止出现重复主键值
		product_template.Id++

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
	pos := position.GetOnePosition(c.GetSession("username").(string))
	if pos != "超级管理员" && pos != "总库管理员" {
		c.Abort("401")
	}

	product_template := models.ProductTemplate{}
	product_template.Id, _ = c.GetInt("template_id")

	o := orm.NewOrm()
	var template_ids orm.ParamsList
	var ids_string string
	if c.GetString("global") == "yes" {
		temp := models.ProductTemplate{}
		o.QueryTable("product_template").Filter("id", product_template.Id).One(&temp, "art_num")
		o.QueryTable("product_template").Filter("art_num", temp.ArtNum).ValuesFlat(&template_ids, "id")

		for index, item := range template_ids {
			if item != product_template.Id {
				if index == 0 {
					ids_string += fmt.Sprintf("%d", item)
				} else {
					ids_string += "," + fmt.Sprintf("%d", item)
				}
			}
		}
		ids_string = "(" + ids_string + ")"
	}

	product_template.Title = c.GetString("title")
	product_template.ArtNum = c.GetString("atr_num")
	product_template.Unit = c.GetString("unit")

	product_template.Brand = GetBrand(c.GetString("brand"))
	product_template.CatNum = GetCategory(c.GetString("three_stage"))

	product_template.Suppliers = c.GetString("supplier_list")

	product_template.Spec = c.GetString("spec")
	if c.GetString("in_price") != "" {
		product_template.InPrice, _ = strconv.ParseFloat(c.GetString("in_price"), 64)
	} else {
		product_template.InPrice = 0
	}

	_, err := o.Update(&product_template)

	if c.GetString("global") == "yes" {
		sql := "UPDATE product_template SET title = ?, brand_id = ?, art_num = ?, cat_num_id = ?, unit = ?, suppliers = ? WHERE id in " + ids_string
		o.Raw(sql, product_template.Title, product_template.Brand.Id, product_template.ArtNum, product_template.CatNum.Id,
			product_template.Unit, product_template.Suppliers).Exec()
	}

	if err != nil {
		c.Data["msg"] = "编辑商品模板失败~"
		c.Data["url"] = "/product_template_add"
		c.TplName = "jump/error.html"
		return
	}
	c.Data["msg"] = "编辑商品模板成功~"
	c.Data["url"] = "/product_template_list"
	c.TplName = "jump/error.html"
}

//删除指定product_template
func (c *ProductController) ProductTemplateDeletePost() {
	pos := position.GetOnePosition(c.GetSession("username").(string))
	if pos != "超级管理员" && pos != "总库管理员" {
		c.Abort("401")
	}

	if c.IsAjax() {
		template := models.ProductTemplate{}
		template.Id, _ = c.GetInt("pid")
		o := orm.NewOrm()
		_, err := o.Delete(&template)
		if err != nil {
			c.Data["json"] = ResponseInfo{
				Code:    "failed",
				Message: "删除失败",
			}
			c.ServeJSON()
		}
		c.Data["json"] = ResponseInfo{
			Code:    "success",
			Message: "删除成功",
		}
		c.ServeJSON()
	}
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

//获取管辖库房列表
func GetStoreSlice(pool_name string) []string {

	var store_slice []string

	if pool_name == "" {
		return []string{}
	}

	if strings.Contains(pool_name, "-") {
		store_slice = append(store_slice, pool_name)
	} else {
		o := orm.NewOrm()
		store := []models.Store{}
		o.QueryTable("store").Filter("pool", pool_name).All(&store, "name")

		for _, item := range store {
			store_slice = append(store_slice, pool_name+"-"+item.Name)
		}
	}

	return store_slice
}

//判断是否属于所管辖库房
func JudgeStore(store_slice []string, store string) bool {
	for _, item := range store_slice {
		if item == store {
			return true
		}
	}
	return false
}

//判断是否属于所管辖库房
func JudgeIsStore(pool_name, input string) bool {
	if pool_name == "" {
		return false
	}

	if pool_name == input {
		return true
	}

	store_slice := GetStoreSlice(pool_name)
	if strings.Contains(pool_name, "-") {
		for _, item := range store_slice {
			if input == item {
				return true
			}
		}
	} else {
		for _, item := range store_slice {
			sub := strings.Split(item, "-")
			if input == sub[0] {
				return true
			}
		}
	}
	return false
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

//根据brand name 获取 brand
func GetBrand(name string) *models.Brand {
	brand := models.Brand{}
	o := orm.NewOrm()
	o.QueryTable("brand").Filter("name", name).One(&brand, "id")
	return &brand
}

//根据dealer name 获取 dealer
func GetDealer(name string) *models.Dealer {
	dealer := models.Dealer{}
	o := orm.NewOrm()
	o.QueryTable("dealer").Filter("name", name).One(&dealer, "id")
	return &dealer
}

//根据三级分类名称获取对应
func GetCategory(three_stage string) *models.Category {
	category := models.Category{}
	o := orm.NewOrm()
	o.QueryTable("category").Filter("three_stage", three_stage).One(&category, "id")
	return &category
}

//根据供应商名称获取供应商
func GetSupplier(name string) *models.Supplier {
	supplier := models.Supplier{}
	o := orm.NewOrm()
	o.QueryTable("supplier").Filter("name", name).One(&supplier, "id")
	return &supplier
}

//根据库房字符串获取库房信息
func GetStore(name string) *models.Store {
	store := models.Store{}
	o := orm.NewOrm()
	store_slice := strings.Split(name, "-")
	o.QueryTable("store").Filter("pool", store_slice[0]).Filter("name", store_slice[1]).One(&store, "id")
	return &store
}

//获取商品模板中的货号列表
func GetArtNumList() string {
	o := orm.NewOrm()
	art_num := []models.ProductTemplate{}
	o.QueryTable("product_template").GroupBy("art_num").All(&art_num, "art_num")
	var art_num_string string
	for _, item := range art_num {
		art_num_string += item.ArtNum + ", "
	}
	return art_num_string
}
