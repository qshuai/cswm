package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id         int
	Username   string      `orm:"size(20);default();"`         //备注：用户名
	Password   string      `orm:"size(50)"`                    //备注：密码
	Name       string      `orm:"size(20);unique"`             //备注：姓名
	Tel        string      `orm:"size(45)"`                    //备注：电话
	Position   string      `orm:"size(10)"`                    //职位名称
	LastLogin  time.Time   `orm:"type(datetime);null"`         //备注：最后一次登录时间
	Ip         string      `orm:"size(15);null"`               //备注：最后一次登录IP
	IsFirst    bool        `orm:"defaul(true)"`                //是否为第一次登陆，第一次登陆可以使用手机号码登陆，在次登陆则不能
	IsActive   bool        `orm:"default(true)"`               //用于用户删除或禁用等操作，不用删除用户信息
	Stage      string      `orm:"size(5);default(on)"`         //状态：on在职，off离职
	PoolName   string      `orm:"size(10);null"`               //所管理库房的名称
	Created    time.Time   `orm:"auto_now_add;type(datetime)"` //备注：用户创建时间
	Updated    time.Time   `orm:"auto_now;type(datetime)"`     //备注：用户更新时间
	Product    []*Product  `orm:"reverse(many)"`
	Message    []*Message  `orm:"reverse(many)"`
	Sale       []*Sale     `orm:"reverse(many)"`
	Permission *Permission `orm:"reverse(one)"`
}

type Brand struct {
	Id              int
	Name            string             `orm:"size(20);unique"`
	ProductTemplate []*ProductTemplate `orm:"reverse(many)"`
	Created         time.Time          `orm:"auto_now_add;type(datetime)"`
}

type Category struct {
	Id              int
	Primary         string             `orm:"size(10)"`            //一级分类；字典：01-试剂， 02-耗材， 03-仪器
	TwoStage        string             `orm:"size(20);null"`       //二级分类
	ThreeStage      string             `orm:"size(50);null;index"` //三级分类
	Is_hidden       bool               //是否隐藏
	ProductTemplate []*ProductTemplate `orm:"reverse(many)"`
}

type Store struct {
	Id      int
	Pool    string     `orm:"size(10)"` //总库房名称, 不能含有`-`
	Name    string     `orm:"size(20)"` //分库房名称, 不能含有`-`
	Product []*Product `orm:"reverse(many)"`
	Sale    []*Sale    `orm:"reverse(many)"`
}

type Supplier struct {
	Id      int
	Name    string `orm:"size(100);unique"`
	Site    string
	Admin   string
	Tel     string
	Created time.Time  `orm:"type(datetime);auto_now_add"`
	Product []*Product `orm:"reverse(many)"`
}

// 经销商
type Dealer struct {
	Id      int
	Name    string     `orm:"size(100);unique"`
	Created time.Time  `orm:"type(datetime);auto_now_add"`
	Product []*Product `orm:"reverse(many)"`
}

type Product struct {
	Id         int
	User       *User     `orm:"rel(fk);on_delete(do_nothing)"` //备注：用户
	Title      string    `orm:"size(100)"`                     //备注：商品名称
	Brand      *Brand    `orm:"rel(fk);on_delete(do_nothing)"` //备注：商标
	ArtNum     string    `orm:"size(20);index"`                //备注：货号
	LotNum     string    `orm:"size(20);null"`                 //备注：批号
	CatNum     *Category `orm:"rel(fk);on_delete(do_nothing)"` //备注：分类号
	Spec       string    `orm:"size(100)"`                     //备注：规格
	Stock      uint32    //备注：库存数量
	Unit       string    `orm:"size(5)"`                       //单位
	Store      *Store    `orm:"rel(fk);on_delete(do_nothing)"` //备注：库房信息
	InTime     time.Time `orm:"type(datetime);auto_now_add"`   //备注：入库时间
	Supplier   *Supplier `orm:"rel(fk);on_delete(do_nothing)"` //备注：供应商
	Dealer     *Dealer   `orm:"rel(fk);null"`                  //经销商
	InPrice    float64   `orm:"digits(10);decimals(2)"`        //备注：进库价格
	HasPay     bool      `orm:"default(false)"`                //备注：是否已经支付货款； 字典：0-否定, 1-肯定
	HasInvoice bool      `orm:"default(false)"`                //备注：是否提供发票； 字典：0-否定, 1-肯定
	GetInvoice time.Time `orm:"type(date);null"`               //备注：发票接收日期
}

type ProductTemplate struct {
	Id        int
	Title     string    `orm:"size(100)"`                     //备注：商品名称
	Brand     *Brand    `orm:"rel(fk);on_delete(do_nothing)"` //备注：商标
	ArtNum    string    `orm:"size(20)"`                      //备注：货号
	CatNum    *Category `orm:"rel(fk);on_delete(do_nothing)"` //备注：分类号
	Spec      string    `orm:"size(100)"`                     //备注：规格
	Unit      string    `orm:"size(5)"`                       //单位
	Suppliers string    //备注：供应商列表(以逗号分隔)
	Dealer    *Dealer   `orm:"rel(fk);null"`                //经销商
	InPrice   float64   `orm:"digits(10);decimals(2);null"` //备注：进库价格
}

