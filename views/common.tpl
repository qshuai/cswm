<!--_meta 作为公共模版分离出去-->
<!DOCTYPE HTML>
<html>
<head>
    <meta charset="utf-8">
    <meta name="renderer" content="webkit|ie-comp|ie-stand">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="viewport" content="width=device-width,initial-scale=1,minimum-scale=1.0,maximum-scale=1.0,user-scalable=no" />
    <meta http-equiv="Cache-Control" content="no-siteapp" />
    <meta name="_xsrf" content="{{.xsrf_token}}" />
    <!--[if lt IE 9]>
    <script type="text/javascript" src="/static/lib/html5.js"></script>
    <script type="text/javascript" src="/static/lib/respond.min.js"></script>
    <![endif]-->
    <link rel="shortcut icon" href="/static/img/icon.png" />
    <link rel="stylesheet" type="text/css" href="/static/static/h-ui/css/H-ui.min.css" />
    <link rel="stylesheet" type="text/css" href="/static/static/h-ui.admin/css/H-ui.admin.css" />
    <link rel="stylesheet" type="text/css" href="/static/lib/Hui-iconfont/1.0.8/iconfont.css" />
    <link rel="stylesheet" type="text/css" href="/static/static/h-ui.admin/skin/default/skin.css" id="skin" />
    <link rel="stylesheet" type="text/css" href="/static/static/h-ui.admin/css/style.css" />

    <!--自动补全插件-->
    <link rel="stylesheet" type="text/css" href="/static/awesomplete/awesomplete.css">

    <!--datepicker插件-->
    <link rel="stylesheet" type="text/css" href="/static/datepicker/pikaday.css">

    <!--自定义css文件-->
    <link rel="stylesheet" type="text/css" href="/static/css/my.css">

    <!--[if IE 6]>
    <script type="text/javascript" src="http://lib.h-ui.net/DD_belatedPNG_0.0.8a-min.js" ></script>
    <script>DD_belatedPNG.fix('*');</script>
    <![endif]-->
    <!--/meta 作为公共模版分离出去-->

    <title>食物链系统</title>
</head>
<body>
<!--_header 作为公共模版分离出去-->
<header class="navbar-wrapper">
    <div class="navbar navbar-fixed-top">
        <div class="container-fluid cl">
            <a class="logo navbar-logo f-l mr-10 hidden-xs" href="/">食物链系统</a>
            <a class="logo navbar-logo-m f-l mr-10 visible-xs" href="/">食物链系统</a>
            <span class="logo navbar-slogan f-l mr-10 hidden-xs">v1.3</span>
            <a aria-hidden="false" class="nav-toggle Hui-iconfont visible-xs" href="javascript:;">&#xe667;</a>
            <nav class="nav navbar-nav">
                <ul class="cl">
                    <li class="dropDown dropDown_hover"><a href="javascript:;" class="dropDown_A"><i class="Hui-iconfont">&#xe600;</i> 新增 <i class="Hui-iconfont">&#xe6d5;</i></a>
                        <ul class="dropDown-menu menu radius box-shadow">
							{{if .authority.AddProduct}}
                            	<li><a href="/product_add"><i class="Hui-iconfont">&#xe620;</i> 产品</a></li>
							{{end}}
							{{if .authority.AddMember}}
                            	<li><a href="/member_add"><i class="Hui-iconfont">&#xe60d;</i> 用户</a></li>
							{{end}}
                            <li><a href="/message_add"><i class="Hui-iconfont">&#xe60d;</i> 消息</a></li>
                        </ul>
                    </li>
                </ul>
            </nav>
            <nav id="Hui-userbar" class="nav navbar-nav navbar-userbar hidden-xs">
                <ul class="cl">
                    <li>{{.grade}}</li>
                    <li class="dropDown dropDown_hover"> <a href="#" class="dropDown_A">{{.username}} <i class="Hui-iconfont">&#xe6d5;</i></a>
                        <ul class="dropDown-menu menu radius box-shadow">
                            <li><a href="/member_edit">个人信息</a></li>
                            <li><a href="/logout">退出</a></li>
                        </ul>
                    </li>
                    <li id="Hui-msg"> <a href="/message_list" title="消息"><span class="badge badge-danger">{{if ne .message_num 0}}{{.message_num}}{{end}}</span><i class="Hui-iconfont" style="font-size:18px">&#xe68a;</i></a> </li>
                    <li id="Hui-skin" class="dropDown right dropDown_hover"> <a href="javascript:;" class="dropDown_A" title="换肤"><i class="Hui-iconfont" style="font-size:18px">&#xe62a;</i></a>
                        <ul class="dropDown-menu menu radius box-shadow">
                            <li><a href="javascript:;" data-val="default" title="默认（黑色）">默认（黑色）</a></li>
                            <li><a href="javascript:;" data-val="blue" title="蓝色">蓝色</a></li>
                            <li><a href="javascript:;" data-val="green" title="绿色">绿色</a></li>
                            <li><a href="javascript:;" data-val="red" title="红色">红色</a></li>
                            <li><a href="javascript:;" data-val="yellow" title="黄色">黄色</a></li>
                            <li><a href="javascript:;" data-val="orange" title="橙色">橙色</a></li>
                        </ul>
                    </li>
                </ul>
            </nav>
        </div>
    </div>
