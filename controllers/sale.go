package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ERP/models"
	"html/template"
	"time"
	"ERP/plugins/permission"
	"encoding/json"
	"strings"
)

type SaleController struct {
	beego.Controller
}

const SaleLimit = 100

type salelist struct {
	Id           int
	Title        string
	No           string
	Pool         string
	StoreName    string
	ArtNum       string
	SalesmanName string
	ConsumerName string
	InPrice      string
	OutPrice     string
	Num          string
	Send         string
	HasInvoice   bool
	InvoiceNum   string
	SendInvoice  string
	GetInvoice   string
	GetMoney     bool
	GetDate      string
	Created      string
}

//获取销售列表数据
func (c *SaleController) Sale_list() {
	username := c.GetSession("username").(string)
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "ViewSale") {
		c.Abort("401")
	}

	user := models.User{}
	o := orm.NewOrm()
	o.QueryTable("user").Filter("username", username).One(&user, "position", "pool_name")
	operate_other_store := !permission.GetOneItemPermission(username, "OperateOtherStore")

	sale := []salelist{}
	qb, _ := orm.NewQueryBuilder("mysql")

	if user.Position != "超级管理员" {
		if operate_other_store {
			if user.PoolName != "" {
				if strings.Contains(user.PoolName, "-") {
					store_slice := strings.Split(user.PoolName, "-")
					qb.Select("sale.id", "sale.out_price", "sale.num", "sale.send", "sale.has_invoice", "sale.invoice_num", "sale.no",
						"sale.send_invoice", "sale.get_invoice", "sale.get_money", "sale.get_date", "sale.created", "store.pool", "store.name as store_name",
						"product.title", "product.art_num", "product.in_price", "user.name as salesman_name", "consumer.name as consumer_name").
						From("sale").
						LeftJoin("product").
						On("product.id = sale.product_id").
						LeftJoin("user").
						On("user.id = sale.salesman_id").
						LeftJoin("consumer").
						On("consumer.id = sale.consumer_id").
						InnerJoin("store").
						On("store.id = sale.store_id AND store.pool = ? AND store.name = ?").
						OrderBy("created").Desc().
						Limit(SaleLimit)
					sql := qb.String()
					o.Raw(sql, store_slice[0], store_slice[1]).QueryRows(&sale)
				}else{
					qb.Select("sale.id", "sale.out_price", "sale.num", "sale.send", "sale.has_invoice", "sale.invoice_num", "sale.no",
						"sale.send_invoice", "sale.get_invoice", "sale.get_money", "sale.get_date", "sale.created", "store.pool", "store.name as store_name",
						"product.title", "product.art_num", "product.in_price", "user.name as salesman_name", "consumer.name as consumer_name").
						From("sale").
						LeftJoin("product").
						On("product.id = sale.product_id").
						LeftJoin("user").
						On("user.id = sale.salesman_id").
						LeftJoin("consumer").
						On("consumer.id = sale.consumer_id").
						InnerJoin("store").
						On("store.id = sale.store_id AND store.pool = ?").
						OrderBy("created").Desc().
						Limit(SaleLimit)
					sql := qb.String()
					o.Raw(sql, user.PoolName).QueryRows(&sale)
				}
			}
		}
	} else if user.Position == "超级管理员" || !operate_other_store {
		qb.Select("sale.id", "sale.out_price", "sale.num", "sale.send", "sale.has_invoice", "sale.invoice_num", "sale.no",
			"sale.send_invoice", "sale.get_invoice", "sale.get_money", "sale.get_date", "sale.created", "store.pool", "store.name as store_name",
			"product.title", "product.art_num", "product.in_price", "user.name as salesman_name", "consumer.name as consumer_name").
			From("sale").
			LeftJoin("product").
			On("product.id = sale.product_id").
			LeftJoin("user").
			On("user.id = sale.salesman_id").
			LeftJoin("consumer").
			On("consumer.id = sale.consumer_id").
			InnerJoin("store").
			On("store.id = sale.store_id").
			OrderBy("created").Desc().
			Limit(SaleLimit)
		sql := qb.String()
		o.Raw(sql).QueryRows(&sale)
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

	sale_byte, _ := json.Marshal(sale)
	c.Data["sale"] = string(sale_byte)

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "sale/sale_list.html"
}

