package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"erp/models"
	"html/template"
	"time"
	"erp/plugins/permission"
	"encoding/json"
	"strings"
	"fmt"
	"strconv"
	"erp/plugins/position"
)

type SaleController struct {
	beego.Controller
}

const SaleLimit = 100

type salelist struct {
	Id           int
	Title        string
	No           string
	Brand        string
	Unit         string
	Spec         string
	Pool         string
	StoreName    string
	ArtNum       string
	SalesmanName string
	ConsumerName string
	InPrice      string
	OutPrice     string
	Num          string
	Send         string
	HasInvoice   bool
	InvoiceNum   string
	SendInvoice  string
	GetInvoice   string
	GetMoney     bool
	GetDate      string
	Comment      string
	HasPrint     bool
	Created      string
}

//获取销售列表数据
func (c *SaleController) Sale_list() {
	username := c.GetSession("username").(string)
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "ViewSale") {
		c.Abort("401")
	}

	user := models.User{}
	o := orm.NewOrm()
	o.QueryTable("user").Filter("username", username).One(&user, "position", "pool_name")
	operate_other_store := !permission.GetOneItemPermission(username, "OperateOtherStore")

	sale := []salelist{}
	qb, _ := orm.NewQueryBuilder("mysql")

	if user.Position != "超级管理员" || operate_other_store {
		if user.PoolName != "" {
			qb.Select("sale.id", "sale.out_price", "sale.num", "sale.send", "sale.has_invoice", "sale.invoice_num", "sale.no", "sale.comment", "sale.has_print",
				"sale.send_invoice", "sale.get_invoice", "sale.get_money", "sale.get_date", "sale.created", "store.pool", "store.name as store_name",
				"product.title", "product.art_num", "product.in_price", "brand.name as brand", "product.unit", "product.spec", "user.name as salesman_name", "consumer.name as consumer_name").
				From("sale").
				LeftJoin("product").
				On("product.id = sale.product_id").
				InnerJoin("brand").
				On("product.brand_id = brand.id").
				LeftJoin("user").
				On("user.id = sale.salesman_id").
				LeftJoin("consumer").
				On("consumer.id = sale.consumer_id").
				InnerJoin("store").
				On("store.id = sale.store_id AND store.name = ?").
				OrderBy("created").Desc().
				Limit(SaleLimit)
			sql := qb.String()
			o.Raw(sql, user.PoolName).QueryRows(&sale)
		}
	} else {
		qb.Select("sale.id", "sale.out_price", "sale.num", "sale.send", "sale.has_invoice", "sale.invoice_num", "sale.no", "sale.comment", "sale.has_print",
			"sale.send_invoice", "sale.get_invoice", "sale.get_money", "sale.get_date", "sale.created", "store.pool", "store.name as store_name",
			"product.title", "product.art_num", "product.in_price", "brand.name as brand", "product.unit", "product.spec", "user.name as salesman_name", "consumer.name as consumer_name").
			From("sale").
			LeftJoin("product").
			On("product.id = sale.product_id").
			InnerJoin("brand").
			On("product.brand_id = brand.id").
			LeftJoin("user").
			On("user.id = sale.salesman_id").
			LeftJoin("consumer").
			On("consumer.id = sale.consumer_id").
			InnerJoin("store").
			On("store.id = sale.store_id").
			OrderBy("created").Desc().
			Limit(SaleLimit)
		sql := qb.String()
		o.Raw(sql).QueryRows(&sale)
	}

	view_consumer := !permission.GetOneItemPermission(username, "ViewSaleConsumer")
	view_in_price := !permission.GetOneItemPermission(username, "ViewSaleInPrice")
	if view_consumer || view_in_price {
		length := len(sale)
		for index := 0; index < length; index++ {
			if view_consumer {
				sale[index].ConsumerName = "***"
			}
			if view_in_price {
				sale[index].InPrice = "***"
			}
		}
	}

	sale_byte, _ := json.Marshal(sale)
	c.Data["sale"] = string(sale_byte)

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "sale/sale_list.html"
}

