package routers

import (
	"erp/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/logout", &controllers.LoginController{}, "get:Logout")

	{
		//添加会员
		beego.Router("/member_add", &controllers.MemberController{}, "get:Member_add")
		beego.Router("/member_add", &controllers.MemberController{}, "post:Member_add_post")
		//完善信息
		beego.Router("/userinfo", &controllers.MemberController{}, "get:UserInfo")
		beego.Router("/userinfo", &controllers.MemberController{}, "post:UserInfo_post")
		//用户列表
		beego.Router("/member_list", &controllers.MemberController{}, "get:Member_list")
		//个人信息修改
		beego.Router("/member_edit", &controllers.MemberController{}, "get:Member_edit")
		beego.Router("/member_edit", &controllers.MemberController{}, "post:Member_edit_post")
		//管理员修改人员信息
		beego.Router("/admin_member_edit/?:uid:int", &controllers.MemberController{}, "get:Admin_member_edit")
		beego.Router("/admin_member_edit", &controllers.MemberController{}, "post:Admin_member_edit_post")  //ajax请求用户检索
		beego.Router("/disable_active_user", &controllers.MemberController{}, "post:Disable_active_member") //ajax激活或禁用账户
		beego.Router("/admin_edit_all", &controllers.MemberController{}, "post:Admin_edit_all")             //修改全部信息
		beego.Router("/disable_member_list", &controllers.MemberController{}, "get:Disable_member_list")    //获取禁用账户列表
		beego.Router("/off_position", &controllers.MemberController{}, "post:OffPosition")                  //获取禁用账户列表
	}

	{
		//商品管理
		beego.Router("/product_list", &controllers.ProductController{}, "get:Get")
		beego.Router("/product_list", &controllers.ProductController{}, "post:Get")
		beego.Router("/product_add", &controllers.ProductController{}, "get:Add_get")
		beego.Router("/product_add", &controllers.ProductController{}, "post:Add_post")

		//ajax
		//通过货号搜索商品
		beego.Router("/searchByCatnum", &controllers.ProductController{}, "post:SearchByCatnum")
		//删除单条商品信息
		beego.Router("/product_item_delete", &controllers.ProductController{}, "post:Product_item_delete")
		//编辑单条商品信息
		beego.Router("/product_item_edit", &controllers.ProductController{}, "post:Product_item_edit")
		//商品跟踪
		beego.Router("/product_track/:pid:int", &controllers.ProductController{}, "get:Product_track")

		beego.Router("/product_template_list", &controllers.ProductController{}, "get:ProductTemplateList")
		beego.Router("/product_template_add", &controllers.ProductController{}, "get:ProductTemplateAdd")
		beego.Router("/product_template_add", &controllers.ProductController{}, "post:ProductTemplateAddPost")
		beego.Router("/product_template_edit", &controllers.ProductController{}, "post:ProductTemplateEditPost")
		beego.Router("/product_template_delete", &controllers.ProductController{}, "post:ProductTemplateDeletePost")
		//加载更多商品
		beego.Router("/product_load_more", &controllers.ProductController{}, "post:ProductLoadMore")
		//加载更多模板
		beego.Router("/template_load_more", &controllers.ProductController{}, "post:TemplateLoadMore")
	}

	{
		//分类管理
		beego.Router("/category_list", &controllers.CategoryController{})
		beego.Router("/category_upload", &controllers.CategoryController{}, "get:Category_upload")
		beego.Router("/category_upload", &controllers.CategoryController{}, "post:Category_upload_post")

		beego.Router("/category_add", &controllers.CategoryController{}, "get:Category_add")       //添加分类
		beego.Router("/category_add", &controllers.CategoryController{}, "post:Category_add_post") //添加分类提交

		beego.Router("/category_edit", &controllers.CategoryController{}, "get:Category_edit")           //编辑分类页面
		beego.Router("/category_edit", &controllers.CategoryController{}, "post:Category_edit_post")     //编辑分类页面
		beego.Router("/category_search_ajax", &controllers.CategoryController{}, "post:Category_search") //ajax请求分类信息
	}

	{
		//供应商
		beego.Router("/supplier_list", &controllers.SupplierController{})
		beego.Router("/supplier_add", &controllers.SupplierController{}, "get:Supplier_add")
		beego.Router("/supplier_add", &controllers.SupplierController{}, "post:Supplier_add_post")
		beego.Router("/supplier_edit", &controllers.SupplierController{}, "post:Supplier_edit_post")
	}

	{
		//经销商
		beego.Router("/dealer_list", &controllers.DealerController{})
		beego.Router("/dealer_add", &controllers.DealerController{}, "get:Dealer_add")
		beego.Router("/dealer_add", &controllers.DealerController{}, "post:Dealer_add_post")
	}

	{
		//品牌
		beego.Router("/brand_list", &controllers.BrandController{})
		beego.Router("/brand_add", &controllers.BrandController{}, "get:Brand_add")
		beego.Router("/brand_add", &controllers.BrandController{}, "post:Brand_add_post")
	}

	{
		//客户
		beego.Router("/consumer_list", &controllers.ConsumerController{})
		beego.Router("/consumer_add", &controllers.ConsumerController{}, "get:Consumer_add")
		beego.Router("/consumer_add", &controllers.ConsumerController{}, "post:Consumer_add_post")
		beego.Router("/consumer_edit", &controllers.ConsumerController{}, "post:Consumer_edit")
	}

	{
		//库房
		beego.Router("/store_list", &controllers.StoreController{})
		beego.Router("/store_add", &controllers.StoreController{}, "get:Store_add")
		beego.Router("/store_add", &controllers.StoreController{}, "post:Store_add_post")

		//出库
		beego.Router("/store_output_action/?:pid:int", &controllers.StoreoutController{}, "get:Store_out_action") //商品出库录入
		beego.Router("/store_output_action", &controllers.StoreoutController{}, "post:Store_out_action_post")     //商品出库录入

		//销售
		beego.Router("/sale_list", &controllers.SaleController{}, "get:Sale_list") //销售记录列表
		beego.Router("/sale_edit", &controllers.SaleController{}, "post:Sale_edit")
		beego.Router("/sale_load_more", &controllers.SaleController{}, "post:SaleLoadMore")
		beego.Router("/print_action/:list/:id", &controllers.SaleController{}, "get:Print")
		beego.Router("/order_list", &controllers.SaleController{}, "get:OrderList")
		beego.Router("/order_close/:oid:int", &controllers.SaleController{}, "get:OrderClose")
		beego.Router("/order_price_edit", &controllers.SaleController{}, "post:OrderPriceEdit")

		beego.Router("/product_sale_info/:art_num", &controllers.SaleController{}, "get:ProductSalInfo")
		beego.Router("/order_list/add", &controllers.SaleController{}, "get:OrderAdd")
		beego.Router("/order_edit", &controllers.SaleController{}, "get:OrderEdit")
		beego.Router("/order_edit_post", &controllers.SaleController{}, "get:OrderEditPost")

		//移库
		beego.Router("/move_request/:pid:int", &controllers.MoveController{}, "get:Move_request") //移库请求页面
		beego.Router("/move_request", &controllers.MoveController{}, "post:Move_request_post")    //移库请求post
		beego.Router("/move_list", &controllers.MoveController{}, "get:Move_list")                //移库列表
		beego.Router("/move_accept", &controllers.MoveController{}, "post:Move_accept")           //移库接受
		beego.Router("/move_deny", &controllers.MoveController{}, "post:Move_deny")               //移库拒绝
		beego.Router("/move_finish", &controllers.MoveController{}, "post:Move_finish")           //移库完成
		beego.Router("/move_info/:mid:int", &controllers.MoveController{}, "get:Move_info")       //移库详情

		//消息
		beego.Router("/message_list", &controllers.MessageController{}, "get:Message_list")          //获取消息列表
		beego.Router("/message_add", &controllers.MessageController{}, "get:Message_add")            //新建消息页面
		beego.Router("/message_add", &controllers.MessageController{}, "post:Message_add_post")      //提交新建消息
		beego.Router("/message_info/:mid:int", &controllers.MessageController{}, "get:Message_info") //消息详情页面

		//权限管理
		beego.Router("/default_permission", &controllers.Permission{}, "get:DefaultPermission")
		beego.Router("/default_permission_edit/:item:int", &controllers.Permission{}, "get:DefaultPermissionEdit")
		beego.Router("/default_permission_edit", &controllers.Permission{}, "post:DefaultPermissionEditPost")
		beego.Router("/permission_member_list", &controllers.Permission{}, "get:PermissionMemberList")
		beego.Router("/permission_member_edit/:uid:int", &controllers.Permission{}, "get:PermissionMemberEdit")
		beego.Router("/permission_member_edit", &controllers.Permission{}, "post:PermissionMemberEditPost")
	}

}
