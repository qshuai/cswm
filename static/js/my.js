/**
 * Created by scrapup on 2017/9/12.
 */

//全局js
	//判断当前页面是否展示侧边栏
var query_url = window.location.pathname;
if (query_url === "/product_list" || query_url === "/sale_list" || query_url === "/product_template_list") {
	$($(".pngfix")).addClass("open");
	$("body").addClass("big-page");
}

//如果左侧导航栏条目下无可选内容，则隐藏改条目
var navlis = $(".nav-left-list");
$.each(navlis, function (index) {
	if (navlis.eq(index).find("li").length === 0) {
		navlis.eq(index).addClass("hide")
	}
});

//product_add.html
//-----------------------------------------------------------------------------------
if (query_url === "/product_add") {
	var disable = false, picker = new Pikaday({
		field: document.getElementById('datepicker'),
		firstDay: 1,
		minDate: new Date(2000, 0, 1),
		maxDate: new Date(),
		yearRange: [2000, 2030],

		showDaysInNextAndPreviousMonths: true,
		enableSelectionDaysInNextAndPreviousMonths: true
	});

//为datepicker初始化为当前日期
	var time = new Date();
	$("#datepicker").val(time.getFullYear() + "-" + (time.getMonth() + 1) + "-" + time.getDate());

//定义变量
	var delete_sku;
	var spec;
	var stock;
	var in_price;

	//product_template_add.html增加sku
	function AddSku() {
		$("#add_spec").parent("div").append('<input type="text" class="input-text mt-10 spec" value="" placeholder="规格" name="spec" style="width: 40%"> <input type="text" class="input-text mt-10 stock" value="" placeholder="数量" id="stock" name="stock" style="width: 20%"> <input type="text" class="input-text mt-10 in_price" value="" placeholder="价格" id="in_price" name="in_price" style="width: 20%"> <a class="btn btn-danger-outline radius delete_sku mt-10"><i class="Hui-iconfont Hui-iconfont-close"></i></a>')

		delete_sku = $(".delete_sku");
		spec = $(".spec");
		stock = $(".stock");
		in_price = $(".in_price");

		$.each(delete_sku, function (index) {
			delete_sku.eq(index).click(function () {
				index += 1;
				spec.eq(index).remove();
				stock.eq(index).remove();
				in_price.eq(index).remove();
				$(this).remove();
			})
		});
	}

	//product_add.html 增加sku
	function AddProductSku() {
		$("#add_spec").append('<input type="text" readonly class="input-text mt-10 spec" value="" placeholder="规格" name="spec" style="width: 40%"> <input type="text" class="input-text mt-10 stock" value="" placeholder="数量" id="stock" name="stock" style="width: 20%"> <input type="text" class="input-text mt-10 in_price" value="" placeholder="价格" id="in_price" name="in_price" style="width: 20%"> <a onclick="DeleteSku(this)" class="btn btn-danger-outline radius mt-10"><i class="Hui-iconfont Hui-iconfont-close"></i></a>')

		delete_sku = $(".delete_sku");
		spec = $(".spec");
		stock = $(".stock");
		in_price = $(".in_price");

		$.each(delete_sku, function (index) {
			delete_sku.eq(index).click(function () {
				index += 1;
				spec.eq(index).remove();
				stock.eq(index).remove();
				in_price.eq(index).remove();
				$(this).remove();
			})
		});
	}

	function DeleteSku(obj) {
		$(obj).prev().prev().prev().remove();
		$(obj).prev().prev().remove();
		$(obj).prev().remove();
		$(obj).remove();
	}


//通过货号快速填充商品信息
	var products = [];
	$("#art_num_search").click(function () {
		$.ajax({
			type: "post",
			url: "/searchByCatnum",
			data: {
				"art_num": $("#art_num").val(),
				"_xsrf": $("input[name=_xsrf]").val()
			},
			success: function (response, status, xhr) {
				products = response;
				var num = products.length;
				if (num > 0) {
					$("#result_art").text("商品名称：" + products[0].Title);
					$("#art_num").attr("readonly", true)
				} else {
					$("#result_art").text("注意：数据库中不存在此货号，请核对后再试~");
				}

			}
		})
	});

	$("#confirm_in").click(function () {
		var num = products.length;
		var radios = $(".radios").find("input");

		$("#title").val(products[0].Title);
		$("#brand").val(products[0].BrandName);
		$("#three_stage").val(products[0].ThreeStage);

		var supplier_array = (products[0].Suppliers).split(",");
		var supplier_len = supplier_array.length;
		var supplier_select = $("select[name=supplier]");
		for (var i = 0; i < supplier_len; i++) {
			supplier_select.append("<option>" + supplier_array[i] + "</option>")
		}

		$("input[name=unit]").val(products[0].Unit);
		$.each(radios, function (index) {
			if (radios.eq(index).val() === products[0].Unit) {
				$(this).attr("checked", true)
			}
		});

		for (var i = 0; i < num; i++) {
			if (i < num - 1) {
				AddProductSku();
			}
			$(".spec").eq(i).val(products[i].Spec);
			$(".in_price").eq(i).val(products[i].InPrice);
			if (products[i].InPrice !== 0){
				$(".in_price").eq(i).attr("type", "password").attr("readonly", true);
			}else{
				$(".in_price").eq(i).val("");
			}
		}

		delete_sku = $(".delete_sku");
		spec = $(".spec");
		stock = $(".stock");
		in_price = $(".in_price");

		$.each(delete_sku, function (index) {
			delete_sku.eq(index).click(function () {
				index += 1;
				spec.eq(index).remove();
				stock.eq(index).remove();
				in_price.eq(index).remove();
				$(this).remove();
			})
		});
		$("#confirm_in").remove();
	});
}

