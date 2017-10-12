package controllers

import (
	"html/template"

	"github.com/astaxie/beego"
	"ERP/models"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)

type SupplierController struct{
	beego.Controller
}

func (c *SupplierController) Get(){
	supplier := []models.Supplier{}
	o := orm.NewOrm()
	o.QueryTable("supplier").All(&supplier)

	c.Data["supplier"] = supplier
	c.Layout = "common.tpl"
	c.TplName = "supplier/supplier_list.html"
}

//添加供应商页面
func (c *SupplierController) Supplier_add(){
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "supplier/supplier_add.html"
}

//添加供应商 post提交
func (c *SupplierController) Supplier_add_post(){
	supplier := models.Supplier{}
	supplier.Name = c.GetString("name")

	o := orm.NewOrm()
	exit := o.QueryTable("supplier").Filter("name", supplier.Name).Exist()
	if exit {
		c.Data["msg"] = "此供应商已经存在，请勿重复添加~"
		c.Data["url"] = "/supplier_add"
		c.TplName = "jump/error.html"
		return
	} else {
		_, err := o.Insert(&supplier)
		if err != nil {
			logs.Error("用户ID:", c.GetSession("uid"), "添加供应商错误：", err)
			c.Data["msg"] = "添加失败，请稍后重试或联系管理员~"
			c.Data["url"] = "/supplier_add"
			c.TplName = "jump/error.html"
			return
		} else {
			c.Data["msg"] = "添加供应商 " + supplier.Name + " 成功~"
			c.Data["url"] = "/supplier_list"
			c.TplName = "jump/success.html"
			return
		}
	}
}
