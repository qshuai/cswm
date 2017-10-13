package controllers

import (
	"github.com/astaxie/beego"
	"html/template"
	"log"
	"github.com/astaxie/beego/orm"
	"ERP/models"
	"crypto/md5"
	"fmt"
	"strconv"
	"ERP/permission"
)

type MemberController struct {
	beego.Controller
}

//添加用户页面
func (c *MemberController) Member_add() {
	c.Data["level"] = beego.AppConfig.Strings("level")
	c.Layout = "common.tpl"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "member/member_add.html"
}

//添加用户提交逻辑
func (c *MemberController) Member_add_post() {
	u := models.User{}

	//md5 crypt password
	password := []byte(beego.AppConfig.String("defaultpassword"))
	passwordMD5 := md5.Sum(password)
	u.Password = fmt.Sprintf("%x", passwordMD5)

	u.Name = c.GetString("name")
	u.Tel = c.GetString("tel")
	u.Position = c.GetString("position")
	u.IsFirst = true
	u.IsActive = true

	o := orm.NewOrm()
	_, err := o.Insert(&u)
	if err != nil {
		log.Fatal("add member", u.Name, " failure: ", err)
	}

	//初始化人员权限
	permission := models.Permission{}
	//查询响应等级人员的默认权限
	defaultPermission := models.DefaultPermission{}
	o.QueryTable("default_permission").Filter("position", u.Position).One(&defaultPermission)

	permission.User = &u
	permission.AddMember = defaultPermission.AddMember
	permission.EditMember = defaultPermission.EditMember
	permission.ActiveMember = defaultPermission.ActiveMember
	permission.AddConsumer = defaultPermission.AddConsumer
	permission.EditConsumer = defaultPermission.EditConsumer
	permission.ViewConsumer = defaultPermission.ViewConsumer
	permission.AddBrand = defaultPermission.AddBrand
	permission.AddDealer = defaultPermission.AddDealer
	permission.ViewDealer = defaultPermission.ViewDealer
	permission.AddSupplier = defaultPermission.AddSupplier
	permission.ViewSupplier = defaultPermission.ViewSupplier
	permission.AddProduct = defaultPermission.AddProduct
	permission.InputInPrice = defaultPermission.InputInPrice
	permission.ViewProductStore = defaultPermission.ViewProductStore
	permission.ViewStock = defaultPermission.ViewStock
	permission.ViewInPrice = defaultPermission.ViewInPrice
	permission.EditProduct = defaultPermission.EditProduct
	permission.DeleteProduct = defaultPermission.DeleteProduct
	permission.OutputProduct = defaultPermission.OutputProduct
	permission.ViewSale = defaultPermission.ViewSale
	permission.ViewSaleConsumer = defaultPermission.ViewSaleConsumer
	permission.ViewSaleInPrice = defaultPermission.ViewSaleInPrice
	permission.EditSale = defaultPermission.EditSale
	permission.OperateCategory = defaultPermission.OperateCategory
	permission.RequestMove = defaultPermission.RequestMove
	permission.ResponseMove = defaultPermission.ResponseMove
	permission.ViewMove = defaultPermission.ViewMove
	permission.AddStore = defaultPermission.AddStore
	permission.ViewStore = defaultPermission.ViewStore

	_, err = o.Insert(&permission)
	if err !=nil {
		log.Fatal("assign permission failure: ", u.Name, "-", err)
	}

	c.Redirect("/member_list", 302)
}

//完善用户信息页面
func (c *MemberController) UserInfo() {
	c.Layout = "common.tpl"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "member/info.html"
}

//完善用户信息提交逻辑
func (c *MemberController) UserInfo_post() {
	u := models.User{}
	u.Id, _ = c.GetSession("uid").(int)
	u.Username = c.GetString("username")

	//同步当前用户的permission数据到redis
	permission.AsyncMysql2RedisOne(u.Id)

	password := []byte(c.GetString("password"))
	passwordMD5 := md5.Sum(password)
	u.Password = fmt.Sprintf("%x", passwordMD5)
	u.IsFirst = false
	o := orm.NewOrm()
	_, err := o.Update(&u, "username", "password", "is_first", "updated")
	if err != nil {
		log.Fatal("完善用户信息错误：", err)
		c.Redirect("/userinfo", 302)
	}
	c.SetSecureCookie("userinfo_secret", "is_first", "false")
	c.Redirect("/", 302)
}