//product_list.html
//-----------------------------------------------------------------------------------
if (query_url === "/product_list") {
	//读取html的script标签中设置的全局变量
	var product = $.parseJSON(product);

	//设置每行的删除按钮集合
	var product_item_delete;

	//默认page_size为10
	var page_size = 10;

	//隐藏某些列的索引数组
	var hidden_index = [];

	var page_size_temp = $.cookie("product_paginator");
	if (page_size_temp !== undefined) {
		page_size = page_size_temp
	}

	$.cookie("product_offset", 0);
	$.cookie("product_current_page", 1);

	var product_rows = $("#product_row");
	product_rows.html("");
	var i = 0;
	product_paginator(product, '#pagination', page_size, product.length, product_rows, hidden_index);

	//用户选择每页显示的条目数，也就是page_size
	var page_size_btn = $(".page_size");
	$.each(page_size_btn, function (index) {
		page_size_btn.eq(index).click(function () {

			//通过hui-ui.js的cookie()方法直接在浏览器设置cookie减少http请求（替代以上ajax请求）
			$.cookie('product_paginator', $(this).attr("data"), {expires: 366});

			//指示为第一页
			var num = 1;

			var page_size_temp = $.cookie("product_paginator");
			if (page_size_temp !== null) {
				page_size = page_size_temp
			}

			product_paginator(product, '#pagination', page_size, product.length, product_rows, hidden_index);
		})
	});

	//隐藏某些列
	var product_item_close = $(".product_item_close");
	var product_title = $(".product_title");
	$.each(product_item_close, function (index) {
		product_item_close.eq(index).click(function () {
			$(this).parent().hide();
			product_title.find("th").eq(index).hide();

			//设置隐藏索引到全局变量
			hidden_index.push(index);
			product_paginator(product, '#pagination', page_size, product.length, product_rows, hidden_index);
		})
	});

	//排序
	var asc = true;
	var product_item_order = $(".product_item_order");
	$.each(product_item_order, function (index) {
		product_item_order.eq(index).click(function () {
			switch (index) {
				case 0:
					product.sort(function (x, y) {
						return asc ? ((x.Title < y.Title) ? -1 : ((x.Title > y.Title) ? 1 : 0)) : ((x.Title < y.Title) ? 1 : ((x.Title > y.Title) ? -1 : 0));
					});
					asc = !asc;
					break;
				case 1:
					product.sort(function (x, y) {
						return asc ? ((x.BrandName < y.BrandName) ? -1 : ((x.BrandName > y.BrandName) ? 1 : 0)) : ((x.BrandName < y.BrandName) ? 1 : ((x.BrandName > y.BrandName) ? -1 : 0));
					});
					asc = !asc;
					break;
				case 2:
					product.sort(function (x, y) {
						return asc ? ((x.SupplierName < y.SupplierName) ? -1 : ((x.SupplierName > y.SupplierName) ? 1 : 0)) : ((x.SupplierName < y.SupplierName) ? 1 : ((x.SupplierName > y.SupplierName) ? -1 : 0));
					});
					asc = !asc;
					break;
				case 3:
					product.sort(function (x, y) {
						return asc ? ((x.ArtNum < y.ArtNum) ? -1 : ((x.ArtNum > y.ArtNum) ? 1 : 0)) : ((x.ArtNum < y.ArtNum) ? 1 : ((x.ArtNum > y.ArtNum) ? -1 : 0));
					});
					asc = !asc;
					break;
				case 4:
					product.sort(function (x, y) {
						return asc ? ((x.LotNum < y.LotNum) ? -1 : ((x.LotNum > y.LotNum) ? 1 : 0)) : ((x.LotNum < y.LotNum) ? 1 : ((x.LotNum > y.LotNum) ? -1 : 0));
					});
					asc = !asc;
					break;
				case 5:
					product.sort(function (x, y) {
						return asc ? ((x.StoreName < y.StoreName) ? -1 : ((x.StoreName > y.StoreName) ? 1 : 0)) : ((x.StoreName < y.StoreName) ? 1 : ((x.StoreName > y.StoreName) ? -1 : 0));
					});
					asc = !asc;
					break;
				case 6:
					product.sort(function (x, y) {
						return asc ? ((x.ThreeStage < y.ThreeStage) ? -1 : ((x.ThreeStage > y.ThreeStage) ? 1 : 0)) : ((x.ThreeStage < y.ThreeStage) ? 1 : ((x.ThreeStage > y.ThreeStage) ? -1 : 0));
					});
					asc = !asc;
					break;
				case 7:
					product.sort(function (x, y) {
						return asc ? ((x.Spec < y.Spec) ? -1 : ((x.Spec > y.Spec) ? 1 : 0)) : ((x.Spec < y.Spec) ? 1 : ((x.Spec > y.Spec) ? -1 : 0));
					});
					asc = !asc;
					break;
				case 8:
					product.sort(function (x, y) {
						return asc ? ((x.Unit < y.Unit) ? -1 : ((x.Unit > y.Unit) ? 1 : 0)) : ((x.Unit < y.Unit) ? 1 : ((x.Unit > y.Unit) ? -1 : 0));
					});
					asc = !asc;
					break;
				case 9:
					product.sort(function (x, y) {
						return asc ? ((x.Stock < y.Stock) ? -1 : ((x.Stock > y.Stock) ? 1 : 0)) : ((x.Stock < y.Stock) ? 1 : ((x.Stock > y.Stock) ? -1 : 0));
					});
					asc = !asc;
					break;
				case 10:
					product.sort(function (x, y) {
						return asc ? ((x.InPrice < y.InPrice) ? -1 : ((x.InPrice > y.InPrice) ? 1 : 0)) : ((x.InPrice < y.InPrice) ? 1 : ((x.InPrice > y.InPrice) ? -1 : 0));
					});
					asc = !asc;
					break;
				case 11:
					product.sort(function (x, y) {
						return asc ? ((x.HasPay < y.HasPay) ? -1 : ((x.HasPay > y.HasPay) ? 1 : 0)) : ((x.HasPay < y.HasPay) ? 1 : ((x.HasPay > y.HasPay) ? -1 : 0));
					});
					asc = !asc;
					break;
				case 12:
					product.sort(function (x, y) {
						return asc ? ((x.HasInvoice < y.HasInvoice) ? -1 : ((x.HasInvoice > y.HasInvoice) ? 1 : 0)) : ((x.HasInvoice < y.HasInvoice) ? 1 : ((x.HasInvoice > y.HasInvoice) ? -1 : 0));
					});
					asc = !asc;
					break;
				case 13:
					product.sort(function (x, y) {
						return asc ? ((x.GetInvoice < y.GetInvoice) ? -1 : ((x.GetInvoice > y.GetInvoice) ? 1 : 0)) : ((x.GetInvoice < y.GetInvoice) ? 1 : ((x.GetInvoice > y.GetInvoice) ? -1 : 0));
					});
					asc = !asc;
					break;
				case 14:
					product.sort(function (x, y) {
						return asc ? ((x.UserName < y.UserName) ? -1 : ((x.UserName > y.UserName) ? 1 : 0)) : ((x.UserName < y.UserName) ? 1 : ((x.UserName > y.UserName) ? -1 : 0));
					});
					asc = !asc;
					break;
				case 15:
					product.sort(function (x, y) {
						return asc ? ((x.InTime < y.InTime) ? -1 : ((x.InTime > y.InTime) ? 1 : 0)) : ((x.InTime < y.InTime) ? 1 : ((x.InTime > y.InTime) ? -1 : 0));
					});
					asc = !asc;
					break;
			}
			product_paginator(product, '#pagination', page_size, product.length, product_rows, hidden_index);
		})
	});

	//对product进行筛选
	var product_copy = product;
	var filter_btn = $(".filter_btn");
	filter_btn.click(function () {
		var splice_array = [];

		var art_num_filter = $("input[name=art_num_filter]").val();
		if (art_num_filter !== "") {
			$.each(product_copy, function (index, item) {
				if (item.ArtNum !== art_num_filter) {
					splice_array.push(index);
				}
			})
		}

		var brand_filter = $("input[name=brand_filter]").val();
		if (brand_filter !== "") {
			$.each(product_copy, function (index, item) {
				if (item.BrandName !== brand_filter) {
					splice_array.push(index);
				}
			})
		}

		var supplier_filter = $("input[name=supplier_filter]").val();
		if (supplier_filter !== "") {
			$.each(product_copy, function (index, item) {
				if (item.SupplierName !== supplier_filter) {
					splice_array.push(index);
				}
			})
		}

		var category_filter = $("input[name=category_filter]").val();
		if (category_filter !== "") {
			$.each(product_copy, function (index, item) {
				if (item.ThreeStage !== category_filter) {
					splice_array.push(index);
				}
			})
		}

		var spec_filter = $("input[name=spec_filter]").val();
		if (spec_filter !== "") {
			$.each(product_copy, function (index, item) {
				if (item.Spec !== spec_filter) {
					splice_array.push(index);
				}
			})
		}

		var user = $("input[name=user]").val();
		if (user !== "") {
			$.each(product_copy, function (index, item) {
				if (item.UserName !== user) {
					splice_array.push(index);
				}
			})
		}

		var store_filter = $("input[name=store_filter]").val();
		if (store_filter !== "") {
			var result = store_filter.split("-");
			$.each(product_copy, function (index, item) {
				if (!(item.Pool === result[0] && item.StoreName === result[1])) {
					splice_array.push(index);
				}
			});
		}

		var has_pay_filter = $("input[name=has_pay_filter]").val();
		switch (has_pay_filter) {
			case "yes":
				has_pay_filter = true;
				break;
			case "no":
				has_pay_filter = false;
				break;
			default:
				has_pay_filter = "";
		}
		if (has_pay_filter !== "") {
			$.each(product_copy, function (index, item) {
				if (item.HasPay !== has_pay_filter) {
					splice_array.push(index);
				}
			})
		}

		var has_invoice_filter = $("input[name=has_invioce_filter]").val();
		switch (has_invoice_filter) {
			case "yes":
				has_invoice_filter = true;
				break;
			case "no":
				has_invoice_filter = false;
				break;
			default:
				has_invoice_filter = "";
		}
		if (has_invoice_filter !== "") {
			$.each(product_copy, function (index, item) {
				if (item.HasInvoice !== has_invoice_filter) {
					splice_array.push(index);
				}
			})
		}

		var splice_array_length = splice_array.length;
		var new_splice_array = [];
		for (var i = 0; i < splice_array_length; i++) {
			if ($.inArray(splice_array[i], new_splice_array) === -1) {
				new_splice_array.push(splice_array[i])
			}
		}

		new_splice_array = new_splice_array.sort(function (x, y) {
			return x - y;
		});

		var ab = 0;
		$.each(new_splice_array, function (index, item) {
			product_copy.splice(item - ab, 1);
			ab++
		});

		product_paginator(product_copy, '#pagination', page_size, product.length, product_rows, hidden_index);
	});

	//加载更多product
	$(".load-more").click(function () {
		$.cookie("product_offset", parseInt($.cookie("product_offset")) + 1);
		$.ajax({
			url : "/product_load_more",
			type : "post",
			data : {
				"offset" : $.cookie("product_offset"),
				"_xsrf" : $("meta[name=_xsrf]").attr("content")
			},
			success : function (response) {
				product = product.concat($.parseJSON(response));
				product_paginator(product, '#pagination', page_size, product.length, product_rows, hidden_index);
			}
		})
	})
}