//获取更多销售记录
func (c *SaleController) SaleLoadMore() {
	username := c.GetSession("username").(string)
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "ViewSale") {
		c.Abort("401")
	}

	user := models.User{}
	o := orm.NewOrm()
	o.QueryTable("user").Filter("username", username).One(&user, "position", "pool_name")
	operate_other_store := !permission.GetOneItemPermission(username, "OperateOtherStore")

	sale := []salelist{}
	qb, _ := orm.NewQueryBuilder("mysql")
	offset, _ := c.GetInt("offset")

	if user.Position != "超级管理员" || operate_other_store {
		if user.PoolName != "" {
			qb.Select("sale.id", "sale.out_price", "sale.num", "sale.send", "sale.has_invoice", "sale.invoice_num", "sale.no", "sale.comment", "sale.has_print",
				"sale.send_invoice", "sale.get_invoice", "sale.get_money", "sale.get_date", "sale.created", "store.pool", "store.name as store_name",
				"product.title", "product.art_num", "product.in_price", "brand.name as brand", "product.unit", "product.spec", "user.name as salesman_name", "consumer.name as consumer_name").
				From("sale").
				LeftJoin("product").
				On("product.id = sale.product_id").
				InnerJoin("brand").
				On("product.brand_id = brand.id").
				LeftJoin("user").
				On("user.id = sale.salesman_id").
				LeftJoin("consumer").
				On("consumer.id = sale.consumer_id").
				InnerJoin("store").
				On("store.id = sale.store_id AND store.name = ?").
				OrderBy("created").Desc().
				Limit(SaleLimit).
				Offset(SaleLimit * offset)

			sql := qb.String()
			o.Raw(sql, user.PoolName).QueryRows(&sale)
		}
	} else {

		qb.Select("sale.id", "sale.out_price", "sale.num", "sale.send", "sale.has_invoice", "sale.invoice_num", "sale.no", "sale.comment", "sale.has_print",
			"sale.send_invoice", "sale.get_invoice", "sale.get_money", "sale.get_date", "sale.created", "store.pool", "store.name as store_name",
			"product.title", "product.art_num", "product.in_price", "brand.name as brand", "product.unit", "product.spec", "user.name as salesman_name", "consumer.name as consumer_name").
			From("sale").
			LeftJoin("product").
			On("product.id = sale.product_id").
			InnerJoin("brand").
			On("product.brand_id = brand.id").
			LeftJoin("user").
			On("user.id = sale.salesman_id").
			LeftJoin("consumer").
			On("consumer.id = sale.consumer_id").
			InnerJoin("store").
			On("store.id = sale.store_id").
			OrderBy("created").Desc().
			Limit(SaleLimit).
			Offset(SaleLimit * offset)
		sql := qb.String()
		o.Raw(sql).QueryRows(&sale)
	}

	view_consumer := !permission.GetOneItemPermission(username, "ViewSaleConsumer")
	view_in_price := !permission.GetOneItemPermission(username, "ViewSaleInPrice")
	if view_consumer || view_in_price {
		length := len(sale)
		for index := 0; index < length; index++ {
			if view_consumer {
				sale[index].ConsumerName = "***"
			}
			if view_in_price {
				sale[index].InPrice = "***"
			}
		}
	}

	sale_byte, _ := json.Marshal(sale)
	c.Data["json"] = string(sale_byte)
	c.ServeJSON()
}

//单条销售记录修改post
func (c *SaleController) Sale_edit() {
	un := c.GetSession("username").(string)
	if !permission.GetOneItemPermission(un, "EditSale") {
		c.Abort("401")
	}

	sale := models.Sale{}
	salesman := models.User{}

	o := orm.NewOrm()
	o.QueryTable("user").Filter("name", c.GetString("salesman")).One(&salesman)
	sale.Salesman = &salesman

	sale.Id, _ = c.GetInt("sale_id")
	sale.OutPrice, _ = c.GetFloat("outprice")
	sale.Num, _ = c.GetUint32("num")
	sale.Send, _ = time.Parse("2006-1-2", c.GetString("send"))
	sale.HasInvoice, _ = c.GetBool("hasinvoice")
	sale.InvoiceNum = c.GetString("invioce_num")
	sale.SendInvoice, _ = time.Parse("2006-1-2", c.GetString("send_invioce"))
	sale.GetInvoice, _ = time.Parse("2006-1-2", c.GetString("get_invioce"))
	sale.GetMoney, _ = c.GetBool("get_money")
	sale.GetDate, _ = time.Parse("2006-1-2", c.GetString("get_date"))
	sale.Comment = c.GetString("comment")
	sale.Store = GetStore(c.GetString("store"))

	_, err := o.Update(&sale, "salesman", "out_price", "num", "send", "has_invoice", "invoice_num", "send_invoice", "get_invoice", "get_money", "get_date", "store_id", "comment")
	if err == nil {
		c.Data["url"] = "/sale_list"
		c.Data["msg"] = "修改销售记录成功"
		c.TplName = "jump/success.html"
	}
}

