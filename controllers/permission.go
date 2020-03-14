package controllers

import (
	"html/template"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/qshuai/cswm/models"
	"github.com/qshuai/cswm/plugins/permission"
	"github.com/qshuai/cswm/plugins/position"
)

type Permission struct {
	beego.Controller
}

//默认权限展示页面
func (c *Permission) DefaultPermission() {
	if position.GetOnePosition(c.GetSession("username").(string)) != "超级管理员" {
		c.Abort("401")
	}

	defaultPermission := []models.DefaultPermission{}
	o := orm.NewOrm()
	o.QueryTable("default_permission").All(&defaultPermission)

	c.Data["defaultPermission"] = defaultPermission
	c.Layout = "common.tpl"
	c.TplName = "permission/default_permission.html"
}

//默认权限编辑页面
func (c *Permission) DefaultPermissionEdit() {
	if position.GetOnePosition(c.GetSession("username").(string)) != "超级管理员" {
		c.Abort("401")
	}

	defaultPermission := models.DefaultPermission{}
	defaultPermission.Id, _ = c.GetInt(":item")

	o := orm.NewOrm()
	o.QueryTable("default_permission").Filter("id", defaultPermission.Id).One(&defaultPermission)

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["defaultPermission"] = defaultPermission
	c.Layout = "common.tpl"
	c.TplName = "permission/default_permission_edit.html"
}

//默认权限编辑post提交
func (c *Permission) DefaultPermissionEditPost() {
	if position.GetOnePosition(c.GetSession("username").(string)) != "超级管理员" {
		c.Abort("401")
	}

	defaultPermission := models.DefaultPermission{}
	defaultPermission.Id, _ = c.GetInt("permission_id")
	defaultPermission.Position = c.GetString("permission_position")
	defaultPermission.AddMember = ConvertPermissionBool(c.GetString("AddMember"))
	defaultPermission.EditMember = ConvertPermissionBool(c.GetString("EditMember"))
	defaultPermission.ActiveMember = ConvertPermissionBool(c.GetString("ActiveMember"))
	defaultPermission.AddConsumer = ConvertPermissionBool(c.GetString("AddConsumer"))
	defaultPermission.EditConsumer = ConvertPermissionBool(c.GetString("EditConsumer"))
	defaultPermission.ViewConsumer = ConvertPermissionBool(c.GetString("ViewConsumer"))
	defaultPermission.AddBrand = ConvertPermissionBool(c.GetString("AddBrand"))
	defaultPermission.AddDealer = ConvertPermissionBool(c.GetString("AddDealer"))
	defaultPermission.ViewDealer = ConvertPermissionBool(c.GetString("ViewDealer"))
	defaultPermission.AddSupplier = ConvertPermissionBool(c.GetString("AddSupplier"))
	defaultPermission.ViewSupplier = ConvertPermissionBool(c.GetString("ViewSupplier"))
	defaultPermission.AddProduct = ConvertPermissionBool(c.GetString("AddProduct"))
	defaultPermission.InputInPrice = ConvertPermissionBool(c.GetString("InputInPrice"))
	defaultPermission.ViewProductStore = ConvertPermissionBool(c.GetString("ViewProductStore"))
	defaultPermission.ViewStock = ConvertPermissionBool(c.GetString("ViewStock"))
	defaultPermission.ViewInPrice = ConvertPermissionBool(c.GetString("ViewInPrice"))
	defaultPermission.EditProduct = ConvertPermissionBool(c.GetString("EditProduct"))
	defaultPermission.DeleteProduct = ConvertPermissionBool(c.GetString("DeleteProduct"))
	defaultPermission.OutputProduct = ConvertPermissionBool(c.GetString("OutputProduct"))
	defaultPermission.ViewSale = ConvertPermissionBool(c.GetString("ViewSale"))
	defaultPermission.ViewSaleConsumer = ConvertPermissionBool(c.GetString("ViewSaleConsumer"))
	defaultPermission.ViewSaleInPrice = ConvertPermissionBool(c.GetString("ViewSaleInPrice"))
	defaultPermission.EditSale = ConvertPermissionBool(c.GetString("EditSale"))
	defaultPermission.OperateCategory = ConvertPermissionBool(c.GetString("OperateCategory"))
	defaultPermission.RequestMove = ConvertPermissionBool(c.GetString("RequestMove"))
	defaultPermission.ResponseMove = ConvertPermissionBool(c.GetString("ResponseMove"))
	defaultPermission.ViewMove = ConvertPermissionBool(c.GetString("ViewMove"))
	defaultPermission.AddStore = ConvertPermissionBool(c.GetString("AddStore"))
	defaultPermission.ViewStore = ConvertPermissionBool(c.GetString("ViewStore"))
	defaultPermission.OperateOtherStore = ConvertPermissionBool(c.GetString("OperateOtherStore"))

	o := orm.NewOrm()
	_, err := o.Update(&defaultPermission)
	if err != nil {
		logs.Error("修改默认权限失败")
		c.Data["url"] = "/default_permission_edit/" + strconv.Itoa(defaultPermission.Id)
		c.Data["msg"] = "权限修改失败~"
		c.TplName = "jump/error.html"
		return
	}
	c.Data["url"] = "/default_permission"
	c.Data["msg"] = "权限修改成功~"
	c.TplName = "jump/success.html"
}