//分页函数（抽象）
function product_paginator(product, paginator_node, page_size, total_item, content_node_obj, hidden_index) {
	//判断数据是否为空
	if (total_item === 0) {
		$(".tip_message").text("Sorry, 商品列表为空~");

		//移除选择分页数量按钮
		$(".page_size_btn").remove();

		//移除表格table
		$(".product_table").remove()
	}

	//计算page_num
	var page_num;
	if (total_item % page_size === 0) {
		page_num = total_item / page_size
	} else {
		page_num = Math.ceil(total_item / page_size)
	}

	var current_page = parseInt($.cookie("product_current_page"));
	if (current_page > page_num){
		current_page = page_num
	}
	$.jqPaginator(paginator_node, {
		totalPages: page_num,
		visiblePages: 10,
		currentPage:  current_page,
		onPageChange: function (num, type) {
			$.cookie("product_current_page", num);
			content_node_obj.html("");
			var is_out = num * page_size;
			if (is_out > total_item) {
				is_out = total_item
			}

			for (var i = page_size * (num - 1); i < is_out; i++) {
				var row = $("<tr product_item_no=''><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td>" +
					"<td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td>" +
					'<td class="text-c">' +
					'<a class="product_item_edit btn size-MINI btn-secondary-outline radius">&nbsp;' +
					'<i class="Hui-iconfont Hui-iconfont-edit"></i>&nbsp;</a> ' +
					' <a class="move_btn btn size-MINI btn-danger-outline radius" href="">&nbsp;' +
					'<i class="Hui-iconfont Hui-iconfont-fabu">' +
					'</i>&nbsp;</a> ' +
					'<a class="product_item_delete btn size-MINI btn-danger-outline radius" onclick=delete_row(this)>&nbsp;' +
					'<i class="Hui-iconfont Hui-iconfont-close">' +
					'</i>&nbsp;</a>' +
					'</td></tr>');

				//为每一行设置id属性，并赋值，便于删除和编辑
				row.attr("product_item_no", product[i].Id);

				var tds = row.find("td");
				tds.eq(0).html('<a href="/product_track/' + product[i].Id + '">' + product[i].Title + '</a>').addClass();
				var sale = $('<a href="/store_output_action/' + product[i].Id + '"target="_blank"><i class="Hui-iconfont Hui-iconfont-daochu"></i></a> ').addClass("c-danger");
				tds.eq(0).prepend(sale);

				tds.eq(16).find(".move_btn").attr("href", "/move_request/" + product[i].Id);

				tds.eq(1).text(product[i].BrandName).addClass("text-c");
				tds.eq(2).text(product[i].SupplierName).addClass("text-c");
				tds.eq(3).text(product[i].ArtNum).addClass("text-c");
				tds.eq(4).text(product[i].LotNum).addClass("text-c");

				tds.eq(5).text(product[i].Pool + "-" + product[i].StoreName).addClass("text-c");
				tds.eq(6).text(product[i].ThreeStage).addClass("text-c");
				tds.eq(7).text(product[i].Spec).addClass("text-c");
				tds.eq(8).text(product[i].Unit).addClass("text-c");
				tds.eq(9).text(product[i].Stock).addClass("text-c");
				tds.eq(10).text(product[i].InPrice).addClass("text-c");
				tds.eq(11).text(product[i].HasPay ? "是" : "否").addClass("text-c");
				tds.eq(12).text(product[i].HasInvoice ? "是" : "否").addClass("text-c");
				tds.eq(13).text((product[i].GetInvoice).substr(0, 10)).addClass("text-c");
				tds.eq(14).text((product[i].UserName)).addClass("text-c");
				tds.eq(15).text((product[i].InTime).substr(0, 19).replace("T", " ")).addClass("text-c");

				//节点追加
				product_rows.append(row);

				if (hidden_index.length > 0) {
					$.each(hidden_index, function (index, value) {
						row.find("td").eq(value).hide()
					})
				}

				//定义每一页的商品删除和编辑按钮，并在分页函数中进行赋值

				var product_item_edit = $(".product_item_edit");

				//编辑单行记录
				$.each(product_item_edit, function (index) {
					product_item_edit.eq(index).click(function () {
						$("input[name=product_id]").val($(this).parent().parent().attr("product_item_no"));

						$("#product_edit_modal").modal("show");

						var row = product_rows.find("tr").eq(index);
						var tds = row.find("td");

						$("#title_edit").val(tds.eq(0).text());
						$("#brand_edit").val(tds.eq(1).text());
						$("#supplier_edit").val(tds.eq(2).text());
						$("#art_num_edit").val(tds.eq(3).text());
						$("#lot_num_edit").val(tds.eq(4).text());
						$("#store_edit").val(tds.eq(5).text());
						$("#three_stage_edit").val(tds.eq(6).text());
						$("#spec_edit").val(tds.eq(7).text());

						var radios = $(".radios_edit").find("input");
						$.each(radios, function (index) {
							if (radios.eq(index).val() === tds.eq(8).text()) {
								radios.eq(index).attr("checked", true)
							}
						});

						$("#stock_edit").val(tds.eq(9).text());
						$("#in_price_edit").val(tds.eq(10).text());

						var has_pay_options = $("select[name=has_pay_edit]").find("option");
						var has_invoice_options = $("select[name=has_invioce_edit]").find("option");

						$.each(has_pay_options, function (index) {
							if (has_pay_options.eq(index).text() === tds.eq(11).text()) {
								has_pay_options.eq(index).attr("selected", true)
							}
						});

						$.each(has_invoice_options, function (index) {
							if (has_invoice_options.eq(index).text() === tds.eq(12).text()) {
								has_invoice_options.eq(index).attr("selected", true)
							}
						});

						$("#get_invioce_edit").val(tds.eq(13).text());
					})
				});
			}
		}
	});
}

