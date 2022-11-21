package Admin

import (
	"aiyun_cloud_srv/app/model/order_device"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"strconv"
	"strings"
)

var Index = indexApi{}

type indexApi struct{}

type OrderRegionData struct {
	Count    uint64 `json:"count"`
	RegionId string `json:"region_id"`
	UseType  int    `json:"use_type"`
}

func (*indexApi) OrderRegionData(r *ghttp.Request) {
	use_type := r.GetQueryInt("use_type")

	region_tree, err := service.RegionService.Tree()
	if err != nil {
		response.ErrorSys(r, err)
	}

	sql := "SELECT COUNT(s) as count,region_id,use_type FROM(" +
		"SELECT O.order_number,D.use_type,H.region_id FROM lpm_order_device AS D " +
		"RIGHT JOIN lpm_order AS O ON O.order_number = D.order_number " +
		"RIGHT JOIN lpm_hospital AS H ON H.id = O.hospital_id " +
		"WHERE O.status in (1,5,9)"
	if use_type > 0 {
		sql += "AND  D.use_type = " + strconv.Itoa(use_type)
	}
	sql += ")as s GROUP BY region_id,use_type"

	res := []*OrderRegionData{}
	if err := g.DB().GetScan(&res, sql); err != nil {
		response.ErrorDb(r, err)
		return
	}

	data := []*UseTypeCount{}
	for k, v := range region_tree {
		data = append(data, &UseTypeCount{Name: strings.TrimRight(v.Label, "省"), UseType: use_type})
		for _, v2 := range res {
			region_ids := strings.Split(v2.RegionId, ",")
			if len(region_ids) > 0 && v.Value == region_ids[0] {
				data[k].Value += v2.Count
			}
		}
	}

	response.Success(r, data)
}

type UseTypeCount struct {
	UseType int    `json:"id"`
	Value   uint64 `json:"value"`
	Name    string `json:"name"`
}

func (*indexApi) DeviceData(r *ghttp.Request) {
	res := []*UseTypeCount{}
	err := order_device.M.Fields("use_type,count(*) as value").Where(g.Map{"status": 1}).Group("use_type").Scan(&res)
	if err != nil {
		response.ErrorDb(r, err)
	}

	data := []*UseTypeCount{
		&UseTypeCount{UseType: 4, Name: "售卖产品"},
		&UseTypeCount{UseType: 1, Name: "临床试验用"},
		&UseTypeCount{UseType: 2, Name: "厂家样机"},
		&UseTypeCount{UseType: 3, Name: "代理商样机"},
	}

	for _, v := range data {
		for _, v2 := range res {
			if v.UseType == v2.UseType {
				v.Value = v2.Value
			}
		}
	}

	response.Success(r, data)
}
