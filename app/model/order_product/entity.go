package order_product

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

// Entity is the golang structure for table order_product
type Entity struct {
	Id           int         `orm:"id,primary,table_comment:'设备管理'" json:"id"`
	OrderNumber  string      `orm:"order_number,size:15,not null,comment:'订单号'" json:"order_number"`
	SerialNumber string      `orm:"serial_number,size:18,not null,comment:'设备系列码'" json:"serial_number"`
	SetMeal      int         `orm:"set_meal,size:2,comment:'设备套餐'" json:"set_meal"`
	Type         int         `orm:"type,size:2,comment:'产品类型'" json:"device_type"`
	ProductId    int         `orm:"product_id,size:2,comment:'产品key'" json:"product_id"`
	MonthCount   int         `orm:"month_count,size:2,comment:'使用时长/月份'" json:"month_count" p:"month_count"`
	DueTime      int         `orm:"due_time,size:8,comment:'过期时间'" json:"due_time"`
	CreateAt     *gtime.Time `orm:"create_at" json:"create_at"`
	UpdateAt     *gtime.Time `orm:"update_at" json:"update_at"`
	Status       int         `orm:"status,size:2,comment:'状态'" json:"status"`
	ProductName  string      `json:"product_name"`
	Children     []*Entity   `json:"children"`
}

var (
	// Table is the table name of order_product.
	Table       = "order_product"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "order_product op"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

var MealOptionDefMap = map[int][]int{
	1: {},
	2: {},
	3: {10001, 10101, 10201, 10202, 10203, 10204, 10205, 10206, 10207, 10208},
	4: {10001, 10101, 10201, 10202, 10203, 10204, 10205, 10206, 10207, 10208, 10301, 10302, 10401},
	5: {10001, 10101, 10201, 10202, 10203, 10204, 10205, 10206, 10207, 10208, 10301, 10302, 10401, 10501, 10601},
	6: {10001, 10101, 10201, 10202, 10203, 10204, 10205, 10206, 10207, 10208, 10301, 10302, 10401, 10501, 10601, 10701, 10801},
	7: {10001, 10101, 10201, 10202, 10203, 10204, 10205, 10206, 10207, 10208, 10301, 10302, 10401, 10501, 10601, 10701, 10801, 10901},
	8: {10001, 10101, 10201, 10202, 10203, 10204, 10205, 10206, 10207, 10208, 10301, 10302, 10401, 10501, 10601, 10701, 10801, 10901, 11001, 11101},
	9: {},
}

var AllMealOptionMap = []int{10001, 10101, 10201, 10202, 10203, 10204, 10205, 10206, 10207, 10208, 10301, 10302, 10401, 10501, 10601, 10701, 10801, 10901, 11001, 11101}

var MealOptionDefMap2 = map[int][]int{
	1: {},
	2: {},
	3: {100, 101, 102},
	4: {100, 101, 102, 103, 104},
	5: {100, 101, 102, 103, 104, 105, 106},
	6: {100, 101, 102, 103, 104, 105, 106, 107, 108},
	7: {100, 101, 102, 103, 104, 105, 106, 107, 108, 109},
	8: {100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111},
	9: {},
}

// 设备套餐类型
var MealMap = map[int]string{
	1: "标准版",
	2: "年费版",
	3: "PUS",
	4: "PUS高配",
	5: "PUS S100",
	6: "PUS ES",
	7: "PUS S200",
	8: "PUS S300",
	9: "试用版",
}

var AllProductList = []*Entity{
	&Entity{ProductId: 100, ProductName: "用户登录模块", Children: []*Entity{
		&Entity{ProductId: 10001, ProductName: "软件基础功能"},
	}},
	&Entity{ProductId: 101, ProductName: "系统管理模块", Children: []*Entity{
		&Entity{ProductId: 10101, ProductName: "软件基础功能"},
	}},
	&Entity{ProductId: 102, ProductName: "影像处理模块", Children: []*Entity{
		&Entity{ProductId: 10201, ProductName: "中晚孕期胎儿系统结构辅助筛选"},
		&Entity{ProductId: 10202, ProductName: "胎儿常见颅脑畸形辅助诊断"},
		&Entity{ProductId: 10203, ProductName: "视频回放诊断"},
		&Entity{ProductId: 10204, ProductName: "单/多胎、病例续接"},
		&Entity{ProductId: 10205, ProductName: "早孕期NT辅助筛查"},
		&Entity{ProductId: 10206, ProductName: "早孕期重大畸形辅助诊断"},
		&Entity{ProductId: 10207, ProductName: "自动取图"},
		&Entity{ProductId: 10208, ProductName: "自动测量"},
	}},
	&Entity{ProductId: 103, ProductName: "图文报告模块", Children: []*Entity{
		&Entity{ProductId: 10301, ProductName: "测量值和图像传输工作站"},
		&Entity{ProductId: 10302, ProductName: "超声报告辅助自动生成"},
	}},
	&Entity{ProductId: 104, ProductName: "知识图谱模块", Children: []*Entity{
		&Entity{ProductId: 10401, ProductName: "胎儿多发畸形诊断思路引导、图文资料、检索、对比鉴别、常用图表等"},
	}},
	&Entity{ProductId: 105, ProductName: "患者管理模块", Children: []*Entity{
		&Entity{ProductId: 10501, ProductName: "增改删患者等功能"},
	}},
	&Entity{ProductId: 106, ProductName: "质控考核模块", Children: []*Entity{
		&Entity{ProductId: 10601, ProductName: "现场留图考核以及导入图片考核，提供静态留图定量评分及考核报告"},
	}},
	&Entity{ProductId: 107, ProductName: "病例管理模块", Children: []*Entity{
		&Entity{ProductId: 10701, ProductName: "3-30秒视频导出、标签功能、病例检索等功能"},
	}},
	&Entity{ProductId: 108, ProductName: "教学规培模块", Children: []*Entity{
		&Entity{ProductId: 10801, ProductName: "胎儿畸形教学考核、标准切面参考图、质控训练模式等功能"},
	}},
	&Entity{ProductId: 109, ProductName: "科室管理模块", Children: []*Entity{
		&Entity{ProductId: 10901, ProductName: "自动排班、病例讨论等"},
	}},
	&Entity{ProductId: 110, ProductName: "全科质控模块", Children: []*Entity{
		&Entity{ProductId: 11001, ProductName: "对全科病例进行质量控制、抽查、定量评分、报告生成、典型病例设置考题等功能，含云点播系统、大数据服务器"},
	}},
	&Entity{ProductId: 111, ProductName: "科研助理模块", Children: []*Entity{
		&Entity{ProductId: 11101, ProductName: "全面支持将搜索出的感兴趣病例数据导出为Excel表格和多段视频导出"},
	}},
}