</header>
<!--/_header 作为公共模版分离出去-->

<!--_menu 作为公共模版分离出去-->
<aside class="Hui-aside">

    <div class="menu_dropdown bk_2">
        <dl id="menu-product" class="nav-left-list">
            <dt><i class="Hui-iconfont">&#xe66a;</i> 库存管理<i class="Hui-iconfont menu_dropdown-arrow">&#xe6d5;</i></dt>
            <dd>
                <ul>
					{{if .authority.ViewStore}}
                    	<li><a href="/store_list" title="库房列表">库房列表</a></li>
					{{end}}
					{{if .authority.AddStore}}
                    	<li><a href="/store_add" title="库房列表">添加库房</a></li>
					{{end}}
                </ul>
            </dd>
        </dl>
		{{if .authority.ViewMove}}
        <dl id="menu-picture" class="nav-left-list">
            <dt><i class="Hui-iconfont">&#xe655;</i> 移库管理<i class="Hui-iconfont menu_dropdown-arrow">&#xe6d5;</i></dt>
            <dd>
                <ul>
                    <li><a href="/move_list" title="图片管理">移库记录</a></li>
                </ul>
            </dd>
        </dl>
		{{end}}

        <dl id="menu-product" class="nav-left-list">
            <dt><i class="Hui-iconfont">&#xe681;</i> 分类管理<i class="Hui-iconfont menu_dropdown-arrow">&#xe6d5;</i></dt>
            <dd>
                <ul>
                    <li><a href="/category_list" title="产品管理">分类列表</a></li>
					{{if .authority.OperateCategory}}
                    	<li><a href="/category_upload" title="产品管理">更新分类</a></li>
                    	<li><a href="/category_add" title="产品管理">添加分类</a></li>
                    	<li><a href="/category_edit" title="产品管理">编辑分类</a></li>
					{{end}}
                </ul>
            </dd>
        </dl>
        <dl id="menu-picture" class="nav-left-list">
            <dt><i class="Hui-iconfont">&#xe620;</i> 产品管理<i class="Hui-iconfont menu_dropdown-arrow">&#xe6d5;</i></dt>
            <dd>
                <ul>
                    <li><a href="/product_list" title="图片管理">产品列表</a></li>
					{{if .authority.AddProduct}}
                    	<li><a href="/product_add" title="图片管理">产品录入</a></li>
					{{end}}
					{{if eq .grade "超级管理员"}}
						<li><a href="/product_template_list">模板列表</a></li>
						<li><a href="/product_template_add">模板录入</a></li>
					{{end}}
					{{if eq .grade "总库管理员"}}
						<li><a href="/product_template_list">模板列表</a></li>
						<li><a href="/product_template_add">模板录入</a></li>
					{{end}}
                </ul>
            </dd>
        </dl>
        <dl id="menu-picture" class="nav-left-list">
            <dt><i class="Hui-iconfont">&#xe61e;</i> 销售管理<i class="Hui-iconfont menu_dropdown-arrow">&#xe6d5;</i></dt>
            <dd>
                <ul>
					{{if .authority.ViewSale}}
                    	<li><a href="/sale_list" title="销售记录">销售列表</a></li>
                    	<li><a href="/order_list" title="销售记录">出库单列表</a></li>
					{{end}}
                </ul>
            </dd>
        </dl>
        <dl id="menu-picture" class="nav-left-list">
            <dt><i class="Hui-iconfont">&#xe60d;</i> 供应商管理<i class="Hui-iconfont menu_dropdown-arrow">&#xe6d5;</i></dt>
            <dd>
                <ul>
					{{if .authority.ViewSupplier}}
                    	<li><a href="/supplier_list" title="供应商列表">供应商列表</a></li>
					{{end}}
					{{if .authority.AddSupplier}}
                    	<li><a href="/supplier_add" title="添加供应商">添加供应商</a></li>
					{{end}}
                </ul>
            </dd>
        </dl>
        <dl id="menu-picture" class="nav-left-list">
            <dt><i class="Hui-iconfont">&#xe62c;</i> 经销商管理<i class="Hui-iconfont menu_dropdown-arrow">&#xe6d5;</i></dt>
            <dd>
                <ul>
					{{if .authority.ViewDealer}}
						<li><a href="/dealer_list" title="经销商列表">经销商列表</a></li>
					{{end}}
					{{if .authority.AddDealer}}
                    	<li><a href="/dealer_add" title="添加经销商">添加经销商</a></li>
					{{end}}
                </ul>
            </dd>
        </dl>
        <dl id="menu-picture" class="nav-left-list">
            <dt><i class="Hui-iconfont">&#xe611;</i> 客户管理<i class="Hui-iconfont menu_dropdown-arrow">&#xe6d5;</i></dt>
            <dd>
                <ul>
					{{if .authority.ViewConsumer}}
                    	<li><a href="/consumer_list" title="客户列表">客户列表</a></li>
					{{end}}
					{{if .authority.AddConsumer}}
                    	<li><a href="/consumer_add" title="添加客户">添加客户</a></li>
					{{end}}
                </ul>
            </dd>
        </dl>
        <dl id="menu-picture" class="nav-left-list">
            <dt><i class="Hui-iconfont">&#xe6d3;</i> 品牌管理<i class="Hui-iconfont menu_dropdown-arrow">&#xe6d5;</i></dt>
            <dd>
                <ul>
                    <li><a href="/brand_list" title="品牌列表">品牌列表</a></li>
					{{if .authority.AddBrand }}
                    	<li><a href="/brand_add" title="添加品牌">添加品牌</a></li>
					{{end}}
                </ul>
            </dd>
        </dl>
        <dl id="menu-member" class="nav-left-list">
            <dt><i class="Hui-iconfont">&#xe602;</i> 员工管理<i class="Hui-iconfont menu_dropdown-arrow">&#xe6d5;</i></dt>
            <dd>
                <ul>
                    <li><a href="/member_list" title="员工列表">员工列表</a></li>
					{{if .authority.AddMember}}
                    	<li><a href="/member_add" title="添加员工">添加员工</a></li>
					{{end}}
					{{if .authority.ActivieMember}}
                    	<li><a href="/disable_member_list" title="禁用账号">禁用账号</a></li>
					{{end}}
                    	<li><a href="/member_edit" title="自我修改">自我修改</a></li>
                </ul>
            </dd>
        </dl>
        <dl id="menu-comments" class="nav-left-list">
            <dt><i class="Hui-iconfont">&#xe622;</i> 消息管理<i class="Hui-iconfont menu_dropdown-arrow">&#xe6d5;</i></dt>
            <dd>
                <ul>
                    <li><a href="/message_list" title="消息列表">消息列表</a></li>
                    <li><a href="/message_add" title="新建消息">新建消息</a></li>
                </ul>
            </dd>
        </dl>

		{{if eq .grade "超级管理员"}}
        <dl id="menu-comments" class="nav-left-list">
            <dt><i class="Hui-iconfont">&#xe61d;</i> 权限管理<i class="Hui-iconfont menu_dropdown-arrow">&#xe6d5;</i></dt>
            <dd>
                <ul>
                    <li><a href="/default_permission" title="默认权限">默认权限</a></li>
					<li><a href="/permission_member_list" title="员工列表">员工列表</a></li>
                </ul>
            </dd>
        </dl>
		{{end}}
    </div>