//打印
func (c *SaleController) Print() {
	print_list := c.GetString(":list")
	print_id, _ := c.GetInt(":id")
	print_slice := strings.Split(print_list, ",")
	new_slice := make([]string, 0)
	length := len(print_slice)
	for i := 0; i < length; i++ {
		for j := 0; j < i; j++ {
			if print_slice[i] == print_slice[j] {
				break
			}
			if j == i-1 {
				new_slice = append(new_slice, print_slice[i])
			}
		}
	}
	new_slice = append(new_slice, print_slice[0])

	str := "(" + strings.Join(new_slice, ",") + ")"

	type printsale struct {
		Num                string
		OutPrice           string
		ConsumerName       string
		ConsumerDepartment string
		SalesmanName       string
		Title              string
		ArtNum             string
		Spec               string
		Unit               string
		BrandName          string
		Pool               string
		StoreName          string
	}

	o := orm.NewOrm()
	order := models.OrderNum{}
	order.Id = print_id
	o.Read(&order)

	//判断是否为fake
	qb, _ := orm.NewQueryBuilder("mysql")
	if order.IsFake {
		qb.Select("sale.num_fake as num", "sale.out_price_fake as out_price", "consumer.name as consumer_name",
			"consumer.department as consumer_department",
			"user.name as salesman_name", "product.title", "product.art_num", "product.spec",
			"product.unit", "brand.name as brand_name", "store.pool", "store.name as store_name").
			From("sale").
			InnerJoin("product").
			On("product.id = sale.product_id").
			InnerJoin("user").
			On("user.id = sale.salesman_id").
			InnerJoin("store").
			On("store.id = sale.store_id").
			InnerJoin("brand").
			On("brand.id = product.brand_id").
			InnerJoin("consumer").
			On("consumer.id = sale.consumer_id").
			Where("sale.id in " + str)
	} else {
		qb.Select("sale.num", "sale.out_price", "consumer.name as consumer_name",
			"consumer.department as consumer_department",
			"user.name as salesman_name", "product.title", "product.art_num", "product.spec",
			"product.unit", "brand.name as brand_name", "store.pool", "store.name as store_name").
			From("sale").
			InnerJoin("product").
			On("product.id = sale.product_id").
			InnerJoin("user").
			On("user.id = sale.salesman_id").
			InnerJoin("store").
			On("store.id = sale.store_id").
			InnerJoin("brand").
			On("brand.id = product.brand_id").
			InnerJoin("consumer").
			On("consumer.id = sale.consumer_id").
			Where("sale.id in " + str)
	}

	sql := qb.String()
	ps := []printsale{}

	//fake的order只能超级管理员打印，正常的order可以有超级管理员、总库管理员或分库管理员打印
	p := position.GetOnePosition(c.GetSession("username").(string))
	lastLetter := order.Asap[len(order.Asap)-1:]
	if lastLetter == "0" {
		if p != "超级管理员" {
			c.Abort("401")
		}
	}else{
		if !(p == "超级管理员" || p == "分库管理员" || p == "总库管理员") {
			c.Abort("401")
		}
	}

	o.Raw(sql).QueryRows(&ps)
	c.Data["order"] = order
	c.Data["print"] = ps
	c.Data["store"] = ps[0].Pool + "-" + ps[0].StoreName
	c.Data["date"] = fmt.Sprint(time.Now())[:19]

	c.TplName = "sale/print.html"
}

//出库单列表
func (c *SaleController) OrderList() {
	order := []models.OrderNum{}
	o := orm.NewOrm()
	o.QueryTable("order_num").OrderBy("-id").All(&order)
	c.Data["order"] = order
	c.Layout = "common.tpl"
	c.TplName = "sale/order_list.html"
}

//作废出库单
func (c *SaleController) OrderClose() {
	oid, _ := c.GetInt(":oid")
	o := orm.NewOrm()
	_, err := o.Raw("UPDATE order_num SET  state = false, updated = ? WHERE id = ?", time.Now(), oid).Exec()
	order_list := models.OrderNum{}
	o.QueryTable("order_num").Filter("id", oid).One(&order_list, "sale_list")
	orders := strings.Split(order_list.SaleList, ",")
	for _, item := range orders {
		o.Raw("UPDATE sale SET has_print = false WHERE id = ?", item).Exec()
	}

	if err != nil {
		c.Data["url"] = "/order_list"
		c.Data["msg"] = "作废出库单失败"
		c.TplName = "jump/success.html"
	} else {
		c.Data["url"] = "/order_list"
		c.Data["msg"] = "作废出库单成功"
		c.TplName = "jump/success.html"
	}
}

