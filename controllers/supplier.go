package controllers

import (
	"html/template"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/qshuai/cswm/models"
	"github.com/qshuai/cswm/plugins/permission"
)

type SupplierController struct {
	beego.Controller
}

func (c *SupplierController) Get() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "ViewSupplier") {
		c.Abort("401")
	}
	supplier := []models.Supplier{}
	o := orm.NewOrm()
	o.QueryTable("supplier").All(&supplier)

	c.Data["supplier"] = supplier
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "supplier/supplier_list.html"
}

// 添加供应商页面
func (c *SupplierController) Supplier_add() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "AddSupplier") {
		c.Abort("401")
	}

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "supplier/supplier_add.html"
}

// 添加供应商 post提交
func (c *SupplierController) Supplier_add_post() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "AddSupplier") {
		c.Abort("401")
	}

	supplier := models.Supplier{}
	supplier.Name = c.GetString("name")
	supplier.Admin = c.GetString("admin")
	supplier.Tel = c.GetString("tel")
	supplier.Site = c.GetString("province") + " " + c.GetString("city") + " " + c.GetString("region") + "-" + c.GetString("introduction")

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

// 编辑供应商信息 ajax
func (c *SupplierController) Supplier_edit_post() {
	sup := models.Supplier{}
	sup.Id, _ = c.GetInt("supplier_id")
	sup.Name = c.GetString("name")
	sup.Admin = c.GetString("admin")
	sup.Tel = c.GetString("tel")
	sup.Site = c.GetString("site")
	o := orm.NewOrm()
	_, err := o.Update(&sup, "name", "site", "tel", "admin")
	if err != nil {
		c.Data["msg"] = "修改供应商失败"
		c.Data["url"] = "/supplier_list"
		c.TplName = "jump/error.html"
		return
	}
	c.Data["msg"] = "修改供应商成功~"
	c.Data["url"] = "/supplier_list"
	c.TplName = "jump/success.html"
}