//人员列表
func (c *Permission) PermissionMemberList() {
	if position.GetOnePosition(c.GetSession("username").(string)) != "超级管理员" {
		c.Abort("401")
	}

	u := []models.User{}
	o := orm.NewOrm()
	o.QueryTable("user").All(&u)
	c.Data["user"] = u
	c.Layout = "common.tpl"
	c.TplName = "permission/permission_member_list.html"
}

//人员权限编辑页面
func (c *Permission) PermissionMemberEdit() {
	if position.GetOnePosition(c.GetSession("username").(string)) != "超级管理员" {
		c.Abort("401")
	}

	o := orm.NewOrm()
	permission := models.Permission{}
	defaultPermission := models.DefaultPermission{}
	user := models.User{}

	uid, _ := c.GetInt(":uid")
	o.QueryTable("user").Filter("id", uid).One(&user, "position")
	o.QueryTable("default_permission").Filter("position", user.Position).One(&defaultPermission)
	o.QueryTable("permission").Filter("User__id", uid).One(&permission)

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["default_permission"] = defaultPermission
	c.Data["permission"] = permission
	c.Layout = "common.tpl"
	c.TplName = "permission/permission_member_edit.html"
}

//人员权限编辑post提交
func (c *Permission) PermissionMemberEditPost() {
	if position.GetOnePosition(c.GetSession("username").(string)) != "超级管理员" {
		c.Abort("401")
	}

	per := models.Permission{}
	per.Id, _ = c.GetInt("permission_id")
	per.AddMember = ConvertPermissionBool(c.GetString("AddMember"))
	per.EditMember = ConvertPermissionBool(c.GetString("EditMember"))
	per.ActiveMember = ConvertPermissionBool(c.GetString("ActiveMember"))
	per.AddConsumer = ConvertPermissionBool(c.GetString("AddConsumer"))
	per.EditConsumer = ConvertPermissionBool(c.GetString("EditConsumer"))
	per.ViewConsumer = ConvertPermissionBool(c.GetString("ViewConsumer"))
	per.AddBrand = ConvertPermissionBool(c.GetString("AddBrand"))
	per.AddDealer = ConvertPermissionBool(c.GetString("AddDealer"))
	per.ViewDealer = ConvertPermissionBool(c.GetString("ViewDealer"))
	per.AddSupplier = ConvertPermissionBool(c.GetString("AddSupplier"))
	per.ViewSupplier = ConvertPermissionBool(c.GetString("ViewSupplier"))
	per.AddProduct = ConvertPermissionBool(c.GetString("AddProduct"))
	per.InputInPrice = ConvertPermissionBool(c.GetString("InputInPrice"))
	per.ViewProductStore = ConvertPermissionBool(c.GetString("ViewProductStore"))
	per.ViewStock = ConvertPermissionBool(c.GetString("ViewStock"))
	per.ViewInPrice = ConvertPermissionBool(c.GetString("ViewInPrice"))
	per.EditProduct = ConvertPermissionBool(c.GetString("EditProduct"))
	per.DeleteProduct = ConvertPermissionBool(c.GetString("DeleteProduct"))
	per.OutputProduct = ConvertPermissionBool(c.GetString("OutputProduct"))
	per.ViewSale = ConvertPermissionBool(c.GetString("ViewSale"))
	per.ViewSaleConsumer = ConvertPermissionBool(c.GetString("ViewSaleConsumer"))
	per.ViewSaleInPrice = ConvertPermissionBool(c.GetString("ViewSaleInPrice"))
	per.EditSale = ConvertPermissionBool(c.GetString("EditSale"))
	per.OperateCategory = ConvertPermissionBool(c.GetString("OperateCategory"))
	per.RequestMove = ConvertPermissionBool(c.GetString("RequestMove"))
	per.ResponseMove = ConvertPermissionBool(c.GetString("ResponseMove"))
	per.ViewMove = ConvertPermissionBool(c.GetString("ViewMove"))
	per.AddStore = ConvertPermissionBool(c.GetString("AddStore"))
	per.ViewStore = ConvertPermissionBool(c.GetString("ViewStore"))
	per.OperateOtherStore = ConvertPermissionBool(c.GetString("OperateOtherStore"))

	o := orm.NewOrm()
	user := models.User{}
	user.Id, _ = c.GetInt("permission_user_id")
	per.User = &user
	permission_id := models.Permission{}
	o.QueryTable("permission").Filter("User__id", user.Id).One(&permission_id, "id")

	per.Id = permission_id.Id
	_, err := o.Update(&per)
	if err != nil {
		logs.Error("修改人员权限失败")
		c.Data["url"] = "/permission_member_edit/" + strconv.Itoa(user.Id)
		c.Data["msg"] = "权限修改失败~"
		c.TplName = "jump/error.html"
		return
	}

	//redis
	o.QueryTable("user").Filter("id", user.Id).One(&user, "id", "username")
	permission.AsyncMysql2RedisOne(user.Username)

	c.Data["url"] = "/permission_member_edit/" + strconv.Itoa(user.Id)
	c.Data["msg"] = "权限修改成功~"
	c.TplName = "jump/success.html"
}

/*
 * 公用函数
 * str 用户输入（可能的值为"on"或为""）
 */
func ConvertPermissionBool(str string) bool {
	if str == "on" {
		return true
	}
	return false
}