//consumer_add.html
//-----------------------------------------------------------------------------------
if (query_url === "/consumer_add") {
	$("#prov").ProvinceCity()
}


//product_template_list.html
//-----------------------------------------------------------------------------------
if (query_url === "/product_template_list") {
	var template = $.parseJSON(template);
	$.cookie("template_offset", 0);
	$.cookie("template_current_page", 1);

	ProductTemplatePaginator(template);

	//加载更多
	$(".load-more-template").click(function () {
		$.cookie("template_offset", parseInt($.cookie("template_offset")) + 1);
		$.ajax({
			url : "/template_load_more",
			type : "post",
			data : {
				"offset" : $.cookie("template_offset"),
				"_xsrf" : $("input[name=_xsrf]").val()
			},
			success : function (response) {
				template = template.concat($.parseJSON(response));
				ProductTemplatePaginator(template)
			}
		})
	});

	//用户选择每页显示的条目数，也就是page_size
	var page_size_btn = $(".page_size");
	$.each(page_size_btn, function (index) {
		page_size_btn.eq(index).click(function () {

			//通过hui-ui.js的cookie()方法直接在浏览器设置cookie减少http请求（替代以上ajax请求）
			$.cookie('template_page_size', $(this).attr("data"), {expires: 366});

			//指示为第一页
			var num = 1;

			var page_size_temp = $.cookie("template_page_size");
			if (page_size_temp !== null) {
				page_size = page_size_temp
			}

			ProductTemplatePaginator(template);
		})
	});
}

