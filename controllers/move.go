package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ERP/models"
	"html/template"
	"strconv"
	"strings"
	"time"
	"github.com/astaxie/beego/logs"
	"ERP/plugins/permission"
)

type MoveController struct{
	beego.Controller
}

//移库请求页面
func (c *MoveController) Move_request() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "RequestMove") {
		c.Abort("401")
	}
	product := models.Product{}
	o := orm.NewOrm()

	pid, _ := c.GetInt(":pid")
	o.QueryTable("product").Filter("id", pid).RelatedSel("store").One(&product)

	store_at := product.Store.Pool + "-" + product.Store.Name

	//获取库房列表
	store := []models.Store{}
	o.QueryTable("store").All(&store, "pool", "name")
	var store_string  string
	for _, item := range store{
		store_string += item.Pool + "-" + item.Name + ", "
	}

	c.Data["store_num"] = product.Store.Id
	c.Data["store_at"] = store_at
	c.Data["store_string"] = store_string
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["product"] = product
	c.Layout = "common.tpl"
	c.TplName = "move/move_request.html"
}

//移库post提交
func (c *MoveController) Move_request_post() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "RequestMove") {
		c.Abort("401")
	}
	pid, _ := c.GetInt("product_id")
	o := orm.NewOrm()
	product := models.Product{}
	o.QueryTable("product").Filter("id", pid).One(&product)

	num, _ := c.GetUint32("num")
	if num > product.Stock {
		c.Data["url"] = "/move_request/" + strconv.Itoa(pid)
		c.Data["msg"] = "移库数量超过当前商品库存"
		c.TplName = "jump/error.html"
		return
	}

	move := models.Move{}
	move.Num = num
	move.Origin = &product

	//起始库房
	store_from := models.Store{}
	store_from.Id, _ = c.GetInt("store_num")
	o.QueryTable("store").Filter("Id", store_from.Id).One(&store_from)
	move.From = &store_from

	//目标库房
	store_to := models.Store{}
	store_slice := strings.Split(c.GetString("move_to"), "-")
	o.QueryTable("store").Filter("pool", store_slice[0]).Filter("name", store_slice[1]).One(&store_to)
	move.To = &store_to

	//发起人
	user := models.User{}
	uid := c.GetSession("uid").(int)
	o.QueryTable("user").Filter("id", uid).One(&user)
	move.Request = &user

	//响应人
	pool_user := []models.User{}
	conf := orm.NewCondition()
	con := conf.Or("pool_name", "S库").Or("pool_name", "J库").Or("pool_name", c.GetString("move_to")).Or("pool_name", c.GetString("store_from"))
	o.QueryTable("user").SetCond(con).All(&pool_user)
	//区分总库和分库管理员
	for _, item := range pool_user {
		//此处有一定的隐患
		if len(item.PoolName) > 2 && item.PoolName != c.GetString("store_from") {
			move.Response = &item
		}
	}

	//初始化操作
	move.Operate = "0"

	//事务
	o.Begin()
	mid, err1 := o.Insert(&move)
	_, err2 := o.QueryTable("product").Filter("id", pid).Update(orm.Params{
		"stock" : orm.ColValue(orm.ColMinus, move.Num),
	})
	if err1 == nil && err2 == nil{
		o.Commit()
		//为目标库房总库管理员和分库管理员发送message
		message := models.Message{}
		for _, item := range pool_user  {
			if item.PoolName != user.PoolName {
				message.From = &user
				message.To = &item
				message.Content = "<span class='c-warning'>"+user.Name + "</span> 将商品： <a href='/product_track/"+ strconv.Itoa(product.Id) +"' class='c-primary'>" + product.Title + "</a> 从 <span class='c-danger'>" + store_from.Pool + "-" + store_from.Name +
					"</span> 移库到 <span class='c-danger'>" + c.GetString("move_to") + "</span><br />具体请查看：<a href='"+ c.Ctx.Input.Site() + ":" +strconv.Itoa(c.Ctx.Input.Port()) +
					"/move_info/" + strconv.FormatInt(mid, 10) + "' target='blank'><u>移库详情</u></a>"
				o.Insert(&message)
			}
		}

		c.Data["url"] = "/move_list"
		c.Data["msg"] = "成功发起移库请求！"
		c.TplName = "jump/success.html"
	} else {
		o.Rollback()
		logs.Error("用户Id： ", uid, "移库失败！ 商品id: ", pid, "; 从库房Id：", product.Stock, " 到 ", c.GetString("move_to"), "原因:", err1, err2)
		c.Data["url"] = "/move_request/" + strconv.Itoa(pid)
		c.Data["msg"] = "发起移库请求失败"
		c.TplName = "jump/error.html"
	}
}

//移库列表
func (c *MoveController) Move_list() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "ViewMove") {
		c.Abort("401")
	}
	uid := c.GetSession("uid").(int)
	user := models.User{}
	o := orm.NewOrm()
	o.QueryTable("user").Filter("id", uid).One(&user)

	move := []models.Move{}
	conf := orm.NewCondition()
	con := conf.Or("request_id", uid).Or("response_id", uid).Or("From__pool", user.PoolName)
	o.QueryTable("move").SetCond(con).RelatedSel().All(&move)

	c.Data["xsrf_token"] = c.XSRFToken()
	c.Data["move"] = move
	c.Layout = "common.tpl"
	c.TplName = "move/move_list.html"
}

//移库接受
func (c *MoveController) Move_accept(){
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "ResponseMove") {
		c.Abort("401")
	}
	if c.IsAjax() {
		o := orm.NewOrm()
		move := models.Move{}

		move.Id, _ = c.GetInt("mid")
		move.Operate = "1"
		move.OperatedTime = time.Now()
		o.Update(&move, "operate", "operated_time")

		c.Data["json"] = ResponseInfo{
			Code : "success",
			Message: "您已同意此次移库操作",
			Data: "",
		}
		c.ServeJSON()
	}
}

//拒绝移库
func (c *MoveController) Move_deny(){
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "ResponseMove") {
		c.Abort("401")
	}
	if c.IsAjax() {
		o := orm.NewOrm()
		move := models.Move{}

		move.Id, _ = c.GetInt("mid")
		move.Operate = "-1"
		move.OperatedTime = time.Now()
		o.Update(&move, "operate", "operated_time")

		c.Data["json"] = ResponseInfo{
			Code : "success",
			Message: "您已拒绝此次移库操作",
			Data: "",
		}
		c.ServeJSON()
	}
}

//移库完成
func (c *MoveController) Move_finish() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "ResponseMove") {
		c.Abort("401")
	}
	if c.IsAjax() {
		o := orm.NewOrm()
		move := models.Move{}

		move.Id, _ = c.GetInt("mid")
		move.Operate = "2"
		move.OperatedTime = time.Now()
		move.Finished = time.Now()
		o.Update(&move, "operate", "operated_time", "finished")

		c.Data["json"] = ResponseInfo{
			Code : "success",
			Message: "您已拒绝此次移库操作",
			Data: "",
		}
		c.ServeJSON()
	}
}

//移库详情页面
func (c *MoveController) Move_info() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "ViewMove") {
		c.Abort("401")
	}
	mid, _ := c.GetInt(":mid")

	o := orm.NewOrm()
	move := models.Move{}
	move.Id = mid
	o.QueryTable("move").Filter("id", mid).RelatedSel().One(&move)

	c.Data["move"] = move
	c.Layout = "common.tpl"
	c.TplName = "move/move_info.html"
}