//从product_list页面点击按钮跳转到相应货号的销售列表
func (c *SaleController) ProductSalInfo() {
	username := c.GetSession("username").(string)
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "ViewSale") {
		c.Abort("401")
	}

	art_num := c.GetString(":art_num")
	user := models.User{}
	o := orm.NewOrm()
	o.QueryTable("user").Filter("username", username).One(&user, "position", "pool_name")
	operate_other_store := !permission.GetOneItemPermission(username, "OperateOtherStore")

	sale := []salelist{}
	qb, _ := orm.NewQueryBuilder("mysql")

	if user.Position != "超级管理员" || operate_other_store {
		if user.PoolName != "" {
			qb.Select("sale.id", "sale.out_price", "sale.num", "sale.send", "sale.has_invoice", "sale.invoice_num", "sale.no", "sale.comment", "sale.has_print",
				"sale.send_invoice", "sale.get_invoice", "sale.get_money", "sale.get_date", "sale.created", "store.pool", "store.name as store_name",
				"product.title", "product.art_num", "product.in_price", "brand.name as brand", "product.unit", "product.spec", "user.name as salesman_name", "consumer.name as consumer_name").
				From("sale").
				LeftJoin("product").
				On("product.id = sale.product_id AND product.art_num = ?").
				InnerJoin("brand").
				On("product.brand_id = brand.id").
				LeftJoin("user").
				On("user.id = sale.salesman_id").
				LeftJoin("consumer").
				On("consumer.id = sale.consumer_id").
				InnerJoin("store").
				On("store.id = sale.store_id AND store.name = ?").
				Where("product.art_num = ?").
				OrderBy("created").Desc()
			sql := qb.String()
			o.Raw(sql, art_num, user.PoolName, art_num).QueryRows(&sale)
		}
	} else {
		qb.Select("sale.id", "sale.out_price", "sale.num", "sale.send", "sale.has_invoice", "sale.invoice_num", "sale.no", "sale.comment", "sale.has_print",
			"sale.send_invoice", "sale.get_invoice", "sale.get_money", "sale.get_date", "sale.created", "store.pool", "store.name as store_name",
			"product.title", "product.art_num", "product.in_price", "brand.name as brand", "product.unit", "product.spec", "user.name as salesman_name", "consumer.name as consumer_name").
			From("sale").
			LeftJoin("product").
			On("product.id = sale.product_id AND product.art_num = ?").
			InnerJoin("brand").
			On("product.brand_id = brand.id").
			LeftJoin("user").
			On("user.id = sale.salesman_id").
			LeftJoin("consumer").
			On("consumer.id = sale.consumer_id").
			InnerJoin("store").
			On("store.id = sale.store_id").
			OrderBy("created").Desc()
		sql := qb.String()
		o.Raw(sql, art_num).QueryRows(&sale)
	}

	view_consumer := !permission.GetOneItemPermission(username, "ViewSaleConsumer")
	view_in_price := !permission.GetOneItemPermission(username, "ViewSaleInPrice")
	if view_consumer || view_in_price {
		length := len(sale)
		for index := 0; index < length; index++ {
			if view_consumer {
				sale[index].ConsumerName = "***"
			}
			if view_in_price {
				sale[index].InPrice = "***"
			}
		}
	}

	sale_byte, _ := json.Marshal(sale)
	c.Data["sale"] = string(sale_byte)

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "sale/sale_list.html"
}

