package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"erp/models"
	"html/template"
	"strconv"
	"github.com/astaxie/beego/logs"
	"erp/plugins/position"
	permission2 "erp/plugins/permission"
)

type Permission struct {
	beego.Controller
}

//默认权限展示页面
func (c *Permission) DefaultPermission(){
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
func (c *Permission) DefaultPermissionEdit(){
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
func (c *Permission) DefaultPermissionEditPost(){
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
func (c *Permission) PermissionMemberList(){
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
func (c *Permission) PermissionMemberEdit(){
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
func (c *Permission) PermissionMemberEditPost(){
	if position.GetOnePosition(c.GetSession("username").(string)) != "超级管理员" {
		c.Abort("401")
	}

	permission := models.Permission{}
	permission.Id, _ = c.GetInt("permission_id")
	permission.AddMember = ConvertPermissionBool(c.GetString("AddMember"))
	permission.EditMember = ConvertPermissionBool(c.GetString("EditMember"))
	permission.ActiveMember = ConvertPermissionBool(c.GetString("ActiveMember"))
	permission.AddConsumer = ConvertPermissionBool(c.GetString("AddConsumer"))
	permission.EditConsumer = ConvertPermissionBool(c.GetString("EditConsumer"))
	permission.ViewConsumer = ConvertPermissionBool(c.GetString("ViewConsumer"))
	permission.AddBrand = ConvertPermissionBool(c.GetString("AddBrand"))
	permission.AddDealer = ConvertPermissionBool(c.GetString("AddDealer"))
	permission.ViewDealer = ConvertPermissionBool(c.GetString("ViewDealer"))
	permission.AddSupplier = ConvertPermissionBool(c.GetString("AddSupplier"))
	permission.ViewSupplier = ConvertPermissionBool(c.GetString("ViewSupplier"))
	permission.AddProduct = ConvertPermissionBool(c.GetString("AddProduct"))
	permission.InputInPrice = ConvertPermissionBool(c.GetString("InputInPrice"))
	permission.ViewProductStore = ConvertPermissionBool(c.GetString("ViewProductStore"))
	permission.ViewStock = ConvertPermissionBool(c.GetString("ViewStock"))
	permission.ViewInPrice = ConvertPermissionBool(c.GetString("ViewInPrice"))
	permission.EditProduct = ConvertPermissionBool(c.GetString("EditProduct"))
	permission.DeleteProduct = ConvertPermissionBool(c.GetString("DeleteProduct"))
	permission.OutputProduct = ConvertPermissionBool(c.GetString("OutputProduct"))
	permission.ViewSale = ConvertPermissionBool(c.GetString("ViewSale"))
	permission.ViewSaleConsumer = ConvertPermissionBool(c.GetString("ViewSaleConsumer"))
	permission.ViewSaleInPrice = ConvertPermissionBool(c.GetString("ViewSaleInPrice"))
	permission.EditSale = ConvertPermissionBool(c.GetString("EditSale"))
	permission.OperateCategory = ConvertPermissionBool(c.GetString("OperateCategory"))
	permission.RequestMove = ConvertPermissionBool(c.GetString("RequestMove"))
	permission.ResponseMove = ConvertPermissionBool(c.GetString("ResponseMove"))
	permission.ViewMove = ConvertPermissionBool(c.GetString("ViewMove"))
	permission.AddStore = ConvertPermissionBool(c.GetString("AddStore"))
	permission.ViewStore = ConvertPermissionBool(c.GetString("ViewStore"))
	permission.OperateOtherStore = ConvertPermissionBool(c.GetString("OperateOtherStore"))

	o := orm.NewOrm()
	user := models.User{}
	user.Id, _ = c.GetInt("permission_user_id")
	permission.User = &user
	permission_id := models.Permission{}
	o.QueryTable("permission").Filter("User__id", user.Id).One(&permission_id, "id")

	permission.Id = permission_id.Id
	_, err := o.Update(&permission)
	if err != nil {
		logs.Error("修改人员权限失败")
		c.Data["url"] = "/permission_member_edit/" + strconv.Itoa(user.Id)
		c.Data["msg"] = "权限修改失败~"
		c.TplName = "jump/error.html"
		return
	}

	//redis
	o.QueryTable("user").Filter("id", user.Id).One(&user, "id", "username")
	permission2.AsyncMysql2RedisOne(user.Username)

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