type Move struct {
	Id           int
	Origin       *Product  `orm:"rel(fk);on_delete(do_nothing)"`      //备注：来源产品Id
	Destination  *Product  `orm:"rel(fk);on_delete(do_nothing);null"` //备注：去处产品Id
	Num          uint32    //备注：移库数量
	From         *Store    `orm:"rel(fk);on_delete(do_nothing)"`      //备注：移出库房
	To           *Store    `orm:"rel(fk);on_delete(do_nothing)"`      //备注：移入库房
	Request      *User     `orm:"rel(fk);on_delete(do_nothing)"`      //备注：发起人
	Response     *User     `orm:"rel(fk);on_delete(do_nothing);null"` //备注：响应人
	Operate      string    `orm:"size(2)"`                            //备注：响应人是否同意 字典：0-未操作，1-同意，-1为拒绝，2-完成移库
	OperatedTime time.Time `orm:"type(datetime);null"`                //备注：响应人是否同意
	Created      time.Time `orm:"auto_now_add;type(datetime)"`        //备注：创建时间
	Finished     time.Time `orm:"type(datetime);null"`                //备注：完成时间
}

type Consumer struct {
	Id           int
	Name         string    `orm:"size(20)"`                    //备注：客户姓名
	Province     string    `orm:"size(10)"`                    //备注：省份
	City         string    `orm:"size(20)"`                    //备注：城市
	Region       string    `orm:"size(20);null"`               //备注：区
	Department   string    `orm:"size(40)"`                    //备注：单位
	Tel          string    `orm:"size(45)"`                    //备注：电话
	Introduction string    `orm:"type(text)"`                  //备注：简介
	Created      time.Time `orm:"type(datetime);auto_now_add"` //备注：创建时间
	Updated      time.Time `orm:"type(datetime);auto_now"`     //备注：更新时间
}

type Sale struct {
	Id           int
	Product      *Product  `orm:"rel(fk);on_delete(do_nothing)"` //备注：商品Id
	Store        *Store    `orm:"rel(fk);on_delete(do_nothing)"`
	No           string    `orm:"size(40)"`                      //备注：订单编号
	Send         time.Time `orm:"type(datetime)"`                //备注：发货时间
	Consumer     *Consumer `orm:"rel(fk);on_delete(do_nothing)"` //备注：客户Id
	Salesman     *User     `orm:"rel(fk);on_delete(do_nothing)"` //备注：销售Id
	Num          uint32    //备注：销售数量
	NumFake      uint32    //备注：伪造销售数量
	OutPrice     float64   `orm:"digits(10);decimals(2)"`      //备注：售出价格
	OutPriceFake float64   `orm:"digits(10);decimals(2)"`      //备注：伪造售出价格
	HasInvoice   bool      `orm:"default(false)"`              //备注：是否已开发票
	SendInvoice  time.Time `orm:"type(datetime);null"`         //备注：开具发票日期
	GetInvoice   time.Time `orm:"type(datetime);null"`         //备注：递交发票日期
	InvoiceNum   string    `orm:"size(10);null"`               //备注：发票编号
	GetMoney     bool      `orm:"default(false)"`              //备注：是否已经接收回款
	GetDate      time.Time `orm:"type(datetime);null"`         //备注：接受回款日期
	HasPrint     bool      `orm:"default(false)"`              //是否已经打印出库单
	Created      time.Time `orm:"type(datetime);auto_now_add"` //备注：订单创建日期
	Updated      time.Time `orm:"type(datetime);auto_now"`     //备注：订单更新日期
	Comment      string    `orm:"size(255);null"`              //备注：备注
}

type Message struct {
	Id      int
	From    *User     `orm:"rel(fk);on_delete(do_nothing)"` //发信人
	To      *User     `orm:"rel(fk);on_delete(do_nothing)"` //收信人
	Content string    `orm:"type(text)"`                    //消息内容
	IsRead  bool      `orm:"default(false)"`                //是否已经读取
	ReadAt  time.Time `orm:"type(datetime);auto_now"`       //读取时间
	Created time.Time `orm:"type(datetime);auto_now_add"`   //创建时间
}

