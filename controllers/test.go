package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"fmt"
)

type TestController struct{
	beego.Controller
}

type Product struct{
	Titlename string
	Username string
	Name string
}

func (c *TestController) Get() {
	product := []Product{}
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("product.title as titlename", "user.username", "brand.name").
	From("product").
	InnerJoin("user").
	On("product.user_id = user.id").
	InnerJoin("brand").
	On("product.brand_id = brand.id")

	sql := qb.String()
	o := orm.NewOrm()
	o.Raw(sql).QueryRows(&product)
	fmt.Printf("%#v", product)
}
