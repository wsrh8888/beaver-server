package conversation

// CityData 城市数据
type CityData struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// GetDefaultCities 获取默认城市列表
func GetDefaultCities() []CityData {
	return []CityData{
		{Code: "ALL", Name: "全国"},
		{Code: "010", Name: "北京"},
		{Code: "021", Name: "上海"},
		{Code: "020", Name: "广州"},
		{Code: "0755", Name: "深圳"},
		{Code: "0571", Name: "杭州"},
		{Code: "028", Name: "成都"},
		{Code: "027", Name: "武汉"},
		{Code: "029", Name: "西安"},
		{Code: "025", Name: "南京"},
		{Code: "023", Name: "重庆"},
		{Code: "022", Name: "天津"},
		{Code: "0512", Name: "苏州"},
		{Code: "0731", Name: "长沙"},
		{Code: "0532", Name: "青岛"},
		{Code: "0510", Name: "无锡"},
		{Code: "0574", Name: "宁波"},
		{Code: "0371", Name: "郑州"},
		{Code: "0757", Name: "佛山"},
		{Code: "0769", Name: "东莞"},
	}
}
