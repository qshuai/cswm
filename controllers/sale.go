package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ERP/models"
	"html/template"
	"time"
)

type SaleController struct {
	beego.Controller
}

//获取销售列表数据
func (c *SaleController) Sale_list()  {
	o := orm.NewOrm()
	sale := []models.Sale{}

	o.QueryTable("sale").RelatedSel().OrderBy("-created").All(&sale)

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["sale"] = sale
	c.Layout = "common.tpl"
	c.TplName = "sale/sale_list.html"
}

//单条销售记录修改post
func (c *SaleController) Sale_edit() {
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
