package controllers

import (
	"github.com/astaxie/beego"
	"ERP/models"
	"github.com/astaxie/beego/orm"
	"html/template"
	"strings"
	"time"
	"github.com/astaxie/beego/logs"
	"strconv"
	"fmt"
	"ERP/plugins/permission"
)

type StoreoutController struct {
	beego.Controller
}

//渲染出库页面
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

//出库post
func (c *StoreoutController) Store_out_action_post() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "OutputProduct") {
		c.Abort("401")
	}

	sale := models.Sale{}
	o := orm.NewOrm()

	product := models.Product{}
	pid, _ := c.GetInt("product_id")
	o.QueryTable("product").Filter("id", pid).One(&product)
	sale.Product = &product

	//判断当前销售数量时候多于库存
	sale.Num, _ = c.GetUint32("num")
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
	sale.HasInvoice, _ = c.GetBool("hasinvoice")
	sale.SendInvoice, _ = time.Parse("2006-1-2", c.GetString("send_invioce"))
	sale.GetInvoice, _ = time.Parse("2006-1-2", c.GetString("get_invioce"))
	sale.GetDate, _ = time.Parse("2006-1-2", c.GetString("get_date"))
	sale.InvoiceNum = c.GetString("invioce_num")
	sale.GetMoney, _ = c.GetBool("get_money")

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
		"stock" : orm.ColValue(orm.ColMinus, sale.Num),
	})
	if err1 != nil || err2 != nil{
		o.Rollback()
		logs.Error(err1, err2)
		return
	} else {
		o.Commit()
		c.Data["url"]= "/sale_list"
		c.Data["msg"] = "出库成功"
		c.TplName = "jump/success.html"
	}
}