function ProductTemplatePaginator(template) {
	//计算page_num
	var page_num;
	var total_item = template.length;
	console.log("total_item", total_item);

	var page_size = $.cookie("template_page_size");
	if (page_size === undefined) {
		page_size = 10
	}else{
		page_size = parseInt(page_size)
	}
	console.log("page_size", page_size);
	if (total_item % page_size === 0) {
		page_num = total_item / page_size
	} else {
		page_num = Math.ceil(total_item / page_size)
	}

	var current_page = parseInt($.cookie("template_current_page"));
	if (current_page > page_num){
		current_page = page_num
	}
	$.jqPaginator("#template_pagination", {
		totalPages: page_num,
		visiblePages: 10,
		currentPage: 1,
		onPageChange: function (num, type) {
			$.cookie("template_current_page", num);
			var template_node = $("#template");
			template_node.html("");
			var is_out = num * page_size;
			if (is_out > total_item) {
				is_out = total_item
			}

			for (var i = page_size * (num - 1); i < is_out; i++) {
				var row = $('<tr class="text-c tds-list"><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td class="hide"></td><td></td></tr>');
				var tds = row.find("td");
				tds.eq(0).text(template[i].Title);
				tds.eq(1).text(template[i].BrandName);
				tds.eq(2).text(template[i].ArtNum);
				tds.eq(3).text(template[i].ThreeStage);
				tds.eq(4).text(template[i].Suppliers);
				tds.eq(5).text(template[i].DealerName);
				tds.eq(6).text(template[i].Spec);
				tds.eq(7).text(template[i].Unit);
				tds.eq(8).text(template[i].InPrice);
				tds.eq(9).text(template[i].Id);
				tds.eq(10).html('<a onclick="ProductTemplateEdit(this)" class="btn size-MINI btn-success-outline radius">&nbsp;<i class="Hui-iconfont Hui-iconfont-edit"></i>&nbsp;</a>' +
				'<a onclick="DeleteTemplateRow(this,'+template[i].Id+')" class="btn size-MINI btn-danger-outline radius">&nbsp;<i class="Hui-iconfont Hui-iconfont-close"></i>&nbsp;</a>');
				template_node.append(row)
			}
		}
	});
}
//admin_member_edit.html
//-----------------------------------------------------------------------------------
//人员检索
$(".member_search").click(function () {
	var tds = $("<tr><td class='text-c'></td><td class='text-c'></td><td class='text-c'></td><td class='text-c'></td><td class='text-c'></td><td class='text-c'></td><td class='text-c'></td><td class='text-c'></td></tr>").find("td");
	$.ajax({
		type: "post",
		url: "/admin_member_edit",
		data: {
			"_xsrf": $("input[name=_xsrf]").val(),
			"search_entry": $("input[name=search_entry]").val()
		},
		success: function (response, status, xhr) {
			tds.eq(0).text(response.Name);
			tds.eq(1).text(response.Tel);
			tds.eq(2).text(response.Position);
			tds.eq(3).text(response.PoolName)
			var control_user = $(".control_user");
			if (response.IsActive) {
				tds.eq(4).text("正常");
				control_user.removeClass("btn-success").addClass("btn-danger").text("禁用账户")
			} else {
				tds.eq(4).text("未激活");
				tds.eq(4).addClass("c-danger");
				control_user.addClass("btn-success").removeClass("btn-danger").text("激活账户")
			}

			var login_raw = response.LastLogin;
			var created = response.Created;

			tds.eq(5).text((login_raw.substring(0, 19)).replace("T", " "));

			if (response.Ip === "") {
				tds.eq(6).text("未登陆过").addClass("c-danger");
			} else {
				tds.eq(6).text(response.Ip);
			}

			tds.eq(7).text((created.substring(0, 19)).replace("T", " "));
			$("tbody").html(tds);

			var input_hidden = $("<input type='hidden' name='userId'>");
			input_hidden.val(response.Id);

			$("tbody").append(input_hidden)
		}
	})
});


