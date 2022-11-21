package service

var AuthMenuService = new(authMenuService)

type authMenuService struct{}

type Menu struct {
	Label    string  `json:"label"`
	Api      string  `json:"api"`
	Path     string  `json:"path"`
	Icon     string  `json:"icon"`
	Redirect string  `json:"redirect"`
	Children []*Menu `json:"children"`
}

var MenuTree []*Menu

//func init() {
//	data, err := ioutil.ReadFile("./data/auth_menu.json")
//	if err != nil {
//		g.Log().Error("菜单配置文件异常：", err)
//		panic(err)
//	}
//
//	err = json.Unmarshal(data, &MenuTree)
//	if err != nil {
//		g.Log().Error("菜单配置文件异常：", err)
//		panic(err)
//	}
//}

//获取全部菜单配置信息
func (s *authMenuService) Tree() []*Menu {
	return MenuTree
}
