package controllers

import (
	"html/template"

	"erp/models"
	"erp/plugins/permission"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type DealerController struct {
	beego.Controller
}

func (c *DealerController) Get() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "ViewDealer") {
		c.Abort("401")
	}
	dealer := []models.Dealer{}
	o := orm.NewOrm()
	o.QueryTable("dealer").All(&dealer)

	c.Data["dealer"] = dealer
	c.Layout = "common.tpl"
	c.TplName = "dealer/dealer_list.html"
}

//添加经销商页面
func (c *DealerController) Dealer_add() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "AddDealer") {
		c.Abort("401")
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "dealer/dealer_add.html"
}

//添加经销商 post提交
func (c *DealerController) Dealer_add_post() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "AddDealer") {
		c.Abort("401")
	}
	dealer := models.Dealer{}
	dealer.Name = c.GetString("name")

	o := orm.NewOrm()
	exit := o.QueryTable("dealer").Filter("name", dealer.Name).Exist()
	if exit {
		c.Data["msg"] = "此经销商已经存在，请勿重复添加~"
		c.Data["url"] = "/dealer_add"
		c.TplName = "jump/error.html"
		return
	} else {
		_, err := o.Insert(&dealer)
		if err != nil {
			logs.Error(c.GetSession("uid"), "添加经销商错误：", err)
			c.Data["msg"] = "添加失败，请稍后重试或联系管理员~"
			c.Data["url"] = "/dealer_add"
			c.TplName = "jump/error.html"
			return
		} else {
			c.Data["msg"] = "添加经销商 " + dealer.Name + " 成功~"
			c.Data["url"] = "/dealer_list"
			c.TplName = "jump/success.html"
			return
		}
	}
}
