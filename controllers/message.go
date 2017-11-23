package controllers

import (
	"html/template"

	"erp/models"
	"erp/plugins/message"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type MessageController struct {
	beego.Controller
}

//message列表
func (c *MessageController) Message_list() {
	current_uid := c.GetSession("uid")
	message := []models.Message{}

	cond := orm.NewCondition()
	cond_build := cond.And("from", current_uid).Or("to", current_uid)

	o := orm.NewOrm()
	o.QueryTable("message").SetCond(cond_build).RelatedSel().OrderBy("-created").All(&message)

	c.Data["message"] = message
	c.Data["xsrftoken"] = c.XSRFToken()
	c.Layout = "common.tpl"
	c.TplName = "message/message_list.html"
}

//get请求单个message详情
func (c *MessageController) Message_info() {
	mid, _ := c.GetInt(":mid")
	o := orm.NewOrm()
	message := models.Message{}

	o.QueryTable("message").Filter("id", mid).RelatedSel().One(&message)

	//不能自己发给自己消息（回复）
	current_uid := c.GetSession("uid").(int)
	if current_uid == message.From.Id {
		c.Data["is_self"] = true
	} else {
		c.Data["is_self"] = false

		//更新消息表中的IsRead字段，设置为true
		message.IsRead = true
		o.Update(&message, "is_read")

		//redis
		msg.DecrOneMessage(c.GetSession("username").(string))
	}

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["message"] = message
	c.Layout = "common.tpl"
	c.TplName = "message/message_info.html"
}

//新建message页面
func (c *MessageController) Message_add() {
	user := []models.User{}
	current_user_id := c.GetSession("uid").(int)

	o := orm.NewOrm()
	o.QueryTable("user").Exclude("id", current_user_id).All(&user, "id", "name")

	var user_string string
	for _, item := range user {
		user_string += item.Name + ", "
	}

	c.Data["user_string"] = user_string

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "message/message_add.html"
}

//新建message提交post
func (c *MessageController) Message_add_post() {
	message_from := c.GetSession("uid").(int)
	user_from := models.User{}
	user_to := models.User{}

	o := orm.NewOrm()
	o.QueryTable("user").Filter("id", message_from).One(&user_from)
	o.QueryTable("user").Filter("name", c.GetString("message_to")).One(&user_to)

	message := models.Message{}
	message.From = &user_from
	message.To = &user_to

	message.Content = c.GetString("message_content")

	_, err := o.Insert(&message)

	//redis
	msg.IncrOneMessage(user_to.Username)
	if err != nil {
		logs.Error("/message_add ", user_from.Name, " 向 ", user_to.Name, " 发送信息失败：", err)
		c.Data["url"] = "/message_add"
		c.Data["msg"] = "发送消息失败~"
		c.TplName = "jump/error.html"
	} else {
		c.Data["url"] = "/message_list"
		c.Data["msg"] = "发送消息成功~"
		c.TplName = "jump/success.html"
	}
}
