package controllers

import (
	"html/template"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/qshuai/cswm/models"
	"github.com/qshuai/cswm/plugins/permission"
)

type CategoryController struct {
	beego.Controller
}

//分类列表数据前台展示
func (c *CategoryController) Get() {
	o := orm.NewOrm()
	category := []models.Category{}
	o.QueryTable("category").Filter("is_hidden", false).OrderBy("primary", "two_stage").All(&category)

	if len(category) == 0 {
		c.Data["msg"] = "分类表数据为空，请联系管理员添加分类信息～"
		c.Data["url"] = "/"
		c.TplName = "jump/error.html"
		return
	}

	//i, j分别为一级分类和二级分类数量
	var i, j int64 = 0, 0
	for _, item := range category {
		if item.TwoStage == "-" {
			i++
		}
		if item.TwoStage != "-" && item.ThreeStage == "-" {
			j++
		}
	}

	//primary, two_stage分别为一级分类和二级分类的以数据库Id为索引，以分类名称为值的map
	primary := make(map[int]string, i)
	two_stage := make(map[int]string, j)
	for _, item := range category {
		if item.TwoStage == "-" {
			primary[item.Id] = item.Primary
		}
		if item.TwoStage != "-" && item.ThreeStage == "-" {
			two_stage[item.Id] = item.TwoStage
		}
	}

	c.Data["category"] = category
	c.Data["primary"] = primary
	c.Data["two_stage"] = two_stage
	c.Layout = "common.tpl"
	c.TplName = "category/category_list.html"
}

//提交分类表excel界面展示
func (c *CategoryController) Category_upload() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "OperateCategory") {
		c.Abort("401")
	}
	c.Layout = "common.tpl"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "category/category_upload.html"
}

//分类表excel文件上传，以及更新数据库分类表
func (c *CategoryController) Category_upload_post() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "OperateCategory") {
		c.Abort("401")
	}
	f, h, err := c.GetFile("category_file")
	if err != nil {
		logs.Error("用户ID：", c.GetSession("uid"), "上传category_file失败，原因:", err)
		c.Data["msg"] = "上传文件出错，请检查后重试～"
		c.Data["url"] = "/category_upload"
		c.TplName = "jump/error.html"
	}
	defer f.Close()
	if h.Header.Get("Content-Type") == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		//文件上传
		filename := strconv.Itoa(c.GetSession("uid").(int)) + "_" + h.Filename
		c.SaveToFile("category_file", "static/upload/"+filename)

		//xlsx文件解析
		xlFile, err := excelize.OpenFile("./static/upload/" + filename)
		if err != nil {
			logs.Error("360EntSecGroup-Skylar/excelize：读取xlsx文件失败->", err)
			c.Data["msg"] = "读取.xlsx文件失败，请检查后重试～"
			c.Data["url"] = "/category_upload"
			c.TplName = "jump/error.html"
			return
		}

		rows := xlFile.GetRows("sheet2")
		rowsnum := len(rows)
		catagory := make([]models.Category, rowsnum)
		var i, j int = 0, 0
		for _, row := range rows {
			temp := make([]string, 5)
			j = 0
			for _, colCell := range row {
				temp[j] = colCell
				j++
			}

			//将Id转化为int类型
			id, _ := strconv.Atoi(temp[0])

			//讲is_hidden转化为bool类型值
			var is_hidden bool = false
			if temp[4] == "1" {
				is_hidden = true
			}

			catagory[i] = models.Category{Id: id, Primary: temp[1], TwoStage: temp[2], ThreeStage: temp[3], Is_hidden: is_hidden}

			i++
		}
		catagory = catagory[1:]

		//数据库操作
		o := orm.NewOrm()
		o.Raw("truncate table category").Exec()
		nums, err := o.InsertMulti(100, catagory)
		if err != nil {
			logs.Error(err)
		} else {
			c.Data["msg"] = "数据库分类表共插入" + strconv.Itoa(int(nums)) + "条数据～"
			c.Data["url"] = "/category_list"
			c.TplName = "jump/success.html"
			return
		}

	} else {
		c.Data["msg"] = "请上传.xlsx为扩展名的excel文件～"
		c.Data["url"] = "/category_upload"
		c.TplName = "jump/error.html"
	}
}

//添加分类
func (c *CategoryController) Category_add() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "OperateCategory") {
		c.Abort("401")
	}
	category := []models.Category{}
	o := orm.NewOrm()

	//将数据一次性查询出来
	o.QueryTable("category").Filter("three_stage", "-").All(&category)

	var primary_string string
	var two_stage_string string
	for _, item := range category {
		if item.TwoStage == "-" {
			primary_string += item.Primary + ", "
		} else {
			two_stage_string += item.TwoStage + ", "
		}
	}

	c.Data["primary_string"] = primary_string
	c.Data["two_stage_string"] = two_stage_string
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "category/category_add.html"
}

