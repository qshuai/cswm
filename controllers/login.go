package controllers

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/qshuai/cswm/models"
)

type LoginController struct {
	beego.Controller
}

// 登陆页面
func (c *LoginController) Get() {
	if _, ok := c.GetSession("uid").(int); ok {
		c.Abort("401")
	}

	beego.ReadFromRequest(&c.Controller)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "login.html"
}

// 登陆提交验证
func (c *LoginController) Post() {
	username := c.GetString("username")
	p := []byte(c.GetString("password"))
	pMD5 := md5.Sum(p)
	password := fmt.Sprintf("%x", pMD5)

	o := orm.NewOrm()
	u := models.User{}
	err := o.QueryTable("user").Filter("username__exact", username).
		Filter("password", password).Filter("is_first", false).Filter("is_active", true).One(&u)
	if err != nil {
		uu := models.User{}
		err := o.QueryTable("user").Filter("tel__exact", username).
			Filter("password", password).Filter("is_first", true).Filter("is_active", true).One(&uu)
		if err != nil {
			flash := beego.NewFlash()
			flash.Error("注意：用户名或密码错误！")
			flash.Store(&c.Controller)
			c.Redirect("/login", 302)
			return
		}
		uu.Ip = c.Ctx.Input.IP()
		uu.LastLogin = time.Now()
		o.Update(&uu, "ip", "last_login")

		//设置session数据，保存user.Id和user.Username
		c.SetSession("uid", uu.Id)

		c.Redirect("/", 302)
	} else {
		u.Ip = c.Ctx.Input.IP()
		u.LastLogin = time.Now()
		o.Update(&u, "ip", "last_login")

		//设置session数据，保存user.Id和user.Username
		c.SetSession("uid", u.Id)
		c.SetSession("username", u.Username)

		c.Redirect("/", 302)
	}
}

// 退出登录
func (c *LoginController) Logout() {
	c.DelSession("uid")
	c.DelSession("username")

	//用户退出是
	c.SetSecureCookie("userinfo_secret", "is_first", "")
	c.Redirect("/login", 302)
}
