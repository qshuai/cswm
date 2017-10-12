package controllers

import (
	"github.com/astaxie/beego"
	"ERP/models"
	"github.com/astaxie/beego/orm"
	"html/template"
	"github.com/astaxie/beego/logs"
)

type BrandController struct {
	beego.Controller
}

//商标列表页面
func (c *BrandController) Get(){
	brand := []models.Brand{}
	o := orm.NewOrm()
	o.QueryTable("brand").All(&brand)

	c.Data["brand"] = brand
	c.Layout = "common.tpl"
	c.TplName = "brand/brand_list.html"
}

//添加商标页面
func (c *BrandController) Brand_add(){
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "brand/brand_add.html"
}

//添加商标 post提交
func (c *BrandController) Brand_add_post(){
	brand := models.Brand{}
	brand.Name = c.GetString("name")

	o := orm.NewOrm()
	exit := o.QueryTable("brand").Filter("name", brand.Name).Exist()
	if exit {
		c.Data["msg"] = "此品牌名称已经存在，请勿重复添加~"
		c.Data["url"] = "/brand_add"
		c.TplName = "jump/error.html"
		return
	} else {
		_, err := o.Insert(&brand)
		if err != nil {
			logs.Error(c.GetSession("uid"), "添加品牌错误：", err)
			c.Data["msg"] = "添加失败，请稍后重试或联系管理员~"
			c.Data["url"] = "/brand_add"
			c.TplName = "jump/error.html"
			return
		} else {
			c.Data["msg"] = "添加品牌 " + brand.Name + " 成功~"
			c.Data["url"] = "/brand_list"
			c.TplName = "jump/success.html"
			return
		}
	}
}