//添加分类提交
func (c *CategoryController) Category_add_post() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "OperateCategory") {
		c.Abort("401")
	}
	primary := c.GetString("primary")
	two_stage := c.GetString("two_stage")
	three_stage := c.GetString("three_stage")

	category := models.Category{}
	category.Primary = primary
	category.TwoStage = two_stage
	category.ThreeStage = three_stage

	o := orm.NewOrm()
	primary_query := models.Category{}
	o.QueryTable("category").Filter("primary", primary).
		One(&primary_query, "id", "is_hidden")
	//如果查询不到将会返回对应类型的零值
	if primary_query.Id != 0 {
		category.Primary = strconv.Itoa(primary_query.Id)
	}

	two_stage_query := models.Category{}
	o.QueryTable("category").Filter("two_stage", two_stage).
		One(&two_stage_query, "id", "is_hidden")
	//如果查询不到将会返回对应类型的零值
	if two_stage_query.Id != 0 {
		category.TwoStage = strconv.Itoa(two_stage_query.Id)
	}

	//判断是否隐藏
	if primary_query.Is_hidden || two_stage_query.Is_hidden {
		category.Is_hidden = true
	} else {
		category.Is_hidden = false
	}

	_, err := o.Insert(&category)
	if err != nil {
		c.Data["url"] = "/category_add"
		c.Data["msg"] = "添加分类失败~"
		c.TplName = "jump/error.html"
		return
	} else {
		c.Data["url"] = "/category_list"
		c.Data["msg"] = "添加分类成功~"
		c.TplName = "jump/success.html"
	}
}

//分类编辑
func (c *CategoryController) Category_edit() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "OperateCategory") {
		c.Abort("401")
	}
	category := []models.Category{}
	o := orm.NewOrm()
	//这里只对常用的三级分类添加检索便利
	o.QueryTable("category").Exclude("three_stage", "-").All(&category, "three_stage")

	var three_stage_string string
	for _, item := range category {
		three_stage_string += item.ThreeStage + ", "
	}

	c.Data["three_stage_string"] = three_stage_string
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "common.tpl"
	c.TplName = "category/category_edit.html"
}

//ajax请求1条分类信息
func (c *CategoryController) Category_search() {
	if !c.IsAjax() {
		return
	}
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "OperateCategory") {
		c.Abort("401")
	}

	o := orm.NewOrm()
	category := models.Category{}
	primary := models.Category{}
	two_stage := models.Category{}
	three_stage := models.Category{}

	item := c.GetString("item")
	stage := c.GetString("stage")

	switch stage {
	case "primary":
		o.QueryTable("category").Filter("primary", item).One(&primary)
		o.QueryTable("category").Filter("id", primary.TwoStage).One(&two_stage, "two_stage")
		o.QueryTable("category").Filter("id", primary.ThreeStage).One(&three_stage, "three_stage")
		category.Id = primary.Id
		category.Is_hidden = primary.Is_hidden
	case "two_stage":
		o.QueryTable("category").Filter("two_stage", item).One(&two_stage)
		o.QueryTable("category").Filter("id", two_stage.Primary).One(&primary, "primary")
		o.QueryTable("category").Filter("id", two_stage.ThreeStage).One(&three_stage, "three_stage")
		category.Id = two_stage.Id
		category.Is_hidden = two_stage.Is_hidden
	case "three_stage":
		o.QueryTable("category").Filter("three_stage", item).One(&three_stage)
		o.QueryTable("category").Filter("id", three_stage.TwoStage).One(&two_stage, "two_stage")
		o.QueryTable("category").Filter("id", three_stage.Primary).One(&primary, "primary")
		category.Id = three_stage.Id
		category.Is_hidden = three_stage.Is_hidden
	}

	category.Primary = primary.Primary
	category.TwoStage = two_stage.TwoStage
	category.ThreeStage = three_stage.ThreeStage
	c.Data["json"] = category
	c.ServeJSON()
}

//分类编辑post
func (c *CategoryController) Category_edit_post() {
	if !permission.GetOneItemPermission(c.GetSession("username").(string), "OperateCategory") {
		c.Abort("401")
	}

	category := models.Category{}
	category.Id, _ = c.GetInt("category_id")
	category.Primary = c.GetString("primary")
	category.TwoStage = c.GetString("two_stage")
	category.ThreeStage = c.GetString("three_stage")

	if category.Primary == "-" {
		c.Data["url"] = "/category_edit"
		c.Data["msg"] = "注意：无此分类信息，请重新检索~"
		c.TplName = "jump/error.html"
		return
	}

	o := orm.NewOrm()
	if category.TwoStage != "-" && category.ThreeStage == "-" {
		category_primary := models.Category{}
		o.QueryTable("category").Filter("primary", category.Primary).
			Filter("two_stage", "-").One(&category_primary, "id")
		category.Primary = strconv.Itoa(category_primary.Id)
	}

	if category.ThreeStage != "-" {
		category_primary := models.Category{}
		o.QueryTable("category").Filter("primary", category.Primary).
			Filter("two_stage", "-").One(&category_primary, "id")
		category.Primary = strconv.Itoa(category_primary.Id)

		category_two_stage := models.Category{}
		o.QueryTable("category").Filter("two_stage", category.TwoStage).
			Filter("three_stage", "-").One(&category_two_stage, "id")
		category.TwoStage = strconv.Itoa(category_two_stage.Id)
	}

	is_hidden := c.GetString("is_hidden")
	if is_hidden == "0" {
		category.Is_hidden = true
	} else {
		category.Is_hidden = false
	}

	_, err := o.Update(&category)
	if err != nil {
		c.Data["url"] = "/category_edit"
		c.Data["msg"] = "更改分类信息失败~"
		c.TplName = "jump/error.html"
		return
	}
	c.Data["url"] = "/category_list"
	c.Data["msg"] = "更改分类信息成功~"
	c.TplName = "jump/success.html"
}