//禁用和激活人员账号
$(".control_user").click(function () {
	var conf;
	var is_active = $(this).hasClass("btn-success");
	var tds = $("td");
	if (is_active) {
		conf = confirm("您确定要激活此账号吗？");
	} else {
		conf = confirm("您确定要禁用此账号吗？");
	}
	if (conf) {
		$.ajax({
			type: "post",
			url: "/disable_active_user",
			data: {
				"_xsrf": $("input[name=_xsrf]").val(),
				"action": is_active ? "active" : "disable",
				"uid": $("input[name=userId]").val()
			},
			success: function (response, status, xhr) {
				if (response.Code === "success") {
					if (is_active) {
						$(".control_user").addClass("btn-danger").removeClass("btn-success").text("禁用账户");
						tds.eq(4).text("正常").removeClass("c-danger")
					} else {
						$(".control_user").addClass("btn-success").removeClass("btn-danger").text("激活账户");
						tds.eq(4).text("未激活").addClass("c-danger")
					}
				}
			}
		})
	}
});

//管理员编辑用户信息，弹窗，并为各个input初始化赋值
$(".edit_user").click(function () {
	$("#member_edit_modal").modal("show");
	var tds = $("td");
	$("input[name=uid]").val($("input[name=userId]").val());
	$("#name").val(tds.eq(0).text());
	$("#tel").val(tds.eq(1).text());
	var options = $("#position").find("option");
	$.each(options, function (index) {
		if (options.eq(index).text() === tds.eq(2).text()) {
			options.eq(index).attr("selected", true);
		}
	})
});


//store_output_action.html
//-----------------------------------------------------------------------------------
//开具发票日期
var disable = false, picker = new Pikaday({
	field: document.getElementById('sendinvioce'),
	firstDay: 1,
	minDate: new Date(2000, 0, 1),
	maxDate: new Date(),
	yearRange: [2000, 2030],

	showDaysInNextAndPreviousMonths: true,
	enableSelectionDaysInNextAndPreviousMonths: true
});

//递交发票日期
var disable = false, picker = new Pikaday({
	field: document.getElementById('getInvoice'),
	firstDay: 1,
	minDate: new Date(2000, 0, 1),
	maxDate: new Date(),
	yearRange: [2000, 2030],

	showDaysInNextAndPreviousMonths: true,
	enableSelectionDaysInNextAndPreviousMonths: true
});

//汇款日期
var disable = false, picker = new Pikaday({
	field: document.getElementById('getdate'),
	firstDay: 1,
	minDate: new Date(2000, 0, 1),
	maxDate: new Date(),
	yearRange: [2000, 2030],

	showDaysInNextAndPreviousMonths: true,
	enableSelectionDaysInNextAndPreviousMonths: true
});

//发货日期
var disable = false, picker = new Pikaday({
	field: document.getElementById('send'),
	firstDay: 1,
	minDate: new Date(2000, 0, 1),
	maxDate: new Date(),
	yearRange: [2000, 2030],

	showDaysInNextAndPreviousMonths: true,
	enableSelectionDaysInNextAndPreviousMonths: true
});

//标记单条销售信息，弹窗，赋值
var sale_item_edit = $(".sale_item_edit");
$.each(sale_item_edit, function (index) {
	sale_item_edit.eq(index).click(function () {
		$("#sale_edit_modal").modal("show");
		var tds = $(this).parent().parent().find("td");
		$("#title").val(tds.eq(0).text());
		$("#artnum").val(tds.eq(1).text());
		$("#salesman").val(tds.eq(2).text());
		$("#consumer").val(tds.eq(3).text());
		$("#inprice").val(tds.eq(4).text());
		$("#outprice").val(tds.eq(5).text());
		$("#num").val(tds.eq(6).text());
		$("#send").val(tds.eq(7).text());

		var hasinvoice = tds.eq(8).text();
		var options = $("select[name=hasinvoice]").find("option");
		$.each(options, function (index) {
			if (options.eq(index).text() === hasinvoice) {
				$(this).attr("selected", true);
			} else {
				$(this).attr("selected", false);
			}
		});

		$("#invioce_num").val(tds.eq(9).text());
		$("#sendinvioce").val(tds.eq(10).text());
		$("#getInvoice").val(tds.eq(11).text());

		var get_money = tds.eq(12).text();
		var option = $("select[name=get_money]").find("option");
		$.each(option, function (index) {
			if (option.eq(index).text() === get_money) {
				option.eq(index).attr("selected", true);
			} else {
				option.eq(index).attr("selected", false);
			}
		});

		$("#getdate").val(tds.eq(13).text());

		$("input[name=sale_id]").val($(".sale_id").eq(index).val());
	})
});

//common functions
//删除单条商品
function delete_row(obj) {
	var conf = confirm("您确定要删除此商品吗？");
	if (conf) {
		$.ajax({
			type: "post",
			url: "/product_item_delete",
			data: {
				"_xsrf": $("meta[name=_xsrf]").attr("content"),
				"product_id": $(obj).parent().parent().attr("product_item_no")
			},
			success: function (response) {
				if (response.Code === "success") {
					$(obj).parent().parent().hide();
					$.each(product, function (index) {
						if (product[index].Id === $(obj).parent().parent().attr("product_item_no")) {
							product.splice(index, 1);
						}
					})
				} else if (response.Code === "error") {
					alert(response.Message)
				} else {
					alert("未知错误，请报告管理员~")
				}
			},
			error: function (response, status, xhr) {
				if (xhr === "Unauthorized") {
					$("#product_edit_modal").modal("show").find(".modal-body p").text("您没有删除商品的权限，如有需要请联系管理员~");
				}
			}
		});
	}
}

