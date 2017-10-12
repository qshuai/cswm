package controllers

import (
	"html/template"

	"github.com/astaxie/beego"
	"ERP/models"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)

type ConsumerController struct {
	beego.Controller
}

//获取客户裂变
func (c *ConsumerController) Get(){
	consumer := []models.Consumer{}
	o := orm.NewOrm()
	o.QueryTable("consumer").All(&consumer)

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["consumer"] = consumer
	c.Layout = "common.tpl"
	c.TplName = "consumer/consumer_list.html"
}

func (c *ConsumerController) Consumer_add(){
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "consumer/consumer_add.html"
}

func (c *ConsumerController) Consumer_add_post(){
	consumer := models.Consumer{}
	consumer.Name = c.GetString("name")
	consumer.Tel = c.GetString("tel")
	consumer.Department = c.GetString("department")
	consumer.Province = c.GetString("province")
	consumer.City = c.GetString("city")
	consumer.Region = c.GetString("region")
	consumer.Introduction = c.GetString("introduction")

	o := orm.NewOrm()
	exit := o.QueryTable("consumer").Filter("tel", consumer.Tel).Exist()
	if exit {
		c.Data["msg"] = "此手机号码已经存在，请勿重复添加~"
		c.Data["url"] = "/consumer_add"
		c.TplName = "jump/error.html"
		return
	} else {
		_, err := o.Insert(&consumer)
		if err != nil {
			logs.Error(c.GetSession("uid"), "添加客户信息错误：", err)
			c.Data["msg"] = "添加失败，请稍后重试或联系管理员~"
			c.Data["url"] = "/consumer_add"
			c.TplName = "jump/error.html"
			return
		} else {
			c.Data["msg"] = "添加客户 " + consumer.Name + " 成功~"
			c.Data["url"] = "/consumer_list"
			c.TplName = "jump/success.html"
			return
		}
	}
}

func (c *ConsumerController) Consumer_edit() {
	consumer := models.Consumer{}
	consumer.Id, _ = c.GetInt("consumer_id")
	consumer.Tel = c.GetString("tel")
	consumer.Department = c.GetString("department")
	consumer.Province = c.GetString("province")
	consumer.City = c.GetString("city")
	consumer.Region = c.GetString("region")
	consumer.Introduction = c.GetString("introduction")

	o := orm.NewOrm()
	_, err := o.Update(&consumer, "tel", "department", "province", "city", "region", "introduction")
	if err != nil{
		logs.Error("/consumer_edit: Id-", c.GetSession("uid").(int), " 更新客户Id-", consumer.Id, "失败")
		c.Data["url"] = "/consumer_list"
		c.Data["msg"] = "更新客户信息十八里~"
		c.TplName = "jump/error.html"
	} else{
		c.Data["url"] = "/consumer_list"
		c.Data["msg"] = "更新客户信息成功~"
		c.TplName = "jump/success.html"
	}
}