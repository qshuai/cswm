package controllers

import (
	"github.com/astaxie/beego"
	"ERP/models"
	"github.com/astaxie/beego/orm"
	"html/template"
	"github.com/astaxie/beego/logs"
)

type StoreController struct {
	beego.Controller
}

//获取库房列表
func (c *StoreController) Get(){
	store := []models.Store{}
	o := orm.NewOrm()
	o.QueryTable("store").OrderBy("pool").All(&store)

	//计算S库和J库数量
	var s, j int
	for _, item := range store {
		switch item.Pool{
		case "S库":
			s++
		case "J库":
			j++
		}
	}

	c.Data["s"] = s
	c.Data["j"] = j
	c.Data["store"] = store
	c.Layout = "common.tpl"
	c.TplName = "store/store_list.html"
}

//添加库房页面
func (c *StoreController) Store_add() {
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "store/store_add.html"
}

//添加库房页面post
func (c *StoreController) Store_add_post() {
	store := models.Store{}
	o := orm.NewOrm()

	store.Pool = c.GetString("pool")
	store.Name = c.GetString("name")
	_, err := o.Insert(&store)
	if err != nil {
		logs.Error("用户Id：", c.GetSession("uid").(int), "添加库房失败： ", err)
		c.Data["url"] = "/store_add"
		c.Data["msg"] = "添加库房失败~"
		c.TplName = "jump/error.html"
		return
	}
	c.Data["url"] = "/store_list"
	c.Data["msg"] = "添加库房 "+ store.Name +" 成功~"
	c.TplName = "jump/error.html"
}
