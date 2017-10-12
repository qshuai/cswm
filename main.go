package main

import (
	"ERP/controllers"
	_ "ERP/models"
	_ "ERP/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"regexp"
	"strconv"
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
	}
}

//验证用户是否完善用户名和密码信息
var FilterUserInfo = func(c *context.Context) {
	o := orm.NewOrm()
	uid := c.Input.Session("uid")

	//is_first变量是为了屏蔽用户登陆后还进行数据库查询操作
	is_first, _ := c.GetSecureCookie("userinfo_secret", "is_first")
	if is_first != "false" && uid != nil {
		exist := o.QueryTable("user").Filter("id", uid).Filter("username", "").Exist()
		if exist && c.Request.RequestURI != "/userinfo" {
			c.Redirect(302, "/userinfo")
		} else if !exist {
			c.SetSecureCookie("userinfo_secret", "is_first", "false")

		}
	}
	//防止用户反复修改（此页面每个用户只能设置一次）
	if is_first == "false" && c.Request.RequestURI == "/userinfo" {
		c.Redirect(302, "/")
	}
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
	//过滤器
	beego.InsertFilter("/*", beego.BeforeRouter, FilterLogin)
	beego.InsertFilter("/*", beego.BeforeRouter, FilterUserInfo)

	//自定义错误界面
	beego.ErrorController(&controllers.ErrorController{})

	beego.AddFuncMap("tostring", ToString)
	beego.AddFuncMap("delete", dd)
	beego.AddFuncMap("stripTags", stripTags)
	beego.AddFuncMap("boolToString", boolToString)
	beego.AddFuncMap("stringToInt", stringToInt)
	beego.Run()
}
