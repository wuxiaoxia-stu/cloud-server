package licence

// client_apply_licence_req 客户端申请证书请求
type ClientApplyLicenceReq struct {
	UkeySerialNumber string `json:"ukey_serial_number"` // ukey唯一编号
	ClientVersion    string `json:"client_version"`     // 客户端版本号
	CpuId            string `json:"cpu_id"`             // cpu_id
	CpuName          string `json:"cpu_name"`           // cpu_name
	DiskId           string `json:"disk_id"`            // disk_id
	BaseboardId      string `json:"baseboard_id"`       // baseboard_id
	Uuid             string `json:"uuid"`               // 主板uuid
	GpuName          string `json:"gpu_name"`           // gpu_name
	Role             int32  `json:"role"`               // 角色 （1-windows 客户端，2-web 端）
	SerialNumber     string `json:"serial_number"`      // 订单设备序列号
	Signature        string `json:"signature"`          // 数据签名，经过ukey签名
}

type ClientAuthStatus struct {
	Uuid string `p:"uuid" v:"required#参数错误,uuid必填"`
}

//配对提交数据
type PairReq struct {
	ClientAuthorNumber string `p:"client_author_number" v:"required#参数错误,客户端授权编号必填"`  // 客户端授权编号
	ClientSerialNumber string `p:"client_serial_number" v:"required#参数错误,客户端设备系列号必填"` // 客户端设备系列号
	ServerAuthorNumber string `p:"service_author_number" v:"required#参数错误,服务端授权编号必填"` // 服务端授权编号
	Signature          string `p:"signature" v:"required#参数错误,服务端签名信息必填"`             // 服务端数据签名
}

//配对提交数据
type BreakPairReq struct {
	ClientAuthorNumber string `p:"client_author_number" v:"required#参数错误,客户端授权编号必填"`  // 客户端授权编号
	ServerAuthorNumber string `p:"service_author_number" v:"required#参数错误,服务端授权编号必填"` // 服务端授权编号
	Signature          string `p:"signature" v:"required#参数错误,服务端签名信息必填"`             // 服务端数据签名
}

//激活查询激活状态
type ActivateStatusReq struct {
	AuthorNumber string `json:"author_number" p:"author_number" v:"required#参数错误：授权编号必填"`  // 授权编号
	SerialNumber string `json:"serial_number" p:"serial_number" v:"required#参数错误：设备系列号必填"` // 设备序列号
}

//激活提交数据验证
type ActivateReq struct {
	AuthorNumber  string `json:"author_number" p:"author_number" v:"required#参数错误：授权编号必填"`        // 授权编号
	SerialNumber  string `json:"serial_number" p:"serial_number" v:"required#参数错误：设备系列号必填"`       // 设备序列号
	MachineBrand  string `json:"machine_brand" p:"machine_brand" v:"length:1,100#超声机品牌不超过100个字符"` // 超声机品牌
	Floor         string `json:"floor" p:"floor" v:"length:1,20#楼层信息不超过20个字符"`                    // 楼层
	Room          string `json:"room" p:"room" v:"length:1,20#房号信息不超过20个字符"`                      // 房号
	MachineNumber string `json:"machine_number" p:"machine_number" v:"length:1,20#超声机简称不超过20个字符"` // 超声机简称
}