//编辑出库单（伪造）
func (c *SaleController) OrderEdit() {
	print_list := c.Ctx.GetCookie("print_sale_list")
	if print_list == "" {
		c.Data["url"] = "/sale_list"
		c.Data["msg"] = "没有指定出库单"
		c.TplName = "jump/error.html"
		return
	}
	print_slice := strings.Split(print_list, "%2C")
	new_slice := make([]string, 0)
	length := len(print_slice)
	for i := 0; i < length; i++ {
		for j := 0; j < i; j++ {
			if print_slice[i] == print_slice[j] {
				break
			}
			if j == i-1 {
				new_slice = append(new_slice, print_slice[i])
			}
		}
	}
	new_slice = append(new_slice, print_slice[0])
	str := "(" + strings.Join(new_slice, ",") + ")"

	//sale := []models.Sale{}
	//o := orm.NewOrm()
	//o.QueryTable("sale").Filter("id__in", new_slice).All(&sale, "id", "title", "out_price_fake", "num_fake")

	type sale_edit struct {
		Id           string
		Title        string
		OutPriceFake string
		NumFake      string
	}
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("sale.id", "sale.out_price_fake", "sale.num_fake", "product.title").
		From("sale").
		InnerJoin("product").
		On("product.id = sale.product_id").
		Where("sale.id in " + str)
	sql := qb.String()
	o := orm.NewOrm()
	sale := []sale_edit{}
	o.Raw(sql).QueryRows(&sale)

	c.Data["sale"] = sale
	c.Layout = "common.tpl"
	c.TplName = "sale/order_edit.html"
}

//编辑出库单（伪造）
func (c *SaleController) OrderEditPost() {
	ids := c.GetStrings("sid")
	outs := c.GetStrings("out_price")
	nums := c.GetStrings("num")

	//和上面代码一样
	print_list := c.Ctx.GetCookie("print_sale_list")
	if print_list == "" {
		c.Data["url"] = "/sale_list"
		c.Data["msg"] = "没有指定出库单"
		c.TplName = "jump/error.html"
		return
	}
	print_slice := strings.Split(print_list, "%2C")
	new_slice := make([]string, 0)
	length := len(print_slice)
	for i := 0; i < length; i++ {
		for j := 0; j < i; j++ {
			if print_slice[i] == print_slice[j] {
				break
			}
			if j == i-1 {
				new_slice = append(new_slice, print_slice[i])
			}
		}
	}
	new_slice = append(new_slice, print_slice[0])
	str := "(" + strings.Join(new_slice, ",") + ")"

	//sale := []models.Sale{}
	//o := orm.NewOrm()
	//o.QueryTable("sale").Filter("id__in", new_slice).All(&sale, "id", "title", "out_price_fake", "num_fake")

	type sale_edit struct {
		Id           string
		Title        string
		OutPriceFake string
		NumFake      string
	}
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("sale.id", "sale.out_price_fake", "sale.num_fake", "product.title").
		From("sale").
		InnerJoin("product").
		On("product.id = sale.product_id").
		Where("sale.id in " + str)
	sql := qb.String()
	o := orm.NewOrm()
	sale := []sale_edit{}
	o.Raw(sql).QueryRows(&sale)
	//------

	isfake := false
	for index, id := range ids {
		for _, item := range sale {
			if id == item.Id {
				if outs[index] != item.OutPriceFake || nums[index] != item.NumFake {
					isfake = true
				}
				o.Raw("UPDATE sale SET out_price_fake = ?, num_fake = ?, has_print = ? where id = ?", outs[index], nums[index], true, id).Exec()
			}
		}
	}

	username := c.GetSession("username").(string)
	if isfake {
		InsertOrder(new_slice, username)
		InsertFakeOrder(new_slice, username)
	} else {
		InsertOrder(new_slice, username)
	}

	//清空cookie
	c.Ctx.SetCookie("print_sale_list", "")

	c.Redirect("/order_list", 302)
}

//显示出库单列表
func (c *SaleController) OrderAdd() {
	o := orm.NewOrm()
	order_list := []models.OrderNum{}
	o.QueryTable("order_num").OrderBy("-id").All(&order_list)

	order_byte, _ := json.Marshal(order_list)
	c.Data["order"] = string(order_byte)

	c.Layout = "common.tpl"
	c.TplName = "sale/order_list.html"
}