//获取更多销售记录
func (c *SaleController) SaleLoadMore() {
	username := c.GetSession("username").(string)
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "ViewSale") {
		c.Abort("401")
	}

	user := models.User{}
	o := orm.NewOrm()
	o.QueryTable("user").Filter("username", username).One(&user, "position", "pool_name")
	operate_other_store := !permission.GetOneItemPermission(username, "OperateOtherStore")

	sale := []salelist{}
	qb, _ := orm.NewQueryBuilder("mysql")
	offset, _ := c.GetInt("offset")

	if user.Position != "超级管理员" {
		if operate_other_store {
			if user.PoolName != "" {
				if strings.Contains(user.PoolName, "-") {
					store_slice := strings.Split(user.PoolName, "-")
					qb.Select("sale.id", "sale.out_price", "sale.num", "sale.send", "sale.has_invoice", "sale.invoice_num", "sale.no",
						"sale.send_invoice", "sale.get_invoice", "sale.get_money", "sale.get_date", "sale.created", "store.pool", "store.name as store_name",
						"product.title", "product.art_num", "product.in_price", "user.name as salesman_name", "consumer.name as consumer_name").
						From("sale").
						LeftJoin("product").
						On("product.id = sale.product_id").
						LeftJoin("user").
						On("user.id = sale.salesman_id").
						LeftJoin("consumer").
						On("consumer.id = sale.consumer_id").
						InnerJoin("store").
						On("store.id = sale.store_id AND store.pool = ? AND store.name = ?").
						OrderBy("created").Desc().
						Limit(SaleLimit).
						Offset(SaleLimit * offset)
					sql := qb.String()
					o.Raw(sql, store_slice[0], store_slice[1]).QueryRows(&sale)
				}else{
					qb.Select("sale.id", "sale.out_price", "sale.num", "sale.send", "sale.has_invoice", "sale.invoice_num", "sale.no",
						"sale.send_invoice", "sale.get_invoice", "sale.get_money", "sale.get_date", "sale.created", "store.pool", "store.name as store_name",
						"product.title", "product.art_num", "product.in_price", "user.name as salesman_name", "consumer.name as consumer_name").
						From("sale").
						LeftJoin("product").
						On("product.id = sale.product_id").
						LeftJoin("user").
						On("user.id = sale.salesman_id").
						LeftJoin("consumer").
						On("consumer.id = sale.consumer_id").
						InnerJoin("store").
						On("store.id = sale.store_id AND store.pool = ?").
						OrderBy("created").Desc().
						Limit(SaleLimit).
						Offset(SaleLimit * offset)

					sql := qb.String()
					o.Raw(sql, user.PoolName).QueryRows(&sale)
				}
			}
		}
	} else if user.Position == "超级管理员" || !operate_other_store {
		qb.Select("sale.id", "sale.out_price", "sale.num", "sale.send", "sale.has_invoice", "sale.invoice_num", "sale.no",
			"sale.send_invoice", "sale.get_invoice", "sale.get_money", "sale.get_date", "sale.created", "store.pool", "store.name as store_name",
			"product.title", "product.art_num", "product.in_price", "user.name as salesman_name", "consumer.name as consumer_name").
			From("sale").
			LeftJoin("product").
			On("product.id = sale.product_id").
			LeftJoin("user").
			On("user.id = sale.salesman_id").
			LeftJoin("consumer").
			On("consumer.id = sale.consumer_id").
			InnerJoin("store").
			On("store.id = sale.store_id").
			OrderBy("created").Desc().
			Limit(SaleLimit).
			Offset(SaleLimit * offset)
		sql := qb.String()
		o.Raw(sql).QueryRows(&sale)
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

	sale_byte, _ := json.Marshal(sale)
	c.Data["json"] = string(sale_byte)
	c.ServeJSON()
}

//单条销售记录修改post
func (c *SaleController) Sale_edit() {
	un := c.GetSession("username").(string)
	if !permission.GetOneItemPermission(un, "EditSale") {
		c.Abort("401")
	}

	sale := models.Sale{}
	salesman := models.User{}

	o := orm.NewOrm()
	o.QueryTable("user").Filter("name", c.GetString("salesman")).One(&salesman)
	sale.Salesman = &salesman

	sale.Id, _ = c.GetInt("sale_id")
	sale.OutPrice, _ = c.GetFloat("outprice")
	sale.Num, _ = c.GetUint32("num")
	sale.Send, _ = time.Parse("2006-1-2", c.GetString("send"))
	sale.HasInvoice, _ = c.GetBool("hasinvoice")
	sale.InvoiceNum = c.GetString("invioce_num")
	sale.SendInvoice, _ = time.Parse("2006-1-2", c.GetString("send_invioce"))
	sale.GetInvoice, _ = time.Parse("2006-1-2", c.GetString("get_invioce"))
	sale.GetMoney, _ = c.GetBool("get_money")
	sale.GetDate, _ = time.Parse("2006-1-2", c.GetString("get_date"))

	_, err := o.Update(&sale, "salesman", "out_price", "num", "send", "has_invoice", "invoice_num", "send_invoice", "get_invoice", "get_money", "get_date")
	if err == nil {
		c.Data["url"] = "/sale_list"
		c.Data["msg"] = "修改销售记录成功"
		c.TplName = "jump/success.html"
	}
}