//获取用户列表
func (c *MemberController) Member_list() {
	u := []models.User{}
	o := orm.NewOrm()
	o.QueryTable("user").All(&u)
	c.Data["user"] = u
	c.Layout = "common.tpl"
	c.TplName = "member/member_list.html"
}

//个人信息修改
func (c *MemberController) Member_edit() {
	u := models.User{}
	u.Id = c.GetSession("uid").(int)

	o := orm.NewOrm()
	o.QueryTable("user").Filter("id", u.Id).One(&u, "username", "tel")

	c.Data["user"] = u
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "member/member_edit.html"
}

//个人信息修改post
func (c *MemberController) Member_edit_post() {
	u := models.User{}
	u.Id = c.GetSession("uid").(int)

	u.Username = c.GetString("username")
	password := c.GetString("password")
	u.Tel = c.GetString("tel")

	o := orm.NewOrm()
	var num int64
	var err error
	if password == "" {
		num, err = o.Update(&u, "username", "tel")
	} else {
		ps := []byte(password)
		psMD5 := md5.Sum(ps)
		u.Password = fmt.Sprintf("%x", psMD5)
		num, err = o.Update(&u, "username", "password", "tel")
	}
	if num == 1 && err == nil {
		c.Data["msg"] = "用户信息更新成功～"
		c.Data["url"] = "/"
		c.TplName = "jump/success.html"
	} else {
		c.Data["msg"] = "用户信息更新失败，如始终无法修改请联系管理员～"
		c.Data["url"] = "/member_edit"
		c.TplName = "jump/error.html"
	}
}

//管理员修改人员信息页面
func (c *MemberController) Admin_member_edit() {
	uid, _ := c.GetInt(":uid")
	if uid != 0 {
		o := orm.NewOrm()
		user := models.User{}

		o.QueryTable("user").Filter("id", uid).One(&user)
		c.Data["user"] = user
	}
	c.Data["level"] = beego.AppConfig.Strings("level")
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "member/admin_edit.html"
}

//管理员检索所要修改的用户
func (c *MemberController) Admin_member_edit_post() {
	if (c.IsAjax()) {
		search_entry := c.GetString("search_entry")
		if search_entry != "" {
			user := models.User{}

			qb, _ := orm.NewQueryBuilder("mysql")
			qb.Select("*").From("user").Where("name = ?").Or("tel = ?").Limit(1)

			sql := qb.String()
			o := orm.NewOrm()
			o.Raw(sql, search_entry, search_entry).QueryRow(&user)
			c.Data["json"] = user
			c.ServeJSON()
		}
	}
}

//管理员禁用或激活用户账号
func (c *MemberController) Disable_active_member() {
	if c.IsAjax() {
		action := c.GetString("action")
		uid, _ := c.GetInt("uid")

		user := models.User{}
		user.Id = uid
		o := orm.NewOrm()

		if action == "disable" {
			user.IsActive = false
			o.Update(&user, "is_active")

			c.Data["json"] = ResponseInfo{
				Code:    "success",
				Message: "禁用用户成功",
				Data:    "",
			}
			c.ServeJSON()
		} else if action == "active" {
			user.IsActive = true
			o.Update(&user, "is_active")

			c.Data["json"] = ResponseInfo{
				Code:    "success",
				Message: "激活用户成功",
				Data:    "",
			}
			c.ServeJSON()
		}
	}
}

//管理员修改账户信息
func (c *MemberController) Admin_edit_all() {
	user := models.User{}
	o := orm.NewOrm()

	user.Id, _ = c.GetInt("uid")
	user.Tel = c.GetString("tel")
	user.Position = c.GetString("position")
	user.PoolName = c.GetString("pool_name")

	_, err := o.Update(&user, "tel", "position", "pool_name")
	if err == nil {
		c.Redirect("/admin_member_edit/"+strconv.Itoa(user.Id), 302)
	}
}

//获取禁用账户列表
func (c *MemberController) Disable_member_list() {
	user := []models.User{}
	o := orm.NewOrm()

	o.QueryTable("user").Filter("is_active", false).All(&user)
	fmt.Println(user)
	c.Data["user"] = user
	c.Layout = "common.tpl"
	c.TplName = "member/disable_member_list.html"
}