//显示消息回复框
function showReplyForm() {
	$(".reply-form").removeClass("hidden");
}

//显示客户信息编辑弹窗
function ConsumerEdit(obj) {
	$("#consumer_edit_modal").modal("show");
	var tds = $(obj).parent().parent().find("td");
	$("input[name=consumer_id]").val($(obj).attr("cid"));
	$("input[name=name]").val(tds.eq(0).text());
	$("input[name=tel]").val(tds.eq(1).text());
	$("input[name=department]").val(tds.eq(2).text());
	$("input[name=province]").val(tds.eq(3).text());
	$("input[name=city]").val(tds.eq(4).text());
	$("input[name=region]").val(tds.eq(5).text());
	$("input[name=introduction]").val(tds.eq(6).text());
}

//返回上一页
function goBack() {
	window.location = document.referrer
}

//move_list.html
//-----------------------------------------------------------------------------------
//同意移库
function agreeMove(obj) {
	var conf = confirm("请在协商完成之后再同意此次移库操作！");
	if (conf) {
		$.ajax({
			type: "post",
			url: "/move_accept",
			data: {
				"_xsrf": $("meta[name=_xsrf]").attr("content"),
				"mid": $(obj).parent().parent().find("input[name=mid]").val()
			},
			success: function (response, status, xhr) {
				if (response.Code === "success") {
					var tds = $(obj).parent().parent().find("td");
					tds.eq(6).removeClass("c-danger").addClass("c-success").text("达成");

					//获取当前时间
					var date = new Date();
					var seperator1 = "-";
					var seperator2 = ":";
					var month = date.getMonth() + 1;
					var strDate = date.getDate();
					if (month >= 1 && month <= 9) {
						month = "0" + month;
					}
					if (strDate >= 0 && strDate <= 9) {
						strDate = "0" + strDate;
					}
					var currentdate = date.getFullYear() + seperator1 + month + seperator1 + strDate
						+ " " + date.getHours() + seperator2 + date.getMinutes()
						+ seperator2 + date.getSeconds();

					tds.eq(8).text(currentdate);

					$(obj).addClass("disabled").next().removeClass("disabled").next().removeClass("disabled")
				} else {
					alert("操作失败")
				}
			}
		})
	}
}

//拒绝移库
function disagreeMove(obj) {
	var conf = confirm("请在协商完成之后再拒绝此次移库操作！");
	if (conf) {
		$.ajax({
			type: "post",
			url: "/move_deny",
			data: {
				"_xsrf": $("meta[name=_xsrf]").attr("content"),
				"mid": $(obj).parent().parent().find("input[name=mid]").val()
			},
			success: function (response, status, xhr) {
				if (response.Code === "success") {
					var tds = $(obj).parent().parent().find("td");
					tds.eq(6).removeClass("c-success").addClass("c-danger").text("拒绝");

					//获取当前时间
					var date = new Date();
					var seperator1 = "-";
					var seperator2 = ":";
					var month = date.getMonth() + 1;
					var strDate = date.getDate();
					if (month >= 1 && month <= 9) {
						month = "0" + month;
					}
					if (strDate >= 0 && strDate <= 9) {
						strDate = "0" + strDate;
					}
					var currentdate = date.getFullYear() + seperator1 + month + seperator1 + strDate
						+ " " + date.getHours() + seperator2 + date.getMinutes()
						+ seperator2 + date.getSeconds();

					tds.eq(8).text(currentdate);

					$(obj).addClass("disabled").prev().removeClass("disabled");
					$(obj).next().addClass("disabled")
				} else {
					alert("操作失败")
				}
			}
		})
	}
}

//完成移库
function finishMove(obj) {
	var conf = confirm("请确定相应人已收到相应数量的货物！");
	if (conf) {
		$.ajax({
			type: "post",
			url: "/move_finish",
			data: {
				"_xsrf": $("meta[name=_xsrf]").attr("content"),
				"mid": $(obj).parent().parent().find("input[name=mid]").val()
			},
			success: function (response, status, xhr) {
				if (response.Code === "success") {
					var tds = $(obj).parent().parent().find("td");
					tds.eq(6).removeClass("c-danger").addClass("c-success").text("完成");

					//获取当前时间
					var date = new Date();
					var seperator1 = "-";
					var seperator2 = ":";
					var month = date.getMonth() + 1;
					var strDate = date.getDate();
					if (month >= 1 && month <= 9) {
						month = "0" + month;
					}
					if (strDate >= 0 && strDate <= 9) {
						strDate = "0" + strDate;
					}
					var currentdate = date.getFullYear() + seperator1 + month + seperator1 + strDate
						+ " " + date.getHours() + seperator2 + date.getMinutes()
						+ seperator2 + date.getSeconds();

					tds.eq(8).text(currentdate);

					$(obj).addClass("disabled").find("i").removeClass("Hui-iconfont-weigouxuan2").addClass("Hui-iconfont-xuanzhong1");
					$(obj).prev().addClass("disabled").prev().addClass("disabled")
				} else {
					alert("操作失败")
				}
			}
		})
	}
}