func InsertOrder(new_slice []string, username string) {
	//更新order_num表
	order := models.OrderNum{}
	//o.QueryTable("order_num").OrderBy("-id").One(&order, "id")
	order.User = username
	order.State = true
	order.SaleList = strings.Join(new_slice, ",")
	str := "(" + strings.Join(new_slice, ",") + ")"

	type printsale struct {
		Num                int
		OutPrice           float64
		ConsumerName       string
		ConsumerDepartment string
		SalesmanName       string
		Title              string
		ArtNum             string
		Spec               string
		Unit               string
		BrandName          string
		Pool               string
		StoreName          string
	}
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("sale.num", "sale.out_price", "consumer.name as consumer_name",
		"consumer.department as consumer_department",
		"user.name as salesman_name", "product.title", "product.art_num", "product.spec",
		"product.unit", "brand.name as brand_name", "store.pool", "store.name as store_name").
		From("sale").
		InnerJoin("product").
		On("product.id = sale.product_id").
		InnerJoin("user").
		On("user.id = sale.salesman_id").
		InnerJoin("store").
		On("store.id = sale.store_id").
		InnerJoin("brand").
		On("brand.id = product.brand_id").
		InnerJoin("consumer").
		On("consumer.id = sale.consumer_id").
		Where("sale.id in " + str)
	sql := qb.String()
	ps := []printsale{}
	o := orm.NewOrm()
	o.Raw(sql).QueryRows(&ps)
	order.Consumer = ps[0].ConsumerName
	order.Department = ps[0].ConsumerDepartment
	order.Salesman = ps[0].SalesmanName

	//计算总额
	var total float64
	for _, item := range ps {
		total = total + item.OutPrice*float64(item.Num)
	}
	order.Sum = fmt.Sprintf("%0.2f", total)

	//ASAP拼凑字符串
	//首先获取order数据库数据，便于取出最新的ASAP单号
	order_list := []models.OrderNum{}
	o.QueryTable("order_num").OrderBy("-id").All(&order_list)

	t := time.Now()
	year := strconv.Itoa(t.Year())[2:]
	month := t.Month().String()
	m2m := map[string]string{
		"January":   "01",
		"February":  "02",
		"March":     "03",
		"April":     "04",
		"May":       "05",
		"June":      "06",
		"July":      "07",
		"August":    "08",
		"September": "09",
		"October":   "10",
		"November":  "11",
		"December":  "12",
	}
	month = m2m[month]

	var old int
	if len(order_list) == 0 {
		old = 0
	} else {
		length := len(order_list[0].Asap)
		old, _ = strconv.Atoi(order_list[0].Asap[4:length-1])
	}
	v := old + 1
	var new_string string
	if v < 10 {
		new_string = "00" + strconv.Itoa(v)
	} else {
		new_string = "0" + strconv.Itoa(v)
	}

	order.IsFake = false
	order.Asap = year + month + new_string + "1"
	o.Insert(&order)
}

func InsertFakeOrder(new_slice []string, username string) {
	//更新order_num表
	order := models.OrderNum{}
	//o.QueryTable("order_num").OrderBy("-id").One(&order, "id")
	order.User = username
	order.State = true
	order.SaleList = strings.Join(new_slice, ",")
	str := "(" + strings.Join(new_slice, ",") + ")"

	type printsale struct {
		Num                int
		OutPrice           float64
		ConsumerName       string
		ConsumerDepartment string
		SalesmanName       string
		Title              string
		ArtNum             string
		Spec               string
		Unit               string
		BrandName          string
		Pool               string
		StoreName          string
	}
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("sale.num_fake as num", "sale.out_price_fake as out_price", "consumer.name as consumer_name",
		"consumer.department as consumer_department",
		"user.name as salesman_name", "product.title", "product.art_num", "product.spec",
		"product.unit", "brand.name as brand_name", "store.pool", "store.name as store_name").
		From("sale").
		InnerJoin("product").
		On("product.id = sale.product_id").
		InnerJoin("user").
		On("user.id = sale.salesman_id").
		InnerJoin("store").
		On("store.id = sale.store_id").
		InnerJoin("brand").
		On("brand.id = product.brand_id").
		InnerJoin("consumer").
		On("consumer.id = sale.consumer_id").
		Where("sale.id in " + str)
	sql := qb.String()
	ps := []printsale{}
	o := orm.NewOrm()
	o.Raw(sql).QueryRows(&ps)
	order.Consumer = ps[0].ConsumerName
	order.Department = ps[0].ConsumerDepartment
	order.Salesman = ps[0].SalesmanName

	//计算总额
	var total float64
	for _, item := range ps {
		total = total + item.OutPrice*float64(item.Num)
	}
	order.Sum = fmt.Sprintf("%0.2f", total)

	//ASAP拼凑字符串
	//首先获取order数据库数据，便于取出最新的ASAP单号
	order_list := []models.OrderNum{}
	o.QueryTable("order_num").OrderBy("-id").All(&order_list)

	length := len(order_list[0].Asap)
	new_string := order_list[0].Asap[:length-1] + "0"

	order.IsFake = true
	order.Asap = new_string
	o.Insert(&order)
}