</aside>
<div class="dislpayArrow hidden-xs"><a class="pngfix" href="javascript:void(0);" onClick="displaynavbar(this)"></a></div>
<!--/_menu 作为公共模版分离出去-->
<div id="fixbug">
{{.LayoutContent}}

<!--_footer 作为公共模版分离出去-->
<script type="text/javascript" src="/static/lib/jquery/1.9.1/jquery.min.js"></script>
<script type="text/javascript" src="/static/lib/layer/2.4/layer.js"></script>
<script type="text/javascript" src="/static/static/h-ui/js/H-ui.min.js"></script>
<script type="text/javascript" src="/static/static/h-ui.admin/js/H-ui.admin.page.js"></script>

<!--城市列表 三级联动-->
<script type="text/javascript" src="/static/province_list/jquery.provincesCity.js"></script>
<script type="text/javascript" src="/static/province_list/provincesData.js"></script>

<!--自动补全插件-->
<script type="text/javascript" src="/static/awesomplete/awesomplete.js"></script>

<!--datepicker插件-->
<script type="text/javascript" src="/static/datepicker/moment.min.js"></script>
<script type="text/javascript" src="/static/datepicker/pikaday.js"></script>

<!--jquery paginator-->
<script type="text/javascript" src="/static/paginator/jqPaginator.js"></script>

<!--验证js-->
<script src="http://static.runoob.com/assets/jquery-validation-1.14.0/dist/jquery.validate.min.js"></script>
<script src="http://static.runoob.com/assets/jquery-validation-1.14.0/dist/localization/messages_zh.js"></script>

<!--自定义js文件-->
<script type="text/javascript" src="/static/js/my.js"></script>

<!--/_footer /作为公共模版分离出去-->

</body>
</html>