//category_edit.html
//-----------------------------------------------------------------------------------
function StageSearch() {
	$.ajax({
		url: "/category_search_ajax",
		type: "post",
		data: {
			"stage": $("select[name=stage]").val(),
			"item": $("input[name=search]").val(),
			"_xsrf": $("input[name=_xsrf]").val()
		},
		success: function (response, status, xhr) {
			$("input[name=category_id]").val(response.Id);
			$("input[name=primary]").val(response.Primary === "" ? "-" : response.Primary);
			$("input[name=two_stage]").val(response.TwoStage === "" ? "-" : response.TwoStage);
			$("input[name=three_stage]").val(response.ThreeStage === "" ? "-" : response.ThreeStage);
			$("select[name=is_hidden]").val(response.Is_hidden ? 0 : 1)
		},
		error: function () {
			alert("请求出错~")
		}
	})
}

//default_permission.html
//-----------------------------------------------------------------------------------
var permission_tds = $("#permission").find("td");
$.each(permission_tds, function (index) {
	if (permission_tds.eq(index).html() === "true") {
		permission_tds.eq(index).html("<i class='Hui-iconfont Hui-iconfont-xuanze'></i>")
	}
	if (permission_tds.eq(index).html() === "false") {
		permission_tds.eq(index).html("<i class='Hui-iconfont Hui-iconfont-close c-danger'></i>")
	}
});

//permission_member_edit.html
//-----------------------------------------------------------------------------------
var permission_member_tds = $("#permission_member").find("td");
$.each(permission_member_tds, function (index) {
	if (permission_member_tds.eq(index).html() === "true") {
		permission_member_tds.eq(index).html("<i class='Hui-iconfont Hui-iconfont-xuanze'></i>")
	}
	if (permission_member_tds.eq(index).html() === "false") {
		permission_member_tds.eq(index).html("<i class='Hui-iconfont Hui-iconfont-close c-danger'></i>")
	}
});

//product_template_add.html
//-----------------------------------------------------------------------------------
//增加sku
function AddTemplateSku() {
	$("#add_spec").parent("div").append('<input type="text" class="input-text mt-10 spec" value="" placeholder="规格" name="spec" style="width: 40%"> <input type="text" class="input-text mt-10 in_price" value="" placeholder="价格" id="in_price" name="in_price" style="width: 20%"> <a class="btn btn-danger-outline radius delete_sku mt-10"><i class="Hui-iconfont Hui-iconfont-close"></i></a>')

	delete_sku = $(".delete_sku");
	spec = $(".spec");
	stock = $(".stock");
	in_price = $(".in_price");

	$.each(delete_sku, function (index) {
		delete_sku.eq(index).click(function () {
			index += 1;
			spec.eq(index).remove();
			stock.eq(index).remove();
			in_price.eq(index).remove();
			$(this).remove();
		})
	});
}

function AppendSupplier() {
	var supplier_input = $("#supplier");
	var supplier_list = $("#supplier-list");
	if (supplier_input.val() !== "") {
		if (supplier_list.val() === "") {
			supplier_list.val(supplier_input.val());
		} else {
			supplier_list.val(supplier_list.val() + "," + supplier_input.val());
		}

		supplier_input.val("");
	}
}

//product_template_list.html
//-----------------------------------------------------------------------------------
//编辑
function ProductTemplateEdit(obj) {
	$("#product_template_edit_modal").modal("show");
	var tds = $(obj).parent().parent().find("td");
	$("#title").val(tds.eq(0).text());
	$("#brand").val(tds.eq(1).text());
	$("#art_num").val(tds.eq(2).text());
	$("#three_stage").val(tds.eq(3).text());
	$("#spec").val(tds.eq(6).text());
	$("#in_price").val(tds.eq(8).text());

	$("input[name=supplier_list]").val(tds.eq(4).text());
	$("input[name=template_id]").val(tds.eq(9).text());

	var radio = $(".radio").find("input");
	$.each(radio, function (index) {
		if (radio.eq(index).val() === tds.eq(7).text()) {
			radio.eq(index).attr("checked", true);
			radio.eq(index).next().click();
		} else {
			radio.eq(index).attr("checked", false)
		}
	});

	var supplier = tds.eq(4).text();
	var supplier_array = supplier.split(",");
	var len = supplier_array.length;

	//判断是否为空
	if (supplier === "") {
		len = 0;
	}
	var supplier_list = $(".supplier_list");
	supplier_list.find(".Huialert").remove();
	for (var i = 0; i < len; i++) {
		supplier_list.append('<div class="Huialert Huialert-success" style="width: 50%; margin: 5px 0 0 0;padding: 5px;"><i class="Hui-iconfont" onclick="DeleteSupplierRow(this)">&#xe6a6;</i><div>' + supplier_array[i] + '</div></div>');
	}
}

function DeleteSupplierRow(obj) {
	var supplier_list = $("input[name=supplier_list]");
	var supplier_array = supplier_list.val().split(",");
	supplier_array.splice($.inArray($(obj).next().text(), supplier_array), 1);
	supplier_list.val(supplier_array.join(","));
	$(obj).parent().remove()
}

function AddSupplierItem() {
	var supplier_list = $("input[name=supplier_list]");
	var supplier_value = $("#supplier").val();

	if (supplier_list.val() === "") {
		supplier_list.val(supplier_value);
	} else {
		supplier_list.val(supplier_list.val() + "," + supplier_value);
	}

	$(".supplier_list").append('<div class="Huialert Huialert-success" style="width: 50%; margin: 5px 0 0 0;padding: 5px;"><i class="Hui-iconfont" onclick="DeleteSupplierRow(this)">&#xe6a6;</i><div>' + supplier_value + '</div></div>');
}

function DeleteTemplateRow(obj, id) {
	$.ajax({
		type: "post",
		url: "/product_template_delete",
		data: {
			"_xsrf": $("input[name=_xsrf]").val(),
			"pid": $.trim(id)
		},
		success: function (response, status, xhr) {
			console.log(response);
			$(obj).parent().parent().remove()
		}
	})
}
