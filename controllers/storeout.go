package controllers

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/qshuai/cswm/models"
	"github.com/qshuai/cswm/plugins/permission"
)

type StoreoutController struct {
	beego.Controller
}

// 渲染出库页面
func (c *StoreoutController) Store_out_action() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "OutputProduct") {
		c.Abort("401")
	}

	o := orm.NewOrm()
	product := models.Product{}
	pid, _ := c.GetInt(":pid")

	o.QueryTable("product").Filter("id", pid).RelatedSel().One(&product)

	//获取业务员列表
	salesman := []models.User{}
	o.QueryTable("user").Filter("position", "业务员").Filter("is_active", true).All(&salesman, "name")
	var salesman_string string
	for _, item := range salesman {
		salesman_string += item.Name + ", "
	}

	//获取客户列表，格式: "姓名-电话"
	consumer := []models.Consumer{}
	o.QueryTable("consumer").All(&consumer)
	var consumer_string string
	for _, item := range consumer {
		consumer_string += item.Name + "-" + item.Tel + ", "
	}

	c.Data["salesman_string"] = salesman_string
	c.Data["consumer_string"] = consumer_string
	c.Data["product_item"] = product
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "storeout/store_output_action.html"
}

// 出库post
func (c *StoreoutController) Store_out_action_post() {
	username := c.GetSession("username").(string)
	if !permission.GetOneItemPermission(username, "OutputProduct") {
		c.Abort("401")
	}

	sale := models.Sale{}
	o := orm.NewOrm()

	product := models.Product{}
	pid, _ := c.GetInt("product_id")
	o.QueryTable("product").Filter("id", pid).One(&product)

	user := models.User{}
	o.QueryTable("user").Filter("username", username).One(&user, "position", "pool_name")

	qb, _ := orm.NewQueryBuilder("mysql")
	if user.Position != "超级管理员" {
		type productcheck struct {
			Id    int
			Stock int
		}
		product_check := productcheck{}
		if user.PoolName != "" {
			if strings.Contains(user.PoolName, "-") {
				store_slice := strings.Split(user.PoolName, "-")
				qb.Select("product.id").
					From("product").
					InnerJoin("store").
					On("store.id = product.store_id AND store.pool = ? AND store.name = ?").
					Where("product.art_num = ? AND product.spec = ?").
					OrderBy("in_time").
					Asc().
					Limit(1)
				sql := qb.String()
				o.Raw(sql, store_slice[0], store_slice[1], product.ArtNum, product.Spec).QueryRow(&product_check)
			} else {
				qb.Select("product.id", "product.stock").
					From("product").
					InnerJoin("store").
					On("store.id = product.store_id AND store.pool = ?").
					Where("product.art_num = ? AND product.spec = ?").
					OrderBy("in_time").
					Asc().
					Limit(1)
				sql := qb.String()
				fmt.Println(o.Raw(sql, user.PoolName, product.ArtNum, product.Spec).QueryRow(&product_check))
			}
		}
		fmt.Println(product_check)
		//检查是否有先录入的商品
		if product_check.Id != product.Id && product_check.Stock != 0 {
			c.Data["url"] = "/store_output_action/" + strconv.Itoa(product_check.Id)
			c.Data["msg"] = "对不起，此规格的商品存在更早录入的，系统已为您自动跳转~"
			c.TplName = "jump/error.html"
			return
		}
	} else {
		product_check := models.Product{}
		o.QueryTable("product").Filter("art_num", product.ArtNum).Filter("spec", product.Spec).OrderBy("in_time").One(&product_check)
		fmt.Println(product_check)
		//检查是否有先录入的商品
		if product_check.Id != product.Id && product_check.Stock != 0 {
			c.Data["url"] = "/store_output_action/" + strconv.Itoa(product_check.Id)
			c.Data["msg"] = "对不起，此规格的商品存在更早录入的，系统已为您自动跳转~"
			c.TplName = "jump/error.html"
			return
		}
	}

	sale.Product = &product
	sale.Store = product.Store

	//判断当前销售数量时候多于库存
	sale.Num, _ = c.GetUint32("num")
	sale.NumFake = sale.Num
	if sale.Num > product.Stock {
		c.Data["url"] = "/store_output_action/" + strconv.Itoa(pid)
		c.Data["msg"] = "对不起，您输入的销售数量多于当前库存数量~"
		c.TplName = "jump/error.html"
		return
	}

	//要求客户姓名中不能含有"-"
	consumer := models.Consumer{}
	o.QueryTable("consumer").Filter("tel", strings.Split(c.GetString("consumer"), "-")[1]).One(&consumer)
	sale.Consumer = &consumer

	//要求user表中的人员姓名不能重复
	salesman := models.User{}
	o.QueryTable("user").Filter("name", c.GetString("salesman")).One(&salesman)
	sale.Salesman = &salesman

	sale.Send, _ = time.Parse("2006-1-2", c.GetString("send"))

	sale.OutPrice, _ = c.GetFloat("outprice")
	sale.OutPriceFake = sale.OutPrice
	sale.HasInvoice, _ = c.GetBool("hasinvoice")
	sale.SendInvoice, _ = time.Parse("2006-1-2", c.GetString("send_invioce"))
	sale.GetInvoice, _ = time.Parse("2006-1-2", c.GetString("get_invioce"))
	sale.GetDate, _ = time.Parse("2006-1-2", c.GetString("get_date"))
	sale.InvoiceNum = c.GetString("invioce_num")
	sale.GetMoney, _ = c.GetBool("get_money")
	sale.Comment = c.GetString("comment")
	sale.HasPrint = false

	//生成唯一订单号
	date_string := time.Now().String()
	date_string = strings.Replace(date_string[:10], "-", "", -1)
	timestamp := time.Now().UnixNano()
	timestamp_string := fmt.Sprintf("%v", timestamp)
	uid := strconv.Itoa(c.GetSession("uid").(int))
	sale.No = date_string + timestamp_string + uid

	//事件处理
	o.Begin()
	_, err1 := o.Insert(&sale)
	_, err2 := o.QueryTable("product").Filter("id", pid).Update(orm.Params{
		"stock": orm.ColValue(orm.ColMinus, sale.Num),
	})
	if err1 != nil || err2 != nil {
		o.Rollback()
		logs.Error(err1, err2)
		return
	} else {
		o.Commit()
		c.Data["url"] = "/sale_list"
		c.Data["msg"] = "出库成功"
		c.TplName = "jump/success.html"
	}
}
