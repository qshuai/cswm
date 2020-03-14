package main

import (
	"regexp"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/qshuai/cswm/controllers"
	_ "github.com/qshuai/cswm/models"
	_ "github.com/qshuai/cswm/modules/redis"
	"github.com/qshuai/cswm/plugins/message"
	"github.com/qshuai/cswm/plugins/permission"
	"github.com/qshuai/cswm/plugins/position"
	_ "github.com/qshuai/cswm/routers"
)

func init() {
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/erp.log", "level":7, "maxlines":0, "maxsize":0, "daily":true, "maxdays":7}`)
	//beego.BeeLogger.DelLogger("console")
	logs.Async(1e3)
}

//验证用户是否登陆
var FilterLogin = func(c *context.Context) {
	_, ok := c.Input.Session("uid").(int)
	if !ok && c.Request.RequestURI != "/login" {
		c.Redirect(302, "/login")
	} else if ok && c.Request.RequestURI == "/login" {
		c.Redirect(302, "/")
	}
}

//验证用户是否完善用户名和密码信息
var FilterUserInfo = func(c *context.Context) {
	o := orm.NewOrm()
	uid, _ := c.Input.Session("uid").(int)

	//is_first变量是为了屏蔽用户登陆后还进行数据库查询操作
	is_first, _ := c.GetSecureCookie(strconv.Itoa(uid), "is_first")
	if is_first != "false" && uid != 0 {
		exist := o.QueryTable("user").Filter("id", uid).Filter("username", "").Exist()
		if exist && c.Request.RequestURI != "/userinfo" {
			c.Redirect(302, "/userinfo")
		} else if !exist {
			c.SetSecureCookie(strconv.Itoa(uid), "is_first", "false")
			c.Redirect(302, "/")

		}
	}
	//防止用户反复修改（此页面每个用户只能设置一次）
	if is_first == "false" && c.Request.RequestURI == "/userinfo" {
		c.Redirect(302, "/")
	}
}

//为每个页面都赋予authority变量（用户权限map）
var PermissionAssign = func(c *context.Context) {
	username, _ := c.Input.Session("username").(string)
	c.Input.SetData("authority", permission.GetOneRowPermission(username))
}

var PositionAssign = func(c *context.Context) {
	username, _ := c.Input.Session("username").(string)
	c.Input.SetData("username", username)
	c.Input.SetData("grade", position.GetOnePosition(username))
}

var MessageAssign = func(c *context.Context) {
	username, _ := c.Input.Session("username").(string)
	c.Input.SetData("message_num", msg.GetOneMessageNum(username))
}

//模板函数
func ToString(s int) string {
	o := strconv.Itoa(s)
	return o
}

func dd(m map[int]string, key int) map[int]string {
	delete(m, key)
	return m
}

//boot convert to string
func boolToString(b bool) string {
	if b {
		return "是"
	}
	return "否"
}

//去除字符串中的html标签
func stripTags(str string) string {
	r, _ := regexp.Compile("(?U)<.+>")
	return r.ReplaceAllString(str, "")
}

//string convert to int
func stringToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

func main() {
	//同步mysql数据表permission到redis
	permission.AsyncMysql2RedisAll()

	//同步mysql数据表position到redis
	position.AsyncAllPosition()

	//同步mysql数据表message到redis
	msg.AsyncAllMessage2Redis()

	//过滤器
	beego.InsertFilter("/*", beego.BeforeRouter, FilterLogin)
	beego.InsertFilter("/*", beego.BeforeRouter, FilterUserInfo)
	beego.InsertFilter("/*", beego.BeforeRouter, PermissionAssign)
	beego.InsertFilter("/*", beego.BeforeRouter, PositionAssign)
	beego.InsertFilter("/*", beego.BeforeRouter, MessageAssign)

	//自定义错误界面
	beego.ErrorController(&controllers.ErrorController{})

	beego.AddFuncMap("tostring", ToString)
	beego.AddFuncMap("delete", dd)
	beego.AddFuncMap("stripTags", stripTags)
	beego.AddFuncMap("boolToString", boolToString)
	beego.AddFuncMap("stringToInt", stringToInt)

	beego.Run()
}