type Permission struct {
	Id                int
	User              *User `orm:"rel(one);unique"`
	AddMember         bool  `orm:"default(false)"` //添加人员
	EditMember        bool  `orm:"default(false)"` //编辑人员信息
	ActiveMember      bool  `orm:"default(false)"` //激活或禁用账户
	AddConsumer       bool  `orm:"default(false)"` //添加客户信息
	EditConsumer      bool  `orm:"default(false)"` //编辑客户信息
	ViewConsumer      bool  `orm:"default(false)"` //查看客户信息
	AddBrand          bool  `orm:"default(false)"` //添加品牌
	AddDealer         bool  `orm:"default(false)"` //添加经销商
	ViewDealer        bool  `orm:"default(false)"` //查看经销商
	AddSupplier       bool  `orm:"default(false)"` //添加供应商
	ViewSupplier      bool  `orm:"default(false)"` //查看供应商
	AddProduct        bool  `orm:"default(false)"` //录入商品
	InputInPrice      bool  `orm:"default(false)"` //录入商品入库价格
	ViewProductStore  bool  `orm:"default(false)"` //查看商品库房
	ViewStock         bool  `orm:"default(false)"` //查看商品库存
	ViewInPrice       bool  `orm:"default(false)"` //查看入库价格
	EditProduct       bool  `orm:"default(false)"` //编辑商品信息
	DeleteProduct     bool  `orm:"default(false)"` //删除商品
	OutputProduct     bool  `orm:"default(false)"` //出库商品
	ViewSale          bool  `orm:"default(false)"` //查看销售记录
	ViewSaleConsumer  bool  `orm:"default(false)"` //查看销售客户
	ViewSaleInPrice   bool  `orm:"default(false)"` //查看销售入库价格
	EditSale          bool  `orm:"default(false)"` //编辑销售信息
	OperateCategory   bool  `orm:"default(false)"` //添加或编辑分类信息
	RequestMove       bool  `orm:"default(false)"` //请求移库
	ResponseMove      bool  `orm:"default(false)"` //响应移库
	ViewMove          bool  `orm:"default(false)"` //查看移库
	AddStore          bool  `orm:"default(false)"` //添加库房
	ViewStore         bool  `orm:"default(false)"` //查看库房
	OperateOtherStore bool  `orm:"default(false)"` //操作非管辖库房
}

type DefaultPermission struct {
	Id                int
	Position          string `orm:"size(20)"`       //人员等级
	AddMember         bool   `orm:"default(false)"` //添加人员
	EditMember        bool   `orm:"default(false)"` //编辑人员信息
	ActiveMember      bool   `orm:"default(false)"` //激活或禁用账户
	AddConsumer       bool   `orm:"default(false)"` //添加客户信息
	EditConsumer      bool   `orm:"default(false)"` //编辑客户信息
	ViewConsumer      bool   `orm:"default(false)"` //查看客户信息
	AddBrand          bool   `orm:"default(false)"` //添加品牌
	AddDealer         bool   `orm:"default(false)"` //添加经销商
	ViewDealer        bool   `orm:"default(false)"` //查看经销商
	AddSupplier       bool   `orm:"default(false)"` //添加供应商
	ViewSupplier      bool   `orm:"default(false)"` //查看供应商
	AddProduct        bool   `orm:"default(false)"` //录入商品
	InputInPrice      bool   `orm:"default(false)"` //录入商品入库价格
	ViewProductStore  bool   `orm:"default(false)"` //查看商品库房
	ViewStock         bool   `orm:"default(false)"` //查看商品库存
	ViewInPrice       bool   `orm:"default(false)"` //查看入库价格
	EditProduct       bool   `orm:"default(false)"` //编辑商品信息
	DeleteProduct     bool   `orm:"default(false)"` //删除商品
	OutputProduct     bool   `orm:"default(false)"` //出库商品
	ViewSale          bool   `orm:"default(false)"` //查看销售记录
	ViewSaleConsumer  bool   `orm:"default(false)"` //查看销售客户
	ViewSaleInPrice   bool   `orm:"default(false)"` //查看销售入库价格
	EditSale          bool   `orm:"default(false)"` //编辑销售信息
	OperateCategory   bool   `orm:"default(false)"` //添加或编辑分类信息
	RequestMove       bool   `orm:"default(false)"` //请求移库
	ResponseMove      bool   `orm:"default(false)"` //响应移库
	ViewMove          bool   `orm:"default(false)"` //查看移库
	AddStore          bool   `orm:"default(false)"` //添加库房
	ViewStore         bool   `orm:"default(false)"` //查看库房
	OperateOtherStore bool   `orm:"default(false)"` //操作非管辖库房
}

// 出库单编号
type OrderNum struct {
	Id         int
	Asap       string
	Consumer   string
	Department string
	Salesman   string
	Sum        string
	SaleList   string
	User       string
	State      bool      `orm:"default(true)"`
	HasPrint   bool      `orm:"default(false)"`
	IsFake     bool      `orm:"default(false)"`
	Created    time.Time `orm:"auto_now_add;type(datetime)"`
	Updated    time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.Debug = true

	//获取配置信息
	username := beego.AppConfig.String("mysql::username")
	password := beego.AppConfig.String("mysql::password")
	host := beego.AppConfig.String("mysql::host")
	port := beego.AppConfig.String("mysql::port")
	database := beego.AppConfig.String("mysql::database")
	err := orm.RegisterDataBase("default", "mysql", username+
		":"+password+"@tcp("+host+":"+port+")/"+database+"?charset=utf8&loc=Asia%2FShanghai")
	if err != nil {
		panic("注册数据库失败" + err.Error())
	}

	orm.RegisterModel(new(User), new(Brand), new(Category), new(Store), new(Supplier), new(Dealer),
		new(Product), new(Move), new(Consumer), new(Sale), new(Message), new(Permission),
		new(DefaultPermission), new(ProductTemplate), new(OrderNum))

	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		panic("同步数据库失败" + err.Error())
	}
}
