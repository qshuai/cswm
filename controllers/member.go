package controllers

import (
	"github.com/astaxie/beego"
	"html/template"
	"log"
	"github.com/astaxie/beego/orm"
	"erp/models"
	"crypto/md5"
	"fmt"
	"strconv"
	"erp/plugins/permission"
	"erp/plugins/position"
	"erp/modules/redis"
	"encoding/json"
)

type MemberController struct {
	beego.Controller
}

//添加用户页面
func (c *MemberController) Member_add() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "AddMember") {
		c.Abort("401")
	}
	c.Data["level"] = modify(beego.AppConfig.Strings("level"), redis_orm.RedisPool.GetOnePosition(c.GetSession("username").(string)))
	c.Layout = "common.tpl"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "member/member_add.html"
}

//添加用户提交逻辑
func (c *MemberController) Member_add_post() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "AddMember") {
		c.Abort("401")
	}
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
	u.Stage = "在职"

	o := orm.NewOrm()
	_, err := o.Insert(&u)
	if err != nil {
		log.Fatal("add member ", u.Name, " failure: ", err)
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
	if err != nil {
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
	o := orm.NewOrm()

	exist := o.QueryTable("user").Filter("username", u.Username).Exist()
	if exist {
		c.Data["url"] = "/userinfo"
		c.Data["msg"] = "此用户名已经存在~"
		c.TplName = "jump/error.html"
	}

	//同步当前用户的permission数据到redis
	permission.AsyncMysql2RedisOne(u.Username)

	password := []byte(c.GetString("password"))
	passwordMD5 := md5.Sum(password)
	u.Password = fmt.Sprintf("%x", passwordMD5)
	u.IsFirst = false
	_, err := o.Update(&u, "username", "password", "is_first", "updated")
	if err != nil {
		log.Fatal("完善用户信息错误：", err)
		c.Redirect("/userinfo", 302)
	}
	//设置session数据，存储user.Username
	c.SetSession("username", u.Username)

	o.QueryTable("user").Filter("id", u.Id).One(&u)
	//同步当前用户的position数据到redis
	position.AsyncOnePosition(u)
	permission.AsyncMysql2RedisOne(u.Username)

	c.SetSecureCookie("userinfo_secret", "is_first", "false")
	c.Redirect("/", 302)
}

//获取用户列表
func (c *MemberController) Member_list() {
	u := []models.User{}
	o := orm.NewOrm()
	o.QueryTable("user").All(&u)
	user_byte, _ := json.Marshal(u)
	c.Data["member"] = string(user_byte)
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

	//判断用户名修改
	o := orm.NewOrm()
	uu := models.User{}
	o.QueryTable("user").Filter("id", u.Id).One(&uu, "username")
	if uu.Username != u.Username {
		c.SetSession("username", u.Username)
		redis_orm.RedisPool.RenameKey(uu.Username, u.Username)
	}

	var err error
	if password == "" {
		_, err = o.Update(&u, "username", "tel")
	} else {
		ps := []byte(password)
		psMD5 := md5.Sum(ps)
		u.Password = fmt.Sprintf("%x", psMD5)
		_, err = o.Update(&u, "username", "password", "tel")
	}
	if err == nil {
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
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "EditMember") {
		c.Abort("401")
	}

	level := beego.AppConfig.Strings("level")
	p := redis_orm.RedisPool.GetOnePosition(c.GetSession("username").(string))

	uid, _ := c.GetInt(":uid")
	if uid != 0 {
		o := orm.NewOrm()
		user := models.User{}

		o.QueryTable("user").Filter("id", uid).One(&user)

		//判断当前用户是否为超级管理员
		if p != level[0] {
			if !judge(level, p, user.Position) {
				c.Abort("401")
			}
		}
		c.Data["user"] = user
	}

	c.Data["level"] = modify(level, p)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "member/admin_edit.html"
}

//管理员检索所要修改的用户
func (c *MemberController) Admin_member_edit_post() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "EditMember") {
		c.Abort("401")
	}

	if (c.IsAjax()) {
		level := beego.AppConfig.Strings("level")
		p := redis_orm.RedisPool.GetOnePosition(c.GetSession("username").(string))

		search_entry := c.GetString("search_entry")
		if search_entry != "" {
			user := models.User{}

			qb, _ := orm.NewQueryBuilder("mysql")
			qb.Select("*").From("user").Where("name = ?").Or("tel = ?").Limit(1)

			sql := qb.String()
			o := orm.NewOrm()
			o.Raw(sql, search_entry, search_entry).QueryRow(&user)
			if !judge(level, p, user.Position) {
				c.Abort("401")
			}
			c.Data["json"] = user
			c.ServeJSON()
		}
	}
}

//管理员禁用或激活用户账号
func (c *MemberController) Disable_active_member() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "ActiveMember") {
		c.Abort("401")
	}

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

func (c *MemberController) OffPosition() {
	if position.GetOnePosition(c.GetSession("username").(string)) != "超级管理员" {
		c.Abort("401")
	}

	if c.IsAjax() {
		action := c.GetString("action")
		uid, _ := c.GetInt("uid")

		user := models.User{}
		user.Id = uid
		o := orm.NewOrm()

		if action == "off" {
			user.Stage = "离职"
			o.Update(&user, "stage")

			c.Data["json"] = ResponseInfo{
				Code:    "success",
				Message: "离职成功",
				Data:    "",
			}
			c.ServeJSON()
		} else if action == "on" {
			user.Stage = "在职"
			o.Update(&user, "stage")

			c.Data["json"] = ResponseInfo{
				Code:    "success",
				Message: "在职成功",
				Data:    "",
			}
			c.ServeJSON()
		}
	}
}

//管理员修改账户信息
func (c *MemberController) Admin_edit_all() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "EditMember") {
		c.Abort("401")
	}
	user := models.User{}
	o := orm.NewOrm()
	user.Id, _ = c.GetInt("uid")

	u := models.User{}
	o.QueryTable("user").Filter("id", user.Id).One(&u, "position", "username")

	user.Tel = c.GetString("tel")
	user.Position = c.GetString("position")
	user.PoolName = c.GetString("pool_name")

	_, err := o.Update(&user, "tel", "position", "pool_name")

	if u.Position != user.Position {
		user.Username = u.Username
		position.AsyncOnePosition(user)
	}
	if err == nil {
		c.Redirect("/admin_member_edit/"+strconv.Itoa(user.Id), 302)
	}
}

//获取禁用账户列表
func (c *MemberController) Disable_member_list() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "ActiveMember") {
		c.Abort("401")
	}
	user := []models.User{}
	o := orm.NewOrm()

	o.QueryTable("user").Filter("is_active", false).All(&user)
	fmt.Println(user)
	c.Data["user"] = user
	c.Layout = "common.tpl"
	c.TplName = "member/disable_member_list.html"
}

//只允许添加比自己等级低的人员
func modify(pp []string, position string) []string {
	if position == pp[0] {
		return pp
	}
	for index, _ := range pp {
		if pp[index] == position {
			return pp[index+1:]
		}
	}
	return []string{}
}

//判断自己的等级是否高于另外一个人
func judge(pp []string, my string, your string) bool {
	var mindex, yindex int
	for index, _ := range pp {
		if my == pp[index] {
			mindex = index
		}
		if your == pp[index] {
			yindex = index
		}
	}
	if mindex < yindex {
		return true
	}
	return false
}
