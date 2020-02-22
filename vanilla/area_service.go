package vanilla

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

var ZONE_NAMES = []string{"直辖市", "华北-东北", "华东地区", "华南-华中", "西北-西南", "其它"}

var PROVINCEID2ZONE = map[int]string{
	1: "直辖市",
	2: "直辖市",
	3: "华北-东北",
	4: "华北-东北",
	5: "华北-东北",
	6: "华北-东北",
	7: "华北-东北",
	8: "华北-东北",
	9: "直辖市",
	10: "华东地区",
	11: "华东地区",
	12: "华东地区",
	13: "华东地区",
	14: "华东地区",
	15: "华东地区",
	16: "华南-华中",
	17: "华南-华中",
	18: "华南-华中",
	19: "华南-华中",
	20: "华南-华中",
	21: "华南-华中",
	22: "直辖市",
	23: "西北-西南",
	24: "西北-西南",
	25: "西北-西南",
	26: "西北-西南",
	27: "西北-西南",
	28: "西北-西南",
	29: "西北-西南",
	30: "西北-西南",
	31: "西北-西南",
	32: "其它",
	33: "其它",
	34: "其它",
}

var AREA = make(map[string][]map[string]interface{})
var provinces = make([]*Province, 0)
var cities = make([]*City, 0)
var districts = make([]*District, 0)
var name2Province = make(map[string]*Province)
var id2Province = make(map[int]*Province)
var name2City = make(map[string]*City)
var id2City = make(map[int]*City)
var id2District = make(map[int]*District)
var name2District = make(map[string]*District)

type Province struct {
	Id int
	Name string
	Zone string
	Cities []*City
}

// IsDGC 是否是直辖市
func (this *Province) IsDGC() bool {
	return len(this.Cities) == 1
}

func NewProvince(data map[string]interface{}) *Province{
	province := new(Province)
	province.Id, _ = strconv.Atoi(data["id"].(string))
	province.Name = data["name"].(string)
	province.Cities = make([]*City, 0)
	return province
}

type City struct {
	Id int
	ProvinceId int
	Name string
	ZipCode string
	Districts []*District
}

func NewCity(data map[string]interface{}) *City{
	city := new(City)
	city.Id, _ = strconv.Atoi(data["id"].(string))
	city.ProvinceId, _ = strconv.Atoi(data["province_id"].(string))
	city.Name = data["name"].(string)
	city.ZipCode = data["zip_code"].(string)
	city.Districts = make([]*District, 0)
	return city
}

type District struct {
	Id int
	CityId int
	Name string
}

func NewDistrict(data map[string]interface{}) *District{
	district := new(District)
	district.Id, _ = strconv.Atoi(data["id"].(string))
	district.CityId, _ = strconv.Atoi(data["city_id"].(string))
	district.Name = data["name"].(string)
	return district
}

type Area struct {
	Province *Province
	City *City
	District *District
}

type AreaService struct {
	ServiceBase
}

func NewAreaService() *AreaService {
	service := new(AreaService)
	return service
}

/*
  Province相关api
*/

func (this *AreaService) GetProvinces() []*Province {
	return provinces
}

func (this *AreaService) GetProvinceByName(name string) *Province{
	return name2Province[name]
}

func (this *AreaService) GetProvincesByNames(names []string) []*Province {
	provinces := make([]*Province, 0)
	for _, name := range names {
		if province, ok := name2Province[name]; ok {
			provinces = append(provinces, province)
		}
	}
	
	return provinces
}


func (this *AreaService) GetProvinceById(id int) *Province {
	return id2Province[id]
}

func (this *AreaService) GetProvincesByIds(ids []int) []*Province {
	provinces := make([]*Province, 0)
	for _, id := range ids {
		if province, ok := id2Province[id]; ok {
			provinces = append(provinces, province)
		}
	}
	
	return provinces
}

/*
  City相关api
*/

func (this *AreaService) GetCityByName(name string) *City{
	return name2City[name]
}

func (this *AreaService) GetCitiesByNames(names []string) []*City {
	cities := make([]*City, 0)
	for _, name := range names {
		if city, ok := name2City[name]; ok {
			cities = append(cities, city)
		}
	}
	
	return cities
}

func (this *AreaService) GetCityById(id int) *City{
	return id2City[id]
}

func (this *AreaService) GetCitiesByIds(ids []int) []*City {
	cities := make([]*City, 0)
	for _, id := range ids {
		if city, ok := id2City[id]; ok {
			cities = append(cities, city)
		}
	}
	
	return cities
}

func (this *AreaService) GetCitiesForProvince(provinceId int) []*City {
	if province, ok := id2Province[provinceId]; ok {
		return province.Cities
	} else {
		return make([]*City, 0)
	}
}

/*
 District相关api
*/
func (this *AreaService) GetDistrictByName(cityId int, name string) *District{
	name = fmt.Sprintf("%d_%s", cityId, name)
	return name2District[name]
}

func (this *AreaService) GetDistrictsByNames(cityId int, names []string) []*District{
	districts := make([]*District, 0)
	for _, name := range names {
		_name := fmt.Sprintf("%d_%s", cityId, name)
		if district, ok := name2District[_name]; ok {
			districts = append(districts, district)
		}
	}
	
	return districts
}

func (this *AreaService) GetDistrictById(id int) *District{
	return id2District[id]
}

func (this *AreaService) GetDistrictsByIds(ids []int) []*District {
	districts := make([]*District, 0)
	for _, id := range ids {
		if district, ok := id2District[id]; ok {
			districts = append(districts, district)
		}
	}
	
	return districts
}

func (this *AreaService) GetDistrictsForCity(cityId int) []*District {
	if city, ok := id2City[cityId]; ok {
		return city.Districts
	} else {
		return make([]*District, 0)
	}
}

/*
 Area相关api
*/

// GetAreaByName 根据area name(北京市 北京市 东城区)获得area
func (this *AreaService) GetAreaByName(name string) *Area {
	var items []string
	if strings.Contains(name, "-"){
		items = strings.Split(name, "-")
	} else {
		items = strings.Split(name, " ")
	}

	province := this.GetProvinceByName(items[0])
	city := this.GetCityByName(items[1])
	district := this.GetDistrictByName(city.Id, items[2])
	
	return &Area{
		Province: province,
		City: city,
		District: district,
	}
}

// GetAreaByCode 根据area code(1_1_1)获取area
func (this *AreaService) GetAreaByCode(code string) *Area {
	items := strings.Split(code, "_")
	
	provinceId, _ := strconv.Atoi(items[0])
	province := this.GetProvinceById(provinceId)
	
	cityId, _ := strconv.Atoi(items[1])
	city := this.GetCityById(cityId)
	
	districtId, _ := strconv.Atoi(items[2])
	district := this.GetDistrictById(districtId)
	
	return &Area{
		Province: province,
		City: city,
		District: district,
	}
}


func init(){
	json.Unmarshal([]byte(areaJSON), &AREA)
	for name, arrs := range AREA {
		switch name {
		case "PROVINCES":
			for _, data := range arrs {
				province := NewProvince(data)
				provinces = append(provinces, province)
				name2Province[province.Name] = province
				id2Province[province.Id] = province
			}
		case "CITIES":
			for _, data := range arrs {
				city := NewCity(data)
				cities = append(cities, city)
			}
		case "DISTRICTS":
			for _, data := range arrs {
				district := NewDistrict(data)
				districts = append(districts, district)
			}
		}
	}

	for _, city := range cities {
		name2City[city.Name] = city
		id2City[city.Id] = city
		//向province.Cities中加入city
		if province, ok := id2Province[city.ProvinceId]; ok {
			province.Cities = append(province.Cities, city)
		}
	}

	for _, district := range districts {
		name := fmt.Sprintf("%d_%s", district.CityId, district.Name)
		name2District[name] = district
		id2District[district.Id] = district

		if city, ok := id2City[district.CityId]; ok {
			city.Districts = append(city.Districts, district)
		}
	}

	//for name, arrs := range AREA{
	//	switch name {
	//	case "PROVINCES":
	//		for _, data := range arrs {
	//			province := NewProvince(data)
	//			provinces = append(provinces, province)
	//			name2Province[province.Name] = province
	//			id2Province[province.Id] = province
	//		}
	//
	//	case "CITIES":
	//		for _, data := range arrs {
	//			city := NewCity(data)
	//			name2City[city.Name] = city
	//			id2City[city.Id] = city
	//			//cities = append(cities, city)
	//
	//			//向province.Cities中加入city
	//			if province, ok := id2Province[city.ProvinceId]; ok {
	//				province.Cities = append(province.Cities, city)
	//			}
	//		}
	//	case "DISTRICTS":
	//		for _, data := range arrs {
	//			district := NewDistrict(data)
	//			name := fmt.Sprintf("%d_%s", district.CityId, district.Name)
	//			name2District[name] = district
	//			id2District[district.Id] = district
	//
	//			if city, ok := id2City[district.CityId]; ok {
	//				city.Districts = append(city.Districts, district)
	//			}
	//		}
	//	}
	//}
}

const areaJSON = `{
    "PROVINCES": [
        {
            "id": "1",
            "name": "北京市"
        },
        {
            "id": "2",
            "name": "天津市"
        },
        {
            "id": "3",
            "name": "河北省"
        },
        {
            "id": "4",
            "name": "山西省"
        },
        {
            "id": "5",
            "name": "内蒙古自治区"
        },
        {
            "id": "6",
            "name": "辽宁省"
        },
        {
            "id": "7",
            "name": "吉林省"
        },
        {
            "id": "8",
            "name": "黑龙江省"
        },
        {
            "id": "9",
            "name": "上海市"
        },
        {
            "id": "10",
            "name": "江苏省"
        },
        {
            "id": "11",
            "name": "浙江省"
        },
        {
            "id": "12",
            "name": "安徽省"
        },
        {
            "id": "13",
            "name": "福建省"
        },
        {
            "id": "14",
            "name": "江西省"
        },
        {
            "id": "15",
            "name": "山东省"
        },
        {
            "id": "16",
            "name": "河南省"
        },
        {
            "id": "17",
            "name": "湖北省"
        },
        {
            "id": "18",
            "name": "湖南省"
        },
        {
            "id": "19",
            "name": "广东省"
        },
        {
            "id": "20",
            "name": "广西壮族自治区"
        },
        {
            "id": "21",
            "name": "海南省"
        },
        {
            "id": "22",
            "name": "重庆市"
        },
        {
            "id": "23",
            "name": "四川省"
        },
        {
            "id": "24",
            "name": "贵州省"
        },
        {
            "id": "25",
            "name": "云南省"
        },
        {
            "id": "26",
            "name": "西藏自治区"
        },
        {
            "id": "27",
            "name": "陕西省"
        },
        {
            "id": "28",
            "name": "甘肃省"
        },
        {
            "id": "29",
            "name": "青海省"
        },
        {
            "id": "30",
            "name": "宁夏回族自治区"
        },
        {
            "id": "31",
            "name": "新疆维吾尔自治区"
        },
        {
            "id": "32",
            "name": "香港特别行政区"
        },
        {
            "id": "33",
            "name": "澳门特别行政区"
        },
        {
            "id": "34",
            "name": "台湾省"
        }
    ],
    "CITIES": [
        {
            "zip_code": "100000",
            "province_id": "1",
            "id": "1",
            "name": "北京市"
        },
        {
            "zip_code": "100000",
            "province_id": "2",
            "id": "2",
            "name": "天津市"
        },
        {
            "zip_code": "050000",
            "province_id": "3",
            "id": "3",
            "name": "石家庄市"
        },
        {
            "zip_code": "063000",
            "province_id": "3",
            "id": "4",
            "name": "唐山市"
        },
        {
            "zip_code": "066000",
            "province_id": "3",
            "id": "5",
            "name": "秦皇岛市"
        },
        {
            "zip_code": "056000",
            "province_id": "3",
            "id": "6",
            "name": "邯郸市"
        },
        {
            "zip_code": "054000",
            "province_id": "3",
            "id": "7",
            "name": "邢台市"
        },
        {
            "zip_code": "071000",
            "province_id": "3",
            "id": "8",
            "name": "保定市"
        },
        {
            "zip_code": "075000",
            "province_id": "3",
            "id": "9",
            "name": "张家口市"
        },
        {
            "zip_code": "067000",
            "province_id": "3",
            "id": "10",
            "name": "承德市"
        },
        {
            "zip_code": "061000",
            "province_id": "3",
            "id": "11",
            "name": "沧州市"
        },
        {
            "zip_code": "065000",
            "province_id": "3",
            "id": "12",
            "name": "廊坊市"
        },
        {
            "zip_code": "053000",
            "province_id": "3",
            "id": "13",
            "name": "衡水市"
        },
        {
            "zip_code": "030000",
            "province_id": "4",
            "id": "14",
            "name": "太原市"
        },
        {
            "zip_code": "037000",
            "province_id": "4",
            "id": "15",
            "name": "大同市"
        },
        {
            "zip_code": "045000",
            "province_id": "4",
            "id": "16",
            "name": "阳泉市"
        },
        {
            "zip_code": "046000",
            "province_id": "4",
            "id": "17",
            "name": "长治市"
        },
        {
            "zip_code": "048000",
            "province_id": "4",
            "id": "18",
            "name": "晋城市"
        },
        {
            "zip_code": "036000",
            "province_id": "4",
            "id": "19",
            "name": "朔州市"
        },
        {
            "zip_code": "030600",
            "province_id": "4",
            "id": "20",
            "name": "晋中市"
        },
        {
            "zip_code": "044000",
            "province_id": "4",
            "id": "21",
            "name": "运城市"
        },
        {
            "zip_code": "034000",
            "province_id": "4",
            "id": "22",
            "name": "忻州市"
        },
        {
            "zip_code": "041000",
            "province_id": "4",
            "id": "23",
            "name": "临汾市"
        },
        {
            "zip_code": "030500",
            "province_id": "4",
            "id": "24",
            "name": "吕梁市"
        },
        {
            "zip_code": "010000",
            "province_id": "5",
            "id": "25",
            "name": "呼和浩特市"
        },
        {
            "zip_code": "014000",
            "province_id": "5",
            "id": "26",
            "name": "包头市"
        },
        {
            "zip_code": "016000",
            "province_id": "5",
            "id": "27",
            "name": "乌海市"
        },
        {
            "zip_code": "024000",
            "province_id": "5",
            "id": "28",
            "name": "赤峰市"
        },
        {
            "zip_code": "028000",
            "province_id": "5",
            "id": "29",
            "name": "通辽市"
        },
        {
            "zip_code": "010300",
            "province_id": "5",
            "id": "30",
            "name": "鄂尔多斯市"
        },
        {
            "zip_code": "021000",
            "province_id": "5",
            "id": "31",
            "name": "呼伦贝尔市"
        },
        {
            "zip_code": "014400",
            "province_id": "5",
            "id": "32",
            "name": "巴彦淖尔市"
        },
        {
            "zip_code": "011800",
            "province_id": "5",
            "id": "33",
            "name": "乌兰察布市"
        },
        {
            "zip_code": "137500",
            "province_id": "5",
            "id": "34",
            "name": "兴安盟"
        },
        {
            "zip_code": "011100",
            "province_id": "5",
            "id": "35",
            "name": "锡林郭勒盟"
        },
        {
            "zip_code": "016000",
            "province_id": "5",
            "id": "36",
            "name": "阿拉善盟"
        },
        {
            "zip_code": "110000",
            "province_id": "6",
            "id": "37",
            "name": "沈阳市"
        },
        {
            "zip_code": "116000",
            "province_id": "6",
            "id": "38",
            "name": "大连市"
        },
        {
            "zip_code": "114000",
            "province_id": "6",
            "id": "39",
            "name": "鞍山市"
        },
        {
            "zip_code": "113000",
            "province_id": "6",
            "id": "40",
            "name": "抚顺市"
        },
        {
            "zip_code": "117000",
            "province_id": "6",
            "id": "41",
            "name": "本溪市"
        },
        {
            "zip_code": "118000",
            "province_id": "6",
            "id": "42",
            "name": "丹东市"
        },
        {
            "zip_code": "121000",
            "province_id": "6",
            "id": "43",
            "name": "锦州市"
        },
        {
            "zip_code": "115000",
            "province_id": "6",
            "id": "44",
            "name": "营口市"
        },
        {
            "zip_code": "123000",
            "province_id": "6",
            "id": "45",
            "name": "阜新市"
        },
        {
            "zip_code": "111000",
            "province_id": "6",
            "id": "46",
            "name": "辽阳市"
        },
        {
            "zip_code": "124000",
            "province_id": "6",
            "id": "47",
            "name": "盘锦市"
        },
        {
            "zip_code": "112000",
            "province_id": "6",
            "id": "48",
            "name": "铁岭市"
        },
        {
            "zip_code": "122000",
            "province_id": "6",
            "id": "49",
            "name": "朝阳市"
        },
        {
            "zip_code": "125000",
            "province_id": "6",
            "id": "50",
            "name": "葫芦岛市"
        },
        {
            "zip_code": "130000",
            "province_id": "7",
            "id": "51",
            "name": "长春市"
        },
        {
            "zip_code": "132000",
            "province_id": "7",
            "id": "52",
            "name": "吉林市"
        },
        {
            "zip_code": "136000",
            "province_id": "7",
            "id": "53",
            "name": "四平市"
        },
        {
            "zip_code": "136200",
            "province_id": "7",
            "id": "54",
            "name": "辽源市"
        },
        {
            "zip_code": "134000",
            "province_id": "7",
            "id": "55",
            "name": "通化市"
        },
        {
            "zip_code": "134300",
            "province_id": "7",
            "id": "56",
            "name": "白山市"
        },
        {
            "zip_code": "131100",
            "province_id": "7",
            "id": "57",
            "name": "松原市"
        },
        {
            "zip_code": "137000",
            "province_id": "7",
            "id": "58",
            "name": "白城市"
        },
        {
            "zip_code": "133000",
            "province_id": "7",
            "id": "59",
            "name": "延边朝鲜族自治州"
        },
        {
            "zip_code": "150000",
            "province_id": "8",
            "id": "60",
            "name": "哈尔滨市"
        },
        {
            "zip_code": "161000",
            "province_id": "8",
            "id": "61",
            "name": "齐齐哈尔市"
        },
        {
            "zip_code": "158100",
            "province_id": "8",
            "id": "62",
            "name": "鸡西市"
        },
        {
            "zip_code": "154100",
            "province_id": "8",
            "id": "63",
            "name": "鹤岗市"
        },
        {
            "zip_code": "155100",
            "province_id": "8",
            "id": "64",
            "name": "双鸭山市"
        },
        {
            "zip_code": "163000",
            "province_id": "8",
            "id": "65",
            "name": "大庆市"
        },
        {
            "zip_code": "152300",
            "province_id": "8",
            "id": "66",
            "name": "伊春市"
        },
        {
            "zip_code": "154000",
            "province_id": "8",
            "id": "67",
            "name": "佳木斯市"
        },
        {
            "zip_code": "154600",
            "province_id": "8",
            "id": "68",
            "name": "七台河市"
        },
        {
            "zip_code": "157000",
            "province_id": "8",
            "id": "69",
            "name": "牡丹江市"
        },
        {
            "zip_code": "164300",
            "province_id": "8",
            "id": "70",
            "name": "黑河市"
        },
        {
            "zip_code": "152000",
            "province_id": "8",
            "id": "71",
            "name": "绥化市"
        },
        {
            "zip_code": "165000",
            "province_id": "8",
            "id": "72",
            "name": "大兴安岭地区"
        },
        {
            "zip_code": "200000",
            "province_id": "9",
            "id": "73",
            "name": "上海市"
        },
        {
            "zip_code": "210000",
            "province_id": "10",
            "id": "74",
            "name": "南京市"
        },
        {
            "zip_code": "214000",
            "province_id": "10",
            "id": "75",
            "name": "无锡市"
        },
        {
            "zip_code": "221000",
            "province_id": "10",
            "id": "76",
            "name": "徐州市"
        },
        {
            "zip_code": "213000",
            "province_id": "10",
            "id": "77",
            "name": "常州市"
        },
        {
            "zip_code": "215000",
            "province_id": "10",
            "id": "78",
            "name": "苏州市"
        },
        {
            "zip_code": "226000",
            "province_id": "10",
            "id": "79",
            "name": "南通市"
        },
        {
            "zip_code": "222000",
            "province_id": "10",
            "id": "80",
            "name": "连云港市"
        },
        {
            "zip_code": "223200",
            "province_id": "10",
            "id": "81",
            "name": "淮安市"
        },
        {
            "zip_code": "224000",
            "province_id": "10",
            "id": "82",
            "name": "盐城市"
        },
        {
            "zip_code": "225000",
            "province_id": "10",
            "id": "83",
            "name": "扬州市"
        },
        {
            "zip_code": "212000",
            "province_id": "10",
            "id": "84",
            "name": "镇江市"
        },
        {
            "zip_code": "225300",
            "province_id": "10",
            "id": "85",
            "name": "泰州市"
        },
        {
            "zip_code": "223800",
            "province_id": "10",
            "id": "86",
            "name": "宿迁市"
        },
        {
            "zip_code": "310000",
            "province_id": "11",
            "id": "87",
            "name": "杭州市"
        },
        {
            "zip_code": "315000",
            "province_id": "11",
            "id": "88",
            "name": "宁波市"
        },
        {
            "zip_code": "325000",
            "province_id": "11",
            "id": "89",
            "name": "温州市"
        },
        {
            "zip_code": "314000",
            "province_id": "11",
            "id": "90",
            "name": "嘉兴市"
        },
        {
            "zip_code": "313000",
            "province_id": "11",
            "id": "91",
            "name": "湖州市"
        },
        {
            "zip_code": "312000",
            "province_id": "11",
            "id": "92",
            "name": "绍兴市"
        },
        {
            "zip_code": "321000",
            "province_id": "11",
            "id": "93",
            "name": "金华市"
        },
        {
            "zip_code": "324000",
            "province_id": "11",
            "id": "94",
            "name": "衢州市"
        },
        {
            "zip_code": "316000",
            "province_id": "11",
            "id": "95",
            "name": "舟山市"
        },
        {
            "zip_code": "318000",
            "province_id": "11",
            "id": "96",
            "name": "台州市"
        },
        {
            "zip_code": "323000",
            "province_id": "11",
            "id": "97",
            "name": "丽水市"
        },
        {
            "zip_code": "230000",
            "province_id": "12",
            "id": "98",
            "name": "合肥市"
        },
        {
            "zip_code": "241000",
            "province_id": "12",
            "id": "99",
            "name": "芜湖市"
        },
        {
            "zip_code": "233000",
            "province_id": "12",
            "id": "100",
            "name": "蚌埠市"
        },
        {
            "zip_code": "232000",
            "province_id": "12",
            "id": "101",
            "name": "淮南市"
        },
        {
            "zip_code": "243000",
            "province_id": "12",
            "id": "102",
            "name": "马鞍山市"
        },
        {
            "zip_code": "235000",
            "province_id": "12",
            "id": "103",
            "name": "淮北市"
        },
        {
            "zip_code": "244000",
            "province_id": "12",
            "id": "104",
            "name": "铜陵市"
        },
        {
            "zip_code": "246000",
            "province_id": "12",
            "id": "105",
            "name": "安庆市"
        },
        {
            "zip_code": "242700",
            "province_id": "12",
            "id": "106",
            "name": "黄山市"
        },
        {
            "zip_code": "239000",
            "province_id": "12",
            "id": "107",
            "name": "滁州市"
        },
        {
            "zip_code": "236100",
            "province_id": "12",
            "id": "108",
            "name": "阜阳市"
        },
        {
            "zip_code": "234100",
            "province_id": "12",
            "id": "109",
            "name": "宿州市"
        },
        {
            "zip_code": "238000",
            "province_id": "12",
            "id": "110",
            "name": "巢湖市"
        },
        {
            "zip_code": "237000",
            "province_id": "12",
            "id": "111",
            "name": "六安市"
        },
        {
            "zip_code": "236800",
            "province_id": "12",
            "id": "112",
            "name": "亳州市"
        },
        {
            "zip_code": "247100",
            "province_id": "12",
            "id": "113",
            "name": "池州市"
        },
        {
            "zip_code": "366000",
            "province_id": "12",
            "id": "114",
            "name": "宣城市"
        },
        {
            "zip_code": "350000",
            "province_id": "13",
            "id": "115",
            "name": "福州市"
        },
        {
            "zip_code": "361000",
            "province_id": "13",
            "id": "116",
            "name": "厦门市"
        },
        {
            "zip_code": "351100",
            "province_id": "13",
            "id": "117",
            "name": "莆田市"
        },
        {
            "zip_code": "365000",
            "province_id": "13",
            "id": "118",
            "name": "三明市"
        },
        {
            "zip_code": "362000",
            "province_id": "13",
            "id": "119",
            "name": "泉州市"
        },
        {
            "zip_code": "363000",
            "province_id": "13",
            "id": "120",
            "name": "漳州市"
        },
        {
            "zip_code": "353000",
            "province_id": "13",
            "id": "121",
            "name": "南平市"
        },
        {
            "zip_code": "364000",
            "province_id": "13",
            "id": "122",
            "name": "龙岩市"
        },
        {
            "zip_code": "352100",
            "province_id": "13",
            "id": "123",
            "name": "宁德市"
        },
        {
            "zip_code": "330000",
            "province_id": "14",
            "id": "124",
            "name": "南昌市"
        },
        {
            "zip_code": "333000",
            "province_id": "14",
            "id": "125",
            "name": "景德镇市"
        },
        {
            "zip_code": "337000",
            "province_id": "14",
            "id": "126",
            "name": "萍乡市"
        },
        {
            "zip_code": "332000",
            "province_id": "14",
            "id": "127",
            "name": "九江市"
        },
        {
            "zip_code": "338000",
            "province_id": "14",
            "id": "128",
            "name": "新余市"
        },
        {
            "zip_code": "335000",
            "province_id": "14",
            "id": "129",
            "name": "鹰潭市"
        },
        {
            "zip_code": "341000",
            "province_id": "14",
            "id": "130",
            "name": "赣州市"
        },
        {
            "zip_code": "343000",
            "province_id": "14",
            "id": "131",
            "name": "吉安市"
        },
        {
            "zip_code": "336000",
            "province_id": "14",
            "id": "132",
            "name": "宜春市"
        },
        {
            "zip_code": "332900",
            "province_id": "14",
            "id": "133",
            "name": "抚州市"
        },
        {
            "zip_code": "334000",
            "province_id": "14",
            "id": "134",
            "name": "上饶市"
        },
        {
            "zip_code": "250000",
            "province_id": "15",
            "id": "135",
            "name": "济南市"
        },
        {
            "zip_code": "266000",
            "province_id": "15",
            "id": "136",
            "name": "青岛市"
        },
        {
            "zip_code": "255000",
            "province_id": "15",
            "id": "137",
            "name": "淄博市"
        },
        {
            "zip_code": "277100",
            "province_id": "15",
            "id": "138",
            "name": "枣庄市"
        },
        {
            "zip_code": "257000",
            "province_id": "15",
            "id": "139",
            "name": "东营市"
        },
        {
            "zip_code": "264000",
            "province_id": "15",
            "id": "140",
            "name": "烟台市"
        },
        {
            "zip_code": "261000",
            "province_id": "15",
            "id": "141",
            "name": "潍坊市"
        },
        {
            "zip_code": "272100",
            "province_id": "15",
            "id": "142",
            "name": "济宁市"
        },
        {
            "zip_code": "271000",
            "province_id": "15",
            "id": "143",
            "name": "泰安市"
        },
        {
            "zip_code": "265700",
            "province_id": "15",
            "id": "144",
            "name": "威海市"
        },
        {
            "zip_code": "276800",
            "province_id": "15",
            "id": "145",
            "name": "日照市"
        },
        {
            "zip_code": "271100",
            "province_id": "15",
            "id": "146",
            "name": "莱芜市"
        },
        {
            "zip_code": "276000",
            "province_id": "15",
            "id": "147",
            "name": "临沂市"
        },
        {
            "zip_code": "253000",
            "province_id": "15",
            "id": "148",
            "name": "德州市"
        },
        {
            "zip_code": "252000",
            "province_id": "15",
            "id": "149",
            "name": "聊城市"
        },
        {
            "zip_code": "256600",
            "province_id": "15",
            "id": "150",
            "name": "滨州市"
        },
        {
            "zip_code": "255000",
            "province_id": "15",
            "id": "151",
            "name": "菏泽市"
        },
        {
            "zip_code": "450000",
            "province_id": "16",
            "id": "152",
            "name": "郑州市"
        },
        {
            "zip_code": "475000",
            "province_id": "16",
            "id": "153",
            "name": "开封市"
        },
        {
            "zip_code": "471000",
            "province_id": "16",
            "id": "154",
            "name": "洛阳市"
        },
        {
            "zip_code": "467000",
            "province_id": "16",
            "id": "155",
            "name": "平顶山市"
        },
        {
            "zip_code": "454900",
            "province_id": "16",
            "id": "156",
            "name": "安阳市"
        },
        {
            "zip_code": "456600",
            "province_id": "16",
            "id": "157",
            "name": "鹤壁市"
        },
        {
            "zip_code": "453000",
            "province_id": "16",
            "id": "158",
            "name": "新乡市"
        },
        {
            "zip_code": "454100",
            "province_id": "16",
            "id": "159",
            "name": "焦作市"
        },
        {
            "zip_code": "457000",
            "province_id": "16",
            "id": "160",
            "name": "濮阳市"
        },
        {
            "zip_code": "461000",
            "province_id": "16",
            "id": "161",
            "name": "许昌市"
        },
        {
            "zip_code": "462000",
            "province_id": "16",
            "id": "162",
            "name": "漯河市"
        },
        {
            "zip_code": "472000",
            "province_id": "16",
            "id": "163",
            "name": "三门峡市"
        },
        {
            "zip_code": "473000",
            "province_id": "16",
            "id": "164",
            "name": "南阳市"
        },
        {
            "zip_code": "476000",
            "province_id": "16",
            "id": "165",
            "name": "商丘市"
        },
        {
            "zip_code": "464000",
            "province_id": "16",
            "id": "166",
            "name": "信阳市"
        },
        {
            "zip_code": "466000",
            "province_id": "16",
            "id": "167",
            "name": "周口市"
        },
        {
            "zip_code": "463000",
            "province_id": "16",
            "id": "168",
            "name": "驻马店市"
        },
        {
            "zip_code": "430000",
            "province_id": "17",
            "id": "169",
            "name": "武汉市"
        },
        {
            "zip_code": "435000",
            "province_id": "17",
            "id": "170",
            "name": "黄石市"
        },
        {
            "zip_code": "442000",
            "province_id": "17",
            "id": "171",
            "name": "十堰市"
        },
        {
            "zip_code": "443000",
            "province_id": "17",
            "id": "172",
            "name": "宜昌市"
        },
        {
            "zip_code": "441000",
            "province_id": "17",
            "id": "173",
            "name": "襄阳市"
        },
        {
            "zip_code": "436000",
            "province_id": "17",
            "id": "174",
            "name": "鄂州市"
        },
        {
            "zip_code": "448000",
            "province_id": "17",
            "id": "175",
            "name": "荆门市"
        },
        {
            "zip_code": "432100",
            "province_id": "17",
            "id": "176",
            "name": "孝感市"
        },
        {
            "zip_code": "434000",
            "province_id": "17",
            "id": "177",
            "name": "荆州市"
        },
        {
            "zip_code": "438000",
            "province_id": "17",
            "id": "178",
            "name": "黄冈市"
        },
        {
            "zip_code": "437000",
            "province_id": "17",
            "id": "179",
            "name": "咸宁市"
        },
        {
            "zip_code": "441300",
            "province_id": "17",
            "id": "180",
            "name": "随州市"
        },
        {
            "zip_code": "445000",
            "province_id": "17",
            "id": "181",
            "name": "恩施土家族苗族自治州"
        },
        {
            "zip_code": "442400",
            "province_id": "17",
            "id": "182",
            "name": "省直辖县"
        },
        {
            "zip_code": "410000",
            "province_id": "18",
            "id": "183",
            "name": "长沙市"
        },
        {
            "zip_code": "412000",
            "province_id": "18",
            "id": "184",
            "name": "株洲市"
        },
        {
            "zip_code": "411100",
            "province_id": "18",
            "id": "185",
            "name": "湘潭市"
        },
        {
            "zip_code": "421000",
            "province_id": "18",
            "id": "186",
            "name": "衡阳市"
        },
        {
            "zip_code": "422000",
            "province_id": "18",
            "id": "187",
            "name": "邵阳市"
        },
        {
            "zip_code": "414000",
            "province_id": "18",
            "id": "188",
            "name": "岳阳市"
        },
        {
            "zip_code": "415000",
            "province_id": "18",
            "id": "189",
            "name": "常德市"
        },
        {
            "zip_code": "427000",
            "province_id": "18",
            "id": "190",
            "name": "张家界市"
        },
        {
            "zip_code": "413000",
            "province_id": "18",
            "id": "191",
            "name": "益阳市"
        },
        {
            "zip_code": "423000",
            "province_id": "18",
            "id": "192",
            "name": "郴州市"
        },
        {
            "zip_code": "425000",
            "province_id": "18",
            "id": "193",
            "name": "永州市"
        },
        {
            "zip_code": "418000",
            "province_id": "18",
            "id": "194",
            "name": "怀化市"
        },
        {
            "zip_code": "417000",
            "province_id": "18",
            "id": "195",
            "name": "娄底市"
        },
        {
            "zip_code": "416000",
            "province_id": "18",
            "id": "196",
            "name": "湘西土家族苗族自治州"
        },
        {
            "zip_code": "510000",
            "province_id": "19",
            "id": "197",
            "name": "广州市"
        },
        {
            "zip_code": "521000",
            "province_id": "19",
            "id": "198",
            "name": "韶关市"
        },
        {
            "zip_code": "518000",
            "province_id": "19",
            "id": "199",
            "name": "深圳市"
        },
        {
            "zip_code": "519000",
            "province_id": "19",
            "id": "200",
            "name": "珠海市"
        },
        {
            "zip_code": "515000",
            "province_id": "19",
            "id": "201",
            "name": "汕头市"
        },
        {
            "zip_code": "528000",
            "province_id": "19",
            "id": "202",
            "name": "佛山市"
        },
        {
            "zip_code": "529000",
            "province_id": "19",
            "id": "203",
            "name": "江门市"
        },
        {
            "zip_code": "524000",
            "province_id": "19",
            "id": "204",
            "name": "湛江市"
        },
        {
            "zip_code": "525000",
            "province_id": "19",
            "id": "205",
            "name": "茂名市"
        },
        {
            "zip_code": "526000",
            "province_id": "19",
            "id": "206",
            "name": "肇庆市"
        },
        {
            "zip_code": "516000",
            "province_id": "19",
            "id": "207",
            "name": "惠州市"
        },
        {
            "zip_code": "514000",
            "province_id": "19",
            "id": "208",
            "name": "梅州市"
        },
        {
            "zip_code": "516600",
            "province_id": "19",
            "id": "209",
            "name": "汕尾市"
        },
        {
            "zip_code": "517000",
            "province_id": "19",
            "id": "210",
            "name": "河源市"
        },
        {
            "zip_code": "529500",
            "province_id": "19",
            "id": "211",
            "name": "阳江市"
        },
        {
            "zip_code": "511500",
            "province_id": "19",
            "id": "212",
            "name": "清远市"
        },
        {
            "zip_code": "511700",
            "province_id": "19",
            "id": "213",
            "name": "东莞市"
        },
        {
            "zip_code": "528400",
            "province_id": "19",
            "id": "214",
            "name": "中山市"
        },
        {
            "zip_code": "515600",
            "province_id": "19",
            "id": "215",
            "name": "潮州市"
        },
        {
            "zip_code": "522000",
            "province_id": "19",
            "id": "216",
            "name": "揭阳市"
        },
        {
            "zip_code": "527300",
            "province_id": "19",
            "id": "217",
            "name": "云浮市"
        },
        {
            "zip_code": "530000",
            "province_id": "20",
            "id": "218",
            "name": "南宁市"
        },
        {
            "zip_code": "545000",
            "province_id": "20",
            "id": "219",
            "name": "柳州市"
        },
        {
            "zip_code": "541000",
            "province_id": "20",
            "id": "220",
            "name": "桂林市"
        },
        {
            "zip_code": "543000",
            "province_id": "20",
            "id": "221",
            "name": "梧州市"
        },
        {
            "zip_code": "536000",
            "province_id": "20",
            "id": "222",
            "name": "北海市"
        },
        {
            "zip_code": "538000",
            "province_id": "20",
            "id": "223",
            "name": "防城港市"
        },
        {
            "zip_code": "535000",
            "province_id": "20",
            "id": "224",
            "name": "钦州市"
        },
        {
            "zip_code": "537100",
            "province_id": "20",
            "id": "225",
            "name": "贵港市"
        },
        {
            "zip_code": "537000",
            "province_id": "20",
            "id": "226",
            "name": "玉林市"
        },
        {
            "zip_code": "533000",
            "province_id": "20",
            "id": "227",
            "name": "百色市"
        },
        {
            "zip_code": "542800",
            "province_id": "20",
            "id": "228",
            "name": "贺州市"
        },
        {
            "zip_code": "547000",
            "province_id": "20",
            "id": "229",
            "name": "河池市"
        },
        {
            "zip_code": "546100",
            "province_id": "20",
            "id": "230",
            "name": "来宾市"
        },
        {
            "zip_code": "532200",
            "province_id": "20",
            "id": "231",
            "name": "崇左市"
        },
        {
            "zip_code": "570000",
            "province_id": "21",
            "id": "232",
            "name": "海口市"
        },
        {
            "zip_code": "572000",
            "province_id": "21",
            "id": "233",
            "name": "三亚市"
        },
        {
            "zip_code": "400000",
            "province_id": "22",
            "id": "234",
            "name": "重庆市"
        },
        {
            "zip_code": "610000",
            "province_id": "23",
            "id": "235",
            "name": "成都市"
        },
        {
            "zip_code": "643000",
            "province_id": "23",
            "id": "236",
            "name": "自贡市"
        },
        {
            "zip_code": "617000",
            "province_id": "23",
            "id": "237",
            "name": "攀枝花市"
        },
        {
            "zip_code": "646100",
            "province_id": "23",
            "id": "238",
            "name": "泸州市"
        },
        {
            "zip_code": "618000",
            "province_id": "23",
            "id": "239",
            "name": "德阳市"
        },
        {
            "zip_code": "621000",
            "province_id": "23",
            "id": "240",
            "name": "绵阳市"
        },
        {
            "zip_code": "628000",
            "province_id": "23",
            "id": "241",
            "name": "广元市"
        },
        {
            "zip_code": "629000",
            "province_id": "23",
            "id": "242",
            "name": "遂宁市"
        },
        {
            "zip_code": "641000",
            "province_id": "23",
            "id": "243",
            "name": "内江市"
        },
        {
            "zip_code": "614000",
            "province_id": "23",
            "id": "244",
            "name": "乐山市"
        },
        {
            "zip_code": "637000",
            "province_id": "23",
            "id": "245",
            "name": "南充市"
        },
        {
            "zip_code": "612100",
            "province_id": "23",
            "id": "246",
            "name": "眉山市"
        },
        {
            "zip_code": "644000",
            "province_id": "23",
            "id": "247",
            "name": "宜宾市"
        },
        {
            "zip_code": "638000",
            "province_id": "23",
            "id": "248",
            "name": "广安市"
        },
        {
            "zip_code": "635000",
            "province_id": "23",
            "id": "249",
            "name": "达州市"
        },
        {
            "zip_code": "625000",
            "province_id": "23",
            "id": "250",
            "name": "雅安市"
        },
        {
            "zip_code": "635500",
            "province_id": "23",
            "id": "251",
            "name": "巴中市"
        },
        {
            "zip_code": "641300",
            "province_id": "23",
            "id": "252",
            "name": "资阳市"
        },
        {
            "zip_code": "624600",
            "province_id": "23",
            "id": "253",
            "name": "阿坝藏族羌族自治州"
        },
        {
            "zip_code": "626000",
            "province_id": "23",
            "id": "254",
            "name": "甘孜藏族自治州"
        },
        {
            "zip_code": "615000",
            "province_id": "23",
            "id": "255",
            "name": "凉山彝族自治州"
        },
        {
            "zip_code": "55000",
            "province_id": "24",
            "id": "256",
            "name": "贵阳市"
        },
        {
            "zip_code": "553000",
            "province_id": "24",
            "id": "257",
            "name": "六盘水市"
        },
        {
            "zip_code": "563000",
            "province_id": "24",
            "id": "258",
            "name": "遵义市"
        },
        {
            "zip_code": "561000",
            "province_id": "24",
            "id": "259",
            "name": "安顺市"
        },
        {
            "zip_code": "554300",
            "province_id": "24",
            "id": "260",
            "name": "铜仁市"
        },
        {
            "zip_code": "551500",
            "province_id": "24",
            "id": "261",
            "name": "黔西南布依族苗族自治州"
        },
        {
            "zip_code": "551700",
            "province_id": "24",
            "id": "262",
            "name": "毕节市"
        },
        {
            "zip_code": "551500",
            "province_id": "24",
            "id": "263",
            "name": "黔东南苗族侗族自治州"
        },
        {
            "zip_code": "550100",
            "province_id": "24",
            "id": "264",
            "name": "黔南布依族苗族自治州"
        },
        {
            "zip_code": "650000",
            "province_id": "25",
            "id": "265",
            "name": "昆明市"
        },
        {
            "zip_code": "655000",
            "province_id": "25",
            "id": "266",
            "name": "曲靖市"
        },
        {
            "zip_code": "653100",
            "province_id": "25",
            "id": "267",
            "name": "玉溪市"
        },
        {
            "zip_code": "678000",
            "province_id": "25",
            "id": "268",
            "name": "保山市"
        },
        {
            "zip_code": "657000",
            "province_id": "25",
            "id": "269",
            "name": "昭通市"
        },
        {
            "zip_code": "674100",
            "province_id": "25",
            "id": "270",
            "name": "丽江市"
        },
        {
            "zip_code": "665000",
            "province_id": "25",
            "id": "271",
            "name": "普洱市"
        },
        {
            "zip_code": "677000",
            "province_id": "25",
            "id": "272",
            "name": "临沧市"
        },
        {
            "zip_code": "675000",
            "province_id": "25",
            "id": "273",
            "name": "楚雄彝族自治州"
        },
        {
            "zip_code": "654400",
            "province_id": "25",
            "id": "274",
            "name": "红河哈尼族彝族自治州"
        },
        {
            "zip_code": "663000",
            "province_id": "25",
            "id": "275",
            "name": "文山壮族苗族自治州"
        },
        {
            "zip_code": "666200",
            "province_id": "25",
            "id": "276",
            "name": "西双版纳傣族自治州"
        },
        {
            "zip_code": "671000",
            "province_id": "25",
            "id": "277",
            "name": "大理白族自治州"
        },
        {
            "zip_code": "678400",
            "province_id": "25",
            "id": "278",
            "name": "德宏傣族景颇族自治州"
        },
        {
            "zip_code": "671400",
            "province_id": "25",
            "id": "279",
            "name": "怒江傈僳族自治州"
        },
        {
            "zip_code": "674400",
            "province_id": "25",
            "id": "280",
            "name": "迪庆藏族自治州"
        },
        {
            "zip_code": "850000",
            "province_id": "26",
            "id": "281",
            "name": "拉萨市"
        },
        {
            "zip_code": "854000",
            "province_id": "26",
            "id": "282",
            "name": "昌都市"
        },
        {
            "zip_code": "856000",
            "province_id": "26",
            "id": "283",
            "name": "山南市"
        },
        {
            "zip_code": "857000",
            "province_id": "26",
            "id": "284",
            "name": "日喀则市"
        },
        {
            "zip_code": "852000",
            "province_id": "26",
            "id": "285",
            "name": "那曲市"
        },
        {
            "zip_code": "859100",
            "province_id": "26",
            "id": "286",
            "name": "阿里地区"
        },
        {
            "zip_code": "860100",
            "province_id": "26",
            "id": "287",
            "name": "林芝市"
        },
        {
            "zip_code": "710000",
            "province_id": "27",
            "id": "288",
            "name": "西安市"
        },
        {
            "zip_code": "727000",
            "province_id": "27",
            "id": "289",
            "name": "铜川市"
        },
        {
            "zip_code": "721000",
            "province_id": "27",
            "id": "290",
            "name": "宝鸡市"
        },
        {
            "zip_code": "712000",
            "province_id": "27",
            "id": "291",
            "name": "咸阳市"
        },
        {
            "zip_code": "714000",
            "province_id": "27",
            "id": "292",
            "name": "渭南市"
        },
        {
            "zip_code": "716000",
            "province_id": "27",
            "id": "293",
            "name": "延安市"
        },
        {
            "zip_code": "723000",
            "province_id": "27",
            "id": "294",
            "name": "汉中市"
        },
        {
            "zip_code": "719000",
            "province_id": "27",
            "id": "295",
            "name": "榆林市"
        },
        {
            "zip_code": "725000",
            "province_id": "27",
            "id": "296",
            "name": "安康市"
        },
        {
            "zip_code": "711500",
            "province_id": "27",
            "id": "297",
            "name": "商洛市"
        },
        {
            "zip_code": "730000",
            "province_id": "28",
            "id": "298",
            "name": "兰州市"
        },
        {
            "zip_code": "735100",
            "province_id": "28",
            "id": "299",
            "name": "嘉峪关市"
        },
        {
            "zip_code": "737100",
            "province_id": "28",
            "id": "300",
            "name": "金昌市"
        },
        {
            "zip_code": "730900",
            "province_id": "28",
            "id": "301",
            "name": "白银市"
        },
        {
            "zip_code": "741000",
            "province_id": "28",
            "id": "302",
            "name": "天水市"
        },
        {
            "zip_code": "733000",
            "province_id": "28",
            "id": "303",
            "name": "武威市"
        },
        {
            "zip_code": "734000",
            "province_id": "28",
            "id": "304",
            "name": "张掖市"
        },
        {
            "zip_code": "744000",
            "province_id": "28",
            "id": "305",
            "name": "平凉市"
        },
        {
            "zip_code": "735000",
            "province_id": "28",
            "id": "306",
            "name": "酒泉市"
        },
        {
            "zip_code": "744500",
            "province_id": "28",
            "id": "307",
            "name": "庆阳市"
        },
        {
            "zip_code": "743000",
            "province_id": "28",
            "id": "308",
            "name": "定西市"
        },
        {
            "zip_code": "742100",
            "province_id": "28",
            "id": "309",
            "name": "陇南市"
        },
        {
            "zip_code": "731100",
            "province_id": "28",
            "id": "310",
            "name": "临夏回族自治州"
        },
        {
            "zip_code": "747000",
            "province_id": "28",
            "id": "311",
            "name": "甘南藏族自治州"
        },
        {
            "zip_code": "810000",
            "province_id": "29",
            "id": "312",
            "name": "西宁市"
        },
        {
            "zip_code": "810600",
            "province_id": "29",
            "id": "313",
            "name": "海东市"
        },
        {
            "zip_code": "810300",
            "province_id": "29",
            "id": "314",
            "name": "海北藏族自治州"
        },
        {
            "zip_code": "811300",
            "province_id": "29",
            "id": "315",
            "name": "黄南藏族自治州"
        },
        {
            "zip_code": "813000",
            "province_id": "29",
            "id": "316",
            "name": "海南藏族自治州"
        },
        {
            "zip_code": "814000",
            "province_id": "29",
            "id": "317",
            "name": "果洛藏族自治州"
        },
        {
            "zip_code": "815000",
            "province_id": "29",
            "id": "318",
            "name": "玉树藏族自治州"
        },
        {
            "zip_code": "817000",
            "province_id": "29",
            "id": "319",
            "name": "海西蒙古族藏族自治州"
        },
        {
            "zip_code": "750000",
            "province_id": "30",
            "id": "320",
            "name": "银川市"
        },
        {
            "zip_code": "753000",
            "province_id": "30",
            "id": "321",
            "name": "石嘴山市"
        },
        {
            "zip_code": "751100",
            "province_id": "30",
            "id": "322",
            "name": "吴忠市"
        },
        {
            "zip_code": "756000",
            "province_id": "30",
            "id": "323",
            "name": "固原市"
        },
        {
            "zip_code": "751700",
            "province_id": "30",
            "id": "324",
            "name": "中卫市"
        },
        {
            "zip_code": "830000",
            "province_id": "31",
            "id": "325",
            "name": "乌鲁木齐市"
        },
        {
            "zip_code": "834000",
            "province_id": "31",
            "id": "326",
            "name": "克拉玛依市"
        },
        {
            "zip_code": "838000",
            "province_id": "31",
            "id": "327",
            "name": "吐鲁番市"
        },
        {
            "zip_code": "839000",
            "province_id": "31",
            "id": "328",
            "name": "哈密市"
        },
        {
            "zip_code": "831100",
            "province_id": "31",
            "id": "329",
            "name": "昌吉回族自治州"
        },
        {
            "zip_code": "833400",
            "province_id": "31",
            "id": "330",
            "name": "博尔塔拉蒙古自治州"
        },
        {
            "zip_code": "841000",
            "province_id": "31",
            "id": "331",
            "name": "巴音郭楞蒙古自治州"
        },
        {
            "zip_code": "843000",
            "province_id": "31",
            "id": "332",
            "name": "阿克苏地区"
        },
        {
            "zip_code": "835600",
            "province_id": "31",
            "id": "333",
            "name": "克孜勒苏柯尔克孜自治州"
        },
        {
            "zip_code": "844000",
            "province_id": "31",
            "id": "334",
            "name": "喀什地区"
        },
        {
            "zip_code": "848000",
            "province_id": "31",
            "id": "335",
            "name": "和田地区"
        },
        {
            "zip_code": "833200",
            "province_id": "31",
            "id": "336",
            "name": "伊犁哈萨克自治州"
        },
        {
            "zip_code": "834700",
            "province_id": "31",
            "id": "337",
            "name": "塔城地区"
        },
        {
            "zip_code": "836500",
            "province_id": "31",
            "id": "338",
            "name": "阿勒泰地区"
        },
        {
            "zip_code": "832000",
            "province_id": "31",
            "id": "339",
            "name": "石河子市"
        },
        {
            "zip_code": "843300",
            "province_id": "31",
            "id": "340",
            "name": "阿拉尔市"
        },
        {
            "zip_code": "843900",
            "province_id": "31",
            "id": "341",
            "name": "图木舒克市"
        },
        {
            "zip_code": "831300",
            "province_id": "31",
            "id": "342",
            "name": "五家渠市"
        },
        {
            "zip_code": "000000",
            "province_id": "32",
            "id": "343",
            "name": "香港特别行政区"
        },
        {
            "zip_code": "000000",
            "province_id": "33",
            "id": "344",
            "name": "澳门特别行政区"
        },
        {
            "zip_code": "000000",
            "province_id": "34",
            "id": "345",
            "name": "台湾省"
        },
        {
            "zip_code": "000000",
            "province_id": "21",
            "id": "346",
            "name": "三沙市"
        },
        {
            "zip_code": "400000",
            "province_id": "22",
            "id": "348",
            "name": "县"
        },
        {
            "zip_code": "572000",
            "province_id": "21",
            "id": "349",
            "name": "省直辖县"
        },
        {
            "zip_code": "000000",
            "province_id": "16",
            "id": "350",
            "name": "省直辖县"
        }
    ],
    "DISTRICTS": [
        {
            "city_id": "1",
            "id": "1",
            "name": "东城区"
        },
        {
            "city_id": "1",
            "id": "2",
            "name": "西城区"
        },
        {
            "city_id": "1",
            "id": "3",
            "name": "崇文区"
        },
        {
            "city_id": "1",
            "id": "4",
            "name": "宣武区"
        },
        {
            "city_id": "1",
            "id": "5",
            "name": "朝阳区"
        },
        {
            "city_id": "1",
            "id": "6",
            "name": "丰台区"
        },
        {
            "city_id": "1",
            "id": "7",
            "name": "石景山区"
        },
        {
            "city_id": "1",
            "id": "8",
            "name": "海淀区"
        },
        {
            "city_id": "1",
            "id": "9",
            "name": "门头沟区"
        },
        {
            "city_id": "1",
            "id": "10",
            "name": "房山区"
        },
        {
            "city_id": "1",
            "id": "11",
            "name": "通州区"
        },
        {
            "city_id": "1",
            "id": "12",
            "name": "顺义区"
        },
        {
            "city_id": "1",
            "id": "13",
            "name": "昌平区"
        },
        {
            "city_id": "1",
            "id": "14",
            "name": "大兴区"
        },
        {
            "city_id": "1",
            "id": "15",
            "name": "怀柔区"
        },
        {
            "city_id": "1",
            "id": "16",
            "name": "平谷区"
        },
        {
            "city_id": "1",
            "id": "17",
            "name": "密云区"
        },
        {
            "city_id": "1",
            "id": "18",
            "name": "延庆区"
        },
        {
            "city_id": "2",
            "id": "19",
            "name": "和平区"
        },
        {
            "city_id": "2",
            "id": "20",
            "name": "河东区"
        },
        {
            "city_id": "2",
            "id": "21",
            "name": "河西区"
        },
        {
            "city_id": "2",
            "id": "22",
            "name": "南开区"
        },
        {
            "city_id": "2",
            "id": "23",
            "name": "河北区"
        },
        {
            "city_id": "2",
            "id": "24",
            "name": "红桥区"
        },
        {
            "city_id": "2",
            "id": "25",
            "name": "塘沽区"
        },
        {
            "city_id": "2",
            "id": "26",
            "name": "汉沽区"
        },
        {
            "city_id": "2",
            "id": "27",
            "name": "大港区"
        },
        {
            "city_id": "2",
            "id": "28",
            "name": "东丽区"
        },
        {
            "city_id": "2",
            "id": "29",
            "name": "西青区"
        },
        {
            "city_id": "2",
            "id": "30",
            "name": "津南区"
        },
        {
            "city_id": "2",
            "id": "31",
            "name": "北辰区"
        },
        {
            "city_id": "2",
            "id": "32",
            "name": "武清区"
        },
        {
            "city_id": "2",
            "id": "33",
            "name": "宝坻区"
        },
        {
            "city_id": "2",
            "id": "34",
            "name": "宁河区"
        },
        {
            "city_id": "2",
            "id": "35",
            "name": "静海区"
        },
        {
            "city_id": "2",
            "id": "36",
            "name": "蓟县"
        },
        {
            "city_id": "3",
            "id": "37",
            "name": "长安区"
        },
        {
            "city_id": "3",
            "id": "38",
            "name": "桥东区"
        },
        {
            "city_id": "3",
            "id": "39",
            "name": "桥西区"
        },
        {
            "city_id": "3",
            "id": "40",
            "name": "新华区"
        },
        {
            "city_id": "3",
            "id": "41",
            "name": "井陉矿区"
        },
        {
            "city_id": "3",
            "id": "42",
            "name": "裕华区"
        },
        {
            "city_id": "3",
            "id": "43",
            "name": "井陉县"
        },
        {
            "city_id": "3",
            "id": "44",
            "name": "正定县"
        },
        {
            "city_id": "3",
            "id": "45",
            "name": "栾城区"
        },
        {
            "city_id": "3",
            "id": "46",
            "name": "行唐县"
        },
        {
            "city_id": "3",
            "id": "47",
            "name": "灵寿县"
        },
        {
            "city_id": "3",
            "id": "48",
            "name": "高邑县"
        },
        {
            "city_id": "3",
            "id": "49",
            "name": "深泽县"
        },
        {
            "city_id": "3",
            "id": "50",
            "name": "赞皇县"
        },
        {
            "city_id": "3",
            "id": "51",
            "name": "无极县"
        },
        {
            "city_id": "3",
            "id": "52",
            "name": "平山县"
        },
        {
            "city_id": "3",
            "id": "53",
            "name": "元氏县"
        },
        {
            "city_id": "3",
            "id": "54",
            "name": "赵县"
        },
        {
            "city_id": "3",
            "id": "55",
            "name": "辛集市"
        },
        {
            "city_id": "3",
            "id": "56",
            "name": "藁城区"
        },
        {
            "city_id": "3",
            "id": "57",
            "name": "晋州市"
        },
        {
            "city_id": "3",
            "id": "58",
            "name": "新乐市"
        },
        {
            "city_id": "3",
            "id": "59",
            "name": "鹿泉区"
        },
        {
            "city_id": "4",
            "id": "60",
            "name": "路南区"
        },
        {
            "city_id": "4",
            "id": "61",
            "name": "路北区"
        },
        {
            "city_id": "4",
            "id": "62",
            "name": "古冶区"
        },
        {
            "city_id": "4",
            "id": "63",
            "name": "开平区"
        },
        {
            "city_id": "4",
            "id": "64",
            "name": "丰南区"
        },
        {
            "city_id": "4",
            "id": "65",
            "name": "丰润区"
        },
        {
            "city_id": "4",
            "id": "66",
            "name": "滦县"
        },
        {
            "city_id": "4",
            "id": "67",
            "name": "滦南县"
        },
        {
            "city_id": "4",
            "id": "68",
            "name": "乐亭县"
        },
        {
            "city_id": "4",
            "id": "69",
            "name": "迁西县"
        },
        {
            "city_id": "4",
            "id": "70",
            "name": "玉田县"
        },
        {
            "city_id": "4",
            "id": "71",
            "name": "唐海县"
        },
        {
            "city_id": "4",
            "id": "72",
            "name": "遵化市"
        },
        {
            "city_id": "4",
            "id": "73",
            "name": "迁安市"
        },
        {
            "city_id": "5",
            "id": "74",
            "name": "海港区"
        },
        {
            "city_id": "5",
            "id": "75",
            "name": "山海关区"
        },
        {
            "city_id": "5",
            "id": "76",
            "name": "北戴河区"
        },
        {
            "city_id": "5",
            "id": "77",
            "name": "青龙满族自治县"
        },
        {
            "city_id": "5",
            "id": "78",
            "name": "昌黎县"
        },
        {
            "city_id": "5",
            "id": "79",
            "name": "抚宁区"
        },
        {
            "city_id": "5",
            "id": "80",
            "name": "卢龙县"
        },
        {
            "city_id": "6",
            "id": "81",
            "name": "邯山区"
        },
        {
            "city_id": "6",
            "id": "82",
            "name": "丛台区"
        },
        {
            "city_id": "6",
            "id": "83",
            "name": "复兴区"
        },
        {
            "city_id": "6",
            "id": "84",
            "name": "峰峰矿区"
        },
        {
            "city_id": "6",
            "id": "85",
            "name": "邯郸县"
        },
        {
            "city_id": "6",
            "id": "86",
            "name": "临漳县"
        },
        {
            "city_id": "6",
            "id": "87",
            "name": "成安县"
        },
        {
            "city_id": "6",
            "id": "88",
            "name": "大名县"
        },
        {
            "city_id": "6",
            "id": "89",
            "name": "涉县"
        },
        {
            "city_id": "6",
            "id": "90",
            "name": "磁县"
        },
        {
            "city_id": "6",
            "id": "91",
            "name": "肥乡区"
        },
        {
            "city_id": "6",
            "id": "92",
            "name": "永年区"
        },
        {
            "city_id": "6",
            "id": "93",
            "name": "邱县"
        },
        {
            "city_id": "6",
            "id": "94",
            "name": "鸡泽县"
        },
        {
            "city_id": "6",
            "id": "95",
            "name": "广平县"
        },
        {
            "city_id": "6",
            "id": "96",
            "name": "馆陶县"
        },
        {
            "city_id": "6",
            "id": "97",
            "name": "魏县"
        },
        {
            "city_id": "6",
            "id": "98",
            "name": "曲周县"
        },
        {
            "city_id": "6",
            "id": "99",
            "name": "武安市"
        },
        {
            "city_id": "7",
            "id": "100",
            "name": "桥东区"
        },
        {
            "city_id": "7",
            "id": "101",
            "name": "桥西区"
        },
        {
            "city_id": "7",
            "id": "102",
            "name": "邢台县"
        },
        {
            "city_id": "7",
            "id": "103",
            "name": "临城县"
        },
        {
            "city_id": "7",
            "id": "104",
            "name": "内丘县"
        },
        {
            "city_id": "7",
            "id": "105",
            "name": "柏乡县"
        },
        {
            "city_id": "7",
            "id": "106",
            "name": "隆尧县"
        },
        {
            "city_id": "7",
            "id": "107",
            "name": "任县"
        },
        {
            "city_id": "7",
            "id": "108",
            "name": "南和县"
        },
        {
            "city_id": "7",
            "id": "109",
            "name": "宁晋县"
        },
        {
            "city_id": "7",
            "id": "110",
            "name": "巨鹿县"
        },
        {
            "city_id": "7",
            "id": "111",
            "name": "新河县"
        },
        {
            "city_id": "7",
            "id": "112",
            "name": "广宗县"
        },
        {
            "city_id": "7",
            "id": "113",
            "name": "平乡县"
        },
        {
            "city_id": "7",
            "id": "114",
            "name": "威县"
        },
        {
            "city_id": "7",
            "id": "115",
            "name": "清河县"
        },
        {
            "city_id": "7",
            "id": "116",
            "name": "临西县"
        },
        {
            "city_id": "7",
            "id": "117",
            "name": "南宫市"
        },
        {
            "city_id": "7",
            "id": "118",
            "name": "沙河市"
        },
        {
            "city_id": "8",
            "id": "119",
            "name": "新市区"
        },
        {
            "city_id": "8",
            "id": "120",
            "name": "北市区"
        },
        {
            "city_id": "8",
            "id": "121",
            "name": "南市区"
        },
        {
            "city_id": "8",
            "id": "122",
            "name": "满城区"
        },
        {
            "city_id": "8",
            "id": "123",
            "name": "清苑区"
        },
        {
            "city_id": "8",
            "id": "124",
            "name": "涞水县"
        },
        {
            "city_id": "8",
            "id": "125",
            "name": "阜平县"
        },
        {
            "city_id": "8",
            "id": "126",
            "name": "徐水区"
        },
        {
            "city_id": "8",
            "id": "127",
            "name": "定兴县"
        },
        {
            "city_id": "8",
            "id": "128",
            "name": "唐县"
        },
        {
            "city_id": "8",
            "id": "129",
            "name": "高阳县"
        },
        {
            "city_id": "8",
            "id": "130",
            "name": "容城县"
        },
        {
            "city_id": "8",
            "id": "131",
            "name": "涞源县"
        },
        {
            "city_id": "8",
            "id": "132",
            "name": "望都县"
        },
        {
            "city_id": "8",
            "id": "133",
            "name": "安新县"
        },
        {
            "city_id": "8",
            "id": "134",
            "name": "易县"
        },
        {
            "city_id": "8",
            "id": "135",
            "name": "曲阳县"
        },
        {
            "city_id": "8",
            "id": "136",
            "name": "蠡县"
        },
        {
            "city_id": "8",
            "id": "137",
            "name": "顺平县"
        },
        {
            "city_id": "8",
            "id": "138",
            "name": "博野县"
        },
        {
            "city_id": "8",
            "id": "139",
            "name": "雄县"
        },
        {
            "city_id": "8",
            "id": "140",
            "name": "涿州市"
        },
        {
            "city_id": "8",
            "id": "141",
            "name": "定州市"
        },
        {
            "city_id": "8",
            "id": "142",
            "name": "安国市"
        },
        {
            "city_id": "8",
            "id": "143",
            "name": "高碑店市"
        },
        {
            "city_id": "9",
            "id": "144",
            "name": "桥东区"
        },
        {
            "city_id": "9",
            "id": "145",
            "name": "桥西区"
        },
        {
            "city_id": "9",
            "id": "146",
            "name": "宣化区"
        },
        {
            "city_id": "9",
            "id": "147",
            "name": "下花园区"
        },
        {
            "city_id": "9",
            "id": "148",
            "name": "宣化县"
        },
        {
            "city_id": "9",
            "id": "149",
            "name": "张北县"
        },
        {
            "city_id": "9",
            "id": "150",
            "name": "康保县"
        },
        {
            "city_id": "9",
            "id": "151",
            "name": "沽源县"
        },
        {
            "city_id": "9",
            "id": "152",
            "name": "尚义县"
        },
        {
            "city_id": "9",
            "id": "153",
            "name": "蔚县"
        },
        {
            "city_id": "9",
            "id": "154",
            "name": "阳原县"
        },
        {
            "city_id": "9",
            "id": "155",
            "name": "怀安县"
        },
        {
            "city_id": "9",
            "id": "156",
            "name": "万全区"
        },
        {
            "city_id": "9",
            "id": "157",
            "name": "怀来县"
        },
        {
            "city_id": "9",
            "id": "158",
            "name": "涿鹿县"
        },
        {
            "city_id": "9",
            "id": "159",
            "name": "赤城县"
        },
        {
            "city_id": "9",
            "id": "160",
            "name": "崇礼区"
        },
        {
            "city_id": "10",
            "id": "161",
            "name": "双桥区"
        },
        {
            "city_id": "10",
            "id": "162",
            "name": "双滦区"
        },
        {
            "city_id": "10",
            "id": "163",
            "name": "鹰手营子矿区"
        },
        {
            "city_id": "10",
            "id": "164",
            "name": "承德县"
        },
        {
            "city_id": "10",
            "id": "165",
            "name": "兴隆县"
        },
        {
            "city_id": "10",
            "id": "166",
            "name": "平泉市"
        },
        {
            "city_id": "10",
            "id": "167",
            "name": "滦平县"
        },
        {
            "city_id": "10",
            "id": "168",
            "name": "隆化县"
        },
        {
            "city_id": "10",
            "id": "169",
            "name": "丰宁满族自治县"
        },
        {
            "city_id": "10",
            "id": "170",
            "name": "宽城满族自治县"
        },
        {
            "city_id": "10",
            "id": "171",
            "name": "围场满族蒙古族自治县"
        },
        {
            "city_id": "11",
            "id": "172",
            "name": "新华区"
        },
        {
            "city_id": "11",
            "id": "173",
            "name": "运河区"
        },
        {
            "city_id": "11",
            "id": "174",
            "name": "沧县"
        },
        {
            "city_id": "11",
            "id": "175",
            "name": "青县"
        },
        {
            "city_id": "11",
            "id": "176",
            "name": "东光县"
        },
        {
            "city_id": "11",
            "id": "177",
            "name": "海兴县"
        },
        {
            "city_id": "11",
            "id": "178",
            "name": "盐山县"
        },
        {
            "city_id": "11",
            "id": "179",
            "name": "肃宁县"
        },
        {
            "city_id": "11",
            "id": "180",
            "name": "南皮县"
        },
        {
            "city_id": "11",
            "id": "181",
            "name": "吴桥县"
        },
        {
            "city_id": "11",
            "id": "182",
            "name": "献县"
        },
        {
            "city_id": "11",
            "id": "183",
            "name": "孟村回族自治县"
        },
        {
            "city_id": "11",
            "id": "184",
            "name": "泊头市"
        },
        {
            "city_id": "11",
            "id": "185",
            "name": "任丘市"
        },
        {
            "city_id": "11",
            "id": "186",
            "name": "黄骅市"
        },
        {
            "city_id": "11",
            "id": "187",
            "name": "河间市"
        },
        {
            "city_id": "12",
            "id": "188",
            "name": "安次区"
        },
        {
            "city_id": "12",
            "id": "189",
            "name": "广阳区"
        },
        {
            "city_id": "12",
            "id": "190",
            "name": "固安县"
        },
        {
            "city_id": "12",
            "id": "191",
            "name": "永清县"
        },
        {
            "city_id": "12",
            "id": "192",
            "name": "香河县"
        },
        {
            "city_id": "12",
            "id": "193",
            "name": "大城县"
        },
        {
            "city_id": "12",
            "id": "194",
            "name": "文安县"
        },
        {
            "city_id": "12",
            "id": "195",
            "name": "大厂回族自治县"
        },
        {
            "city_id": "12",
            "id": "196",
            "name": "霸州市"
        },
        {
            "city_id": "12",
            "id": "197",
            "name": "三河市"
        },
        {
            "city_id": "13",
            "id": "198",
            "name": "桃城区"
        },
        {
            "city_id": "13",
            "id": "199",
            "name": "枣强县"
        },
        {
            "city_id": "13",
            "id": "200",
            "name": "武邑县"
        },
        {
            "city_id": "13",
            "id": "201",
            "name": "武强县"
        },
        {
            "city_id": "13",
            "id": "202",
            "name": "饶阳县"
        },
        {
            "city_id": "13",
            "id": "203",
            "name": "安平县"
        },
        {
            "city_id": "13",
            "id": "204",
            "name": "故城县"
        },
        {
            "city_id": "13",
            "id": "205",
            "name": "景县"
        },
        {
            "city_id": "13",
            "id": "206",
            "name": "阜城县"
        },
        {
            "city_id": "13",
            "id": "207",
            "name": "冀州区"
        },
        {
            "city_id": "13",
            "id": "208",
            "name": "深州市"
        },
        {
            "city_id": "14",
            "id": "209",
            "name": "小店区"
        },
        {
            "city_id": "14",
            "id": "210",
            "name": "迎泽区"
        },
        {
            "city_id": "14",
            "id": "211",
            "name": "杏花岭区"
        },
        {
            "city_id": "14",
            "id": "212",
            "name": "尖草坪区"
        },
        {
            "city_id": "14",
            "id": "213",
            "name": "万柏林区"
        },
        {
            "city_id": "14",
            "id": "214",
            "name": "晋源区"
        },
        {
            "city_id": "14",
            "id": "215",
            "name": "清徐县"
        },
        {
            "city_id": "14",
            "id": "216",
            "name": "阳曲县"
        },
        {
            "city_id": "14",
            "id": "217",
            "name": "娄烦县"
        },
        {
            "city_id": "14",
            "id": "218",
            "name": "古交市"
        },
        {
            "city_id": "15",
            "id": "219",
            "name": "城区"
        },
        {
            "city_id": "15",
            "id": "220",
            "name": "矿区"
        },
        {
            "city_id": "15",
            "id": "221",
            "name": "南郊区"
        },
        {
            "city_id": "15",
            "id": "222",
            "name": "新荣区"
        },
        {
            "city_id": "15",
            "id": "223",
            "name": "阳高县"
        },
        {
            "city_id": "15",
            "id": "224",
            "name": "天镇县"
        },
        {
            "city_id": "15",
            "id": "225",
            "name": "广灵县"
        },
        {
            "city_id": "15",
            "id": "226",
            "name": "灵丘县"
        },
        {
            "city_id": "15",
            "id": "227",
            "name": "浑源县"
        },
        {
            "city_id": "15",
            "id": "228",
            "name": "左云县"
        },
        {
            "city_id": "15",
            "id": "229",
            "name": "大同县"
        },
        {
            "city_id": "16",
            "id": "230",
            "name": "城区"
        },
        {
            "city_id": "16",
            "id": "231",
            "name": "矿区"
        },
        {
            "city_id": "16",
            "id": "232",
            "name": "郊区"
        },
        {
            "city_id": "16",
            "id": "233",
            "name": "平定县"
        },
        {
            "city_id": "16",
            "id": "234",
            "name": "盂县"
        },
        {
            "city_id": "17",
            "id": "235",
            "name": "城区"
        },
        {
            "city_id": "17",
            "id": "236",
            "name": "郊区"
        },
        {
            "city_id": "17",
            "id": "237",
            "name": "长治县"
        },
        {
            "city_id": "17",
            "id": "238",
            "name": "襄垣县"
        },
        {
            "city_id": "17",
            "id": "239",
            "name": "屯留区"
        },
        {
            "city_id": "17",
            "id": "240",
            "name": "平顺县"
        },
        {
            "city_id": "17",
            "id": "241",
            "name": "黎城县"
        },
        {
            "city_id": "17",
            "id": "242",
            "name": "壶关县"
        },
        {
            "city_id": "17",
            "id": "243",
            "name": "长子县"
        },
        {
            "city_id": "17",
            "id": "244",
            "name": "武乡县"
        },
        {
            "city_id": "17",
            "id": "245",
            "name": "沁县"
        },
        {
            "city_id": "17",
            "id": "246",
            "name": "沁源县"
        },
        {
            "city_id": "17",
            "id": "247",
            "name": "潞城区"
        },
        {
            "city_id": "18",
            "id": "248",
            "name": "城区"
        },
        {
            "city_id": "18",
            "id": "249",
            "name": "沁水县"
        },
        {
            "city_id": "18",
            "id": "250",
            "name": "阳城县"
        },
        {
            "city_id": "18",
            "id": "251",
            "name": "陵川县"
        },
        {
            "city_id": "18",
            "id": "252",
            "name": "泽州县"
        },
        {
            "city_id": "18",
            "id": "253",
            "name": "高平市"
        },
        {
            "city_id": "19",
            "id": "254",
            "name": "朔城区"
        },
        {
            "city_id": "19",
            "id": "255",
            "name": "平鲁区"
        },
        {
            "city_id": "19",
            "id": "256",
            "name": "山阴县"
        },
        {
            "city_id": "19",
            "id": "257",
            "name": "应县"
        },
        {
            "city_id": "19",
            "id": "258",
            "name": "右玉县"
        },
        {
            "city_id": "19",
            "id": "259",
            "name": "怀仁市"
        },
        {
            "city_id": "20",
            "id": "260",
            "name": "榆次区"
        },
        {
            "city_id": "20",
            "id": "261",
            "name": "榆社县"
        },
        {
            "city_id": "20",
            "id": "262",
            "name": "左权县"
        },
        {
            "city_id": "20",
            "id": "263",
            "name": "和顺县"
        },
        {
            "city_id": "20",
            "id": "264",
            "name": "昔阳县"
        },
        {
            "city_id": "20",
            "id": "265",
            "name": "寿阳县"
        },
        {
            "city_id": "20",
            "id": "266",
            "name": "太谷县"
        },
        {
            "city_id": "20",
            "id": "267",
            "name": "祁县"
        },
        {
            "city_id": "20",
            "id": "268",
            "name": "平遥县"
        },
        {
            "city_id": "20",
            "id": "269",
            "name": "灵石县"
        },
        {
            "city_id": "20",
            "id": "270",
            "name": "介休市"
        },
        {
            "city_id": "21",
            "id": "271",
            "name": "盐湖区"
        },
        {
            "city_id": "21",
            "id": "272",
            "name": "临猗县"
        },
        {
            "city_id": "21",
            "id": "273",
            "name": "万荣县"
        },
        {
            "city_id": "21",
            "id": "274",
            "name": "闻喜县"
        },
        {
            "city_id": "21",
            "id": "275",
            "name": "稷山县"
        },
        {
            "city_id": "21",
            "id": "276",
            "name": "新绛县"
        },
        {
            "city_id": "21",
            "id": "277",
            "name": "绛县"
        },
        {
            "city_id": "21",
            "id": "278",
            "name": "垣曲县"
        },
        {
            "city_id": "21",
            "id": "279",
            "name": "夏县"
        },
        {
            "city_id": "21",
            "id": "280",
            "name": "平陆县"
        },
        {
            "city_id": "21",
            "id": "281",
            "name": "芮城县"
        },
        {
            "city_id": "21",
            "id": "282",
            "name": "永济市"
        },
        {
            "city_id": "21",
            "id": "283",
            "name": "河津市"
        },
        {
            "city_id": "22",
            "id": "284",
            "name": "忻府区"
        },
        {
            "city_id": "22",
            "id": "285",
            "name": "定襄县"
        },
        {
            "city_id": "22",
            "id": "286",
            "name": "五台县"
        },
        {
            "city_id": "22",
            "id": "287",
            "name": "代县"
        },
        {
            "city_id": "22",
            "id": "288",
            "name": "繁峙县"
        },
        {
            "city_id": "22",
            "id": "289",
            "name": "宁武县"
        },
        {
            "city_id": "22",
            "id": "290",
            "name": "静乐县"
        },
        {
            "city_id": "22",
            "id": "291",
            "name": "神池县"
        },
        {
            "city_id": "22",
            "id": "292",
            "name": "五寨县"
        },
        {
            "city_id": "22",
            "id": "293",
            "name": "岢岚县"
        },
        {
            "city_id": "22",
            "id": "294",
            "name": "河曲县"
        },
        {
            "city_id": "22",
            "id": "295",
            "name": "保德县"
        },
        {
            "city_id": "22",
            "id": "296",
            "name": "偏关县"
        },
        {
            "city_id": "22",
            "id": "297",
            "name": "原平市"
        },
        {
            "city_id": "23",
            "id": "298",
            "name": "尧都区"
        },
        {
            "city_id": "23",
            "id": "299",
            "name": "曲沃县"
        },
        {
            "city_id": "23",
            "id": "300",
            "name": "翼城县"
        },
        {
            "city_id": "23",
            "id": "301",
            "name": "襄汾县"
        },
        {
            "city_id": "23",
            "id": "302",
            "name": "洪洞县"
        },
        {
            "city_id": "23",
            "id": "303",
            "name": "古县"
        },
        {
            "city_id": "23",
            "id": "304",
            "name": "安泽县"
        },
        {
            "city_id": "23",
            "id": "305",
            "name": "浮山县"
        },
        {
            "city_id": "23",
            "id": "306",
            "name": "吉县"
        },
        {
            "city_id": "23",
            "id": "307",
            "name": "乡宁县"
        },
        {
            "city_id": "23",
            "id": "308",
            "name": "大宁县"
        },
        {
            "city_id": "23",
            "id": "309",
            "name": "隰县"
        },
        {
            "city_id": "23",
            "id": "310",
            "name": "永和县"
        },
        {
            "city_id": "23",
            "id": "311",
            "name": "蒲县"
        },
        {
            "city_id": "23",
            "id": "312",
            "name": "汾西县"
        },
        {
            "city_id": "23",
            "id": "313",
            "name": "侯马市"
        },
        {
            "city_id": "23",
            "id": "314",
            "name": "霍州市"
        },
        {
            "city_id": "24",
            "id": "315",
            "name": "离石区"
        },
        {
            "city_id": "24",
            "id": "316",
            "name": "文水县"
        },
        {
            "city_id": "24",
            "id": "317",
            "name": "交城县"
        },
        {
            "city_id": "24",
            "id": "318",
            "name": "兴县"
        },
        {
            "city_id": "24",
            "id": "319",
            "name": "临县"
        },
        {
            "city_id": "24",
            "id": "320",
            "name": "柳林县"
        },
        {
            "city_id": "24",
            "id": "321",
            "name": "石楼县"
        },
        {
            "city_id": "24",
            "id": "322",
            "name": "岚县"
        },
        {
            "city_id": "24",
            "id": "323",
            "name": "方山县"
        },
        {
            "city_id": "24",
            "id": "324",
            "name": "中阳县"
        },
        {
            "city_id": "24",
            "id": "325",
            "name": "交口县"
        },
        {
            "city_id": "24",
            "id": "326",
            "name": "孝义市"
        },
        {
            "city_id": "24",
            "id": "327",
            "name": "汾阳市"
        },
        {
            "city_id": "25",
            "id": "328",
            "name": "新城区"
        },
        {
            "city_id": "25",
            "id": "329",
            "name": "回民区"
        },
        {
            "city_id": "25",
            "id": "330",
            "name": "玉泉区"
        },
        {
            "city_id": "25",
            "id": "331",
            "name": "赛罕区"
        },
        {
            "city_id": "25",
            "id": "332",
            "name": "土默特左旗"
        },
        {
            "city_id": "25",
            "id": "333",
            "name": "托克托县"
        },
        {
            "city_id": "25",
            "id": "334",
            "name": "和林格尔县"
        },
        {
            "city_id": "25",
            "id": "335",
            "name": "清水河县"
        },
        {
            "city_id": "25",
            "id": "336",
            "name": "武川县"
        },
        {
            "city_id": "26",
            "id": "337",
            "name": "东河区"
        },
        {
            "city_id": "26",
            "id": "338",
            "name": "昆都仑区"
        },
        {
            "city_id": "26",
            "id": "339",
            "name": "青山区"
        },
        {
            "city_id": "26",
            "id": "340",
            "name": "石拐区"
        },
        {
            "city_id": "26",
            "id": "341",
            "name": "白云矿区"
        },
        {
            "city_id": "26",
            "id": "342",
            "name": "九原区"
        },
        {
            "city_id": "26",
            "id": "343",
            "name": "土默特右旗"
        },
        {
            "city_id": "26",
            "id": "344",
            "name": "固阳县"
        },
        {
            "city_id": "26",
            "id": "345",
            "name": "达尔罕茂明安联合旗"
        },
        {
            "city_id": "27",
            "id": "346",
            "name": "海勃湾区"
        },
        {
            "city_id": "27",
            "id": "347",
            "name": "海南区"
        },
        {
            "city_id": "27",
            "id": "348",
            "name": "乌达区"
        },
        {
            "city_id": "28",
            "id": "349",
            "name": "红山区"
        },
        {
            "city_id": "28",
            "id": "350",
            "name": "元宝山区"
        },
        {
            "city_id": "28",
            "id": "351",
            "name": "松山区"
        },
        {
            "city_id": "28",
            "id": "352",
            "name": "阿鲁科尔沁旗"
        },
        {
            "city_id": "28",
            "id": "353",
            "name": "巴林左旗"
        },
        {
            "city_id": "28",
            "id": "354",
            "name": "巴林右旗"
        },
        {
            "city_id": "28",
            "id": "355",
            "name": "林西县"
        },
        {
            "city_id": "28",
            "id": "356",
            "name": "克什克腾旗"
        },
        {
            "city_id": "28",
            "id": "357",
            "name": "翁牛特旗"
        },
        {
            "city_id": "28",
            "id": "358",
            "name": "喀喇沁旗"
        },
        {
            "city_id": "28",
            "id": "359",
            "name": "宁城县"
        },
        {
            "city_id": "28",
            "id": "360",
            "name": "敖汉旗"
        },
        {
            "city_id": "29",
            "id": "361",
            "name": "科尔沁区"
        },
        {
            "city_id": "29",
            "id": "362",
            "name": "科尔沁左翼中旗"
        },
        {
            "city_id": "29",
            "id": "363",
            "name": "科尔沁左翼后旗"
        },
        {
            "city_id": "29",
            "id": "364",
            "name": "开鲁县"
        },
        {
            "city_id": "29",
            "id": "365",
            "name": "库伦旗"
        },
        {
            "city_id": "29",
            "id": "366",
            "name": "奈曼旗"
        },
        {
            "city_id": "29",
            "id": "367",
            "name": "扎鲁特旗"
        },
        {
            "city_id": "29",
            "id": "368",
            "name": "霍林郭勒市"
        },
        {
            "city_id": "30",
            "id": "369",
            "name": "东胜区"
        },
        {
            "city_id": "30",
            "id": "370",
            "name": "达拉特旗"
        },
        {
            "city_id": "30",
            "id": "371",
            "name": "准格尔旗"
        },
        {
            "city_id": "30",
            "id": "372",
            "name": "鄂托克前旗"
        },
        {
            "city_id": "30",
            "id": "373",
            "name": "鄂托克旗"
        },
        {
            "city_id": "30",
            "id": "374",
            "name": "杭锦旗"
        },
        {
            "city_id": "30",
            "id": "375",
            "name": "乌审旗"
        },
        {
            "city_id": "30",
            "id": "376",
            "name": "伊金霍洛旗"
        },
        {
            "city_id": "31",
            "id": "377",
            "name": "海拉尔区"
        },
        {
            "city_id": "31",
            "id": "378",
            "name": "阿荣旗"
        },
        {
            "city_id": "31",
            "id": "379",
            "name": "莫力达瓦达斡尔族自治旗"
        },
        {
            "city_id": "31",
            "id": "380",
            "name": "鄂伦春自治旗"
        },
        {
            "city_id": "31",
            "id": "381",
            "name": "鄂温克族自治旗"
        },
        {
            "city_id": "31",
            "id": "382",
            "name": "陈巴尔虎旗"
        },
        {
            "city_id": "31",
            "id": "383",
            "name": "新巴尔虎左旗"
        },
        {
            "city_id": "31",
            "id": "384",
            "name": "新巴尔虎右旗"
        },
        {
            "city_id": "31",
            "id": "385",
            "name": "满洲里市"
        },
        {
            "city_id": "31",
            "id": "386",
            "name": "牙克石市"
        },
        {
            "city_id": "31",
            "id": "387",
            "name": "扎兰屯市"
        },
        {
            "city_id": "31",
            "id": "388",
            "name": "额尔古纳市"
        },
        {
            "city_id": "31",
            "id": "389",
            "name": "根河市"
        },
        {
            "city_id": "32",
            "id": "390",
            "name": "临河区"
        },
        {
            "city_id": "32",
            "id": "391",
            "name": "五原县"
        },
        {
            "city_id": "32",
            "id": "392",
            "name": "磴口县"
        },
        {
            "city_id": "32",
            "id": "393",
            "name": "乌拉特前旗"
        },
        {
            "city_id": "32",
            "id": "394",
            "name": "乌拉特中旗"
        },
        {
            "city_id": "32",
            "id": "395",
            "name": "乌拉特后旗"
        },
        {
            "city_id": "32",
            "id": "396",
            "name": "杭锦后旗"
        },
        {
            "city_id": "33",
            "id": "397",
            "name": "集宁区"
        },
        {
            "city_id": "33",
            "id": "398",
            "name": "卓资县"
        },
        {
            "city_id": "33",
            "id": "399",
            "name": "化德县"
        },
        {
            "city_id": "33",
            "id": "400",
            "name": "商都县"
        },
        {
            "city_id": "33",
            "id": "401",
            "name": "兴和县"
        },
        {
            "city_id": "33",
            "id": "402",
            "name": "凉城县"
        },
        {
            "city_id": "33",
            "id": "403",
            "name": "察哈尔右翼前旗"
        },
        {
            "city_id": "33",
            "id": "404",
            "name": "察哈尔右翼中旗"
        },
        {
            "city_id": "33",
            "id": "405",
            "name": "察哈尔右翼后旗"
        },
        {
            "city_id": "33",
            "id": "406",
            "name": "四子王旗"
        },
        {
            "city_id": "33",
            "id": "407",
            "name": "丰镇市"
        },
        {
            "city_id": "34",
            "id": "408",
            "name": "乌兰浩特市"
        },
        {
            "city_id": "34",
            "id": "409",
            "name": "阿尔山市"
        },
        {
            "city_id": "34",
            "id": "410",
            "name": "科尔沁右翼前旗"
        },
        {
            "city_id": "34",
            "id": "411",
            "name": "科尔沁右翼中旗"
        },
        {
            "city_id": "34",
            "id": "412",
            "name": "扎赉特旗"
        },
        {
            "city_id": "34",
            "id": "413",
            "name": "突泉县"
        },
        {
            "city_id": "35",
            "id": "414",
            "name": "二连浩特市"
        },
        {
            "city_id": "35",
            "id": "415",
            "name": "锡林浩特市"
        },
        {
            "city_id": "35",
            "id": "416",
            "name": "阿巴嘎旗"
        },
        {
            "city_id": "35",
            "id": "417",
            "name": "苏尼特左旗"
        },
        {
            "city_id": "35",
            "id": "418",
            "name": "苏尼特右旗"
        },
        {
            "city_id": "35",
            "id": "419",
            "name": "东乌珠穆沁旗"
        },
        {
            "city_id": "35",
            "id": "420",
            "name": "西乌珠穆沁旗"
        },
        {
            "city_id": "35",
            "id": "421",
            "name": "太仆寺旗"
        },
        {
            "city_id": "35",
            "id": "422",
            "name": "镶黄旗"
        },
        {
            "city_id": "35",
            "id": "423",
            "name": "正镶白旗"
        },
        {
            "city_id": "35",
            "id": "424",
            "name": "正蓝旗"
        },
        {
            "city_id": "35",
            "id": "425",
            "name": "多伦县"
        },
        {
            "city_id": "36",
            "id": "426",
            "name": "阿拉善左旗"
        },
        {
            "city_id": "36",
            "id": "427",
            "name": "阿拉善右旗"
        },
        {
            "city_id": "36",
            "id": "428",
            "name": "额济纳旗"
        },
        {
            "city_id": "37",
            "id": "429",
            "name": "和平区"
        },
        {
            "city_id": "37",
            "id": "430",
            "name": "沈河区"
        },
        {
            "city_id": "37",
            "id": "431",
            "name": "大东区"
        },
        {
            "city_id": "37",
            "id": "432",
            "name": "皇姑区"
        },
        {
            "city_id": "37",
            "id": "433",
            "name": "铁西区"
        },
        {
            "city_id": "37",
            "id": "434",
            "name": "苏家屯区"
        },
        {
            "city_id": "37",
            "id": "435",
            "name": "东陵区"
        },
        {
            "city_id": "37",
            "id": "436",
            "name": "新城子区"
        },
        {
            "city_id": "37",
            "id": "437",
            "name": "于洪区"
        },
        {
            "city_id": "37",
            "id": "438",
            "name": "辽中区"
        },
        {
            "city_id": "37",
            "id": "439",
            "name": "康平县"
        },
        {
            "city_id": "37",
            "id": "440",
            "name": "法库县"
        },
        {
            "city_id": "37",
            "id": "441",
            "name": "新民市"
        },
        {
            "city_id": "38",
            "id": "442",
            "name": "中山区"
        },
        {
            "city_id": "38",
            "id": "443",
            "name": "西岗区"
        },
        {
            "city_id": "38",
            "id": "444",
            "name": "沙河口区"
        },
        {
            "city_id": "38",
            "id": "445",
            "name": "甘井子区"
        },
        {
            "city_id": "38",
            "id": "446",
            "name": "旅顺口区"
        },
        {
            "city_id": "38",
            "id": "447",
            "name": "金州区"
        },
        {
            "city_id": "38",
            "id": "448",
            "name": "长海县"
        },
        {
            "city_id": "38",
            "id": "449",
            "name": "瓦房店市"
        },
        {
            "city_id": "38",
            "id": "450",
            "name": "普兰店区"
        },
        {
            "city_id": "38",
            "id": "451",
            "name": "庄河市"
        },
        {
            "city_id": "39",
            "id": "452",
            "name": "铁东区"
        },
        {
            "city_id": "39",
            "id": "453",
            "name": "铁西区"
        },
        {
            "city_id": "39",
            "id": "454",
            "name": "立山区"
        },
        {
            "city_id": "39",
            "id": "455",
            "name": "千山区"
        },
        {
            "city_id": "39",
            "id": "456",
            "name": "台安县"
        },
        {
            "city_id": "39",
            "id": "457",
            "name": "岫岩满族自治县"
        },
        {
            "city_id": "39",
            "id": "458",
            "name": "海城市"
        },
        {
            "city_id": "40",
            "id": "459",
            "name": "新抚区"
        },
        {
            "city_id": "40",
            "id": "460",
            "name": "东洲区"
        },
        {
            "city_id": "40",
            "id": "461",
            "name": "望花区"
        },
        {
            "city_id": "40",
            "id": "462",
            "name": "顺城区"
        },
        {
            "city_id": "40",
            "id": "463",
            "name": "抚顺县"
        },
        {
            "city_id": "40",
            "id": "464",
            "name": "新宾满族自治县"
        },
        {
            "city_id": "40",
            "id": "465",
            "name": "清原满族自治县"
        },
        {
            "city_id": "41",
            "id": "466",
            "name": "平山区"
        },
        {
            "city_id": "41",
            "id": "467",
            "name": "溪湖区"
        },
        {
            "city_id": "41",
            "id": "468",
            "name": "明山区"
        },
        {
            "city_id": "41",
            "id": "469",
            "name": "南芬区"
        },
        {
            "city_id": "41",
            "id": "470",
            "name": "本溪满族自治县"
        },
        {
            "city_id": "41",
            "id": "471",
            "name": "桓仁满族自治县"
        },
        {
            "city_id": "42",
            "id": "472",
            "name": "元宝区"
        },
        {
            "city_id": "42",
            "id": "473",
            "name": "振兴区"
        },
        {
            "city_id": "42",
            "id": "474",
            "name": "振安区"
        },
        {
            "city_id": "42",
            "id": "475",
            "name": "宽甸满族自治县"
        },
        {
            "city_id": "42",
            "id": "476",
            "name": "东港市"
        },
        {
            "city_id": "42",
            "id": "477",
            "name": "凤城市"
        },
        {
            "city_id": "43",
            "id": "478",
            "name": "古塔区"
        },
        {
            "city_id": "43",
            "id": "479",
            "name": "凌河区"
        },
        {
            "city_id": "43",
            "id": "480",
            "name": "太和区"
        },
        {
            "city_id": "43",
            "id": "481",
            "name": "黑山县"
        },
        {
            "city_id": "43",
            "id": "482",
            "name": "义县"
        },
        {
            "city_id": "43",
            "id": "483",
            "name": "凌海市"
        },
        {
            "city_id": "43",
            "id": "484",
            "name": "北宁市"
        },
        {
            "city_id": "44",
            "id": "485",
            "name": "站前区"
        },
        {
            "city_id": "44",
            "id": "486",
            "name": "西市区"
        },
        {
            "city_id": "44",
            "id": "487",
            "name": "鲅鱼圈区"
        },
        {
            "city_id": "44",
            "id": "488",
            "name": "老边区"
        },
        {
            "city_id": "44",
            "id": "489",
            "name": "盖州市"
        },
        {
            "city_id": "44",
            "id": "490",
            "name": "大石桥市"
        },
        {
            "city_id": "45",
            "id": "491",
            "name": "海州区"
        },
        {
            "city_id": "45",
            "id": "492",
            "name": "新邱区"
        },
        {
            "city_id": "45",
            "id": "493",
            "name": "太平区"
        },
        {
            "city_id": "45",
            "id": "494",
            "name": "清河门区"
        },
        {
            "city_id": "45",
            "id": "495",
            "name": "细河区"
        },
        {
            "city_id": "45",
            "id": "496",
            "name": "阜新蒙古族自治县"
        },
        {
            "city_id": "45",
            "id": "497",
            "name": "彰武县"
        },
        {
            "city_id": "46",
            "id": "498",
            "name": "白塔区"
        },
        {
            "city_id": "46",
            "id": "499",
            "name": "文圣区"
        },
        {
            "city_id": "46",
            "id": "500",
            "name": "宏伟区"
        },
        {
            "city_id": "46",
            "id": "501",
            "name": "弓长岭区"
        },
        {
            "city_id": "46",
            "id": "502",
            "name": "太子河区"
        },
        {
            "city_id": "46",
            "id": "503",
            "name": "辽阳县"
        },
        {
            "city_id": "46",
            "id": "504",
            "name": "灯塔市"
        },
        {
            "city_id": "47",
            "id": "505",
            "name": "双台子区"
        },
        {
            "city_id": "47",
            "id": "506",
            "name": "兴隆台区"
        },
        {
            "city_id": "47",
            "id": "507",
            "name": "大洼区"
        },
        {
            "city_id": "47",
            "id": "508",
            "name": "盘山县"
        },
        {
            "city_id": "48",
            "id": "509",
            "name": "银州区"
        },
        {
            "city_id": "48",
            "id": "510",
            "name": "清河区"
        },
        {
            "city_id": "48",
            "id": "511",
            "name": "铁岭县"
        },
        {
            "city_id": "48",
            "id": "512",
            "name": "西丰县"
        },
        {
            "city_id": "48",
            "id": "513",
            "name": "昌图县"
        },
        {
            "city_id": "48",
            "id": "514",
            "name": "调兵山市"
        },
        {
            "city_id": "48",
            "id": "515",
            "name": "开原市"
        },
        {
            "city_id": "49",
            "id": "516",
            "name": "双塔区"
        },
        {
            "city_id": "49",
            "id": "517",
            "name": "龙城区"
        },
        {
            "city_id": "49",
            "id": "518",
            "name": "朝阳县"
        },
        {
            "city_id": "49",
            "id": "519",
            "name": "建平县"
        },
        {
            "city_id": "49",
            "id": "520",
            "name": "喀喇沁左翼蒙古族自治县"
        },
        {
            "city_id": "49",
            "id": "521",
            "name": "北票市"
        },
        {
            "city_id": "49",
            "id": "522",
            "name": "凌源市"
        },
        {
            "city_id": "50",
            "id": "523",
            "name": "连山区"
        },
        {
            "city_id": "50",
            "id": "524",
            "name": "龙港区"
        },
        {
            "city_id": "50",
            "id": "525",
            "name": "南票区"
        },
        {
            "city_id": "50",
            "id": "526",
            "name": "绥中县"
        },
        {
            "city_id": "50",
            "id": "527",
            "name": "建昌县"
        },
        {
            "city_id": "50",
            "id": "528",
            "name": "兴城市"
        },
        {
            "city_id": "51",
            "id": "529",
            "name": "南关区"
        },
        {
            "city_id": "51",
            "id": "530",
            "name": "宽城区"
        },
        {
            "city_id": "51",
            "id": "531",
            "name": "朝阳区"
        },
        {
            "city_id": "51",
            "id": "532",
            "name": "二道区"
        },
        {
            "city_id": "51",
            "id": "533",
            "name": "绿园区"
        },
        {
            "city_id": "51",
            "id": "534",
            "name": "双阳区"
        },
        {
            "city_id": "51",
            "id": "535",
            "name": "农安县"
        },
        {
            "city_id": "51",
            "id": "536",
            "name": "九台区"
        },
        {
            "city_id": "51",
            "id": "537",
            "name": "榆树市"
        },
        {
            "city_id": "51",
            "id": "538",
            "name": "德惠市"
        },
        {
            "city_id": "52",
            "id": "539",
            "name": "昌邑区"
        },
        {
            "city_id": "52",
            "id": "540",
            "name": "龙潭区"
        },
        {
            "city_id": "52",
            "id": "541",
            "name": "船营区"
        },
        {
            "city_id": "52",
            "id": "542",
            "name": "丰满区"
        },
        {
            "city_id": "52",
            "id": "543",
            "name": "永吉县"
        },
        {
            "city_id": "52",
            "id": "544",
            "name": "蛟河市"
        },
        {
            "city_id": "52",
            "id": "545",
            "name": "桦甸市"
        },
        {
            "city_id": "52",
            "id": "546",
            "name": "舒兰市"
        },
        {
            "city_id": "52",
            "id": "547",
            "name": "磐石市"
        },
        {
            "city_id": "53",
            "id": "548",
            "name": "铁西区"
        },
        {
            "city_id": "53",
            "id": "549",
            "name": "铁东区"
        },
        {
            "city_id": "53",
            "id": "550",
            "name": "梨树县"
        },
        {
            "city_id": "53",
            "id": "551",
            "name": "伊通满族自治县"
        },
        {
            "city_id": "53",
            "id": "552",
            "name": "公主岭市"
        },
        {
            "city_id": "53",
            "id": "553",
            "name": "双辽市"
        },
        {
            "city_id": "54",
            "id": "554",
            "name": "龙山区"
        },
        {
            "city_id": "54",
            "id": "555",
            "name": "西安区"
        },
        {
            "city_id": "54",
            "id": "556",
            "name": "东丰县"
        },
        {
            "city_id": "54",
            "id": "557",
            "name": "东辽县"
        },
        {
            "city_id": "55",
            "id": "558",
            "name": "东昌区"
        },
        {
            "city_id": "55",
            "id": "559",
            "name": "二道江区"
        },
        {
            "city_id": "55",
            "id": "560",
            "name": "通化县"
        },
        {
            "city_id": "55",
            "id": "561",
            "name": "辉南县"
        },
        {
            "city_id": "55",
            "id": "562",
            "name": "柳河县"
        },
        {
            "city_id": "55",
            "id": "563",
            "name": "梅河口市"
        },
        {
            "city_id": "55",
            "id": "564",
            "name": "集安市"
        },
        {
            "city_id": "56",
            "id": "565",
            "name": "八道江区"
        },
        {
            "city_id": "56",
            "id": "566",
            "name": "抚松县"
        },
        {
            "city_id": "56",
            "id": "567",
            "name": "靖宇县"
        },
        {
            "city_id": "56",
            "id": "568",
            "name": "长白朝鲜族自治县"
        },
        {
            "city_id": "56",
            "id": "569",
            "name": "江源区"
        },
        {
            "city_id": "56",
            "id": "570",
            "name": "临江市"
        },
        {
            "city_id": "57",
            "id": "571",
            "name": "宁江区"
        },
        {
            "city_id": "57",
            "id": "572",
            "name": "前郭尔罗斯蒙古族自治县"
        },
        {
            "city_id": "57",
            "id": "573",
            "name": "长岭县"
        },
        {
            "city_id": "57",
            "id": "574",
            "name": "乾安县"
        },
        {
            "city_id": "57",
            "id": "575",
            "name": "扶余市"
        },
        {
            "city_id": "58",
            "id": "576",
            "name": "洮北区"
        },
        {
            "city_id": "58",
            "id": "577",
            "name": "镇赉县"
        },
        {
            "city_id": "58",
            "id": "578",
            "name": "通榆县"
        },
        {
            "city_id": "58",
            "id": "579",
            "name": "洮南市"
        },
        {
            "city_id": "58",
            "id": "580",
            "name": "大安市"
        },
        {
            "city_id": "59",
            "id": "581",
            "name": "延吉市"
        },
        {
            "city_id": "59",
            "id": "582",
            "name": "图们市"
        },
        {
            "city_id": "59",
            "id": "583",
            "name": "敦化市"
        },
        {
            "city_id": "59",
            "id": "584",
            "name": "珲春市"
        },
        {
            "city_id": "59",
            "id": "585",
            "name": "龙井市"
        },
        {
            "city_id": "59",
            "id": "586",
            "name": "和龙市"
        },
        {
            "city_id": "59",
            "id": "587",
            "name": "汪清县"
        },
        {
            "city_id": "59",
            "id": "588",
            "name": "安图县"
        },
        {
            "city_id": "60",
            "id": "589",
            "name": "道里区"
        },
        {
            "city_id": "60",
            "id": "590",
            "name": "南岗区"
        },
        {
            "city_id": "60",
            "id": "591",
            "name": "道外区"
        },
        {
            "city_id": "60",
            "id": "592",
            "name": "香坊区"
        },
        {
            "city_id": "60",
            "id": "593",
            "name": "动力区"
        },
        {
            "city_id": "60",
            "id": "594",
            "name": "平房区"
        },
        {
            "city_id": "60",
            "id": "595",
            "name": "松北区"
        },
        {
            "city_id": "60",
            "id": "596",
            "name": "呼兰区"
        },
        {
            "city_id": "60",
            "id": "597",
            "name": "依兰县"
        },
        {
            "city_id": "60",
            "id": "598",
            "name": "方正县"
        },
        {
            "city_id": "60",
            "id": "599",
            "name": "宾县"
        },
        {
            "city_id": "60",
            "id": "600",
            "name": "巴彦县"
        },
        {
            "city_id": "60",
            "id": "601",
            "name": "木兰县"
        },
        {
            "city_id": "60",
            "id": "602",
            "name": "通河县"
        },
        {
            "city_id": "60",
            "id": "603",
            "name": "延寿县"
        },
        {
            "city_id": "60",
            "id": "604",
            "name": "阿城区"
        },
        {
            "city_id": "60",
            "id": "605",
            "name": "双城区"
        },
        {
            "city_id": "60",
            "id": "606",
            "name": "尚志市"
        },
        {
            "city_id": "60",
            "id": "607",
            "name": "五常市"
        },
        {
            "city_id": "61",
            "id": "608",
            "name": "龙沙区"
        },
        {
            "city_id": "61",
            "id": "609",
            "name": "建华区"
        },
        {
            "city_id": "61",
            "id": "610",
            "name": "铁锋区"
        },
        {
            "city_id": "61",
            "id": "611",
            "name": "昂昂溪区"
        },
        {
            "city_id": "61",
            "id": "612",
            "name": "富拉尔基区"
        },
        {
            "city_id": "61",
            "id": "613",
            "name": "碾子山区"
        },
        {
            "city_id": "61",
            "id": "614",
            "name": "梅里斯达斡尔族区"
        },
        {
            "city_id": "61",
            "id": "615",
            "name": "龙江县"
        },
        {
            "city_id": "61",
            "id": "616",
            "name": "依安县"
        },
        {
            "city_id": "61",
            "id": "617",
            "name": "泰来县"
        },
        {
            "city_id": "61",
            "id": "618",
            "name": "甘南县"
        },
        {
            "city_id": "61",
            "id": "619",
            "name": "富裕县"
        },
        {
            "city_id": "61",
            "id": "620",
            "name": "克山县"
        },
        {
            "city_id": "61",
            "id": "621",
            "name": "克东县"
        },
        {
            "city_id": "61",
            "id": "622",
            "name": "拜泉县"
        },
        {
            "city_id": "61",
            "id": "623",
            "name": "讷河市"
        },
        {
            "city_id": "62",
            "id": "624",
            "name": "鸡冠区"
        },
        {
            "city_id": "62",
            "id": "625",
            "name": "恒山区"
        },
        {
            "city_id": "62",
            "id": "626",
            "name": "滴道区"
        },
        {
            "city_id": "62",
            "id": "627",
            "name": "梨树区"
        },
        {
            "city_id": "62",
            "id": "628",
            "name": "城子河区"
        },
        {
            "city_id": "62",
            "id": "629",
            "name": "麻山区"
        },
        {
            "city_id": "62",
            "id": "630",
            "name": "鸡东县"
        },
        {
            "city_id": "62",
            "id": "631",
            "name": "虎林市"
        },
        {
            "city_id": "62",
            "id": "632",
            "name": "密山市"
        },
        {
            "city_id": "63",
            "id": "633",
            "name": "向阳区"
        },
        {
            "city_id": "63",
            "id": "634",
            "name": "工农区"
        },
        {
            "city_id": "63",
            "id": "635",
            "name": "南山区"
        },
        {
            "city_id": "63",
            "id": "636",
            "name": "兴安区"
        },
        {
            "city_id": "63",
            "id": "637",
            "name": "东山区"
        },
        {
            "city_id": "63",
            "id": "638",
            "name": "兴山区"
        },
        {
            "city_id": "63",
            "id": "639",
            "name": "萝北县"
        },
        {
            "city_id": "63",
            "id": "640",
            "name": "绥滨县"
        },
        {
            "city_id": "64",
            "id": "641",
            "name": "尖山区"
        },
        {
            "city_id": "64",
            "id": "642",
            "name": "岭东区"
        },
        {
            "city_id": "64",
            "id": "643",
            "name": "四方台区"
        },
        {
            "city_id": "64",
            "id": "644",
            "name": "宝山区"
        },
        {
            "city_id": "64",
            "id": "645",
            "name": "集贤县"
        },
        {
            "city_id": "64",
            "id": "646",
            "name": "友谊县"
        },
        {
            "city_id": "64",
            "id": "647",
            "name": "宝清县"
        },
        {
            "city_id": "64",
            "id": "648",
            "name": "饶河县"
        },
        {
            "city_id": "65",
            "id": "649",
            "name": "萨尔图区"
        },
        {
            "city_id": "65",
            "id": "650",
            "name": "龙凤区"
        },
        {
            "city_id": "65",
            "id": "651",
            "name": "让胡路区"
        },
        {
            "city_id": "65",
            "id": "652",
            "name": "红岗区"
        },
        {
            "city_id": "65",
            "id": "653",
            "name": "大同区"
        },
        {
            "city_id": "65",
            "id": "654",
            "name": "肇州县"
        },
        {
            "city_id": "65",
            "id": "655",
            "name": "肇源县"
        },
        {
            "city_id": "65",
            "id": "656",
            "name": "林甸县"
        },
        {
            "city_id": "65",
            "id": "657",
            "name": "杜尔伯特蒙古族自治县"
        },
        {
            "city_id": "66",
            "id": "658",
            "name": "伊春区"
        },
        {
            "city_id": "66",
            "id": "659",
            "name": "南岔区"
        },
        {
            "city_id": "66",
            "id": "660",
            "name": "友好区"
        },
        {
            "city_id": "66",
            "id": "661",
            "name": "西林区"
        },
        {
            "city_id": "66",
            "id": "662",
            "name": "翠峦区"
        },
        {
            "city_id": "66",
            "id": "663",
            "name": "新青区"
        },
        {
            "city_id": "66",
            "id": "664",
            "name": "美溪区"
        },
        {
            "city_id": "66",
            "id": "665",
            "name": "金山屯区"
        },
        {
            "city_id": "66",
            "id": "666",
            "name": "五营区"
        },
        {
            "city_id": "66",
            "id": "667",
            "name": "乌马河区"
        },
        {
            "city_id": "66",
            "id": "668",
            "name": "汤旺河区"
        },
        {
            "city_id": "66",
            "id": "669",
            "name": "带岭区"
        },
        {
            "city_id": "66",
            "id": "670",
            "name": "乌伊岭区"
        },
        {
            "city_id": "66",
            "id": "671",
            "name": "红星区"
        },
        {
            "city_id": "66",
            "id": "672",
            "name": "上甘岭区"
        },
        {
            "city_id": "66",
            "id": "673",
            "name": "嘉荫县"
        },
        {
            "city_id": "66",
            "id": "674",
            "name": "铁力市"
        },
        {
            "city_id": "67",
            "id": "675",
            "name": "永红区"
        },
        {
            "city_id": "67",
            "id": "676",
            "name": "向阳区"
        },
        {
            "city_id": "67",
            "id": "677",
            "name": "前进区"
        },
        {
            "city_id": "67",
            "id": "678",
            "name": "东风区"
        },
        {
            "city_id": "67",
            "id": "679",
            "name": "郊区"
        },
        {
            "city_id": "67",
            "id": "680",
            "name": "桦南县"
        },
        {
            "city_id": "67",
            "id": "681",
            "name": "桦川县"
        },
        {
            "city_id": "67",
            "id": "682",
            "name": "汤原县"
        },
        {
            "city_id": "67",
            "id": "683",
            "name": "抚远市"
        },
        {
            "city_id": "67",
            "id": "684",
            "name": "同江市"
        },
        {
            "city_id": "67",
            "id": "685",
            "name": "富锦市"
        },
        {
            "city_id": "68",
            "id": "686",
            "name": "新兴区"
        },
        {
            "city_id": "68",
            "id": "687",
            "name": "桃山区"
        },
        {
            "city_id": "68",
            "id": "688",
            "name": "茄子河区"
        },
        {
            "city_id": "68",
            "id": "689",
            "name": "勃利县"
        },
        {
            "city_id": "69",
            "id": "690",
            "name": "东安区"
        },
        {
            "city_id": "69",
            "id": "691",
            "name": "阳明区"
        },
        {
            "city_id": "69",
            "id": "692",
            "name": "爱民区"
        },
        {
            "city_id": "69",
            "id": "693",
            "name": "西安区"
        },
        {
            "city_id": "69",
            "id": "694",
            "name": "东宁市"
        },
        {
            "city_id": "69",
            "id": "695",
            "name": "林口县"
        },
        {
            "city_id": "69",
            "id": "696",
            "name": "绥芬河市"
        },
        {
            "city_id": "69",
            "id": "697",
            "name": "海林市"
        },
        {
            "city_id": "69",
            "id": "698",
            "name": "宁安市"
        },
        {
            "city_id": "69",
            "id": "699",
            "name": "穆棱市"
        },
        {
            "city_id": "70",
            "id": "700",
            "name": "爱辉区"
        },
        {
            "city_id": "70",
            "id": "701",
            "name": "嫩江县"
        },
        {
            "city_id": "70",
            "id": "702",
            "name": "逊克县"
        },
        {
            "city_id": "70",
            "id": "703",
            "name": "孙吴县"
        },
        {
            "city_id": "70",
            "id": "704",
            "name": "北安市"
        },
        {
            "city_id": "70",
            "id": "705",
            "name": "五大连池市"
        },
        {
            "city_id": "71",
            "id": "706",
            "name": "北林区"
        },
        {
            "city_id": "71",
            "id": "707",
            "name": "望奎县"
        },
        {
            "city_id": "71",
            "id": "708",
            "name": "兰西县"
        },
        {
            "city_id": "71",
            "id": "709",
            "name": "青冈县"
        },
        {
            "city_id": "71",
            "id": "710",
            "name": "庆安县"
        },
        {
            "city_id": "71",
            "id": "711",
            "name": "明水县"
        },
        {
            "city_id": "71",
            "id": "712",
            "name": "绥棱县"
        },
        {
            "city_id": "71",
            "id": "713",
            "name": "安达市"
        },
        {
            "city_id": "71",
            "id": "714",
            "name": "肇东市"
        },
        {
            "city_id": "71",
            "id": "715",
            "name": "海伦市"
        },
        {
            "city_id": "72",
            "id": "716",
            "name": "呼玛县"
        },
        {
            "city_id": "72",
            "id": "717",
            "name": "塔河县"
        },
        {
            "city_id": "72",
            "id": "718",
            "name": "漠河市"
        },
        {
            "city_id": "73",
            "id": "719",
            "name": "黄浦区"
        },
        {
            "city_id": "73",
            "id": "720",
            "name": "卢湾区"
        },
        {
            "city_id": "73",
            "id": "721",
            "name": "徐汇区"
        },
        {
            "city_id": "73",
            "id": "722",
            "name": "长宁区"
        },
        {
            "city_id": "73",
            "id": "723",
            "name": "静安区"
        },
        {
            "city_id": "73",
            "id": "724",
            "name": "普陀区"
        },
        {
            "city_id": "73",
            "id": "725",
            "name": "闸北区"
        },
        {
            "city_id": "73",
            "id": "726",
            "name": "虹口区"
        },
        {
            "city_id": "73",
            "id": "727",
            "name": "杨浦区"
        },
        {
            "city_id": "73",
            "id": "728",
            "name": "闵行区"
        },
        {
            "city_id": "73",
            "id": "729",
            "name": "宝山区"
        },
        {
            "city_id": "73",
            "id": "730",
            "name": "嘉定区"
        },
        {
            "city_id": "73",
            "id": "731",
            "name": "浦东新区"
        },
        {
            "city_id": "73",
            "id": "732",
            "name": "金山区"
        },
        {
            "city_id": "73",
            "id": "733",
            "name": "松江区"
        },
        {
            "city_id": "73",
            "id": "734",
            "name": "青浦区"
        },
        {
            "city_id": "73",
            "id": "735",
            "name": "南汇区"
        },
        {
            "city_id": "73",
            "id": "736",
            "name": "奉贤区"
        },
        {
            "city_id": "73",
            "id": "737",
            "name": "崇明区"
        },
        {
            "city_id": "74",
            "id": "738",
            "name": "玄武区"
        },
        {
            "city_id": "74",
            "id": "740",
            "name": "秦淮区"
        },
        {
            "city_id": "74",
            "id": "741",
            "name": "建邺区"
        },
        {
            "city_id": "74",
            "id": "742",
            "name": "鼓楼区"
        },
        {
            "city_id": "74",
            "id": "743",
            "name": "下关区"
        },
        {
            "city_id": "74",
            "id": "744",
            "name": "浦口区"
        },
        {
            "city_id": "74",
            "id": "745",
            "name": "栖霞区"
        },
        {
            "city_id": "74",
            "id": "746",
            "name": "雨花台区"
        },
        {
            "city_id": "74",
            "id": "747",
            "name": "江宁区"
        },
        {
            "city_id": "74",
            "id": "748",
            "name": "六合区"
        },
        {
            "city_id": "74",
            "id": "749",
            "name": "溧水区"
        },
        {
            "city_id": "74",
            "id": "750",
            "name": "高淳区"
        },
        {
            "city_id": "75",
            "id": "751",
            "name": "崇安区"
        },
        {
            "city_id": "75",
            "id": "752",
            "name": "南长区"
        },
        {
            "city_id": "75",
            "id": "753",
            "name": "北塘区"
        },
        {
            "city_id": "75",
            "id": "754",
            "name": "锡山区"
        },
        {
            "city_id": "75",
            "id": "755",
            "name": "惠山区"
        },
        {
            "city_id": "75",
            "id": "756",
            "name": "滨湖区"
        },
        {
            "city_id": "75",
            "id": "757",
            "name": "江阴市"
        },
        {
            "city_id": "75",
            "id": "758",
            "name": "宜兴市"
        },
        {
            "city_id": "76",
            "id": "759",
            "name": "鼓楼区"
        },
        {
            "city_id": "76",
            "id": "760",
            "name": "云龙区"
        },
        {
            "city_id": "76",
            "id": "761",
            "name": "九里区"
        },
        {
            "city_id": "76",
            "id": "762",
            "name": "贾汪区"
        },
        {
            "city_id": "76",
            "id": "763",
            "name": "泉山区"
        },
        {
            "city_id": "76",
            "id": "764",
            "name": "丰县"
        },
        {
            "city_id": "76",
            "id": "765",
            "name": "沛县"
        },
        {
            "city_id": "76",
            "id": "766",
            "name": "铜山区"
        },
        {
            "city_id": "76",
            "id": "767",
            "name": "睢宁县"
        },
        {
            "city_id": "76",
            "id": "768",
            "name": "新沂市"
        },
        {
            "city_id": "76",
            "id": "769",
            "name": "邳州市"
        },
        {
            "city_id": "77",
            "id": "770",
            "name": "天宁区"
        },
        {
            "city_id": "77",
            "id": "771",
            "name": "钟楼区"
        },
        {
            "city_id": "77",
            "id": "772",
            "name": "戚墅堰区"
        },
        {
            "city_id": "77",
            "id": "773",
            "name": "新北区"
        },
        {
            "city_id": "77",
            "id": "774",
            "name": "武进区"
        },
        {
            "city_id": "77",
            "id": "775",
            "name": "溧阳市"
        },
        {
            "city_id": "77",
            "id": "776",
            "name": "金坛区"
        },
        {
            "city_id": "78",
            "id": "777",
            "name": "沧浪区"
        },
        {
            "city_id": "78",
            "id": "778",
            "name": "平江区"
        },
        {
            "city_id": "78",
            "id": "779",
            "name": "金阊区"
        },
        {
            "city_id": "78",
            "id": "780",
            "name": "虎丘区"
        },
        {
            "city_id": "78",
            "id": "781",
            "name": "吴中区"
        },
        {
            "city_id": "78",
            "id": "782",
            "name": "相城区"
        },
        {
            "city_id": "78",
            "id": "783",
            "name": "常熟市"
        },
        {
            "city_id": "78",
            "id": "784",
            "name": "张家港市"
        },
        {
            "city_id": "78",
            "id": "785",
            "name": "昆山市"
        },
        {
            "city_id": "78",
            "id": "786",
            "name": "吴江区"
        },
        {
            "city_id": "78",
            "id": "787",
            "name": "太仓市"
        },
        {
            "city_id": "79",
            "id": "788",
            "name": "崇川区"
        },
        {
            "city_id": "79",
            "id": "789",
            "name": "港闸区"
        },
        {
            "city_id": "79",
            "id": "790",
            "name": "海安市"
        },
        {
            "city_id": "79",
            "id": "791",
            "name": "如东县"
        },
        {
            "city_id": "79",
            "id": "792",
            "name": "启东市"
        },
        {
            "city_id": "79",
            "id": "793",
            "name": "如皋市"
        },
        {
            "city_id": "79",
            "id": "794",
            "name": "通州区"
        },
        {
            "city_id": "79",
            "id": "795",
            "name": "海门市"
        },
        {
            "city_id": "80",
            "id": "796",
            "name": "连云区"
        },
        {
            "city_id": "80",
            "id": "797",
            "name": "新浦区"
        },
        {
            "city_id": "80",
            "id": "798",
            "name": "海州区"
        },
        {
            "city_id": "80",
            "id": "799",
            "name": "赣榆区"
        },
        {
            "city_id": "80",
            "id": "800",
            "name": "东海县"
        },
        {
            "city_id": "80",
            "id": "801",
            "name": "灌云县"
        },
        {
            "city_id": "80",
            "id": "802",
            "name": "灌南县"
        },
        {
            "city_id": "81",
            "id": "803",
            "name": "清河区"
        },
        {
            "city_id": "81",
            "id": "804",
            "name": "楚州区"
        },
        {
            "city_id": "81",
            "id": "805",
            "name": "淮阴区"
        },
        {
            "city_id": "81",
            "id": "806",
            "name": "清浦区"
        },
        {
            "city_id": "81",
            "id": "807",
            "name": "涟水县"
        },
        {
            "city_id": "81",
            "id": "808",
            "name": "洪泽区"
        },
        {
            "city_id": "81",
            "id": "809",
            "name": "盱眙县"
        },
        {
            "city_id": "81",
            "id": "810",
            "name": "金湖县"
        },
        {
            "city_id": "82",
            "id": "811",
            "name": "亭湖区"
        },
        {
            "city_id": "82",
            "id": "812",
            "name": "盐都区"
        },
        {
            "city_id": "82",
            "id": "813",
            "name": "响水县"
        },
        {
            "city_id": "82",
            "id": "814",
            "name": "滨海县"
        },
        {
            "city_id": "82",
            "id": "815",
            "name": "阜宁县"
        },
        {
            "city_id": "82",
            "id": "816",
            "name": "射阳县"
        },
        {
            "city_id": "82",
            "id": "817",
            "name": "建湖县"
        },
        {
            "city_id": "82",
            "id": "818",
            "name": "东台市"
        },
        {
            "city_id": "82",
            "id": "819",
            "name": "大丰区"
        },
        {
            "city_id": "83",
            "id": "820",
            "name": "广陵区"
        },
        {
            "city_id": "83",
            "id": "821",
            "name": "邗江区"
        },
        {
            "city_id": "83",
            "id": "822",
            "name": "维扬区"
        },
        {
            "city_id": "83",
            "id": "823",
            "name": "宝应县"
        },
        {
            "city_id": "83",
            "id": "824",
            "name": "仪征市"
        },
        {
            "city_id": "83",
            "id": "825",
            "name": "高邮市"
        },
        {
            "city_id": "83",
            "id": "826",
            "name": "江都区"
        },
        {
            "city_id": "84",
            "id": "827",
            "name": "京口区"
        },
        {
            "city_id": "84",
            "id": "828",
            "name": "润州区"
        },
        {
            "city_id": "84",
            "id": "829",
            "name": "丹徒区"
        },
        {
            "city_id": "84",
            "id": "830",
            "name": "丹阳市"
        },
        {
            "city_id": "84",
            "id": "831",
            "name": "扬中市"
        },
        {
            "city_id": "84",
            "id": "832",
            "name": "句容市"
        },
        {
            "city_id": "85",
            "id": "833",
            "name": "海陵区"
        },
        {
            "city_id": "85",
            "id": "834",
            "name": "高港区"
        },
        {
            "city_id": "85",
            "id": "835",
            "name": "兴化市"
        },
        {
            "city_id": "85",
            "id": "836",
            "name": "靖江市"
        },
        {
            "city_id": "85",
            "id": "837",
            "name": "泰兴市"
        },
        {
            "city_id": "85",
            "id": "838",
            "name": "姜堰区"
        },
        {
            "city_id": "86",
            "id": "839",
            "name": "宿城区"
        },
        {
            "city_id": "86",
            "id": "840",
            "name": "宿豫区"
        },
        {
            "city_id": "86",
            "id": "841",
            "name": "沭阳县"
        },
        {
            "city_id": "86",
            "id": "842",
            "name": "泗阳县"
        },
        {
            "city_id": "86",
            "id": "843",
            "name": "泗洪县"
        },
        {
            "city_id": "87",
            "id": "844",
            "name": "上城区"
        },
        {
            "city_id": "87",
            "id": "845",
            "name": "下城区"
        },
        {
            "city_id": "87",
            "id": "846",
            "name": "江干区"
        },
        {
            "city_id": "87",
            "id": "847",
            "name": "拱墅区"
        },
        {
            "city_id": "87",
            "id": "848",
            "name": "西湖区"
        },
        {
            "city_id": "87",
            "id": "849",
            "name": "滨江区"
        },
        {
            "city_id": "87",
            "id": "850",
            "name": "萧山区"
        },
        {
            "city_id": "87",
            "id": "851",
            "name": "余杭区"
        },
        {
            "city_id": "87",
            "id": "852",
            "name": "桐庐县"
        },
        {
            "city_id": "87",
            "id": "853",
            "name": "淳安县"
        },
        {
            "city_id": "87",
            "id": "854",
            "name": "建德市"
        },
        {
            "city_id": "87",
            "id": "855",
            "name": "富阳区"
        },
        {
            "city_id": "87",
            "id": "856",
            "name": "临安区"
        },
        {
            "city_id": "88",
            "id": "857",
            "name": "海曙区"
        },
        {
            "city_id": "88",
            "id": "858",
            "name": "江东区"
        },
        {
            "city_id": "88",
            "id": "859",
            "name": "江北区"
        },
        {
            "city_id": "88",
            "id": "860",
            "name": "北仑区"
        },
        {
            "city_id": "88",
            "id": "861",
            "name": "镇海区"
        },
        {
            "city_id": "88",
            "id": "862",
            "name": "鄞州区"
        },
        {
            "city_id": "88",
            "id": "863",
            "name": "象山县"
        },
        {
            "city_id": "88",
            "id": "864",
            "name": "宁海县"
        },
        {
            "city_id": "88",
            "id": "865",
            "name": "余姚市"
        },
        {
            "city_id": "88",
            "id": "866",
            "name": "慈溪市"
        },
        {
            "city_id": "88",
            "id": "867",
            "name": "奉化区"
        },
        {
            "city_id": "89",
            "id": "868",
            "name": "鹿城区"
        },
        {
            "city_id": "89",
            "id": "869",
            "name": "龙湾区"
        },
        {
            "city_id": "89",
            "id": "870",
            "name": "瓯海区"
        },
        {
            "city_id": "89",
            "id": "871",
            "name": "洞头区"
        },
        {
            "city_id": "89",
            "id": "872",
            "name": "永嘉县"
        },
        {
            "city_id": "89",
            "id": "873",
            "name": "平阳县"
        },
        {
            "city_id": "89",
            "id": "874",
            "name": "苍南县"
        },
        {
            "city_id": "89",
            "id": "875",
            "name": "文成县"
        },
        {
            "city_id": "89",
            "id": "876",
            "name": "泰顺县"
        },
        {
            "city_id": "89",
            "id": "877",
            "name": "瑞安市"
        },
        {
            "city_id": "89",
            "id": "878",
            "name": "乐清市"
        },
        {
            "city_id": "90",
            "id": "879",
            "name": "秀城区"
        },
        {
            "city_id": "90",
            "id": "880",
            "name": "秀洲区"
        },
        {
            "city_id": "90",
            "id": "881",
            "name": "嘉善县"
        },
        {
            "city_id": "90",
            "id": "882",
            "name": "海盐县"
        },
        {
            "city_id": "90",
            "id": "883",
            "name": "海宁市"
        },
        {
            "city_id": "90",
            "id": "884",
            "name": "平湖市"
        },
        {
            "city_id": "90",
            "id": "885",
            "name": "桐乡市"
        },
        {
            "city_id": "91",
            "id": "886",
            "name": "吴兴区"
        },
        {
            "city_id": "91",
            "id": "887",
            "name": "南浔区"
        },
        {
            "city_id": "91",
            "id": "888",
            "name": "德清县"
        },
        {
            "city_id": "91",
            "id": "889",
            "name": "长兴县"
        },
        {
            "city_id": "91",
            "id": "890",
            "name": "安吉县"
        },
        {
            "city_id": "92",
            "id": "891",
            "name": "越城区"
        },
        {
            "city_id": "92",
            "id": "892",
            "name": "绍兴县"
        },
        {
            "city_id": "92",
            "id": "893",
            "name": "新昌县"
        },
        {
            "city_id": "92",
            "id": "894",
            "name": "诸暨市"
        },
        {
            "city_id": "92",
            "id": "895",
            "name": "上虞区"
        },
        {
            "city_id": "92",
            "id": "896",
            "name": "嵊州市"
        },
        {
            "city_id": "93",
            "id": "897",
            "name": "婺城区"
        },
        {
            "city_id": "93",
            "id": "898",
            "name": "金东区"
        },
        {
            "city_id": "93",
            "id": "899",
            "name": "武义县"
        },
        {
            "city_id": "93",
            "id": "900",
            "name": "浦江县"
        },
        {
            "city_id": "93",
            "id": "901",
            "name": "磐安县"
        },
        {
            "city_id": "93",
            "id": "902",
            "name": "兰溪市"
        },
        {
            "city_id": "93",
            "id": "903",
            "name": "义乌市"
        },
        {
            "city_id": "93",
            "id": "904",
            "name": "东阳市"
        },
        {
            "city_id": "93",
            "id": "905",
            "name": "永康市"
        },
        {
            "city_id": "94",
            "id": "906",
            "name": "柯城区"
        },
        {
            "city_id": "94",
            "id": "907",
            "name": "衢江区"
        },
        {
            "city_id": "94",
            "id": "908",
            "name": "常山县"
        },
        {
            "city_id": "94",
            "id": "909",
            "name": "开化县"
        },
        {
            "city_id": "94",
            "id": "910",
            "name": "龙游县"
        },
        {
            "city_id": "94",
            "id": "911",
            "name": "江山市"
        },
        {
            "city_id": "95",
            "id": "912",
            "name": "定海区"
        },
        {
            "city_id": "95",
            "id": "913",
            "name": "普陀区"
        },
        {
            "city_id": "95",
            "id": "914",
            "name": "岱山县"
        },
        {
            "city_id": "95",
            "id": "915",
            "name": "嵊泗县"
        },
        {
            "city_id": "96",
            "id": "916",
            "name": "椒江区"
        },
        {
            "city_id": "96",
            "id": "917",
            "name": "黄岩区"
        },
        {
            "city_id": "96",
            "id": "918",
            "name": "路桥区"
        },
        {
            "city_id": "96",
            "id": "919",
            "name": "玉环市"
        },
        {
            "city_id": "96",
            "id": "920",
            "name": "三门县"
        },
        {
            "city_id": "96",
            "id": "921",
            "name": "天台县"
        },
        {
            "city_id": "96",
            "id": "922",
            "name": "仙居县"
        },
        {
            "city_id": "96",
            "id": "923",
            "name": "温岭市"
        },
        {
            "city_id": "96",
            "id": "924",
            "name": "临海市"
        },
        {
            "city_id": "97",
            "id": "925",
            "name": "莲都区"
        },
        {
            "city_id": "97",
            "id": "926",
            "name": "青田县"
        },
        {
            "city_id": "97",
            "id": "927",
            "name": "缙云县"
        },
        {
            "city_id": "97",
            "id": "928",
            "name": "遂昌县"
        },
        {
            "city_id": "97",
            "id": "929",
            "name": "松阳县"
        },
        {
            "city_id": "97",
            "id": "930",
            "name": "云和县"
        },
        {
            "city_id": "97",
            "id": "931",
            "name": "庆元县"
        },
        {
            "city_id": "97",
            "id": "932",
            "name": "景宁畲族自治县"
        },
        {
            "city_id": "97",
            "id": "933",
            "name": "龙泉市"
        },
        {
            "city_id": "98",
            "id": "934",
            "name": "瑶海区"
        },
        {
            "city_id": "98",
            "id": "935",
            "name": "庐阳区"
        },
        {
            "city_id": "98",
            "id": "936",
            "name": "蜀山区"
        },
        {
            "city_id": "98",
            "id": "937",
            "name": "包河区"
        },
        {
            "city_id": "98",
            "id": "938",
            "name": "长丰县"
        },
        {
            "city_id": "98",
            "id": "939",
            "name": "肥东县"
        },
        {
            "city_id": "98",
            "id": "940",
            "name": "肥西县"
        },
        {
            "city_id": "99",
            "id": "941",
            "name": "镜湖区"
        },
        {
            "city_id": "99",
            "id": "942",
            "name": "马塘区"
        },
        {
            "city_id": "99",
            "id": "943",
            "name": "新芜区"
        },
        {
            "city_id": "99",
            "id": "944",
            "name": "鸠江区"
        },
        {
            "city_id": "99",
            "id": "945",
            "name": "芜湖县"
        },
        {
            "city_id": "99",
            "id": "946",
            "name": "繁昌县"
        },
        {
            "city_id": "99",
            "id": "947",
            "name": "南陵县"
        },
        {
            "city_id": "100",
            "id": "948",
            "name": "龙子湖区"
        },
        {
            "city_id": "100",
            "id": "949",
            "name": "蚌山区"
        },
        {
            "city_id": "100",
            "id": "950",
            "name": "禹会区"
        },
        {
            "city_id": "100",
            "id": "951",
            "name": "淮上区"
        },
        {
            "city_id": "100",
            "id": "952",
            "name": "怀远县"
        },
        {
            "city_id": "100",
            "id": "953",
            "name": "五河县"
        },
        {
            "city_id": "100",
            "id": "954",
            "name": "固镇县"
        },
        {
            "city_id": "101",
            "id": "955",
            "name": "大通区"
        },
        {
            "city_id": "101",
            "id": "956",
            "name": "田家庵区"
        },
        {
            "city_id": "101",
            "id": "957",
            "name": "谢家集区"
        },
        {
            "city_id": "101",
            "id": "958",
            "name": "八公山区"
        },
        {
            "city_id": "101",
            "id": "959",
            "name": "潘集区"
        },
        {
            "city_id": "101",
            "id": "960",
            "name": "凤台县"
        },
        {
            "city_id": "102",
            "id": "961",
            "name": "金家庄区"
        },
        {
            "city_id": "102",
            "id": "962",
            "name": "花山区"
        },
        {
            "city_id": "102",
            "id": "963",
            "name": "雨山区"
        },
        {
            "city_id": "102",
            "id": "964",
            "name": "当涂县"
        },
        {
            "city_id": "103",
            "id": "965",
            "name": "杜集区"
        },
        {
            "city_id": "103",
            "id": "966",
            "name": "相山区"
        },
        {
            "city_id": "103",
            "id": "967",
            "name": "烈山区"
        },
        {
            "city_id": "103",
            "id": "968",
            "name": "濉溪县"
        },
        {
            "city_id": "104",
            "id": "969",
            "name": "铜官区"
        },
        {
            "city_id": "104",
            "id": "970",
            "name": "狮子山区"
        },
        {
            "city_id": "104",
            "id": "971",
            "name": "郊区"
        },
        {
            "city_id": "104",
            "id": "972",
            "name": "铜陵县"
        },
        {
            "city_id": "105",
            "id": "973",
            "name": "迎江区"
        },
        {
            "city_id": "105",
            "id": "974",
            "name": "大观区"
        },
        {
            "city_id": "105",
            "id": "975",
            "name": "郊区"
        },
        {
            "city_id": "105",
            "id": "976",
            "name": "怀宁县"
        },
        {
            "city_id": "105",
            "id": "977",
            "name": "枞阳县"
        },
        {
            "city_id": "105",
            "id": "978",
            "name": "潜山县"
        },
        {
            "city_id": "105",
            "id": "979",
            "name": "太湖县"
        },
        {
            "city_id": "105",
            "id": "980",
            "name": "宿松县"
        },
        {
            "city_id": "105",
            "id": "981",
            "name": "望江县"
        },
        {
            "city_id": "105",
            "id": "982",
            "name": "岳西县"
        },
        {
            "city_id": "105",
            "id": "983",
            "name": "桐城市"
        },
        {
            "city_id": "106",
            "id": "984",
            "name": "屯溪区"
        },
        {
            "city_id": "106",
            "id": "985",
            "name": "黄山区"
        },
        {
            "city_id": "106",
            "id": "986",
            "name": "徽州区"
        },
        {
            "city_id": "106",
            "id": "987",
            "name": "歙县"
        },
        {
            "city_id": "106",
            "id": "988",
            "name": "休宁县"
        },
        {
            "city_id": "106",
            "id": "989",
            "name": "黟县"
        },
        {
            "city_id": "106",
            "id": "990",
            "name": "祁门县"
        },
        {
            "city_id": "107",
            "id": "991",
            "name": "琅琊区"
        },
        {
            "city_id": "107",
            "id": "992",
            "name": "南谯区"
        },
        {
            "city_id": "107",
            "id": "993",
            "name": "来安县"
        },
        {
            "city_id": "107",
            "id": "994",
            "name": "全椒县"
        },
        {
            "city_id": "107",
            "id": "995",
            "name": "定远县"
        },
        {
            "city_id": "107",
            "id": "996",
            "name": "凤阳县"
        },
        {
            "city_id": "107",
            "id": "997",
            "name": "天长市"
        },
        {
            "city_id": "107",
            "id": "998",
            "name": "明光市"
        },
        {
            "city_id": "108",
            "id": "999",
            "name": "颍州区"
        },
        {
            "city_id": "108",
            "id": "1000",
            "name": "颍东区"
        },
        {
            "city_id": "108",
            "id": "1001",
            "name": "颍泉区"
        },
        {
            "city_id": "108",
            "id": "1002",
            "name": "临泉县"
        },
        {
            "city_id": "108",
            "id": "1003",
            "name": "太和县"
        },
        {
            "city_id": "108",
            "id": "1004",
            "name": "阜南县"
        },
        {
            "city_id": "108",
            "id": "1005",
            "name": "颍上县"
        },
        {
            "city_id": "108",
            "id": "1006",
            "name": "界首市"
        },
        {
            "city_id": "109",
            "id": "1007",
            "name": "埇桥区"
        },
        {
            "city_id": "109",
            "id": "1008",
            "name": "砀山县"
        },
        {
            "city_id": "109",
            "id": "1009",
            "name": "萧县"
        },
        {
            "city_id": "109",
            "id": "1010",
            "name": "灵璧县"
        },
        {
            "city_id": "109",
            "id": "1011",
            "name": "泗县"
        },
        {
            "city_id": "110",
            "id": "1012",
            "name": "居巢区"
        },
        {
            "city_id": "110",
            "id": "1013",
            "name": "庐江县"
        },
        {
            "city_id": "110",
            "id": "1014",
            "name": "无为县"
        },
        {
            "city_id": "110",
            "id": "1015",
            "name": "含山县"
        },
        {
            "city_id": "110",
            "id": "1016",
            "name": "和县"
        },
        {
            "city_id": "111",
            "id": "1017",
            "name": "金安区"
        },
        {
            "city_id": "111",
            "id": "1018",
            "name": "裕安区"
        },
        {
            "city_id": "111",
            "id": "1019",
            "name": "寿县"
        },
        {
            "city_id": "111",
            "id": "1020",
            "name": "霍邱县"
        },
        {
            "city_id": "111",
            "id": "1021",
            "name": "舒城县"
        },
        {
            "city_id": "111",
            "id": "1022",
            "name": "金寨县"
        },
        {
            "city_id": "111",
            "id": "1023",
            "name": "霍山县"
        },
        {
            "city_id": "112",
            "id": "1024",
            "name": "谯城区"
        },
        {
            "city_id": "112",
            "id": "1025",
            "name": "涡阳县"
        },
        {
            "city_id": "112",
            "id": "1026",
            "name": "蒙城县"
        },
        {
            "city_id": "112",
            "id": "1027",
            "name": "利辛县"
        },
        {
            "city_id": "113",
            "id": "1028",
            "name": "贵池区"
        },
        {
            "city_id": "113",
            "id": "1029",
            "name": "东至县"
        },
        {
            "city_id": "113",
            "id": "1030",
            "name": "石台县"
        },
        {
            "city_id": "113",
            "id": "1031",
            "name": "青阳县"
        },
        {
            "city_id": "114",
            "id": "1032",
            "name": "宣州区"
        },
        {
            "city_id": "114",
            "id": "1033",
            "name": "郎溪县"
        },
        {
            "city_id": "114",
            "id": "1034",
            "name": "广德县"
        },
        {
            "city_id": "114",
            "id": "1035",
            "name": "泾县"
        },
        {
            "city_id": "114",
            "id": "1036",
            "name": "绩溪县"
        },
        {
            "city_id": "114",
            "id": "1037",
            "name": "旌德县"
        },
        {
            "city_id": "114",
            "id": "1038",
            "name": "宁国市"
        },
        {
            "city_id": "115",
            "id": "1039",
            "name": "鼓楼区"
        },
        {
            "city_id": "115",
            "id": "1040",
            "name": "台江区"
        },
        {
            "city_id": "115",
            "id": "1041",
            "name": "仓山区"
        },
        {
            "city_id": "115",
            "id": "1042",
            "name": "马尾区"
        },
        {
            "city_id": "115",
            "id": "1043",
            "name": "晋安区"
        },
        {
            "city_id": "115",
            "id": "1044",
            "name": "闽侯县"
        },
        {
            "city_id": "115",
            "id": "1045",
            "name": "连江县"
        },
        {
            "city_id": "115",
            "id": "1046",
            "name": "罗源县"
        },
        {
            "city_id": "115",
            "id": "1047",
            "name": "闽清县"
        },
        {
            "city_id": "115",
            "id": "1048",
            "name": "永泰县"
        },
        {
            "city_id": "115",
            "id": "1049",
            "name": "平潭县"
        },
        {
            "city_id": "115",
            "id": "1050",
            "name": "福清市"
        },
        {
            "city_id": "115",
            "id": "1051",
            "name": "长乐区"
        },
        {
            "city_id": "116",
            "id": "1052",
            "name": "思明区"
        },
        {
            "city_id": "116",
            "id": "1053",
            "name": "海沧区"
        },
        {
            "city_id": "116",
            "id": "1054",
            "name": "湖里区"
        },
        {
            "city_id": "116",
            "id": "1055",
            "name": "集美区"
        },
        {
            "city_id": "116",
            "id": "1056",
            "name": "同安区"
        },
        {
            "city_id": "116",
            "id": "1057",
            "name": "翔安区"
        },
        {
            "city_id": "117",
            "id": "1058",
            "name": "城厢区"
        },
        {
            "city_id": "117",
            "id": "1059",
            "name": "涵江区"
        },
        {
            "city_id": "117",
            "id": "1060",
            "name": "荔城区"
        },
        {
            "city_id": "117",
            "id": "1061",
            "name": "秀屿区"
        },
        {
            "city_id": "117",
            "id": "1062",
            "name": "仙游县"
        },
        {
            "city_id": "118",
            "id": "1063",
            "name": "梅列区"
        },
        {
            "city_id": "118",
            "id": "1064",
            "name": "三元区"
        },
        {
            "city_id": "118",
            "id": "1065",
            "name": "明溪县"
        },
        {
            "city_id": "118",
            "id": "1066",
            "name": "清流县"
        },
        {
            "city_id": "118",
            "id": "1067",
            "name": "宁化县"
        },
        {
            "city_id": "118",
            "id": "1068",
            "name": "大田县"
        },
        {
            "city_id": "118",
            "id": "1069",
            "name": "尤溪县"
        },
        {
            "city_id": "118",
            "id": "1070",
            "name": "沙县"
        },
        {
            "city_id": "118",
            "id": "1071",
            "name": "将乐县"
        },
        {
            "city_id": "118",
            "id": "1072",
            "name": "泰宁县"
        },
        {
            "city_id": "118",
            "id": "1073",
            "name": "建宁县"
        },
        {
            "city_id": "118",
            "id": "1074",
            "name": "永安市"
        },
        {
            "city_id": "119",
            "id": "1075",
            "name": "鲤城区"
        },
        {
            "city_id": "119",
            "id": "1076",
            "name": "丰泽区"
        },
        {
            "city_id": "119",
            "id": "1077",
            "name": "洛江区"
        },
        {
            "city_id": "119",
            "id": "1078",
            "name": "泉港区"
        },
        {
            "city_id": "119",
            "id": "1079",
            "name": "惠安县"
        },
        {
            "city_id": "119",
            "id": "1080",
            "name": "安溪县"
        },
        {
            "city_id": "119",
            "id": "1081",
            "name": "永春县"
        },
        {
            "city_id": "119",
            "id": "1082",
            "name": "德化县"
        },
        {
            "city_id": "119",
            "id": "1083",
            "name": "金门县"
        },
        {
            "city_id": "119",
            "id": "1084",
            "name": "石狮市"
        },
        {
            "city_id": "119",
            "id": "1085",
            "name": "晋江市"
        },
        {
            "city_id": "119",
            "id": "1086",
            "name": "南安市"
        },
        {
            "city_id": "120",
            "id": "1087",
            "name": "芗城区"
        },
        {
            "city_id": "120",
            "id": "1088",
            "name": "龙文区"
        },
        {
            "city_id": "120",
            "id": "1089",
            "name": "云霄县"
        },
        {
            "city_id": "120",
            "id": "1090",
            "name": "漳浦县"
        },
        {
            "city_id": "120",
            "id": "1091",
            "name": "诏安县"
        },
        {
            "city_id": "120",
            "id": "1092",
            "name": "长泰县"
        },
        {
            "city_id": "120",
            "id": "1093",
            "name": "东山县"
        },
        {
            "city_id": "120",
            "id": "1094",
            "name": "南靖县"
        },
        {
            "city_id": "120",
            "id": "1095",
            "name": "平和县"
        },
        {
            "city_id": "120",
            "id": "1096",
            "name": "华安县"
        },
        {
            "city_id": "120",
            "id": "1097",
            "name": "龙海市"
        },
        {
            "city_id": "121",
            "id": "1098",
            "name": "延平区"
        },
        {
            "city_id": "121",
            "id": "1099",
            "name": "顺昌县"
        },
        {
            "city_id": "121",
            "id": "1100",
            "name": "浦城县"
        },
        {
            "city_id": "121",
            "id": "1101",
            "name": "光泽县"
        },
        {
            "city_id": "121",
            "id": "1102",
            "name": "松溪县"
        },
        {
            "city_id": "121",
            "id": "1103",
            "name": "政和县"
        },
        {
            "city_id": "121",
            "id": "1104",
            "name": "邵武市"
        },
        {
            "city_id": "121",
            "id": "1105",
            "name": "武夷山市"
        },
        {
            "city_id": "121",
            "id": "1106",
            "name": "建瓯市"
        },
        {
            "city_id": "121",
            "id": "1107",
            "name": "建阳区"
        },
        {
            "city_id": "122",
            "id": "1108",
            "name": "新罗区"
        },
        {
            "city_id": "122",
            "id": "1109",
            "name": "长汀县"
        },
        {
            "city_id": "122",
            "id": "1110",
            "name": "永定区"
        },
        {
            "city_id": "122",
            "id": "1111",
            "name": "上杭县"
        },
        {
            "city_id": "122",
            "id": "1112",
            "name": "武平县"
        },
        {
            "city_id": "122",
            "id": "1113",
            "name": "连城县"
        },
        {
            "city_id": "122",
            "id": "1114",
            "name": "漳平市"
        },
        {
            "city_id": "123",
            "id": "1115",
            "name": "蕉城区"
        },
        {
            "city_id": "123",
            "id": "1116",
            "name": "霞浦县"
        },
        {
            "city_id": "123",
            "id": "1117",
            "name": "古田县"
        },
        {
            "city_id": "123",
            "id": "1118",
            "name": "屏南县"
        },
        {
            "city_id": "123",
            "id": "1119",
            "name": "寿宁县"
        },
        {
            "city_id": "123",
            "id": "1120",
            "name": "周宁县"
        },
        {
            "city_id": "123",
            "id": "1121",
            "name": "柘荣县"
        },
        {
            "city_id": "123",
            "id": "1122",
            "name": "福安市"
        },
        {
            "city_id": "123",
            "id": "1123",
            "name": "福鼎市"
        },
        {
            "city_id": "124",
            "id": "1124",
            "name": "东湖区"
        },
        {
            "city_id": "124",
            "id": "1125",
            "name": "西湖区"
        },
        {
            "city_id": "124",
            "id": "1126",
            "name": "青云谱区"
        },
        {
            "city_id": "124",
            "id": "1127",
            "name": "湾里区"
        },
        {
            "city_id": "124",
            "id": "1128",
            "name": "青山湖区"
        },
        {
            "city_id": "124",
            "id": "1129",
            "name": "南昌县"
        },
        {
            "city_id": "124",
            "id": "1130",
            "name": "新建区"
        },
        {
            "city_id": "124",
            "id": "1131",
            "name": "安义县"
        },
        {
            "city_id": "124",
            "id": "1132",
            "name": "进贤县"
        },
        {
            "city_id": "125",
            "id": "1133",
            "name": "昌江区"
        },
        {
            "city_id": "125",
            "id": "1134",
            "name": "珠山区"
        },
        {
            "city_id": "125",
            "id": "1135",
            "name": "浮梁县"
        },
        {
            "city_id": "125",
            "id": "1136",
            "name": "乐平市"
        },
        {
            "city_id": "126",
            "id": "1137",
            "name": "安源区"
        },
        {
            "city_id": "126",
            "id": "1138",
            "name": "湘东区"
        },
        {
            "city_id": "126",
            "id": "1139",
            "name": "莲花县"
        },
        {
            "city_id": "126",
            "id": "1140",
            "name": "上栗县"
        },
        {
            "city_id": "126",
            "id": "1141",
            "name": "芦溪县"
        },
        {
            "city_id": "127",
            "id": "1142",
            "name": "庐山市"
        },
        {
            "city_id": "127",
            "id": "1143",
            "name": "浔阳区"
        },
        {
            "city_id": "127",
            "id": "1144",
            "name": "九江县"
        },
        {
            "city_id": "127",
            "id": "1145",
            "name": "武宁县"
        },
        {
            "city_id": "127",
            "id": "1146",
            "name": "修水县"
        },
        {
            "city_id": "127",
            "id": "1147",
            "name": "永修县"
        },
        {
            "city_id": "127",
            "id": "1148",
            "name": "德安县"
        },
        {
            "city_id": "127",
            "id": "1149",
            "name": "星子县"
        },
        {
            "city_id": "127",
            "id": "1150",
            "name": "都昌县"
        },
        {
            "city_id": "127",
            "id": "1151",
            "name": "湖口县"
        },
        {
            "city_id": "127",
            "id": "1152",
            "name": "彭泽县"
        },
        {
            "city_id": "127",
            "id": "1153",
            "name": "瑞昌市"
        },
        {
            "city_id": "128",
            "id": "1154",
            "name": "渝水区"
        },
        {
            "city_id": "128",
            "id": "1155",
            "name": "分宜县"
        },
        {
            "city_id": "129",
            "id": "1156",
            "name": "月湖区"
        },
        {
            "city_id": "129",
            "id": "1157",
            "name": "余江区"
        },
        {
            "city_id": "129",
            "id": "1158",
            "name": "贵溪市"
        },
        {
            "city_id": "130",
            "id": "1159",
            "name": "章贡区"
        },
        {
            "city_id": "130",
            "id": "1160",
            "name": "赣县区"
        },
        {
            "city_id": "130",
            "id": "1161",
            "name": "信丰县"
        },
        {
            "city_id": "130",
            "id": "1162",
            "name": "大余县"
        },
        {
            "city_id": "130",
            "id": "1163",
            "name": "上犹县"
        },
        {
            "city_id": "130",
            "id": "1164",
            "name": "崇义县"
        },
        {
            "city_id": "130",
            "id": "1165",
            "name": "安远县"
        },
        {
            "city_id": "130",
            "id": "1166",
            "name": "龙南县"
        },
        {
            "city_id": "130",
            "id": "1167",
            "name": "定南县"
        },
        {
            "city_id": "130",
            "id": "1168",
            "name": "全南县"
        },
        {
            "city_id": "130",
            "id": "1169",
            "name": "宁都县"
        },
        {
            "city_id": "130",
            "id": "1170",
            "name": "于都县"
        },
        {
            "city_id": "130",
            "id": "1171",
            "name": "兴国县"
        },
        {
            "city_id": "130",
            "id": "1172",
            "name": "会昌县"
        },
        {
            "city_id": "130",
            "id": "1173",
            "name": "寻乌县"
        },
        {
            "city_id": "130",
            "id": "1174",
            "name": "石城县"
        },
        {
            "city_id": "130",
            "id": "1175",
            "name": "瑞金市"
        },
        {
            "city_id": "130",
            "id": "1176",
            "name": "南康区"
        },
        {
            "city_id": "131",
            "id": "1177",
            "name": "吉州区"
        },
        {
            "city_id": "131",
            "id": "1178",
            "name": "青原区"
        },
        {
            "city_id": "131",
            "id": "1179",
            "name": "吉安县"
        },
        {
            "city_id": "131",
            "id": "1180",
            "name": "吉水县"
        },
        {
            "city_id": "131",
            "id": "1181",
            "name": "峡江县"
        },
        {
            "city_id": "131",
            "id": "1182",
            "name": "新干县"
        },
        {
            "city_id": "131",
            "id": "1183",
            "name": "永丰县"
        },
        {
            "city_id": "131",
            "id": "1184",
            "name": "泰和县"
        },
        {
            "city_id": "131",
            "id": "1185",
            "name": "遂川县"
        },
        {
            "city_id": "131",
            "id": "1186",
            "name": "万安县"
        },
        {
            "city_id": "131",
            "id": "1187",
            "name": "安福县"
        },
        {
            "city_id": "131",
            "id": "1188",
            "name": "永新县"
        },
        {
            "city_id": "131",
            "id": "1189",
            "name": "井冈山市"
        },
        {
            "city_id": "132",
            "id": "1190",
            "name": "袁州区"
        },
        {
            "city_id": "132",
            "id": "1191",
            "name": "奉新县"
        },
        {
            "city_id": "132",
            "id": "1192",
            "name": "万载县"
        },
        {
            "city_id": "132",
            "id": "1193",
            "name": "上高县"
        },
        {
            "city_id": "132",
            "id": "1194",
            "name": "宜丰县"
        },
        {
            "city_id": "132",
            "id": "1195",
            "name": "靖安县"
        },
        {
            "city_id": "132",
            "id": "1196",
            "name": "铜鼓县"
        },
        {
            "city_id": "132",
            "id": "1197",
            "name": "丰城市"
        },
        {
            "city_id": "132",
            "id": "1198",
            "name": "樟树市"
        },
        {
            "city_id": "132",
            "id": "1199",
            "name": "高安市"
        },
        {
            "city_id": "133",
            "id": "1200",
            "name": "临川区"
        },
        {
            "city_id": "133",
            "id": "1201",
            "name": "南城县"
        },
        {
            "city_id": "133",
            "id": "1202",
            "name": "黎川县"
        },
        {
            "city_id": "133",
            "id": "1203",
            "name": "南丰县"
        },
        {
            "city_id": "133",
            "id": "1204",
            "name": "崇仁县"
        },
        {
            "city_id": "133",
            "id": "1205",
            "name": "乐安县"
        },
        {
            "city_id": "133",
            "id": "1206",
            "name": "宜黄县"
        },
        {
            "city_id": "133",
            "id": "1207",
            "name": "金溪县"
        },
        {
            "city_id": "133",
            "id": "1208",
            "name": "资溪县"
        },
        {
            "city_id": "133",
            "id": "1209",
            "name": "东乡区"
        },
        {
            "city_id": "133",
            "id": "1210",
            "name": "广昌县"
        },
        {
            "city_id": "134",
            "id": "1211",
            "name": "信州区"
        },
        {
            "city_id": "134",
            "id": "1212",
            "name": "上饶县"
        },
        {
            "city_id": "134",
            "id": "1213",
            "name": "广丰区"
        },
        {
            "city_id": "134",
            "id": "1214",
            "name": "玉山县"
        },
        {
            "city_id": "134",
            "id": "1215",
            "name": "铅山县"
        },
        {
            "city_id": "134",
            "id": "1216",
            "name": "横峰县"
        },
        {
            "city_id": "134",
            "id": "1217",
            "name": "弋阳县"
        },
        {
            "city_id": "134",
            "id": "1218",
            "name": "余干县"
        },
        {
            "city_id": "134",
            "id": "1219",
            "name": "鄱阳县"
        },
        {
            "city_id": "134",
            "id": "1220",
            "name": "万年县"
        },
        {
            "city_id": "134",
            "id": "1221",
            "name": "婺源县"
        },
        {
            "city_id": "134",
            "id": "1222",
            "name": "德兴市"
        },
        {
            "city_id": "135",
            "id": "1223",
            "name": "历下区"
        },
        {
            "city_id": "135",
            "id": "1224",
            "name": "市中区"
        },
        {
            "city_id": "135",
            "id": "1225",
            "name": "槐荫区"
        },
        {
            "city_id": "135",
            "id": "1226",
            "name": "天桥区"
        },
        {
            "city_id": "135",
            "id": "1227",
            "name": "历城区"
        },
        {
            "city_id": "135",
            "id": "1228",
            "name": "长清区"
        },
        {
            "city_id": "135",
            "id": "1229",
            "name": "平阴县"
        },
        {
            "city_id": "135",
            "id": "1230",
            "name": "济阳区"
        },
        {
            "city_id": "135",
            "id": "1231",
            "name": "商河县"
        },
        {
            "city_id": "135",
            "id": "1232",
            "name": "章丘区"
        },
        {
            "city_id": "136",
            "id": "1233",
            "name": "市南区"
        },
        {
            "city_id": "136",
            "id": "1234",
            "name": "市北区"
        },
        {
            "city_id": "136",
            "id": "1235",
            "name": "四方区"
        },
        {
            "city_id": "136",
            "id": "1236",
            "name": "黄岛区"
        },
        {
            "city_id": "136",
            "id": "1237",
            "name": "崂山区"
        },
        {
            "city_id": "136",
            "id": "1238",
            "name": "李沧区"
        },
        {
            "city_id": "136",
            "id": "1239",
            "name": "城阳区"
        },
        {
            "city_id": "136",
            "id": "1240",
            "name": "胶州市"
        },
        {
            "city_id": "136",
            "id": "1241",
            "name": "即墨区"
        },
        {
            "city_id": "136",
            "id": "1242",
            "name": "平度市"
        },
        {
            "city_id": "136",
            "id": "1243",
            "name": "胶南市"
        },
        {
            "city_id": "136",
            "id": "1244",
            "name": "莱西市"
        },
        {
            "city_id": "137",
            "id": "1245",
            "name": "淄川区"
        },
        {
            "city_id": "137",
            "id": "1246",
            "name": "张店区"
        },
        {
            "city_id": "137",
            "id": "1247",
            "name": "博山区"
        },
        {
            "city_id": "137",
            "id": "1248",
            "name": "临淄区"
        },
        {
            "city_id": "137",
            "id": "1249",
            "name": "周村区"
        },
        {
            "city_id": "137",
            "id": "1250",
            "name": "桓台县"
        },
        {
            "city_id": "137",
            "id": "1251",
            "name": "高青县"
        },
        {
            "city_id": "137",
            "id": "1252",
            "name": "沂源县"
        },
        {
            "city_id": "138",
            "id": "1253",
            "name": "市中区"
        },
        {
            "city_id": "138",
            "id": "1254",
            "name": "薛城区"
        },
        {
            "city_id": "138",
            "id": "1255",
            "name": "峄城区"
        },
        {
            "city_id": "138",
            "id": "1256",
            "name": "台儿庄区"
        },
        {
            "city_id": "138",
            "id": "1257",
            "name": "山亭区"
        },
        {
            "city_id": "138",
            "id": "1258",
            "name": "滕州市"
        },
        {
            "city_id": "139",
            "id": "1259",
            "name": "东营区"
        },
        {
            "city_id": "139",
            "id": "1260",
            "name": "河口区"
        },
        {
            "city_id": "139",
            "id": "1261",
            "name": "垦利区"
        },
        {
            "city_id": "139",
            "id": "1262",
            "name": "利津县"
        },
        {
            "city_id": "139",
            "id": "1263",
            "name": "广饶县"
        },
        {
            "city_id": "140",
            "id": "1264",
            "name": "芝罘区"
        },
        {
            "city_id": "140",
            "id": "1265",
            "name": "福山区"
        },
        {
            "city_id": "140",
            "id": "1266",
            "name": "牟平区"
        },
        {
            "city_id": "140",
            "id": "1267",
            "name": "莱山区"
        },
        {
            "city_id": "140",
            "id": "1268",
            "name": "长岛县"
        },
        {
            "city_id": "140",
            "id": "1269",
            "name": "龙口市"
        },
        {
            "city_id": "140",
            "id": "1270",
            "name": "莱阳市"
        },
        {
            "city_id": "140",
            "id": "1271",
            "name": "莱州市"
        },
        {
            "city_id": "140",
            "id": "1272",
            "name": "蓬莱市"
        },
        {
            "city_id": "140",
            "id": "1273",
            "name": "招远市"
        },
        {
            "city_id": "140",
            "id": "1274",
            "name": "栖霞市"
        },
        {
            "city_id": "140",
            "id": "1275",
            "name": "海阳市"
        },
        {
            "city_id": "141",
            "id": "1276",
            "name": "潍城区"
        },
        {
            "city_id": "141",
            "id": "1277",
            "name": "寒亭区"
        },
        {
            "city_id": "141",
            "id": "1278",
            "name": "坊子区"
        },
        {
            "city_id": "141",
            "id": "1279",
            "name": "奎文区"
        },
        {
            "city_id": "141",
            "id": "1280",
            "name": "临朐县"
        },
        {
            "city_id": "141",
            "id": "1281",
            "name": "昌乐县"
        },
        {
            "city_id": "141",
            "id": "1282",
            "name": "青州市"
        },
        {
            "city_id": "141",
            "id": "1283",
            "name": "诸城市"
        },
        {
            "city_id": "141",
            "id": "1284",
            "name": "寿光市"
        },
        {
            "city_id": "141",
            "id": "1285",
            "name": "安丘市"
        },
        {
            "city_id": "141",
            "id": "1286",
            "name": "高密市"
        },
        {
            "city_id": "141",
            "id": "1287",
            "name": "昌邑市"
        },
        {
            "city_id": "142",
            "id": "1288",
            "name": "市中区"
        },
        {
            "city_id": "142",
            "id": "1289",
            "name": "任城区"
        },
        {
            "city_id": "142",
            "id": "1290",
            "name": "微山县"
        },
        {
            "city_id": "142",
            "id": "1291",
            "name": "鱼台县"
        },
        {
            "city_id": "142",
            "id": "1292",
            "name": "金乡县"
        },
        {
            "city_id": "142",
            "id": "1293",
            "name": "嘉祥县"
        },
        {
            "city_id": "142",
            "id": "1294",
            "name": "汶上县"
        },
        {
            "city_id": "142",
            "id": "1295",
            "name": "泗水县"
        },
        {
            "city_id": "142",
            "id": "1296",
            "name": "梁山县"
        },
        {
            "city_id": "142",
            "id": "1297",
            "name": "曲阜市"
        },
        {
            "city_id": "142",
            "id": "1298",
            "name": "兖州区"
        },
        {
            "city_id": "142",
            "id": "1299",
            "name": "邹城市"
        },
        {
            "city_id": "143",
            "id": "1300",
            "name": "泰山区"
        },
        {
            "city_id": "143",
            "id": "1301",
            "name": "岱岳区"
        },
        {
            "city_id": "143",
            "id": "1302",
            "name": "宁阳县"
        },
        {
            "city_id": "143",
            "id": "1303",
            "name": "东平县"
        },
        {
            "city_id": "143",
            "id": "1304",
            "name": "新泰市"
        },
        {
            "city_id": "143",
            "id": "1305",
            "name": "肥城市"
        },
        {
            "city_id": "144",
            "id": "1306",
            "name": "环翠区"
        },
        {
            "city_id": "144",
            "id": "1307",
            "name": "文登区"
        },
        {
            "city_id": "144",
            "id": "1308",
            "name": "荣成市"
        },
        {
            "city_id": "144",
            "id": "1309",
            "name": "乳山市"
        },
        {
            "city_id": "145",
            "id": "1310",
            "name": "东港区"
        },
        {
            "city_id": "145",
            "id": "1311",
            "name": "岚山区"
        },
        {
            "city_id": "145",
            "id": "1312",
            "name": "五莲县"
        },
        {
            "city_id": "145",
            "id": "1313",
            "name": "莒县"
        },
        {
            "city_id": "146",
            "id": "1314",
            "name": "莱城区"
        },
        {
            "city_id": "146",
            "id": "1315",
            "name": "钢城区"
        },
        {
            "city_id": "147",
            "id": "1316",
            "name": "兰山区"
        },
        {
            "city_id": "147",
            "id": "1317",
            "name": "罗庄区"
        },
        {
            "city_id": "147",
            "id": "1318",
            "name": "河东区"
        },
        {
            "city_id": "147",
            "id": "1319",
            "name": "沂南县"
        },
        {
            "city_id": "147",
            "id": "1320",
            "name": "郯城县"
        },
        {
            "city_id": "147",
            "id": "1321",
            "name": "沂水县"
        },
        {
            "city_id": "147",
            "id": "1322",
            "name": "苍山县"
        },
        {
            "city_id": "147",
            "id": "1323",
            "name": "费县"
        },
        {
            "city_id": "147",
            "id": "1324",
            "name": "平邑县"
        },
        {
            "city_id": "147",
            "id": "1325",
            "name": "莒南县"
        },
        {
            "city_id": "147",
            "id": "1326",
            "name": "蒙阴县"
        },
        {
            "city_id": "147",
            "id": "1327",
            "name": "临沭县"
        },
        {
            "city_id": "148",
            "id": "1328",
            "name": "德城区"
        },
        {
            "city_id": "148",
            "id": "1329",
            "name": "陵县"
        },
        {
            "city_id": "148",
            "id": "1330",
            "name": "宁津县"
        },
        {
            "city_id": "148",
            "id": "1331",
            "name": "庆云县"
        },
        {
            "city_id": "148",
            "id": "1332",
            "name": "临邑县"
        },
        {
            "city_id": "148",
            "id": "1333",
            "name": "齐河县"
        },
        {
            "city_id": "148",
            "id": "1334",
            "name": "平原县"
        },
        {
            "city_id": "148",
            "id": "1335",
            "name": "夏津县"
        },
        {
            "city_id": "148",
            "id": "1336",
            "name": "武城县"
        },
        {
            "city_id": "148",
            "id": "1337",
            "name": "乐陵市"
        },
        {
            "city_id": "148",
            "id": "1338",
            "name": "禹城市"
        },
        {
            "city_id": "149",
            "id": "1339",
            "name": "东昌府区"
        },
        {
            "city_id": "149",
            "id": "1340",
            "name": "阳谷县"
        },
        {
            "city_id": "149",
            "id": "1341",
            "name": "莘县"
        },
        {
            "city_id": "149",
            "id": "1342",
            "name": "茌平县"
        },
        {
            "city_id": "149",
            "id": "1343",
            "name": "东阿县"
        },
        {
            "city_id": "149",
            "id": "1344",
            "name": "冠县"
        },
        {
            "city_id": "149",
            "id": "1345",
            "name": "高唐县"
        },
        {
            "city_id": "149",
            "id": "1346",
            "name": "临清市"
        },
        {
            "city_id": "150",
            "id": "1347",
            "name": "滨城区"
        },
        {
            "city_id": "150",
            "id": "1348",
            "name": "惠民县"
        },
        {
            "city_id": "150",
            "id": "1349",
            "name": "阳信县"
        },
        {
            "city_id": "150",
            "id": "1350",
            "name": "无棣县"
        },
        {
            "city_id": "150",
            "id": "1351",
            "name": "沾化区"
        },
        {
            "city_id": "150",
            "id": "1352",
            "name": "博兴县"
        },
        {
            "city_id": "150",
            "id": "1353",
            "name": "邹平市"
        },
        {
            "city_id": "151",
            "id": "1354",
            "name": "牡丹区"
        },
        {
            "city_id": "151",
            "id": "1355",
            "name": "曹县"
        },
        {
            "city_id": "151",
            "id": "1356",
            "name": "单县"
        },
        {
            "city_id": "151",
            "id": "1357",
            "name": "成武县"
        },
        {
            "city_id": "151",
            "id": "1358",
            "name": "巨野县"
        },
        {
            "city_id": "151",
            "id": "1359",
            "name": "郓城县"
        },
        {
            "city_id": "151",
            "id": "1360",
            "name": "鄄城县"
        },
        {
            "city_id": "151",
            "id": "1361",
            "name": "定陶区"
        },
        {
            "city_id": "151",
            "id": "1362",
            "name": "东明县"
        },
        {
            "city_id": "152",
            "id": "1363",
            "name": "中原区"
        },
        {
            "city_id": "152",
            "id": "1364",
            "name": "二七区"
        },
        {
            "city_id": "152",
            "id": "1365",
            "name": "管城回族区"
        },
        {
            "city_id": "152",
            "id": "1366",
            "name": "金水区"
        },
        {
            "city_id": "152",
            "id": "1367",
            "name": "上街区"
        },
        {
            "city_id": "152",
            "id": "1368",
            "name": "惠济区"
        },
        {
            "city_id": "152",
            "id": "1369",
            "name": "中牟县"
        },
        {
            "city_id": "152",
            "id": "1370",
            "name": "巩义市"
        },
        {
            "city_id": "152",
            "id": "1371",
            "name": "荥阳市"
        },
        {
            "city_id": "152",
            "id": "1372",
            "name": "新密市"
        },
        {
            "city_id": "152",
            "id": "1373",
            "name": "新郑市"
        },
        {
            "city_id": "152",
            "id": "1374",
            "name": "登封市"
        },
        {
            "city_id": "153",
            "id": "1375",
            "name": "龙亭区"
        },
        {
            "city_id": "153",
            "id": "1376",
            "name": "顺河回族区"
        },
        {
            "city_id": "153",
            "id": "1377",
            "name": "鼓楼区"
        },
        {
            "city_id": "153",
            "id": "1378",
            "name": "南关区"
        },
        {
            "city_id": "153",
            "id": "1379",
            "name": "郊区"
        },
        {
            "city_id": "153",
            "id": "1380",
            "name": "杞县"
        },
        {
            "city_id": "153",
            "id": "1381",
            "name": "通许县"
        },
        {
            "city_id": "153",
            "id": "1382",
            "name": "尉氏县"
        },
        {
            "city_id": "153",
            "id": "1383",
            "name": "开封县"
        },
        {
            "city_id": "153",
            "id": "1384",
            "name": "兰考县"
        },
        {
            "city_id": "154",
            "id": "1385",
            "name": "老城区"
        },
        {
            "city_id": "154",
            "id": "1386",
            "name": "西工区"
        },
        {
            "city_id": "154",
            "id": "1387",
            "name": "廛河回族区"
        },
        {
            "city_id": "154",
            "id": "1388",
            "name": "涧西区"
        },
        {
            "city_id": "154",
            "id": "1389",
            "name": "吉利区"
        },
        {
            "city_id": "154",
            "id": "1390",
            "name": "洛龙区"
        },
        {
            "city_id": "154",
            "id": "1391",
            "name": "孟津县"
        },
        {
            "city_id": "154",
            "id": "1392",
            "name": "新安县"
        },
        {
            "city_id": "154",
            "id": "1393",
            "name": "栾川县"
        },
        {
            "city_id": "154",
            "id": "1394",
            "name": "嵩县"
        },
        {
            "city_id": "154",
            "id": "1395",
            "name": "汝阳县"
        },
        {
            "city_id": "154",
            "id": "1396",
            "name": "宜阳县"
        },
        {
            "city_id": "154",
            "id": "1397",
            "name": "洛宁县"
        },
        {
            "city_id": "154",
            "id": "1398",
            "name": "伊川县"
        },
        {
            "city_id": "154",
            "id": "1399",
            "name": "偃师市"
        },
        {
            "city_id": "155",
            "id": "1400",
            "name": "新华区"
        },
        {
            "city_id": "155",
            "id": "1401",
            "name": "卫东区"
        },
        {
            "city_id": "155",
            "id": "1402",
            "name": "石龙区"
        },
        {
            "city_id": "155",
            "id": "1403",
            "name": "湛河区"
        },
        {
            "city_id": "155",
            "id": "1404",
            "name": "宝丰县"
        },
        {
            "city_id": "155",
            "id": "1405",
            "name": "叶县"
        },
        {
            "city_id": "155",
            "id": "1406",
            "name": "鲁山县"
        },
        {
            "city_id": "155",
            "id": "1407",
            "name": "郏县"
        },
        {
            "city_id": "155",
            "id": "1408",
            "name": "舞钢市"
        },
        {
            "city_id": "155",
            "id": "1409",
            "name": "汝州市"
        },
        {
            "city_id": "156",
            "id": "1410",
            "name": "文峰区"
        },
        {
            "city_id": "156",
            "id": "1411",
            "name": "北关区"
        },
        {
            "city_id": "156",
            "id": "1412",
            "name": "殷都区"
        },
        {
            "city_id": "156",
            "id": "1413",
            "name": "龙安区"
        },
        {
            "city_id": "156",
            "id": "1414",
            "name": "安阳县"
        },
        {
            "city_id": "156",
            "id": "1415",
            "name": "汤阴县"
        },
        {
            "city_id": "156",
            "id": "1416",
            "name": "滑县"
        },
        {
            "city_id": "156",
            "id": "1417",
            "name": "内黄县"
        },
        {
            "city_id": "156",
            "id": "1418",
            "name": "林州市"
        },
        {
            "city_id": "157",
            "id": "1419",
            "name": "鹤山区"
        },
        {
            "city_id": "157",
            "id": "1420",
            "name": "山城区"
        },
        {
            "city_id": "157",
            "id": "1421",
            "name": "淇滨区"
        },
        {
            "city_id": "157",
            "id": "1422",
            "name": "浚县"
        },
        {
            "city_id": "157",
            "id": "1423",
            "name": "淇县"
        },
        {
            "city_id": "158",
            "id": "1424",
            "name": "红旗区"
        },
        {
            "city_id": "158",
            "id": "1425",
            "name": "卫滨区"
        },
        {
            "city_id": "158",
            "id": "1426",
            "name": "凤泉区"
        },
        {
            "city_id": "158",
            "id": "1427",
            "name": "牧野区"
        },
        {
            "city_id": "158",
            "id": "1428",
            "name": "新乡县"
        },
        {
            "city_id": "158",
            "id": "1429",
            "name": "获嘉县"
        },
        {
            "city_id": "158",
            "id": "1430",
            "name": "原阳县"
        },
        {
            "city_id": "158",
            "id": "1431",
            "name": "延津县"
        },
        {
            "city_id": "158",
            "id": "1432",
            "name": "封丘县"
        },
        {
            "city_id": "158",
            "id": "1433",
            "name": "长垣县"
        },
        {
            "city_id": "158",
            "id": "1434",
            "name": "卫辉市"
        },
        {
            "city_id": "158",
            "id": "1435",
            "name": "辉县市"
        },
        {
            "city_id": "159",
            "id": "1436",
            "name": "解放区"
        },
        {
            "city_id": "159",
            "id": "1437",
            "name": "中站区"
        },
        {
            "city_id": "159",
            "id": "1438",
            "name": "马村区"
        },
        {
            "city_id": "159",
            "id": "1439",
            "name": "山阳区"
        },
        {
            "city_id": "159",
            "id": "1440",
            "name": "修武县"
        },
        {
            "city_id": "159",
            "id": "1441",
            "name": "博爱县"
        },
        {
            "city_id": "159",
            "id": "1442",
            "name": "武陟县"
        },
        {
            "city_id": "159",
            "id": "1443",
            "name": "温县"
        },
        {
            "city_id": "350",
            "id": "1444",
            "name": "济源市"
        },
        {
            "city_id": "159",
            "id": "1445",
            "name": "沁阳市"
        },
        {
            "city_id": "159",
            "id": "1446",
            "name": "孟州市"
        },
        {
            "city_id": "160",
            "id": "1447",
            "name": "华龙区"
        },
        {
            "city_id": "160",
            "id": "1448",
            "name": "清丰县"
        },
        {
            "city_id": "160",
            "id": "1449",
            "name": "南乐县"
        },
        {
            "city_id": "160",
            "id": "1450",
            "name": "范县"
        },
        {
            "city_id": "160",
            "id": "1451",
            "name": "台前县"
        },
        {
            "city_id": "160",
            "id": "1452",
            "name": "濮阳县"
        },
        {
            "city_id": "161",
            "id": "1453",
            "name": "魏都区"
        },
        {
            "city_id": "161",
            "id": "1454",
            "name": "许昌县"
        },
        {
            "city_id": "161",
            "id": "1455",
            "name": "鄢陵县"
        },
        {
            "city_id": "161",
            "id": "1456",
            "name": "襄城县"
        },
        {
            "city_id": "161",
            "id": "1457",
            "name": "禹州市"
        },
        {
            "city_id": "161",
            "id": "1458",
            "name": "长葛市"
        },
        {
            "city_id": "162",
            "id": "1459",
            "name": "源汇区"
        },
        {
            "city_id": "162",
            "id": "1460",
            "name": "郾城区"
        },
        {
            "city_id": "162",
            "id": "1461",
            "name": "召陵区"
        },
        {
            "city_id": "162",
            "id": "1462",
            "name": "舞阳县"
        },
        {
            "city_id": "162",
            "id": "1463",
            "name": "临颍县"
        },
        {
            "city_id": "163",
            "id": "1464",
            "name": "市辖区"
        },
        {
            "city_id": "163",
            "id": "1465",
            "name": "湖滨区"
        },
        {
            "city_id": "163",
            "id": "1466",
            "name": "渑池县"
        },
        {
            "city_id": "163",
            "id": "1467",
            "name": "陕县"
        },
        {
            "city_id": "163",
            "id": "1468",
            "name": "卢氏县"
        },
        {
            "city_id": "163",
            "id": "1469",
            "name": "义马市"
        },
        {
            "city_id": "163",
            "id": "1470",
            "name": "灵宝市"
        },
        {
            "city_id": "164",
            "id": "1471",
            "name": "宛城区"
        },
        {
            "city_id": "164",
            "id": "1472",
            "name": "卧龙区"
        },
        {
            "city_id": "164",
            "id": "1473",
            "name": "南召县"
        },
        {
            "city_id": "164",
            "id": "1474",
            "name": "方城县"
        },
        {
            "city_id": "164",
            "id": "1475",
            "name": "西峡县"
        },
        {
            "city_id": "164",
            "id": "1476",
            "name": "镇平县"
        },
        {
            "city_id": "164",
            "id": "1477",
            "name": "内乡县"
        },
        {
            "city_id": "164",
            "id": "1478",
            "name": "淅川县"
        },
        {
            "city_id": "164",
            "id": "1479",
            "name": "社旗县"
        },
        {
            "city_id": "164",
            "id": "1480",
            "name": "唐河县"
        },
        {
            "city_id": "164",
            "id": "1481",
            "name": "新野县"
        },
        {
            "city_id": "164",
            "id": "1482",
            "name": "桐柏县"
        },
        {
            "city_id": "164",
            "id": "1483",
            "name": "邓州市"
        },
        {
            "city_id": "165",
            "id": "1484",
            "name": "梁园区"
        },
        {
            "city_id": "165",
            "id": "1485",
            "name": "睢阳区"
        },
        {
            "city_id": "165",
            "id": "1486",
            "name": "民权县"
        },
        {
            "city_id": "165",
            "id": "1487",
            "name": "睢县"
        },
        {
            "city_id": "165",
            "id": "1488",
            "name": "宁陵县"
        },
        {
            "city_id": "165",
            "id": "1489",
            "name": "柘城县"
        },
        {
            "city_id": "165",
            "id": "1490",
            "name": "虞城县"
        },
        {
            "city_id": "165",
            "id": "1491",
            "name": "夏邑县"
        },
        {
            "city_id": "165",
            "id": "1492",
            "name": "永城市"
        },
        {
            "city_id": "166",
            "id": "1493",
            "name": "浉河区"
        },
        {
            "city_id": "166",
            "id": "1494",
            "name": "平桥区"
        },
        {
            "city_id": "166",
            "id": "1495",
            "name": "罗山县"
        },
        {
            "city_id": "166",
            "id": "1496",
            "name": "光山县"
        },
        {
            "city_id": "166",
            "id": "1497",
            "name": "新县"
        },
        {
            "city_id": "166",
            "id": "1498",
            "name": "商城县"
        },
        {
            "city_id": "166",
            "id": "1499",
            "name": "固始县"
        },
        {
            "city_id": "166",
            "id": "1500",
            "name": "潢川县"
        },
        {
            "city_id": "166",
            "id": "1501",
            "name": "淮滨县"
        },
        {
            "city_id": "166",
            "id": "1502",
            "name": "息县"
        },
        {
            "city_id": "167",
            "id": "1503",
            "name": "川汇区"
        },
        {
            "city_id": "167",
            "id": "1504",
            "name": "扶沟县"
        },
        {
            "city_id": "167",
            "id": "1505",
            "name": "西华县"
        },
        {
            "city_id": "167",
            "id": "1506",
            "name": "商水县"
        },
        {
            "city_id": "167",
            "id": "1507",
            "name": "沈丘县"
        },
        {
            "city_id": "167",
            "id": "1508",
            "name": "郸城县"
        },
        {
            "city_id": "167",
            "id": "1509",
            "name": "淮阳县"
        },
        {
            "city_id": "167",
            "id": "1510",
            "name": "太康县"
        },
        {
            "city_id": "167",
            "id": "1511",
            "name": "鹿邑县"
        },
        {
            "city_id": "167",
            "id": "1512",
            "name": "项城市"
        },
        {
            "city_id": "168",
            "id": "1513",
            "name": "驿城区"
        },
        {
            "city_id": "168",
            "id": "1514",
            "name": "西平县"
        },
        {
            "city_id": "168",
            "id": "1515",
            "name": "上蔡县"
        },
        {
            "city_id": "168",
            "id": "1516",
            "name": "平舆县"
        },
        {
            "city_id": "168",
            "id": "1517",
            "name": "正阳县"
        },
        {
            "city_id": "168",
            "id": "1518",
            "name": "确山县"
        },
        {
            "city_id": "168",
            "id": "1519",
            "name": "泌阳县"
        },
        {
            "city_id": "168",
            "id": "1520",
            "name": "汝南县"
        },
        {
            "city_id": "168",
            "id": "1521",
            "name": "遂平县"
        },
        {
            "city_id": "168",
            "id": "1522",
            "name": "新蔡县"
        },
        {
            "city_id": "169",
            "id": "1523",
            "name": "江岸区"
        },
        {
            "city_id": "169",
            "id": "1524",
            "name": "江汉区"
        },
        {
            "city_id": "169",
            "id": "1525",
            "name": "硚口区"
        },
        {
            "city_id": "169",
            "id": "1526",
            "name": "汉阳区"
        },
        {
            "city_id": "169",
            "id": "1527",
            "name": "武昌区"
        },
        {
            "city_id": "169",
            "id": "1528",
            "name": "青山区"
        },
        {
            "city_id": "169",
            "id": "1529",
            "name": "洪山区"
        },
        {
            "city_id": "169",
            "id": "1530",
            "name": "东西湖区"
        },
        {
            "city_id": "169",
            "id": "1531",
            "name": "汉南区"
        },
        {
            "city_id": "169",
            "id": "1532",
            "name": "蔡甸区"
        },
        {
            "city_id": "169",
            "id": "1533",
            "name": "江夏区"
        },
        {
            "city_id": "169",
            "id": "1534",
            "name": "黄陂区"
        },
        {
            "city_id": "169",
            "id": "1535",
            "name": "新洲区"
        },
        {
            "city_id": "170",
            "id": "1536",
            "name": "黄石港区"
        },
        {
            "city_id": "170",
            "id": "1537",
            "name": "西塞山区"
        },
        {
            "city_id": "170",
            "id": "1538",
            "name": "下陆区"
        },
        {
            "city_id": "170",
            "id": "1539",
            "name": "铁山区"
        },
        {
            "city_id": "170",
            "id": "1540",
            "name": "阳新县"
        },
        {
            "city_id": "170",
            "id": "1541",
            "name": "大冶市"
        },
        {
            "city_id": "171",
            "id": "1542",
            "name": "茅箭区"
        },
        {
            "city_id": "171",
            "id": "1543",
            "name": "张湾区"
        },
        {
            "city_id": "171",
            "id": "1544",
            "name": "郧县"
        },
        {
            "city_id": "171",
            "id": "1545",
            "name": "郧西县"
        },
        {
            "city_id": "171",
            "id": "1546",
            "name": "竹山县"
        },
        {
            "city_id": "171",
            "id": "1547",
            "name": "竹溪县"
        },
        {
            "city_id": "171",
            "id": "1548",
            "name": "房县"
        },
        {
            "city_id": "171",
            "id": "1549",
            "name": "丹江口市"
        },
        {
            "city_id": "172",
            "id": "1550",
            "name": "西陵区"
        },
        {
            "city_id": "172",
            "id": "1551",
            "name": "伍家岗区"
        },
        {
            "city_id": "172",
            "id": "1552",
            "name": "点军区"
        },
        {
            "city_id": "172",
            "id": "1553",
            "name": "猇亭区"
        },
        {
            "city_id": "172",
            "id": "1554",
            "name": "夷陵区"
        },
        {
            "city_id": "172",
            "id": "1555",
            "name": "远安县"
        },
        {
            "city_id": "172",
            "id": "1556",
            "name": "兴山县"
        },
        {
            "city_id": "172",
            "id": "1557",
            "name": "秭归县"
        },
        {
            "city_id": "172",
            "id": "1558",
            "name": "长阳土家族自治县"
        },
        {
            "city_id": "172",
            "id": "1559",
            "name": "五峰土家族自治县"
        },
        {
            "city_id": "172",
            "id": "1560",
            "name": "宜都市"
        },
        {
            "city_id": "172",
            "id": "1561",
            "name": "当阳市"
        },
        {
            "city_id": "172",
            "id": "1562",
            "name": "枝江市"
        },
        {
            "city_id": "173",
            "id": "1563",
            "name": "襄城区"
        },
        {
            "city_id": "173",
            "id": "1564",
            "name": "樊城区"
        },
        {
            "city_id": "173",
            "id": "1565",
            "name": "襄阳区"
        },
        {
            "city_id": "173",
            "id": "1566",
            "name": "南漳县"
        },
        {
            "city_id": "173",
            "id": "1567",
            "name": "谷城县"
        },
        {
            "city_id": "173",
            "id": "1568",
            "name": "保康县"
        },
        {
            "city_id": "173",
            "id": "1569",
            "name": "老河口市"
        },
        {
            "city_id": "173",
            "id": "1570",
            "name": "枣阳市"
        },
        {
            "city_id": "173",
            "id": "1571",
            "name": "宜城市"
        },
        {
            "city_id": "174",
            "id": "1572",
            "name": "梁子湖区"
        },
        {
            "city_id": "174",
            "id": "1573",
            "name": "华容区"
        },
        {
            "city_id": "174",
            "id": "1574",
            "name": "鄂城区"
        },
        {
            "city_id": "175",
            "id": "1575",
            "name": "东宝区"
        },
        {
            "city_id": "175",
            "id": "1576",
            "name": "掇刀区"
        },
        {
            "city_id": "175",
            "id": "1577",
            "name": "京山市"
        },
        {
            "city_id": "175",
            "id": "1578",
            "name": "沙洋县"
        },
        {
            "city_id": "175",
            "id": "1579",
            "name": "钟祥市"
        },
        {
            "city_id": "176",
            "id": "1580",
            "name": "孝南区"
        },
        {
            "city_id": "176",
            "id": "1581",
            "name": "孝昌县"
        },
        {
            "city_id": "176",
            "id": "1582",
            "name": "大悟县"
        },
        {
            "city_id": "176",
            "id": "1583",
            "name": "云梦县"
        },
        {
            "city_id": "176",
            "id": "1584",
            "name": "应城市"
        },
        {
            "city_id": "176",
            "id": "1585",
            "name": "安陆市"
        },
        {
            "city_id": "176",
            "id": "1586",
            "name": "汉川市"
        },
        {
            "city_id": "177",
            "id": "1587",
            "name": "沙市区"
        },
        {
            "city_id": "177",
            "id": "1588",
            "name": "荆州区"
        },
        {
            "city_id": "177",
            "id": "1589",
            "name": "公安县"
        },
        {
            "city_id": "177",
            "id": "1590",
            "name": "监利县"
        },
        {
            "city_id": "177",
            "id": "1591",
            "name": "江陵县"
        },
        {
            "city_id": "177",
            "id": "1592",
            "name": "石首市"
        },
        {
            "city_id": "177",
            "id": "1593",
            "name": "洪湖市"
        },
        {
            "city_id": "177",
            "id": "1594",
            "name": "松滋市"
        },
        {
            "city_id": "178",
            "id": "1595",
            "name": "黄州区"
        },
        {
            "city_id": "178",
            "id": "1596",
            "name": "团风县"
        },
        {
            "city_id": "178",
            "id": "1597",
            "name": "红安县"
        },
        {
            "city_id": "178",
            "id": "1598",
            "name": "罗田县"
        },
        {
            "city_id": "178",
            "id": "1599",
            "name": "英山县"
        },
        {
            "city_id": "178",
            "id": "1600",
            "name": "浠水县"
        },
        {
            "city_id": "178",
            "id": "1601",
            "name": "蕲春县"
        },
        {
            "city_id": "178",
            "id": "1602",
            "name": "黄梅县"
        },
        {
            "city_id": "178",
            "id": "1603",
            "name": "麻城市"
        },
        {
            "city_id": "178",
            "id": "1604",
            "name": "武穴市"
        },
        {
            "city_id": "179",
            "id": "1605",
            "name": "咸安区"
        },
        {
            "city_id": "179",
            "id": "1606",
            "name": "嘉鱼县"
        },
        {
            "city_id": "179",
            "id": "1607",
            "name": "通城县"
        },
        {
            "city_id": "179",
            "id": "1608",
            "name": "崇阳县"
        },
        {
            "city_id": "179",
            "id": "1609",
            "name": "通山县"
        },
        {
            "city_id": "179",
            "id": "1610",
            "name": "赤壁市"
        },
        {
            "city_id": "180",
            "id": "1611",
            "name": "曾都区"
        },
        {
            "city_id": "180",
            "id": "1612",
            "name": "广水市"
        },
        {
            "city_id": "181",
            "id": "1613",
            "name": "恩施市"
        },
        {
            "city_id": "181",
            "id": "1614",
            "name": "利川市"
        },
        {
            "city_id": "181",
            "id": "1615",
            "name": "建始县"
        },
        {
            "city_id": "181",
            "id": "1616",
            "name": "巴东县"
        },
        {
            "city_id": "181",
            "id": "1617",
            "name": "宣恩县"
        },
        {
            "city_id": "181",
            "id": "1618",
            "name": "咸丰县"
        },
        {
            "city_id": "181",
            "id": "1619",
            "name": "来凤县"
        },
        {
            "city_id": "181",
            "id": "1620",
            "name": "鹤峰县"
        },
        {
            "city_id": "182",
            "id": "1621",
            "name": "仙桃市"
        },
        {
            "city_id": "182",
            "id": "1622",
            "name": "潜江市"
        },
        {
            "city_id": "182",
            "id": "1623",
            "name": "天门市"
        },
        {
            "city_id": "182",
            "id": "1624",
            "name": "神农架林区"
        },
        {
            "city_id": "183",
            "id": "1625",
            "name": "芙蓉区"
        },
        {
            "city_id": "183",
            "id": "1626",
            "name": "天心区"
        },
        {
            "city_id": "183",
            "id": "1627",
            "name": "岳麓区"
        },
        {
            "city_id": "183",
            "id": "1628",
            "name": "开福区"
        },
        {
            "city_id": "183",
            "id": "1629",
            "name": "雨花区"
        },
        {
            "city_id": "183",
            "id": "1630",
            "name": "长沙县"
        },
        {
            "city_id": "183",
            "id": "1631",
            "name": "望城区"
        },
        {
            "city_id": "183",
            "id": "1632",
            "name": "宁乡市"
        },
        {
            "city_id": "183",
            "id": "1633",
            "name": "浏阳市"
        },
        {
            "city_id": "184",
            "id": "1634",
            "name": "荷塘区"
        },
        {
            "city_id": "184",
            "id": "1635",
            "name": "芦淞区"
        },
        {
            "city_id": "184",
            "id": "1636",
            "name": "石峰区"
        },
        {
            "city_id": "184",
            "id": "1637",
            "name": "天元区"
        },
        {
            "city_id": "184",
            "id": "1638",
            "name": "株洲县"
        },
        {
            "city_id": "184",
            "id": "1639",
            "name": "攸县"
        },
        {
            "city_id": "184",
            "id": "1640",
            "name": "茶陵县"
        },
        {
            "city_id": "184",
            "id": "1641",
            "name": "炎陵县"
        },
        {
            "city_id": "184",
            "id": "1642",
            "name": "醴陵市"
        },
        {
            "city_id": "185",
            "id": "1643",
            "name": "雨湖区"
        },
        {
            "city_id": "185",
            "id": "1644",
            "name": "岳塘区"
        },
        {
            "city_id": "185",
            "id": "1645",
            "name": "湘潭县"
        },
        {
            "city_id": "185",
            "id": "1646",
            "name": "湘乡市"
        },
        {
            "city_id": "185",
            "id": "1647",
            "name": "韶山市"
        },
        {
            "city_id": "186",
            "id": "1648",
            "name": "珠晖区"
        },
        {
            "city_id": "186",
            "id": "1649",
            "name": "雁峰区"
        },
        {
            "city_id": "186",
            "id": "1650",
            "name": "石鼓区"
        },
        {
            "city_id": "186",
            "id": "1651",
            "name": "蒸湘区"
        },
        {
            "city_id": "186",
            "id": "1652",
            "name": "南岳区"
        },
        {
            "city_id": "186",
            "id": "1653",
            "name": "衡阳县"
        },
        {
            "city_id": "186",
            "id": "1654",
            "name": "衡南县"
        },
        {
            "city_id": "186",
            "id": "1655",
            "name": "衡山县"
        },
        {
            "city_id": "186",
            "id": "1656",
            "name": "衡东县"
        },
        {
            "city_id": "186",
            "id": "1657",
            "name": "祁东县"
        },
        {
            "city_id": "186",
            "id": "1658",
            "name": "耒阳市"
        },
        {
            "city_id": "186",
            "id": "1659",
            "name": "常宁市"
        },
        {
            "city_id": "187",
            "id": "1660",
            "name": "双清区"
        },
        {
            "city_id": "187",
            "id": "1661",
            "name": "大祥区"
        },
        {
            "city_id": "187",
            "id": "1662",
            "name": "北塔区"
        },
        {
            "city_id": "187",
            "id": "1663",
            "name": "邵东县"
        },
        {
            "city_id": "187",
            "id": "1664",
            "name": "新邵县"
        },
        {
            "city_id": "187",
            "id": "1665",
            "name": "邵阳县"
        },
        {
            "city_id": "187",
            "id": "1666",
            "name": "隆回县"
        },
        {
            "city_id": "187",
            "id": "1667",
            "name": "洞口县"
        },
        {
            "city_id": "187",
            "id": "1668",
            "name": "绥宁县"
        },
        {
            "city_id": "187",
            "id": "1669",
            "name": "新宁县"
        },
        {
            "city_id": "187",
            "id": "1670",
            "name": "城步苗族自治县"
        },
        {
            "city_id": "187",
            "id": "1671",
            "name": "武冈市"
        },
        {
            "city_id": "188",
            "id": "1672",
            "name": "岳阳楼区"
        },
        {
            "city_id": "188",
            "id": "1673",
            "name": "云溪区"
        },
        {
            "city_id": "188",
            "id": "1674",
            "name": "君山区"
        },
        {
            "city_id": "188",
            "id": "1675",
            "name": "岳阳县"
        },
        {
            "city_id": "188",
            "id": "1676",
            "name": "华容县"
        },
        {
            "city_id": "188",
            "id": "1677",
            "name": "湘阴县"
        },
        {
            "city_id": "188",
            "id": "1678",
            "name": "平江县"
        },
        {
            "city_id": "188",
            "id": "1679",
            "name": "汨罗市"
        },
        {
            "city_id": "188",
            "id": "1680",
            "name": "临湘市"
        },
        {
            "city_id": "189",
            "id": "1681",
            "name": "武陵区"
        },
        {
            "city_id": "189",
            "id": "1682",
            "name": "鼎城区"
        },
        {
            "city_id": "189",
            "id": "1683",
            "name": "安乡县"
        },
        {
            "city_id": "189",
            "id": "1684",
            "name": "汉寿县"
        },
        {
            "city_id": "189",
            "id": "1685",
            "name": "澧县"
        },
        {
            "city_id": "189",
            "id": "1686",
            "name": "临澧县"
        },
        {
            "city_id": "189",
            "id": "1687",
            "name": "桃源县"
        },
        {
            "city_id": "189",
            "id": "1688",
            "name": "石门县"
        },
        {
            "city_id": "189",
            "id": "1689",
            "name": "津市市"
        },
        {
            "city_id": "190",
            "id": "1690",
            "name": "永定区"
        },
        {
            "city_id": "190",
            "id": "1691",
            "name": "武陵源区"
        },
        {
            "city_id": "190",
            "id": "1692",
            "name": "慈利县"
        },
        {
            "city_id": "190",
            "id": "1693",
            "name": "桑植县"
        },
        {
            "city_id": "191",
            "id": "1694",
            "name": "资阳区"
        },
        {
            "city_id": "191",
            "id": "1695",
            "name": "赫山区"
        },
        {
            "city_id": "191",
            "id": "1696",
            "name": "南县"
        },
        {
            "city_id": "191",
            "id": "1697",
            "name": "桃江县"
        },
        {
            "city_id": "191",
            "id": "1698",
            "name": "安化县"
        },
        {
            "city_id": "191",
            "id": "1699",
            "name": "沅江市"
        },
        {
            "city_id": "192",
            "id": "1700",
            "name": "北湖区"
        },
        {
            "city_id": "192",
            "id": "1701",
            "name": "苏仙区"
        },
        {
            "city_id": "192",
            "id": "1702",
            "name": "桂阳县"
        },
        {
            "city_id": "192",
            "id": "1703",
            "name": "宜章县"
        },
        {
            "city_id": "192",
            "id": "1704",
            "name": "永兴县"
        },
        {
            "city_id": "192",
            "id": "1705",
            "name": "嘉禾县"
        },
        {
            "city_id": "192",
            "id": "1706",
            "name": "临武县"
        },
        {
            "city_id": "192",
            "id": "1707",
            "name": "汝城县"
        },
        {
            "city_id": "192",
            "id": "1708",
            "name": "桂东县"
        },
        {
            "city_id": "192",
            "id": "1709",
            "name": "安仁县"
        },
        {
            "city_id": "192",
            "id": "1710",
            "name": "资兴市"
        },
        {
            "city_id": "193",
            "id": "1711",
            "name": "芝山区"
        },
        {
            "city_id": "193",
            "id": "1712",
            "name": "冷水滩区"
        },
        {
            "city_id": "193",
            "id": "1713",
            "name": "祁阳县"
        },
        {
            "city_id": "193",
            "id": "1714",
            "name": "东安县"
        },
        {
            "city_id": "193",
            "id": "1715",
            "name": "双牌县"
        },
        {
            "city_id": "193",
            "id": "1716",
            "name": "道县"
        },
        {
            "city_id": "193",
            "id": "1717",
            "name": "江永县"
        },
        {
            "city_id": "193",
            "id": "1718",
            "name": "宁远县"
        },
        {
            "city_id": "193",
            "id": "1719",
            "name": "蓝山县"
        },
        {
            "city_id": "193",
            "id": "1720",
            "name": "新田县"
        },
        {
            "city_id": "193",
            "id": "1721",
            "name": "江华瑶族自治县"
        },
        {
            "city_id": "194",
            "id": "1722",
            "name": "鹤城区"
        },
        {
            "city_id": "194",
            "id": "1723",
            "name": "中方县"
        },
        {
            "city_id": "194",
            "id": "1724",
            "name": "沅陵县"
        },
        {
            "city_id": "194",
            "id": "1725",
            "name": "辰溪县"
        },
        {
            "city_id": "194",
            "id": "1726",
            "name": "溆浦县"
        },
        {
            "city_id": "194",
            "id": "1727",
            "name": "会同县"
        },
        {
            "city_id": "194",
            "id": "1728",
            "name": "麻阳苗族自治县"
        },
        {
            "city_id": "194",
            "id": "1729",
            "name": "新晃侗族自治县"
        },
        {
            "city_id": "194",
            "id": "1730",
            "name": "芷江侗族自治县"
        },
        {
            "city_id": "194",
            "id": "1731",
            "name": "靖州苗族侗族自治县"
        },
        {
            "city_id": "194",
            "id": "1732",
            "name": "通道侗族自治县"
        },
        {
            "city_id": "194",
            "id": "1733",
            "name": "洪江市"
        },
        {
            "city_id": "195",
            "id": "1734",
            "name": "娄星区"
        },
        {
            "city_id": "195",
            "id": "1735",
            "name": "双峰县"
        },
        {
            "city_id": "195",
            "id": "1736",
            "name": "新化县"
        },
        {
            "city_id": "195",
            "id": "1737",
            "name": "冷水江市"
        },
        {
            "city_id": "195",
            "id": "1738",
            "name": "涟源市"
        },
        {
            "city_id": "196",
            "id": "1739",
            "name": "吉首市"
        },
        {
            "city_id": "196",
            "id": "1740",
            "name": "泸溪县"
        },
        {
            "city_id": "196",
            "id": "1741",
            "name": "凤凰县"
        },
        {
            "city_id": "196",
            "id": "1742",
            "name": "花垣县"
        },
        {
            "city_id": "196",
            "id": "1743",
            "name": "保靖县"
        },
        {
            "city_id": "196",
            "id": "1744",
            "name": "古丈县"
        },
        {
            "city_id": "196",
            "id": "1745",
            "name": "永顺县"
        },
        {
            "city_id": "196",
            "id": "1746",
            "name": "龙山县"
        },
        {
            "city_id": "197",
            "id": "1747",
            "name": "东山区"
        },
        {
            "city_id": "197",
            "id": "1748",
            "name": "荔湾区"
        },
        {
            "city_id": "197",
            "id": "1749",
            "name": "越秀区"
        },
        {
            "city_id": "197",
            "id": "1750",
            "name": "海珠区"
        },
        {
            "city_id": "197",
            "id": "1751",
            "name": "天河区"
        },
        {
            "city_id": "197",
            "id": "1752",
            "name": "芳村区"
        },
        {
            "city_id": "197",
            "id": "1753",
            "name": "白云区"
        },
        {
            "city_id": "197",
            "id": "1754",
            "name": "黄埔区"
        },
        {
            "city_id": "197",
            "id": "1755",
            "name": "番禺区"
        },
        {
            "city_id": "197",
            "id": "1756",
            "name": "花都区"
        },
        {
            "city_id": "197",
            "id": "1757",
            "name": "增城区"
        },
        {
            "city_id": "197",
            "id": "1758",
            "name": "从化区"
        },
        {
            "city_id": "198",
            "id": "1759",
            "name": "武江区"
        },
        {
            "city_id": "198",
            "id": "1760",
            "name": "浈江区"
        },
        {
            "city_id": "198",
            "id": "1761",
            "name": "曲江区"
        },
        {
            "city_id": "198",
            "id": "1762",
            "name": "始兴县"
        },
        {
            "city_id": "198",
            "id": "1763",
            "name": "仁化县"
        },
        {
            "city_id": "198",
            "id": "1764",
            "name": "翁源县"
        },
        {
            "city_id": "198",
            "id": "1765",
            "name": "乳源瑶族自治县"
        },
        {
            "city_id": "198",
            "id": "1766",
            "name": "新丰县"
        },
        {
            "city_id": "198",
            "id": "1767",
            "name": "乐昌市"
        },
        {
            "city_id": "198",
            "id": "1768",
            "name": "南雄市"
        },
        {
            "city_id": "199",
            "id": "1769",
            "name": "罗湖区"
        },
        {
            "city_id": "199",
            "id": "1770",
            "name": "福田区"
        },
        {
            "city_id": "199",
            "id": "1771",
            "name": "南山区"
        },
        {
            "city_id": "199",
            "id": "1772",
            "name": "宝安区"
        },
        {
            "city_id": "199",
            "id": "1773",
            "name": "龙岗区"
        },
        {
            "city_id": "199",
            "id": "1774",
            "name": "盐田区"
        },
        {
            "city_id": "200",
            "id": "1775",
            "name": "香洲区"
        },
        {
            "city_id": "200",
            "id": "1776",
            "name": "斗门区"
        },
        {
            "city_id": "200",
            "id": "1777",
            "name": "金湾区"
        },
        {
            "city_id": "201",
            "id": "1778",
            "name": "龙湖区"
        },
        {
            "city_id": "201",
            "id": "1779",
            "name": "金平区"
        },
        {
            "city_id": "201",
            "id": "1780",
            "name": "濠江区"
        },
        {
            "city_id": "201",
            "id": "1781",
            "name": "潮阳区"
        },
        {
            "city_id": "201",
            "id": "1782",
            "name": "潮南区"
        },
        {
            "city_id": "201",
            "id": "1783",
            "name": "澄海区"
        },
        {
            "city_id": "201",
            "id": "1784",
            "name": "南澳县"
        },
        {
            "city_id": "202",
            "id": "1785",
            "name": "禅城区"
        },
        {
            "city_id": "202",
            "id": "1786",
            "name": "南海区"
        },
        {
            "city_id": "202",
            "id": "1787",
            "name": "顺德区"
        },
        {
            "city_id": "202",
            "id": "1788",
            "name": "三水区"
        },
        {
            "city_id": "202",
            "id": "1789",
            "name": "高明区"
        },
        {
            "city_id": "203",
            "id": "1790",
            "name": "蓬江区"
        },
        {
            "city_id": "203",
            "id": "1791",
            "name": "江海区"
        },
        {
            "city_id": "203",
            "id": "1792",
            "name": "新会区"
        },
        {
            "city_id": "203",
            "id": "1793",
            "name": "台山市"
        },
        {
            "city_id": "203",
            "id": "1794",
            "name": "开平市"
        },
        {
            "city_id": "203",
            "id": "1795",
            "name": "鹤山市"
        },
        {
            "city_id": "203",
            "id": "1796",
            "name": "恩平市"
        },
        {
            "city_id": "204",
            "id": "1797",
            "name": "赤坎区"
        },
        {
            "city_id": "204",
            "id": "1798",
            "name": "霞山区"
        },
        {
            "city_id": "204",
            "id": "1799",
            "name": "坡头区"
        },
        {
            "city_id": "204",
            "id": "1800",
            "name": "麻章区"
        },
        {
            "city_id": "204",
            "id": "1801",
            "name": "遂溪县"
        },
        {
            "city_id": "204",
            "id": "1802",
            "name": "徐闻县"
        },
        {
            "city_id": "204",
            "id": "1803",
            "name": "廉江市"
        },
        {
            "city_id": "204",
            "id": "1804",
            "name": "雷州市"
        },
        {
            "city_id": "204",
            "id": "1805",
            "name": "吴川市"
        },
        {
            "city_id": "205",
            "id": "1806",
            "name": "茂南区"
        },
        {
            "city_id": "205",
            "id": "1807",
            "name": "茂港区"
        },
        {
            "city_id": "205",
            "id": "1808",
            "name": "电白区"
        },
        {
            "city_id": "205",
            "id": "1809",
            "name": "高州市"
        },
        {
            "city_id": "205",
            "id": "1810",
            "name": "化州市"
        },
        {
            "city_id": "205",
            "id": "1811",
            "name": "信宜市"
        },
        {
            "city_id": "206",
            "id": "1812",
            "name": "端州区"
        },
        {
            "city_id": "206",
            "id": "1813",
            "name": "鼎湖区"
        },
        {
            "city_id": "206",
            "id": "1814",
            "name": "广宁县"
        },
        {
            "city_id": "206",
            "id": "1815",
            "name": "怀集县"
        },
        {
            "city_id": "206",
            "id": "1816",
            "name": "封开县"
        },
        {
            "city_id": "206",
            "id": "1817",
            "name": "德庆县"
        },
        {
            "city_id": "206",
            "id": "1818",
            "name": "高要区"
        },
        {
            "city_id": "206",
            "id": "1819",
            "name": "四会市"
        },
        {
            "city_id": "207",
            "id": "1820",
            "name": "惠城区"
        },
        {
            "city_id": "207",
            "id": "1821",
            "name": "惠阳区"
        },
        {
            "city_id": "207",
            "id": "1822",
            "name": "博罗县"
        },
        {
            "city_id": "207",
            "id": "1823",
            "name": "惠东县"
        },
        {
            "city_id": "207",
            "id": "1824",
            "name": "龙门县"
        },
        {
            "city_id": "208",
            "id": "1825",
            "name": "梅江区"
        },
        {
            "city_id": "208",
            "id": "1826",
            "name": "梅县区"
        },
        {
            "city_id": "208",
            "id": "1827",
            "name": "大埔县"
        },
        {
            "city_id": "208",
            "id": "1828",
            "name": "丰顺县"
        },
        {
            "city_id": "208",
            "id": "1829",
            "name": "五华县"
        },
        {
            "city_id": "208",
            "id": "1830",
            "name": "平远县"
        },
        {
            "city_id": "208",
            "id": "1831",
            "name": "蕉岭县"
        },
        {
            "city_id": "208",
            "id": "1832",
            "name": "兴宁市"
        },
        {
            "city_id": "209",
            "id": "1833",
            "name": "城区"
        },
        {
            "city_id": "209",
            "id": "1834",
            "name": "海丰县"
        },
        {
            "city_id": "209",
            "id": "1835",
            "name": "陆河县"
        },
        {
            "city_id": "209",
            "id": "1836",
            "name": "陆丰市"
        },
        {
            "city_id": "210",
            "id": "1837",
            "name": "源城区"
        },
        {
            "city_id": "210",
            "id": "1838",
            "name": "紫金县"
        },
        {
            "city_id": "210",
            "id": "1839",
            "name": "龙川县"
        },
        {
            "city_id": "210",
            "id": "1840",
            "name": "连平县"
        },
        {
            "city_id": "210",
            "id": "1841",
            "name": "和平县"
        },
        {
            "city_id": "210",
            "id": "1842",
            "name": "东源县"
        },
        {
            "city_id": "211",
            "id": "1843",
            "name": "江城区"
        },
        {
            "city_id": "211",
            "id": "1844",
            "name": "阳西县"
        },
        {
            "city_id": "211",
            "id": "1845",
            "name": "阳东区"
        },
        {
            "city_id": "211",
            "id": "1846",
            "name": "阳春市"
        },
        {
            "city_id": "212",
            "id": "1847",
            "name": "清城区"
        },
        {
            "city_id": "212",
            "id": "1848",
            "name": "佛冈县"
        },
        {
            "city_id": "212",
            "id": "1849",
            "name": "阳山县"
        },
        {
            "city_id": "212",
            "id": "1850",
            "name": "连山壮族瑶族自治县"
        },
        {
            "city_id": "212",
            "id": "1851",
            "name": "连南瑶族自治县"
        },
        {
            "city_id": "212",
            "id": "1852",
            "name": "清新区"
        },
        {
            "city_id": "212",
            "id": "1853",
            "name": "英德市"
        },
        {
            "city_id": "212",
            "id": "1854",
            "name": "连州市"
        },
        {
            "city_id": "215",
            "id": "1855",
            "name": "湘桥区"
        },
        {
            "city_id": "215",
            "id": "1856",
            "name": "潮安区"
        },
        {
            "city_id": "215",
            "id": "1857",
            "name": "饶平县"
        },
        {
            "city_id": "216",
            "id": "1858",
            "name": "榕城区"
        },
        {
            "city_id": "216",
            "id": "1859",
            "name": "揭东区"
        },
        {
            "city_id": "216",
            "id": "1860",
            "name": "揭西县"
        },
        {
            "city_id": "216",
            "id": "1861",
            "name": "惠来县"
        },
        {
            "city_id": "216",
            "id": "1862",
            "name": "普宁市"
        },
        {
            "city_id": "217",
            "id": "1863",
            "name": "云城区"
        },
        {
            "city_id": "217",
            "id": "1864",
            "name": "新兴县"
        },
        {
            "city_id": "217",
            "id": "1865",
            "name": "郁南县"
        },
        {
            "city_id": "217",
            "id": "1866",
            "name": "云安区"
        },
        {
            "city_id": "217",
            "id": "1867",
            "name": "罗定市"
        },
        {
            "city_id": "218",
            "id": "1868",
            "name": "兴宁区"
        },
        {
            "city_id": "218",
            "id": "1869",
            "name": "青秀区"
        },
        {
            "city_id": "218",
            "id": "1870",
            "name": "江南区"
        },
        {
            "city_id": "218",
            "id": "1871",
            "name": "西乡塘区"
        },
        {
            "city_id": "218",
            "id": "1872",
            "name": "良庆区"
        },
        {
            "city_id": "218",
            "id": "1873",
            "name": "邕宁区"
        },
        {
            "city_id": "218",
            "id": "1874",
            "name": "武鸣区"
        },
        {
            "city_id": "218",
            "id": "1875",
            "name": "隆安县"
        },
        {
            "city_id": "218",
            "id": "1876",
            "name": "马山县"
        },
        {
            "city_id": "218",
            "id": "1877",
            "name": "上林县"
        },
        {
            "city_id": "218",
            "id": "1878",
            "name": "宾阳县"
        },
        {
            "city_id": "218",
            "id": "1879",
            "name": "横县"
        },
        {
            "city_id": "219",
            "id": "1880",
            "name": "城中区"
        },
        {
            "city_id": "219",
            "id": "1881",
            "name": "鱼峰区"
        },
        {
            "city_id": "219",
            "id": "1882",
            "name": "柳南区"
        },
        {
            "city_id": "219",
            "id": "1883",
            "name": "柳北区"
        },
        {
            "city_id": "219",
            "id": "1884",
            "name": "柳江区"
        },
        {
            "city_id": "219",
            "id": "1885",
            "name": "柳城县"
        },
        {
            "city_id": "219",
            "id": "1886",
            "name": "鹿寨县"
        },
        {
            "city_id": "219",
            "id": "1887",
            "name": "融安县"
        },
        {
            "city_id": "219",
            "id": "1888",
            "name": "融水苗族自治县"
        },
        {
            "city_id": "219",
            "id": "1889",
            "name": "三江侗族自治县"
        },
        {
            "city_id": "220",
            "id": "1890",
            "name": "秀峰区"
        },
        {
            "city_id": "220",
            "id": "1891",
            "name": "叠彩区"
        },
        {
            "city_id": "220",
            "id": "1892",
            "name": "象山区"
        },
        {
            "city_id": "220",
            "id": "1893",
            "name": "七星区"
        },
        {
            "city_id": "220",
            "id": "1894",
            "name": "雁山区"
        },
        {
            "city_id": "220",
            "id": "1895",
            "name": "阳朔县"
        },
        {
            "city_id": "220",
            "id": "1896",
            "name": "临桂区"
        },
        {
            "city_id": "220",
            "id": "1897",
            "name": "灵川县"
        },
        {
            "city_id": "220",
            "id": "1898",
            "name": "全州县"
        },
        {
            "city_id": "220",
            "id": "1899",
            "name": "兴安县"
        },
        {
            "city_id": "220",
            "id": "1900",
            "name": "永福县"
        },
        {
            "city_id": "220",
            "id": "1901",
            "name": "灌阳县"
        },
        {
            "city_id": "220",
            "id": "1902",
            "name": "龙胜各族自治县"
        },
        {
            "city_id": "220",
            "id": "1903",
            "name": "资源县"
        },
        {
            "city_id": "220",
            "id": "1904",
            "name": "平乐县"
        },
        {
            "city_id": "220",
            "id": "1905",
            "name": "荔蒲县"
        },
        {
            "city_id": "220",
            "id": "1906",
            "name": "恭城瑶族自治县"
        },
        {
            "city_id": "221",
            "id": "1907",
            "name": "万秀区"
        },
        {
            "city_id": "221",
            "id": "1908",
            "name": "蝶山区"
        },
        {
            "city_id": "221",
            "id": "1909",
            "name": "长洲区"
        },
        {
            "city_id": "221",
            "id": "1910",
            "name": "苍梧县"
        },
        {
            "city_id": "221",
            "id": "1911",
            "name": "藤县"
        },
        {
            "city_id": "221",
            "id": "1912",
            "name": "蒙山县"
        },
        {
            "city_id": "221",
            "id": "1913",
            "name": "岑溪市"
        },
        {
            "city_id": "222",
            "id": "1914",
            "name": "海城区"
        },
        {
            "city_id": "222",
            "id": "1915",
            "name": "银海区"
        },
        {
            "city_id": "222",
            "id": "1916",
            "name": "铁山港区"
        },
        {
            "city_id": "222",
            "id": "1917",
            "name": "合浦县"
        },
        {
            "city_id": "223",
            "id": "1918",
            "name": "港口区"
        },
        {
            "city_id": "223",
            "id": "1919",
            "name": "防城区"
        },
        {
            "city_id": "223",
            "id": "1920",
            "name": "上思县"
        },
        {
            "city_id": "223",
            "id": "1921",
            "name": "东兴市"
        },
        {
            "city_id": "224",
            "id": "1922",
            "name": "钦南区"
        },
        {
            "city_id": "224",
            "id": "1923",
            "name": "钦北区"
        },
        {
            "city_id": "224",
            "id": "1924",
            "name": "灵山县"
        },
        {
            "city_id": "224",
            "id": "1925",
            "name": "浦北县"
        },
        {
            "city_id": "225",
            "id": "1926",
            "name": "港北区"
        },
        {
            "city_id": "225",
            "id": "1927",
            "name": "港南区"
        },
        {
            "city_id": "225",
            "id": "1928",
            "name": "覃塘区"
        },
        {
            "city_id": "225",
            "id": "1929",
            "name": "平南县"
        },
        {
            "city_id": "225",
            "id": "1930",
            "name": "桂平市"
        },
        {
            "city_id": "226",
            "id": "1931",
            "name": "玉州区"
        },
        {
            "city_id": "226",
            "id": "1932",
            "name": "容县"
        },
        {
            "city_id": "226",
            "id": "1933",
            "name": "陆川县"
        },
        {
            "city_id": "226",
            "id": "1934",
            "name": "博白县"
        },
        {
            "city_id": "226",
            "id": "1935",
            "name": "兴业县"
        },
        {
            "city_id": "226",
            "id": "1936",
            "name": "北流市"
        },
        {
            "city_id": "227",
            "id": "1937",
            "name": "右江区"
        },
        {
            "city_id": "227",
            "id": "1938",
            "name": "田阳县"
        },
        {
            "city_id": "227",
            "id": "1939",
            "name": "田东县"
        },
        {
            "city_id": "227",
            "id": "1940",
            "name": "平果县"
        },
        {
            "city_id": "227",
            "id": "1941",
            "name": "德保县"
        },
        {
            "city_id": "227",
            "id": "1942",
            "name": "靖西市"
        },
        {
            "city_id": "227",
            "id": "1943",
            "name": "那坡县"
        },
        {
            "city_id": "227",
            "id": "1944",
            "name": "凌云县"
        },
        {
            "city_id": "227",
            "id": "1945",
            "name": "乐业县"
        },
        {
            "city_id": "227",
            "id": "1946",
            "name": "田林县"
        },
        {
            "city_id": "227",
            "id": "1947",
            "name": "西林县"
        },
        {
            "city_id": "227",
            "id": "1948",
            "name": "隆林各族自治县"
        },
        {
            "city_id": "228",
            "id": "1949",
            "name": "八步区"
        },
        {
            "city_id": "228",
            "id": "1950",
            "name": "昭平县"
        },
        {
            "city_id": "228",
            "id": "1951",
            "name": "钟山县"
        },
        {
            "city_id": "228",
            "id": "1952",
            "name": "富川瑶族自治县"
        },
        {
            "city_id": "229",
            "id": "1953",
            "name": "金城江区"
        },
        {
            "city_id": "229",
            "id": "1954",
            "name": "南丹县"
        },
        {
            "city_id": "229",
            "id": "1955",
            "name": "天峨县"
        },
        {
            "city_id": "229",
            "id": "1956",
            "name": "凤山县"
        },
        {
            "city_id": "229",
            "id": "1957",
            "name": "东兰县"
        },
        {
            "city_id": "229",
            "id": "1958",
            "name": "罗城仫佬族自治县"
        },
        {
            "city_id": "229",
            "id": "1959",
            "name": "环江毛南族自治县"
        },
        {
            "city_id": "229",
            "id": "1960",
            "name": "巴马瑶族自治县"
        },
        {
            "city_id": "229",
            "id": "1961",
            "name": "都安瑶族自治县"
        },
        {
            "city_id": "229",
            "id": "1962",
            "name": "大化瑶族自治县"
        },
        {
            "city_id": "229",
            "id": "1963",
            "name": "宜州区"
        },
        {
            "city_id": "230",
            "id": "1964",
            "name": "兴宾区"
        },
        {
            "city_id": "230",
            "id": "1965",
            "name": "忻城县"
        },
        {
            "city_id": "230",
            "id": "1966",
            "name": "象州县"
        },
        {
            "city_id": "230",
            "id": "1967",
            "name": "武宣县"
        },
        {
            "city_id": "230",
            "id": "1968",
            "name": "金秀瑶族自治县"
        },
        {
            "city_id": "230",
            "id": "1969",
            "name": "合山市"
        },
        {
            "city_id": "231",
            "id": "1970",
            "name": "江洲区"
        },
        {
            "city_id": "231",
            "id": "1971",
            "name": "扶绥县"
        },
        {
            "city_id": "231",
            "id": "1972",
            "name": "宁明县"
        },
        {
            "city_id": "231",
            "id": "1973",
            "name": "龙州县"
        },
        {
            "city_id": "231",
            "id": "1974",
            "name": "大新县"
        },
        {
            "city_id": "231",
            "id": "1975",
            "name": "天等县"
        },
        {
            "city_id": "231",
            "id": "1976",
            "name": "凭祥市"
        },
        {
            "city_id": "232",
            "id": "1977",
            "name": "秀英区"
        },
        {
            "city_id": "232",
            "id": "1978",
            "name": "龙华区"
        },
        {
            "city_id": "232",
            "id": "1979",
            "name": "琼山区"
        },
        {
            "city_id": "232",
            "id": "1980",
            "name": "美兰区"
        },
        {
            "city_id": "349",
            "id": "1981",
            "name": "五指山市"
        },
        {
            "city_id": "349",
            "id": "1982",
            "name": "琼海市"
        },
        {
            "city_id": "233",
            "id": "1983",
            "name": "儋州市"
        },
        {
            "city_id": "349",
            "id": "1984",
            "name": "文昌市"
        },
        {
            "city_id": "349",
            "id": "1985",
            "name": "万宁市"
        },
        {
            "city_id": "349",
            "id": "1986",
            "name": "东方市"
        },
        {
            "city_id": "349",
            "id": "1987",
            "name": "定安县"
        },
        {
            "city_id": "349",
            "id": "1988",
            "name": "屯昌县"
        },
        {
            "city_id": "349",
            "id": "1989",
            "name": "澄迈县"
        },
        {
            "city_id": "349",
            "id": "1990",
            "name": "临高县"
        },
        {
            "city_id": "349",
            "id": "1991",
            "name": "白沙黎族自治县"
        },
        {
            "city_id": "349",
            "id": "1992",
            "name": "昌江黎族自治县"
        },
        {
            "city_id": "349",
            "id": "1993",
            "name": "乐东黎族自治县"
        },
        {
            "city_id": "349",
            "id": "1994",
            "name": "陵水黎族自治县"
        },
        {
            "city_id": "349",
            "id": "1995",
            "name": "保亭黎族苗族自治县"
        },
        {
            "city_id": "349",
            "id": "1996",
            "name": "琼中黎族苗族自治县"
        },
        {
            "city_id": "233",
            "id": "1997",
            "name": "西沙群岛"
        },
        {
            "city_id": "233",
            "id": "1998",
            "name": "南沙群岛"
        },
        {
            "city_id": "233",
            "id": "1999",
            "name": "中沙群岛的岛礁及其海域"
        },
        {
            "city_id": "234",
            "id": "2000",
            "name": "万州区"
        },
        {
            "city_id": "234",
            "id": "2001",
            "name": "涪陵区"
        },
        {
            "city_id": "234",
            "id": "2002",
            "name": "渝中区"
        },
        {
            "city_id": "234",
            "id": "2003",
            "name": "大渡口区"
        },
        {
            "city_id": "234",
            "id": "2004",
            "name": "江北区"
        },
        {
            "city_id": "234",
            "id": "2005",
            "name": "沙坪坝区"
        },
        {
            "city_id": "234",
            "id": "2006",
            "name": "九龙坡区"
        },
        {
            "city_id": "234",
            "id": "2007",
            "name": "南岸区"
        },
        {
            "city_id": "234",
            "id": "2008",
            "name": "北碚区"
        },
        {
            "city_id": "234",
            "id": "2009",
            "name": "万盛区"
        },
        {
            "city_id": "234",
            "id": "2010",
            "name": "双桥区"
        },
        {
            "city_id": "234",
            "id": "2011",
            "name": "渝北区"
        },
        {
            "city_id": "234",
            "id": "2012",
            "name": "巴南区"
        },
        {
            "city_id": "234",
            "id": "2013",
            "name": "黔江区"
        },
        {
            "city_id": "234",
            "id": "2014",
            "name": "长寿区"
        },
        {
            "city_id": "234",
            "id": "2015",
            "name": "綦江区"
        },
        {
            "city_id": "234",
            "id": "2016",
            "name": "潼南区"
        },
        {
            "city_id": "234",
            "id": "2017",
            "name": "铜梁区"
        },
        {
            "city_id": "234",
            "id": "2018",
            "name": "大足区"
        },
        {
            "city_id": "234",
            "id": "2019",
            "name": "荣昌区"
        },
        {
            "city_id": "234",
            "id": "2020",
            "name": "璧山区"
        },
        {
            "city_id": "234",
            "id": "2021",
            "name": "梁平区"
        },
        {
            "city_id": "348",
            "id": "2022",
            "name": "城口县"
        },
        {
            "city_id": "348",
            "id": "2023",
            "name": "丰都县"
        },
        {
            "city_id": "348",
            "id": "2024",
            "name": "垫江县"
        },
        {
            "city_id": "234",
            "id": "2025",
            "name": "武隆区"
        },
        {
            "city_id": "348",
            "id": "2026",
            "name": "忠县"
        },
        {
            "city_id": "234",
            "id": "2027",
            "name": "开县"
        },
        {
            "city_id": "348",
            "id": "2028",
            "name": "云阳县"
        },
        {
            "city_id": "348",
            "id": "2029",
            "name": "奉节县"
        },
        {
            "city_id": "348",
            "id": "2030",
            "name": "巫山县"
        },
        {
            "city_id": "348",
            "id": "2031",
            "name": "巫溪县"
        },
        {
            "city_id": "348",
            "id": "2032",
            "name": "石柱土家族自治县"
        },
        {
            "city_id": "348",
            "id": "2033",
            "name": "秀山土家族苗族自治县"
        },
        {
            "city_id": "348",
            "id": "2034",
            "name": "酉阳土家族苗族自治县"
        },
        {
            "city_id": "348",
            "id": "2035",
            "name": "彭水苗族土家族自治县"
        },
        {
            "city_id": "234",
            "id": "2036",
            "name": "江津区"
        },
        {
            "city_id": "234",
            "id": "2037",
            "name": "合川区"
        },
        {
            "city_id": "234",
            "id": "2038",
            "name": "永川区"
        },
        {
            "city_id": "234",
            "id": "2039",
            "name": "南川区"
        },
        {
            "city_id": "235",
            "id": "2040",
            "name": "锦江区"
        },
        {
            "city_id": "235",
            "id": "2041",
            "name": "青羊区"
        },
        {
            "city_id": "235",
            "id": "2042",
            "name": "金牛区"
        },
        {
            "city_id": "235",
            "id": "2043",
            "name": "武侯区"
        },
        {
            "city_id": "235",
            "id": "2044",
            "name": "成华区"
        },
        {
            "city_id": "235",
            "id": "2045",
            "name": "龙泉驿区"
        },
        {
            "city_id": "235",
            "id": "2046",
            "name": "青白江区"
        },
        {
            "city_id": "235",
            "id": "2047",
            "name": "新都区"
        },
        {
            "city_id": "235",
            "id": "2048",
            "name": "温江区"
        },
        {
            "city_id": "235",
            "id": "2049",
            "name": "金堂县"
        },
        {
            "city_id": "235",
            "id": "2050",
            "name": "双流区"
        },
        {
            "city_id": "235",
            "id": "2051",
            "name": "郫县"
        },
        {
            "city_id": "235",
            "id": "2052",
            "name": "大邑县"
        },
        {
            "city_id": "235",
            "id": "2053",
            "name": "蒲江县"
        },
        {
            "city_id": "235",
            "id": "2054",
            "name": "新津县"
        },
        {
            "city_id": "235",
            "id": "2055",
            "name": "都江堰市"
        },
        {
            "city_id": "235",
            "id": "2056",
            "name": "彭州市"
        },
        {
            "city_id": "235",
            "id": "2057",
            "name": "邛崃市"
        },
        {
            "city_id": "235",
            "id": "2058",
            "name": "崇州市"
        },
        {
            "city_id": "236",
            "id": "2059",
            "name": "自流井区"
        },
        {
            "city_id": "236",
            "id": "2060",
            "name": "贡井区"
        },
        {
            "city_id": "236",
            "id": "2061",
            "name": "大安区"
        },
        {
            "city_id": "236",
            "id": "2062",
            "name": "沿滩区"
        },
        {
            "city_id": "236",
            "id": "2063",
            "name": "荣县"
        },
        {
            "city_id": "236",
            "id": "2064",
            "name": "富顺县"
        },
        {
            "city_id": "237",
            "id": "2065",
            "name": "东区"
        },
        {
            "city_id": "237",
            "id": "2066",
            "name": "西区"
        },
        {
            "city_id": "237",
            "id": "2067",
            "name": "仁和区"
        },
        {
            "city_id": "237",
            "id": "2068",
            "name": "米易县"
        },
        {
            "city_id": "237",
            "id": "2069",
            "name": "盐边县"
        },
        {
            "city_id": "238",
            "id": "2070",
            "name": "江阳区"
        },
        {
            "city_id": "238",
            "id": "2071",
            "name": "纳溪区"
        },
        {
            "city_id": "238",
            "id": "2072",
            "name": "龙马潭区"
        },
        {
            "city_id": "238",
            "id": "2073",
            "name": "泸县"
        },
        {
            "city_id": "238",
            "id": "2074",
            "name": "合江县"
        },
        {
            "city_id": "238",
            "id": "2075",
            "name": "叙永县"
        },
        {
            "city_id": "238",
            "id": "2076",
            "name": "古蔺县"
        },
        {
            "city_id": "239",
            "id": "2077",
            "name": "旌阳区"
        },
        {
            "city_id": "239",
            "id": "2078",
            "name": "中江县"
        },
        {
            "city_id": "239",
            "id": "2079",
            "name": "罗江区"
        },
        {
            "city_id": "239",
            "id": "2080",
            "name": "广汉市"
        },
        {
            "city_id": "239",
            "id": "2081",
            "name": "什邡市"
        },
        {
            "city_id": "239",
            "id": "2082",
            "name": "绵竹市"
        },
        {
            "city_id": "240",
            "id": "2083",
            "name": "涪城区"
        },
        {
            "city_id": "240",
            "id": "2084",
            "name": "游仙区"
        },
        {
            "city_id": "240",
            "id": "2085",
            "name": "三台县"
        },
        {
            "city_id": "240",
            "id": "2086",
            "name": "盐亭县"
        },
        {
            "city_id": "240",
            "id": "2087",
            "name": "安县"
        },
        {
            "city_id": "240",
            "id": "2088",
            "name": "梓潼县"
        },
        {
            "city_id": "240",
            "id": "2089",
            "name": "北川羌族自治县"
        },
        {
            "city_id": "240",
            "id": "2090",
            "name": "平武县"
        },
        {
            "city_id": "240",
            "id": "2091",
            "name": "江油市"
        },
        {
            "city_id": "241",
            "id": "2092",
            "name": "市中区"
        },
        {
            "city_id": "241",
            "id": "2093",
            "name": "元坝区"
        },
        {
            "city_id": "241",
            "id": "2094",
            "name": "朝天区"
        },
        {
            "city_id": "241",
            "id": "2095",
            "name": "旺苍县"
        },
        {
            "city_id": "241",
            "id": "2096",
            "name": "青川县"
        },
        {
            "city_id": "241",
            "id": "2097",
            "name": "剑阁县"
        },
        {
            "city_id": "241",
            "id": "2098",
            "name": "苍溪县"
        },
        {
            "city_id": "242",
            "id": "2099",
            "name": "船山区"
        },
        {
            "city_id": "242",
            "id": "2100",
            "name": "安居区"
        },
        {
            "city_id": "242",
            "id": "2101",
            "name": "蓬溪县"
        },
        {
            "city_id": "242",
            "id": "2102",
            "name": "射洪县"
        },
        {
            "city_id": "242",
            "id": "2103",
            "name": "大英县"
        },
        {
            "city_id": "243",
            "id": "2104",
            "name": "市中区"
        },
        {
            "city_id": "243",
            "id": "2105",
            "name": "东兴区"
        },
        {
            "city_id": "243",
            "id": "2106",
            "name": "威远县"
        },
        {
            "city_id": "243",
            "id": "2107",
            "name": "资中县"
        },
        {
            "city_id": "243",
            "id": "2108",
            "name": "隆昌市"
        },
        {
            "city_id": "244",
            "id": "2109",
            "name": "市中区"
        },
        {
            "city_id": "244",
            "id": "2110",
            "name": "沙湾区"
        },
        {
            "city_id": "244",
            "id": "2111",
            "name": "五通桥区"
        },
        {
            "city_id": "244",
            "id": "2112",
            "name": "金口河区"
        },
        {
            "city_id": "244",
            "id": "2113",
            "name": "犍为县"
        },
        {
            "city_id": "244",
            "id": "2114",
            "name": "井研县"
        },
        {
            "city_id": "244",
            "id": "2115",
            "name": "夹江县"
        },
        {
            "city_id": "244",
            "id": "2116",
            "name": "沐川县"
        },
        {
            "city_id": "244",
            "id": "2117",
            "name": "峨边彝族自治县"
        },
        {
            "city_id": "244",
            "id": "2118",
            "name": "马边彝族自治县"
        },
        {
            "city_id": "244",
            "id": "2119",
            "name": "峨眉山市"
        },
        {
            "city_id": "245",
            "id": "2120",
            "name": "顺庆区"
        },
        {
            "city_id": "245",
            "id": "2121",
            "name": "高坪区"
        },
        {
            "city_id": "245",
            "id": "2122",
            "name": "嘉陵区"
        },
        {
            "city_id": "245",
            "id": "2123",
            "name": "南部县"
        },
        {
            "city_id": "245",
            "id": "2124",
            "name": "营山县"
        },
        {
            "city_id": "245",
            "id": "2125",
            "name": "蓬安县"
        },
        {
            "city_id": "245",
            "id": "2126",
            "name": "仪陇县"
        },
        {
            "city_id": "245",
            "id": "2127",
            "name": "西充县"
        },
        {
            "city_id": "245",
            "id": "2128",
            "name": "阆中市"
        },
        {
            "city_id": "246",
            "id": "2129",
            "name": "东坡区"
        },
        {
            "city_id": "246",
            "id": "2130",
            "name": "仁寿县"
        },
        {
            "city_id": "246",
            "id": "2131",
            "name": "彭山区"
        },
        {
            "city_id": "246",
            "id": "2132",
            "name": "洪雅县"
        },
        {
            "city_id": "246",
            "id": "2133",
            "name": "丹棱县"
        },
        {
            "city_id": "246",
            "id": "2134",
            "name": "青神县"
        },
        {
            "city_id": "247",
            "id": "2135",
            "name": "翠屏区"
        },
        {
            "city_id": "247",
            "id": "2136",
            "name": "宜宾县"
        },
        {
            "city_id": "247",
            "id": "2137",
            "name": "南溪区"
        },
        {
            "city_id": "247",
            "id": "2138",
            "name": "江安县"
        },
        {
            "city_id": "247",
            "id": "2139",
            "name": "长宁县"
        },
        {
            "city_id": "247",
            "id": "2140",
            "name": "高县"
        },
        {
            "city_id": "247",
            "id": "2141",
            "name": "珙县"
        },
        {
            "city_id": "247",
            "id": "2142",
            "name": "筠连县"
        },
        {
            "city_id": "247",
            "id": "2143",
            "name": "兴文县"
        },
        {
            "city_id": "247",
            "id": "2144",
            "name": "屏山县"
        },
        {
            "city_id": "248",
            "id": "2145",
            "name": "广安区"
        },
        {
            "city_id": "248",
            "id": "2146",
            "name": "岳池县"
        },
        {
            "city_id": "248",
            "id": "2147",
            "name": "武胜县"
        },
        {
            "city_id": "248",
            "id": "2148",
            "name": "邻水县"
        },
        {
            "city_id": "248",
            "id": "2149",
            "name": "华蓥市"
        },
        {
            "city_id": "249",
            "id": "2150",
            "name": "通川区"
        },
        {
            "city_id": "249",
            "id": "2151",
            "name": "达县"
        },
        {
            "city_id": "249",
            "id": "2152",
            "name": "宣汉县"
        },
        {
            "city_id": "249",
            "id": "2153",
            "name": "开江县"
        },
        {
            "city_id": "249",
            "id": "2154",
            "name": "大竹县"
        },
        {
            "city_id": "249",
            "id": "2155",
            "name": "渠县"
        },
        {
            "city_id": "249",
            "id": "2156",
            "name": "万源市"
        },
        {
            "city_id": "250",
            "id": "2157",
            "name": "雨城区"
        },
        {
            "city_id": "250",
            "id": "2158",
            "name": "名山区"
        },
        {
            "city_id": "250",
            "id": "2159",
            "name": "荥经县"
        },
        {
            "city_id": "250",
            "id": "2160",
            "name": "汉源县"
        },
        {
            "city_id": "250",
            "id": "2161",
            "name": "石棉县"
        },
        {
            "city_id": "250",
            "id": "2162",
            "name": "天全县"
        },
        {
            "city_id": "250",
            "id": "2163",
            "name": "芦山县"
        },
        {
            "city_id": "250",
            "id": "2164",
            "name": "宝兴县"
        },
        {
            "city_id": "251",
            "id": "2165",
            "name": "巴州区"
        },
        {
            "city_id": "251",
            "id": "2166",
            "name": "通江县"
        },
        {
            "city_id": "251",
            "id": "2167",
            "name": "南江县"
        },
        {
            "city_id": "251",
            "id": "2168",
            "name": "平昌县"
        },
        {
            "city_id": "252",
            "id": "2169",
            "name": "雁江区"
        },
        {
            "city_id": "252",
            "id": "2170",
            "name": "安岳县"
        },
        {
            "city_id": "252",
            "id": "2171",
            "name": "乐至县"
        },
        {
            "city_id": "252",
            "id": "2172",
            "name": "简阳市"
        },
        {
            "city_id": "253",
            "id": "2173",
            "name": "汶川县"
        },
        {
            "city_id": "253",
            "id": "2174",
            "name": "理县"
        },
        {
            "city_id": "253",
            "id": "2175",
            "name": "茂县"
        },
        {
            "city_id": "253",
            "id": "2176",
            "name": "松潘县"
        },
        {
            "city_id": "253",
            "id": "2177",
            "name": "九寨沟县"
        },
        {
            "city_id": "253",
            "id": "2178",
            "name": "金川县"
        },
        {
            "city_id": "253",
            "id": "2179",
            "name": "小金县"
        },
        {
            "city_id": "253",
            "id": "2180",
            "name": "黑水县"
        },
        {
            "city_id": "253",
            "id": "2181",
            "name": "马尔康市"
        },
        {
            "city_id": "253",
            "id": "2182",
            "name": "壤塘县"
        },
        {
            "city_id": "253",
            "id": "2183",
            "name": "阿坝县"
        },
        {
            "city_id": "253",
            "id": "2184",
            "name": "若尔盖县"
        },
        {
            "city_id": "253",
            "id": "2185",
            "name": "红原县"
        },
        {
            "city_id": "254",
            "id": "2186",
            "name": "康定市"
        },
        {
            "city_id": "254",
            "id": "2187",
            "name": "泸定县"
        },
        {
            "city_id": "254",
            "id": "2188",
            "name": "丹巴县"
        },
        {
            "city_id": "254",
            "id": "2189",
            "name": "九龙县"
        },
        {
            "city_id": "254",
            "id": "2190",
            "name": "雅江县"
        },
        {
            "city_id": "254",
            "id": "2191",
            "name": "道孚县"
        },
        {
            "city_id": "254",
            "id": "2192",
            "name": "炉霍县"
        },
        {
            "city_id": "254",
            "id": "2193",
            "name": "甘孜县"
        },
        {
            "city_id": "254",
            "id": "2194",
            "name": "新龙县"
        },
        {
            "city_id": "254",
            "id": "2195",
            "name": "德格县"
        },
        {
            "city_id": "254",
            "id": "2196",
            "name": "白玉县"
        },
        {
            "city_id": "254",
            "id": "2197",
            "name": "石渠县"
        },
        {
            "city_id": "254",
            "id": "2198",
            "name": "色达县"
        },
        {
            "city_id": "254",
            "id": "2199",
            "name": "理塘县"
        },
        {
            "city_id": "254",
            "id": "2200",
            "name": "巴塘县"
        },
        {
            "city_id": "254",
            "id": "2201",
            "name": "乡城县"
        },
        {
            "city_id": "254",
            "id": "2202",
            "name": "稻城县"
        },
        {
            "city_id": "254",
            "id": "2203",
            "name": "得荣县"
        },
        {
            "city_id": "255",
            "id": "2204",
            "name": "西昌市"
        },
        {
            "city_id": "255",
            "id": "2205",
            "name": "木里藏族自治县"
        },
        {
            "city_id": "255",
            "id": "2206",
            "name": "盐源县"
        },
        {
            "city_id": "255",
            "id": "2207",
            "name": "德昌县"
        },
        {
            "city_id": "255",
            "id": "2208",
            "name": "会理县"
        },
        {
            "city_id": "255",
            "id": "2209",
            "name": "会东县"
        },
        {
            "city_id": "255",
            "id": "2210",
            "name": "宁南县"
        },
        {
            "city_id": "255",
            "id": "2211",
            "name": "普格县"
        },
        {
            "city_id": "255",
            "id": "2212",
            "name": "布拖县"
        },
        {
            "city_id": "255",
            "id": "2213",
            "name": "金阳县"
        },
        {
            "city_id": "255",
            "id": "2214",
            "name": "昭觉县"
        },
        {
            "city_id": "255",
            "id": "2215",
            "name": "喜德县"
        },
        {
            "city_id": "255",
            "id": "2216",
            "name": "冕宁县"
        },
        {
            "city_id": "255",
            "id": "2217",
            "name": "越西县"
        },
        {
            "city_id": "255",
            "id": "2218",
            "name": "甘洛县"
        },
        {
            "city_id": "255",
            "id": "2219",
            "name": "美姑县"
        },
        {
            "city_id": "255",
            "id": "2220",
            "name": "雷波县"
        },
        {
            "city_id": "256",
            "id": "2221",
            "name": "南明区"
        },
        {
            "city_id": "256",
            "id": "2222",
            "name": "云岩区"
        },
        {
            "city_id": "256",
            "id": "2223",
            "name": "花溪区"
        },
        {
            "city_id": "256",
            "id": "2224",
            "name": "乌当区"
        },
        {
            "city_id": "256",
            "id": "2225",
            "name": "白云区"
        },
        {
            "city_id": "256",
            "id": "2226",
            "name": "小河区"
        },
        {
            "city_id": "256",
            "id": "2227",
            "name": "开阳县"
        },
        {
            "city_id": "256",
            "id": "2228",
            "name": "息烽县"
        },
        {
            "city_id": "256",
            "id": "2229",
            "name": "修文县"
        },
        {
            "city_id": "256",
            "id": "2230",
            "name": "清镇市"
        },
        {
            "city_id": "257",
            "id": "2231",
            "name": "钟山区"
        },
        {
            "city_id": "257",
            "id": "2232",
            "name": "六枝特区"
        },
        {
            "city_id": "257",
            "id": "2233",
            "name": "水城县"
        },
        {
            "city_id": "257",
            "id": "2234",
            "name": "盘县"
        },
        {
            "city_id": "258",
            "id": "2235",
            "name": "红花岗区"
        },
        {
            "city_id": "258",
            "id": "2236",
            "name": "汇川区"
        },
        {
            "city_id": "258",
            "id": "2237",
            "name": "遵义县"
        },
        {
            "city_id": "258",
            "id": "2238",
            "name": "桐梓县"
        },
        {
            "city_id": "258",
            "id": "2239",
            "name": "绥阳县"
        },
        {
            "city_id": "258",
            "id": "2240",
            "name": "正安县"
        },
        {
            "city_id": "258",
            "id": "2241",
            "name": "道真仡佬族苗族自治县"
        },
        {
            "city_id": "258",
            "id": "2242",
            "name": "务川仡佬族苗族自治县"
        },
        {
            "city_id": "258",
            "id": "2243",
            "name": "凤冈县"
        },
        {
            "city_id": "258",
            "id": "2244",
            "name": "湄潭县"
        },
        {
            "city_id": "258",
            "id": "2245",
            "name": "余庆县"
        },
        {
            "city_id": "258",
            "id": "2246",
            "name": "习水县"
        },
        {
            "city_id": "258",
            "id": "2247",
            "name": "赤水市"
        },
        {
            "city_id": "258",
            "id": "2248",
            "name": "仁怀市"
        },
        {
            "city_id": "259",
            "id": "2249",
            "name": "西秀区"
        },
        {
            "city_id": "259",
            "id": "2250",
            "name": "平坝区"
        },
        {
            "city_id": "259",
            "id": "2251",
            "name": "普定县"
        },
        {
            "city_id": "259",
            "id": "2252",
            "name": "镇宁布依族苗族自治县"
        },
        {
            "city_id": "259",
            "id": "2253",
            "name": "关岭布依族苗族自治县"
        },
        {
            "city_id": "259",
            "id": "2254",
            "name": "紫云苗族布依族自治县"
        },
        {
            "city_id": "260",
            "id": "2255",
            "name": "铜仁市"
        },
        {
            "city_id": "260",
            "id": "2256",
            "name": "江口县"
        },
        {
            "city_id": "260",
            "id": "2257",
            "name": "玉屏侗族自治县"
        },
        {
            "city_id": "260",
            "id": "2258",
            "name": "石阡县"
        },
        {
            "city_id": "260",
            "id": "2259",
            "name": "思南县"
        },
        {
            "city_id": "260",
            "id": "2260",
            "name": "印江土家族苗族自治县"
        },
        {
            "city_id": "260",
            "id": "2261",
            "name": "德江县"
        },
        {
            "city_id": "260",
            "id": "2262",
            "name": "沿河土家族自治县"
        },
        {
            "city_id": "260",
            "id": "2263",
            "name": "松桃苗族自治县"
        },
        {
            "city_id": "260",
            "id": "2264",
            "name": "万山区"
        },
        {
            "city_id": "261",
            "id": "2265",
            "name": "兴义市"
        },
        {
            "city_id": "261",
            "id": "2266",
            "name": "兴仁市"
        },
        {
            "city_id": "261",
            "id": "2267",
            "name": "普安县"
        },
        {
            "city_id": "261",
            "id": "2268",
            "name": "晴隆县"
        },
        {
            "city_id": "261",
            "id": "2269",
            "name": "贞丰县"
        },
        {
            "city_id": "261",
            "id": "2270",
            "name": "望谟县"
        },
        {
            "city_id": "261",
            "id": "2271",
            "name": "册亨县"
        },
        {
            "city_id": "261",
            "id": "2272",
            "name": "安龙县"
        },
        {
            "city_id": "262",
            "id": "2273",
            "name": "毕节市"
        },
        {
            "city_id": "262",
            "id": "2274",
            "name": "大方县"
        },
        {
            "city_id": "262",
            "id": "2275",
            "name": "黔西县"
        },
        {
            "city_id": "262",
            "id": "2276",
            "name": "金沙县"
        },
        {
            "city_id": "262",
            "id": "2277",
            "name": "织金县"
        },
        {
            "city_id": "262",
            "id": "2278",
            "name": "纳雍县"
        },
        {
            "city_id": "262",
            "id": "2279",
            "name": "威宁彝族回族苗族自治县"
        },
        {
            "city_id": "262",
            "id": "2280",
            "name": "赫章县"
        },
        {
            "city_id": "263",
            "id": "2281",
            "name": "凯里市"
        },
        {
            "city_id": "263",
            "id": "2282",
            "name": "黄平县"
        },
        {
            "city_id": "263",
            "id": "2283",
            "name": "施秉县"
        },
        {
            "city_id": "263",
            "id": "2284",
            "name": "三穗县"
        },
        {
            "city_id": "263",
            "id": "2285",
            "name": "镇远县"
        },
        {
            "city_id": "263",
            "id": "2286",
            "name": "岑巩县"
        },
        {
            "city_id": "263",
            "id": "2287",
            "name": "天柱县"
        },
        {
            "city_id": "263",
            "id": "2288",
            "name": "锦屏县"
        },
        {
            "city_id": "263",
            "id": "2289",
            "name": "剑河县"
        },
        {
            "city_id": "263",
            "id": "2290",
            "name": "台江县"
        },
        {
            "city_id": "263",
            "id": "2291",
            "name": "黎平县"
        },
        {
            "city_id": "263",
            "id": "2292",
            "name": "榕江县"
        },
        {
            "city_id": "263",
            "id": "2293",
            "name": "从江县"
        },
        {
            "city_id": "263",
            "id": "2294",
            "name": "雷山县"
        },
        {
            "city_id": "263",
            "id": "2295",
            "name": "麻江县"
        },
        {
            "city_id": "263",
            "id": "2296",
            "name": "丹寨县"
        },
        {
            "city_id": "264",
            "id": "2297",
            "name": "都匀市"
        },
        {
            "city_id": "264",
            "id": "2298",
            "name": "福泉市"
        },
        {
            "city_id": "264",
            "id": "2299",
            "name": "荔波县"
        },
        {
            "city_id": "264",
            "id": "2300",
            "name": "贵定县"
        },
        {
            "city_id": "264",
            "id": "2301",
            "name": "瓮安县"
        },
        {
            "city_id": "264",
            "id": "2302",
            "name": "独山县"
        },
        {
            "city_id": "264",
            "id": "2303",
            "name": "平塘县"
        },
        {
            "city_id": "264",
            "id": "2304",
            "name": "罗甸县"
        },
        {
            "city_id": "264",
            "id": "2305",
            "name": "长顺县"
        },
        {
            "city_id": "264",
            "id": "2306",
            "name": "龙里县"
        },
        {
            "city_id": "264",
            "id": "2307",
            "name": "惠水县"
        },
        {
            "city_id": "264",
            "id": "2308",
            "name": "三都水族自治县"
        },
        {
            "city_id": "265",
            "id": "2309",
            "name": "五华区"
        },
        {
            "city_id": "265",
            "id": "2310",
            "name": "盘龙区"
        },
        {
            "city_id": "265",
            "id": "2311",
            "name": "官渡区"
        },
        {
            "city_id": "265",
            "id": "2312",
            "name": "西山区"
        },
        {
            "city_id": "265",
            "id": "2313",
            "name": "东川区"
        },
        {
            "city_id": "265",
            "id": "2314",
            "name": "呈贡区"
        },
        {
            "city_id": "265",
            "id": "2315",
            "name": "晋宁区"
        },
        {
            "city_id": "265",
            "id": "2316",
            "name": "富民县"
        },
        {
            "city_id": "265",
            "id": "2317",
            "name": "宜良县"
        },
        {
            "city_id": "265",
            "id": "2318",
            "name": "石林彝族自治县"
        },
        {
            "city_id": "265",
            "id": "2319",
            "name": "嵩明县"
        },
        {
            "city_id": "265",
            "id": "2320",
            "name": "禄劝彝族苗族自治县"
        },
        {
            "city_id": "265",
            "id": "2321",
            "name": "寻甸回族彝族自治县"
        },
        {
            "city_id": "265",
            "id": "2322",
            "name": "安宁市"
        },
        {
            "city_id": "266",
            "id": "2323",
            "name": "麒麟区"
        },
        {
            "city_id": "266",
            "id": "2324",
            "name": "马龙区"
        },
        {
            "city_id": "266",
            "id": "2325",
            "name": "陆良县"
        },
        {
            "city_id": "266",
            "id": "2326",
            "name": "师宗县"
        },
        {
            "city_id": "266",
            "id": "2327",
            "name": "罗平县"
        },
        {
            "city_id": "266",
            "id": "2328",
            "name": "富源县"
        },
        {
            "city_id": "266",
            "id": "2329",
            "name": "会泽县"
        },
        {
            "city_id": "266",
            "id": "2330",
            "name": "沾益区"
        },
        {
            "city_id": "266",
            "id": "2331",
            "name": "宣威市"
        },
        {
            "city_id": "267",
            "id": "2332",
            "name": "红塔区"
        },
        {
            "city_id": "267",
            "id": "2333",
            "name": "江川区"
        },
        {
            "city_id": "267",
            "id": "2334",
            "name": "澄江县"
        },
        {
            "city_id": "267",
            "id": "2335",
            "name": "通海县"
        },
        {
            "city_id": "267",
            "id": "2336",
            "name": "华宁县"
        },
        {
            "city_id": "267",
            "id": "2337",
            "name": "易门县"
        },
        {
            "city_id": "267",
            "id": "2338",
            "name": "峨山彝族自治县"
        },
        {
            "city_id": "267",
            "id": "2339",
            "name": "新平彝族傣族自治县"
        },
        {
            "city_id": "267",
            "id": "2340",
            "name": "元江哈尼族彝族傣族自治县"
        },
        {
            "city_id": "268",
            "id": "2341",
            "name": "隆阳区"
        },
        {
            "city_id": "268",
            "id": "2342",
            "name": "施甸县"
        },
        {
            "city_id": "268",
            "id": "2343",
            "name": "腾冲市"
        },
        {
            "city_id": "268",
            "id": "2344",
            "name": "龙陵县"
        },
        {
            "city_id": "268",
            "id": "2345",
            "name": "昌宁县"
        },
        {
            "city_id": "269",
            "id": "2346",
            "name": "昭阳区"
        },
        {
            "city_id": "269",
            "id": "2347",
            "name": "鲁甸县"
        },
        {
            "city_id": "269",
            "id": "2348",
            "name": "巧家县"
        },
        {
            "city_id": "269",
            "id": "2349",
            "name": "盐津县"
        },
        {
            "city_id": "269",
            "id": "2350",
            "name": "大关县"
        },
        {
            "city_id": "269",
            "id": "2351",
            "name": "永善县"
        },
        {
            "city_id": "269",
            "id": "2352",
            "name": "绥江县"
        },
        {
            "city_id": "269",
            "id": "2353",
            "name": "镇雄县"
        },
        {
            "city_id": "269",
            "id": "2354",
            "name": "彝良县"
        },
        {
            "city_id": "269",
            "id": "2355",
            "name": "威信县"
        },
        {
            "city_id": "269",
            "id": "2356",
            "name": "水富市"
        },
        {
            "city_id": "270",
            "id": "2357",
            "name": "古城区"
        },
        {
            "city_id": "270",
            "id": "2358",
            "name": "玉龙纳西族自治县"
        },
        {
            "city_id": "270",
            "id": "2359",
            "name": "永胜县"
        },
        {
            "city_id": "270",
            "id": "2360",
            "name": "华坪县"
        },
        {
            "city_id": "270",
            "id": "2361",
            "name": "宁蒗彝族自治县"
        },
        {
            "city_id": "271",
            "id": "2362",
            "name": "翠云区"
        },
        {
            "city_id": "271",
            "id": "2363",
            "name": "普洱哈尼族彝族自治县"
        },
        {
            "city_id": "271",
            "id": "2364",
            "name": "墨江哈尼族自治县"
        },
        {
            "city_id": "271",
            "id": "2365",
            "name": "景东彝族自治县"
        },
        {
            "city_id": "271",
            "id": "2366",
            "name": "景谷傣族彝族自治县"
        },
        {
            "city_id": "271",
            "id": "2367",
            "name": "镇沅彝族哈尼族拉祜族自治县"
        },
        {
            "city_id": "271",
            "id": "2368",
            "name": "江城哈尼族彝族自治县"
        },
        {
            "city_id": "271",
            "id": "2369",
            "name": "孟连傣族拉祜族佤族自治县"
        },
        {
            "city_id": "271",
            "id": "2370",
            "name": "澜沧拉祜族自治县"
        },
        {
            "city_id": "271",
            "id": "2371",
            "name": "西盟佤族自治县"
        },
        {
            "city_id": "272",
            "id": "2372",
            "name": "临翔区"
        },
        {
            "city_id": "272",
            "id": "2373",
            "name": "凤庆县"
        },
        {
            "city_id": "272",
            "id": "2374",
            "name": "云县"
        },
        {
            "city_id": "272",
            "id": "2375",
            "name": "永德县"
        },
        {
            "city_id": "272",
            "id": "2376",
            "name": "镇康县"
        },
        {
            "city_id": "272",
            "id": "2377",
            "name": "双江拉祜族佤族布朗族傣族自治县"
        },
        {
            "city_id": "272",
            "id": "2378",
            "name": "耿马傣族佤族自治县"
        },
        {
            "city_id": "272",
            "id": "2379",
            "name": "沧源佤族自治县"
        },
        {
            "city_id": "273",
            "id": "2380",
            "name": "楚雄市"
        },
        {
            "city_id": "273",
            "id": "2381",
            "name": "双柏县"
        },
        {
            "city_id": "273",
            "id": "2382",
            "name": "牟定县"
        },
        {
            "city_id": "273",
            "id": "2383",
            "name": "南华县"
        },
        {
            "city_id": "273",
            "id": "2384",
            "name": "姚安县"
        },
        {
            "city_id": "273",
            "id": "2385",
            "name": "大姚县"
        },
        {
            "city_id": "273",
            "id": "2386",
            "name": "永仁县"
        },
        {
            "city_id": "273",
            "id": "2387",
            "name": "元谋县"
        },
        {
            "city_id": "273",
            "id": "2388",
            "name": "武定县"
        },
        {
            "city_id": "273",
            "id": "2389",
            "name": "禄丰县"
        },
        {
            "city_id": "274",
            "id": "2390",
            "name": "个旧市"
        },
        {
            "city_id": "274",
            "id": "2391",
            "name": "开远市"
        },
        {
            "city_id": "274",
            "id": "2392",
            "name": "蒙自市"
        },
        {
            "city_id": "274",
            "id": "2393",
            "name": "屏边苗族自治县"
        },
        {
            "city_id": "274",
            "id": "2394",
            "name": "建水县"
        },
        {
            "city_id": "274",
            "id": "2395",
            "name": "石屏县"
        },
        {
            "city_id": "274",
            "id": "2396",
            "name": "弥勒市"
        },
        {
            "city_id": "274",
            "id": "2397",
            "name": "泸西县"
        },
        {
            "city_id": "274",
            "id": "2398",
            "name": "元阳县"
        },
        {
            "city_id": "274",
            "id": "2399",
            "name": "红河县"
        },
        {
            "city_id": "274",
            "id": "2400",
            "name": "金平苗族瑶族傣族自治县"
        },
        {
            "city_id": "274",
            "id": "2401",
            "name": "绿春县"
        },
        {
            "city_id": "274",
            "id": "2402",
            "name": "河口瑶族自治县"
        },
        {
            "city_id": "275",
            "id": "2403",
            "name": "文山市"
        },
        {
            "city_id": "275",
            "id": "2404",
            "name": "砚山县"
        },
        {
            "city_id": "275",
            "id": "2405",
            "name": "西畴县"
        },
        {
            "city_id": "275",
            "id": "2406",
            "name": "麻栗坡县"
        },
        {
            "city_id": "275",
            "id": "2407",
            "name": "马关县"
        },
        {
            "city_id": "275",
            "id": "2408",
            "name": "丘北县"
        },
        {
            "city_id": "275",
            "id": "2409",
            "name": "广南县"
        },
        {
            "city_id": "275",
            "id": "2410",
            "name": "富宁县"
        },
        {
            "city_id": "276",
            "id": "2411",
            "name": "景洪市"
        },
        {
            "city_id": "276",
            "id": "2412",
            "name": "勐海县"
        },
        {
            "city_id": "276",
            "id": "2413",
            "name": "勐腊县"
        },
        {
            "city_id": "277",
            "id": "2414",
            "name": "大理市"
        },
        {
            "city_id": "277",
            "id": "2415",
            "name": "漾濞彝族自治县"
        },
        {
            "city_id": "277",
            "id": "2416",
            "name": "祥云县"
        },
        {
            "city_id": "277",
            "id": "2417",
            "name": "宾川县"
        },
        {
            "city_id": "277",
            "id": "2418",
            "name": "弥渡县"
        },
        {
            "city_id": "277",
            "id": "2419",
            "name": "南涧彝族自治县"
        },
        {
            "city_id": "277",
            "id": "2420",
            "name": "巍山彝族回族自治县"
        },
        {
            "city_id": "277",
            "id": "2421",
            "name": "永平县"
        },
        {
            "city_id": "277",
            "id": "2422",
            "name": "云龙县"
        },
        {
            "city_id": "277",
            "id": "2423",
            "name": "洱源县"
        },
        {
            "city_id": "277",
            "id": "2424",
            "name": "剑川县"
        },
        {
            "city_id": "277",
            "id": "2425",
            "name": "鹤庆县"
        },
        {
            "city_id": "278",
            "id": "2426",
            "name": "瑞丽市"
        },
        {
            "city_id": "278",
            "id": "2427",
            "name": "潞西市"
        },
        {
            "city_id": "278",
            "id": "2428",
            "name": "梁河县"
        },
        {
            "city_id": "278",
            "id": "2429",
            "name": "盈江县"
        },
        {
            "city_id": "278",
            "id": "2430",
            "name": "陇川县"
        },
        {
            "city_id": "279",
            "id": "2431",
            "name": "泸水市"
        },
        {
            "city_id": "279",
            "id": "2432",
            "name": "福贡县"
        },
        {
            "city_id": "279",
            "id": "2433",
            "name": "贡山独龙族怒族自治县"
        },
        {
            "city_id": "279",
            "id": "2434",
            "name": "兰坪白族普米族自治县"
        },
        {
            "city_id": "280",
            "id": "2435",
            "name": "香格里拉市"
        },
        {
            "city_id": "280",
            "id": "2436",
            "name": "德钦县"
        },
        {
            "city_id": "280",
            "id": "2437",
            "name": "维西傈僳族自治县"
        },
        {
            "city_id": "281",
            "id": "2438",
            "name": "城关区"
        },
        {
            "city_id": "281",
            "id": "2439",
            "name": "林周县"
        },
        {
            "city_id": "281",
            "id": "2440",
            "name": "当雄县"
        },
        {
            "city_id": "281",
            "id": "2441",
            "name": "尼木县"
        },
        {
            "city_id": "281",
            "id": "2442",
            "name": "曲水县"
        },
        {
            "city_id": "281",
            "id": "2443",
            "name": "堆龙德庆区"
        },
        {
            "city_id": "281",
            "id": "2444",
            "name": "达孜区"
        },
        {
            "city_id": "281",
            "id": "2445",
            "name": "墨竹工卡县"
        },
        {
            "city_id": "282",
            "id": "2446",
            "name": "昌都县"
        },
        {
            "city_id": "282",
            "id": "2447",
            "name": "江达县"
        },
        {
            "city_id": "282",
            "id": "2448",
            "name": "贡觉县"
        },
        {
            "city_id": "282",
            "id": "2449",
            "name": "类乌齐县"
        },
        {
            "city_id": "282",
            "id": "2450",
            "name": "丁青县"
        },
        {
            "city_id": "282",
            "id": "2451",
            "name": "察雅县"
        },
        {
            "city_id": "282",
            "id": "2452",
            "name": "八宿县"
        },
        {
            "city_id": "282",
            "id": "2453",
            "name": "左贡县"
        },
        {
            "city_id": "282",
            "id": "2454",
            "name": "芒康县"
        },
        {
            "city_id": "282",
            "id": "2455",
            "name": "洛隆县"
        },
        {
            "city_id": "282",
            "id": "2456",
            "name": "边坝县"
        },
        {
            "city_id": "283",
            "id": "2457",
            "name": "乃东区"
        },
        {
            "city_id": "283",
            "id": "2458",
            "name": "扎囊县"
        },
        {
            "city_id": "283",
            "id": "2459",
            "name": "贡嘎县"
        },
        {
            "city_id": "283",
            "id": "2460",
            "name": "桑日县"
        },
        {
            "city_id": "283",
            "id": "2461",
            "name": "琼结县"
        },
        {
            "city_id": "283",
            "id": "2462",
            "name": "曲松县"
        },
        {
            "city_id": "283",
            "id": "2463",
            "name": "措美县"
        },
        {
            "city_id": "283",
            "id": "2464",
            "name": "洛扎县"
        },
        {
            "city_id": "283",
            "id": "2465",
            "name": "加查县"
        },
        {
            "city_id": "283",
            "id": "2466",
            "name": "隆子县"
        },
        {
            "city_id": "283",
            "id": "2467",
            "name": "错那县"
        },
        {
            "city_id": "283",
            "id": "2468",
            "name": "浪卡子县"
        },
        {
            "city_id": "284",
            "id": "2469",
            "name": "日喀则市"
        },
        {
            "city_id": "284",
            "id": "2470",
            "name": "南木林县"
        },
        {
            "city_id": "284",
            "id": "2471",
            "name": "江孜县"
        },
        {
            "city_id": "284",
            "id": "2472",
            "name": "定日县"
        },
        {
            "city_id": "284",
            "id": "2473",
            "name": "萨迦县"
        },
        {
            "city_id": "284",
            "id": "2474",
            "name": "拉孜县"
        },
        {
            "city_id": "284",
            "id": "2475",
            "name": "昂仁县"
        },
        {
            "city_id": "284",
            "id": "2476",
            "name": "谢通门县"
        },
        {
            "city_id": "284",
            "id": "2477",
            "name": "白朗县"
        },
        {
            "city_id": "284",
            "id": "2478",
            "name": "仁布县"
        },
        {
            "city_id": "284",
            "id": "2479",
            "name": "康马县"
        },
        {
            "city_id": "284",
            "id": "2480",
            "name": "定结县"
        },
        {
            "city_id": "284",
            "id": "2481",
            "name": "仲巴县"
        },
        {
            "city_id": "284",
            "id": "2482",
            "name": "亚东县"
        },
        {
            "city_id": "284",
            "id": "2483",
            "name": "吉隆县"
        },
        {
            "city_id": "284",
            "id": "2484",
            "name": "聂拉木县"
        },
        {
            "city_id": "284",
            "id": "2485",
            "name": "萨嘎县"
        },
        {
            "city_id": "284",
            "id": "2486",
            "name": "岗巴县"
        },
        {
            "city_id": "285",
            "id": "2487",
            "name": "那曲县"
        },
        {
            "city_id": "285",
            "id": "2488",
            "name": "嘉黎县"
        },
        {
            "city_id": "285",
            "id": "2489",
            "name": "比如县"
        },
        {
            "city_id": "285",
            "id": "2490",
            "name": "聂荣县"
        },
        {
            "city_id": "285",
            "id": "2491",
            "name": "安多县"
        },
        {
            "city_id": "285",
            "id": "2492",
            "name": "申扎县"
        },
        {
            "city_id": "285",
            "id": "2493",
            "name": "索县"
        },
        {
            "city_id": "285",
            "id": "2494",
            "name": "班戈县"
        },
        {
            "city_id": "285",
            "id": "2495",
            "name": "巴青县"
        },
        {
            "city_id": "285",
            "id": "2496",
            "name": "尼玛县"
        },
        {
            "city_id": "286",
            "id": "2497",
            "name": "普兰县"
        },
        {
            "city_id": "286",
            "id": "2498",
            "name": "札达县"
        },
        {
            "city_id": "286",
            "id": "2499",
            "name": "噶尔县"
        },
        {
            "city_id": "286",
            "id": "2500",
            "name": "日土县"
        },
        {
            "city_id": "286",
            "id": "2501",
            "name": "革吉县"
        },
        {
            "city_id": "286",
            "id": "2502",
            "name": "改则县"
        },
        {
            "city_id": "286",
            "id": "2503",
            "name": "措勤县"
        },
        {
            "city_id": "287",
            "id": "2504",
            "name": "林芝县"
        },
        {
            "city_id": "287",
            "id": "2505",
            "name": "工布江达县"
        },
        {
            "city_id": "287",
            "id": "2506",
            "name": "米林县"
        },
        {
            "city_id": "287",
            "id": "2507",
            "name": "墨脱县"
        },
        {
            "city_id": "287",
            "id": "2508",
            "name": "波密县"
        },
        {
            "city_id": "287",
            "id": "2509",
            "name": "察隅县"
        },
        {
            "city_id": "287",
            "id": "2510",
            "name": "朗县"
        },
        {
            "city_id": "288",
            "id": "2511",
            "name": "新城区"
        },
        {
            "city_id": "288",
            "id": "2512",
            "name": "碑林区"
        },
        {
            "city_id": "288",
            "id": "2513",
            "name": "莲湖区"
        },
        {
            "city_id": "288",
            "id": "2514",
            "name": "灞桥区"
        },
        {
            "city_id": "288",
            "id": "2515",
            "name": "未央区"
        },
        {
            "city_id": "288",
            "id": "2516",
            "name": "雁塔区"
        },
        {
            "city_id": "288",
            "id": "2517",
            "name": "阎良区"
        },
        {
            "city_id": "288",
            "id": "2518",
            "name": "临潼区"
        },
        {
            "city_id": "288",
            "id": "2519",
            "name": "长安区"
        },
        {
            "city_id": "288",
            "id": "2520",
            "name": "蓝田县"
        },
        {
            "city_id": "288",
            "id": "2521",
            "name": "周至县"
        },
        {
            "city_id": "288",
            "id": "2522",
            "name": "户县"
        },
        {
            "city_id": "288",
            "id": "2523",
            "name": "高陵区"
        },
        {
            "city_id": "289",
            "id": "2524",
            "name": "王益区"
        },
        {
            "city_id": "289",
            "id": "2525",
            "name": "印台区"
        },
        {
            "city_id": "289",
            "id": "2526",
            "name": "耀州区"
        },
        {
            "city_id": "289",
            "id": "2527",
            "name": "宜君县"
        },
        {
            "city_id": "290",
            "id": "2528",
            "name": "渭滨区"
        },
        {
            "city_id": "290",
            "id": "2529",
            "name": "金台区"
        },
        {
            "city_id": "290",
            "id": "2530",
            "name": "陈仓区"
        },
        {
            "city_id": "290",
            "id": "2531",
            "name": "凤翔县"
        },
        {
            "city_id": "290",
            "id": "2532",
            "name": "岐山县"
        },
        {
            "city_id": "290",
            "id": "2533",
            "name": "扶风县"
        },
        {
            "city_id": "290",
            "id": "2534",
            "name": "眉县"
        },
        {
            "city_id": "290",
            "id": "2535",
            "name": "陇县"
        },
        {
            "city_id": "290",
            "id": "2536",
            "name": "千阳县"
        },
        {
            "city_id": "290",
            "id": "2537",
            "name": "麟游县"
        },
        {
            "city_id": "290",
            "id": "2538",
            "name": "凤县"
        },
        {
            "city_id": "290",
            "id": "2539",
            "name": "太白县"
        },
        {
            "city_id": "291",
            "id": "2540",
            "name": "秦都区"
        },
        {
            "city_id": "291",
            "id": "2541",
            "name": "杨凌区"
        },
        {
            "city_id": "291",
            "id": "2542",
            "name": "渭城区"
        },
        {
            "city_id": "291",
            "id": "2543",
            "name": "三原县"
        },
        {
            "city_id": "291",
            "id": "2544",
            "name": "泾阳县"
        },
        {
            "city_id": "291",
            "id": "2545",
            "name": "乾县"
        },
        {
            "city_id": "291",
            "id": "2546",
            "name": "礼泉县"
        },
        {
            "city_id": "291",
            "id": "2547",
            "name": "永寿县"
        },
        {
            "city_id": "291",
            "id": "2548",
            "name": "彬县"
        },
        {
            "city_id": "291",
            "id": "2549",
            "name": "长武县"
        },
        {
            "city_id": "291",
            "id": "2550",
            "name": "旬邑县"
        },
        {
            "city_id": "291",
            "id": "2551",
            "name": "淳化县"
        },
        {
            "city_id": "291",
            "id": "2552",
            "name": "武功县"
        },
        {
            "city_id": "291",
            "id": "2553",
            "name": "兴平市"
        },
        {
            "city_id": "292",
            "id": "2554",
            "name": "临渭区"
        },
        {
            "city_id": "292",
            "id": "2555",
            "name": "华县"
        },
        {
            "city_id": "292",
            "id": "2556",
            "name": "潼关县"
        },
        {
            "city_id": "292",
            "id": "2557",
            "name": "大荔县"
        },
        {
            "city_id": "292",
            "id": "2558",
            "name": "合阳县"
        },
        {
            "city_id": "292",
            "id": "2559",
            "name": "澄城县"
        },
        {
            "city_id": "292",
            "id": "2560",
            "name": "蒲城县"
        },
        {
            "city_id": "292",
            "id": "2561",
            "name": "白水县"
        },
        {
            "city_id": "292",
            "id": "2562",
            "name": "富平县"
        },
        {
            "city_id": "292",
            "id": "2563",
            "name": "韩城市"
        },
        {
            "city_id": "292",
            "id": "2564",
            "name": "华阴市"
        },
        {
            "city_id": "293",
            "id": "2565",
            "name": "宝塔区"
        },
        {
            "city_id": "293",
            "id": "2566",
            "name": "延长县"
        },
        {
            "city_id": "293",
            "id": "2567",
            "name": "延川县"
        },
        {
            "city_id": "293",
            "id": "2568",
            "name": "子长县"
        },
        {
            "city_id": "293",
            "id": "2569",
            "name": "安塞区"
        },
        {
            "city_id": "293",
            "id": "2570",
            "name": "志丹县"
        },
        {
            "city_id": "293",
            "id": "2571",
            "name": "吴旗县"
        },
        {
            "city_id": "293",
            "id": "2572",
            "name": "甘泉县"
        },
        {
            "city_id": "293",
            "id": "2573",
            "name": "富县"
        },
        {
            "city_id": "293",
            "id": "2574",
            "name": "洛川县"
        },
        {
            "city_id": "293",
            "id": "2575",
            "name": "宜川县"
        },
        {
            "city_id": "293",
            "id": "2576",
            "name": "黄龙县"
        },
        {
            "city_id": "293",
            "id": "2577",
            "name": "黄陵县"
        },
        {
            "city_id": "294",
            "id": "2578",
            "name": "汉台区"
        },
        {
            "city_id": "294",
            "id": "2579",
            "name": "南郑区"
        },
        {
            "city_id": "294",
            "id": "2580",
            "name": "城固县"
        },
        {
            "city_id": "294",
            "id": "2581",
            "name": "洋县"
        },
        {
            "city_id": "294",
            "id": "2582",
            "name": "西乡县"
        },
        {
            "city_id": "294",
            "id": "2583",
            "name": "勉县"
        },
        {
            "city_id": "294",
            "id": "2584",
            "name": "宁强县"
        },
        {
            "city_id": "294",
            "id": "2585",
            "name": "略阳县"
        },
        {
            "city_id": "294",
            "id": "2586",
            "name": "镇巴县"
        },
        {
            "city_id": "294",
            "id": "2587",
            "name": "留坝县"
        },
        {
            "city_id": "294",
            "id": "2588",
            "name": "佛坪县"
        },
        {
            "city_id": "295",
            "id": "2589",
            "name": "榆阳区"
        },
        {
            "city_id": "295",
            "id": "2590",
            "name": "神木市"
        },
        {
            "city_id": "295",
            "id": "2591",
            "name": "府谷县"
        },
        {
            "city_id": "295",
            "id": "2592",
            "name": "横山区"
        },
        {
            "city_id": "295",
            "id": "2593",
            "name": "靖边县"
        },
        {
            "city_id": "295",
            "id": "2594",
            "name": "定边县"
        },
        {
            "city_id": "295",
            "id": "2595",
            "name": "绥德县"
        },
        {
            "city_id": "295",
            "id": "2596",
            "name": "米脂县"
        },
        {
            "city_id": "295",
            "id": "2597",
            "name": "佳县"
        },
        {
            "city_id": "295",
            "id": "2598",
            "name": "吴堡县"
        },
        {
            "city_id": "295",
            "id": "2599",
            "name": "清涧县"
        },
        {
            "city_id": "295",
            "id": "2600",
            "name": "子洲县"
        },
        {
            "city_id": "296",
            "id": "2601",
            "name": "汉滨区"
        },
        {
            "city_id": "296",
            "id": "2602",
            "name": "汉阴县"
        },
        {
            "city_id": "296",
            "id": "2603",
            "name": "石泉县"
        },
        {
            "city_id": "296",
            "id": "2604",
            "name": "宁陕县"
        },
        {
            "city_id": "296",
            "id": "2605",
            "name": "紫阳县"
        },
        {
            "city_id": "296",
            "id": "2606",
            "name": "岚皋县"
        },
        {
            "city_id": "296",
            "id": "2607",
            "name": "平利县"
        },
        {
            "city_id": "296",
            "id": "2608",
            "name": "镇坪县"
        },
        {
            "city_id": "296",
            "id": "2609",
            "name": "旬阳县"
        },
        {
            "city_id": "296",
            "id": "2610",
            "name": "白河县"
        },
        {
            "city_id": "297",
            "id": "2611",
            "name": "商州区"
        },
        {
            "city_id": "297",
            "id": "2612",
            "name": "洛南县"
        },
        {
            "city_id": "297",
            "id": "2613",
            "name": "丹凤县"
        },
        {
            "city_id": "297",
            "id": "2614",
            "name": "商南县"
        },
        {
            "city_id": "297",
            "id": "2615",
            "name": "山阳县"
        },
        {
            "city_id": "297",
            "id": "2616",
            "name": "镇安县"
        },
        {
            "city_id": "297",
            "id": "2617",
            "name": "柞水县"
        },
        {
            "city_id": "298",
            "id": "2618",
            "name": "城关区"
        },
        {
            "city_id": "298",
            "id": "2619",
            "name": "七里河区"
        },
        {
            "city_id": "298",
            "id": "2620",
            "name": "西固区"
        },
        {
            "city_id": "298",
            "id": "2621",
            "name": "安宁区"
        },
        {
            "city_id": "298",
            "id": "2622",
            "name": "红古区"
        },
        {
            "city_id": "298",
            "id": "2623",
            "name": "永登县"
        },
        {
            "city_id": "298",
            "id": "2624",
            "name": "皋兰县"
        },
        {
            "city_id": "298",
            "id": "2625",
            "name": "榆中县"
        },
        {
            "city_id": "300",
            "id": "2626",
            "name": "金川区"
        },
        {
            "city_id": "300",
            "id": "2627",
            "name": "永昌县"
        },
        {
            "city_id": "301",
            "id": "2628",
            "name": "白银区"
        },
        {
            "city_id": "301",
            "id": "2629",
            "name": "平川区"
        },
        {
            "city_id": "301",
            "id": "2630",
            "name": "靖远县"
        },
        {
            "city_id": "301",
            "id": "2631",
            "name": "会宁县"
        },
        {
            "city_id": "301",
            "id": "2632",
            "name": "景泰县"
        },
        {
            "city_id": "302",
            "id": "2633",
            "name": "秦城区"
        },
        {
            "city_id": "302",
            "id": "2634",
            "name": "北道区"
        },
        {
            "city_id": "302",
            "id": "2635",
            "name": "清水县"
        },
        {
            "city_id": "302",
            "id": "2636",
            "name": "秦安县"
        },
        {
            "city_id": "302",
            "id": "2637",
            "name": "甘谷县"
        },
        {
            "city_id": "302",
            "id": "2638",
            "name": "武山县"
        },
        {
            "city_id": "302",
            "id": "2639",
            "name": "张家川回族自治县"
        },
        {
            "city_id": "303",
            "id": "2640",
            "name": "凉州区"
        },
        {
            "city_id": "303",
            "id": "2641",
            "name": "民勤县"
        },
        {
            "city_id": "303",
            "id": "2642",
            "name": "古浪县"
        },
        {
            "city_id": "303",
            "id": "2643",
            "name": "天祝藏族自治县"
        },
        {
            "city_id": "304",
            "id": "2644",
            "name": "甘州区"
        },
        {
            "city_id": "304",
            "id": "2645",
            "name": "肃南裕固族自治县"
        },
        {
            "city_id": "304",
            "id": "2646",
            "name": "民乐县"
        },
        {
            "city_id": "304",
            "id": "2647",
            "name": "临泽县"
        },
        {
            "city_id": "304",
            "id": "2648",
            "name": "高台县"
        },
        {
            "city_id": "304",
            "id": "2649",
            "name": "山丹县"
        },
        {
            "city_id": "305",
            "id": "2650",
            "name": "崆峒区"
        },
        {
            "city_id": "305",
            "id": "2651",
            "name": "泾川县"
        },
        {
            "city_id": "305",
            "id": "2652",
            "name": "灵台县"
        },
        {
            "city_id": "305",
            "id": "2653",
            "name": "崇信县"
        },
        {
            "city_id": "305",
            "id": "2654",
            "name": "华亭市"
        },
        {
            "city_id": "305",
            "id": "2655",
            "name": "庄浪县"
        },
        {
            "city_id": "305",
            "id": "2656",
            "name": "静宁县"
        },
        {
            "city_id": "306",
            "id": "2657",
            "name": "肃州区"
        },
        {
            "city_id": "306",
            "id": "2658",
            "name": "金塔县"
        },
        {
            "city_id": "306",
            "id": "2659",
            "name": "安西县"
        },
        {
            "city_id": "306",
            "id": "2660",
            "name": "肃北蒙古族自治县"
        },
        {
            "city_id": "306",
            "id": "2661",
            "name": "阿克塞哈萨克族自治县"
        },
        {
            "city_id": "306",
            "id": "2662",
            "name": "玉门市"
        },
        {
            "city_id": "306",
            "id": "2663",
            "name": "敦煌市"
        },
        {
            "city_id": "307",
            "id": "2664",
            "name": "西峰区"
        },
        {
            "city_id": "307",
            "id": "2665",
            "name": "庆城县"
        },
        {
            "city_id": "307",
            "id": "2666",
            "name": "环县"
        },
        {
            "city_id": "307",
            "id": "2667",
            "name": "华池县"
        },
        {
            "city_id": "307",
            "id": "2668",
            "name": "合水县"
        },
        {
            "city_id": "307",
            "id": "2669",
            "name": "正宁县"
        },
        {
            "city_id": "307",
            "id": "2670",
            "name": "宁县"
        },
        {
            "city_id": "307",
            "id": "2671",
            "name": "镇原县"
        },
        {
            "city_id": "308",
            "id": "2672",
            "name": "安定区"
        },
        {
            "city_id": "308",
            "id": "2673",
            "name": "通渭县"
        },
        {
            "city_id": "308",
            "id": "2674",
            "name": "陇西县"
        },
        {
            "city_id": "308",
            "id": "2675",
            "name": "渭源县"
        },
        {
            "city_id": "308",
            "id": "2676",
            "name": "临洮县"
        },
        {
            "city_id": "308",
            "id": "2677",
            "name": "漳县"
        },
        {
            "city_id": "308",
            "id": "2678",
            "name": "岷县"
        },
        {
            "city_id": "309",
            "id": "2679",
            "name": "武都区"
        },
        {
            "city_id": "309",
            "id": "2680",
            "name": "成县"
        },
        {
            "city_id": "309",
            "id": "2681",
            "name": "文县"
        },
        {
            "city_id": "309",
            "id": "2682",
            "name": "宕昌县"
        },
        {
            "city_id": "309",
            "id": "2683",
            "name": "康县"
        },
        {
            "city_id": "309",
            "id": "2684",
            "name": "西和县"
        },
        {
            "city_id": "309",
            "id": "2685",
            "name": "礼县"
        },
        {
            "city_id": "309",
            "id": "2686",
            "name": "徽县"
        },
        {
            "city_id": "309",
            "id": "2687",
            "name": "两当县"
        },
        {
            "city_id": "310",
            "id": "2688",
            "name": "临夏市"
        },
        {
            "city_id": "310",
            "id": "2689",
            "name": "临夏县"
        },
        {
            "city_id": "310",
            "id": "2690",
            "name": "康乐县"
        },
        {
            "city_id": "310",
            "id": "2691",
            "name": "永靖县"
        },
        {
            "city_id": "310",
            "id": "2692",
            "name": "广河县"
        },
        {
            "city_id": "310",
            "id": "2693",
            "name": "和政县"
        },
        {
            "city_id": "310",
            "id": "2694",
            "name": "东乡族自治县"
        },
        {
            "city_id": "310",
            "id": "2695",
            "name": "积石山保安族东乡族撒拉族自治县"
        },
        {
            "city_id": "311",
            "id": "2696",
            "name": "合作市"
        },
        {
            "city_id": "311",
            "id": "2697",
            "name": "临潭县"
        },
        {
            "city_id": "311",
            "id": "2698",
            "name": "卓尼县"
        },
        {
            "city_id": "311",
            "id": "2699",
            "name": "舟曲县"
        },
        {
            "city_id": "311",
            "id": "2700",
            "name": "迭部县"
        },
        {
            "city_id": "311",
            "id": "2701",
            "name": "玛曲县"
        },
        {
            "city_id": "311",
            "id": "2702",
            "name": "碌曲县"
        },
        {
            "city_id": "311",
            "id": "2703",
            "name": "夏河县"
        },
        {
            "city_id": "312",
            "id": "2704",
            "name": "城东区"
        },
        {
            "city_id": "312",
            "id": "2705",
            "name": "城中区"
        },
        {
            "city_id": "312",
            "id": "2706",
            "name": "城西区"
        },
        {
            "city_id": "312",
            "id": "2707",
            "name": "城北区"
        },
        {
            "city_id": "312",
            "id": "2708",
            "name": "大通回族土族自治县"
        },
        {
            "city_id": "312",
            "id": "2709",
            "name": "湟中县"
        },
        {
            "city_id": "312",
            "id": "2710",
            "name": "湟源县"
        },
        {
            "city_id": "313",
            "id": "2711",
            "name": "平安区"
        },
        {
            "city_id": "313",
            "id": "2712",
            "name": "民和回族土族自治县"
        },
        {
            "city_id": "313",
            "id": "2713",
            "name": "乐都区"
        },
        {
            "city_id": "313",
            "id": "2714",
            "name": "互助土族自治县"
        },
        {
            "city_id": "313",
            "id": "2715",
            "name": "化隆回族自治县"
        },
        {
            "city_id": "313",
            "id": "2716",
            "name": "循化撒拉族自治县"
        },
        {
            "city_id": "314",
            "id": "2717",
            "name": "门源回族自治县"
        },
        {
            "city_id": "314",
            "id": "2718",
            "name": "祁连县"
        },
        {
            "city_id": "314",
            "id": "2719",
            "name": "海晏县"
        },
        {
            "city_id": "314",
            "id": "2720",
            "name": "刚察县"
        },
        {
            "city_id": "315",
            "id": "2721",
            "name": "同仁县"
        },
        {
            "city_id": "315",
            "id": "2722",
            "name": "尖扎县"
        },
        {
            "city_id": "315",
            "id": "2723",
            "name": "泽库县"
        },
        {
            "city_id": "315",
            "id": "2724",
            "name": "河南蒙古族自治县"
        },
        {
            "city_id": "316",
            "id": "2725",
            "name": "共和县"
        },
        {
            "city_id": "316",
            "id": "2726",
            "name": "同德县"
        },
        {
            "city_id": "316",
            "id": "2727",
            "name": "贵德县"
        },
        {
            "city_id": "316",
            "id": "2728",
            "name": "兴海县"
        },
        {
            "city_id": "316",
            "id": "2729",
            "name": "贵南县"
        },
        {
            "city_id": "317",
            "id": "2730",
            "name": "玛沁县"
        },
        {
            "city_id": "317",
            "id": "2731",
            "name": "班玛县"
        },
        {
            "city_id": "317",
            "id": "2732",
            "name": "甘德县"
        },
        {
            "city_id": "317",
            "id": "2733",
            "name": "达日县"
        },
        {
            "city_id": "317",
            "id": "2734",
            "name": "久治县"
        },
        {
            "city_id": "317",
            "id": "2735",
            "name": "玛多县"
        },
        {
            "city_id": "318",
            "id": "2736",
            "name": "玉树市"
        },
        {
            "city_id": "318",
            "id": "2737",
            "name": "杂多县"
        },
        {
            "city_id": "318",
            "id": "2738",
            "name": "称多县"
        },
        {
            "city_id": "318",
            "id": "2739",
            "name": "治多县"
        },
        {
            "city_id": "318",
            "id": "2740",
            "name": "囊谦县"
        },
        {
            "city_id": "318",
            "id": "2741",
            "name": "曲麻莱县"
        },
        {
            "city_id": "319",
            "id": "2742",
            "name": "格尔木市"
        },
        {
            "city_id": "319",
            "id": "2743",
            "name": "德令哈市"
        },
        {
            "city_id": "319",
            "id": "2744",
            "name": "乌兰县"
        },
        {
            "city_id": "319",
            "id": "2745",
            "name": "都兰县"
        },
        {
            "city_id": "319",
            "id": "2746",
            "name": "天峻县"
        },
        {
            "city_id": "320",
            "id": "2747",
            "name": "兴庆区"
        },
        {
            "city_id": "320",
            "id": "2748",
            "name": "西夏区"
        },
        {
            "city_id": "320",
            "id": "2749",
            "name": "金凤区"
        },
        {
            "city_id": "320",
            "id": "2750",
            "name": "永宁县"
        },
        {
            "city_id": "320",
            "id": "2751",
            "name": "贺兰县"
        },
        {
            "city_id": "320",
            "id": "2752",
            "name": "灵武市"
        },
        {
            "city_id": "321",
            "id": "2753",
            "name": "大武口区"
        },
        {
            "city_id": "321",
            "id": "2754",
            "name": "惠农区"
        },
        {
            "city_id": "321",
            "id": "2755",
            "name": "平罗县"
        },
        {
            "city_id": "322",
            "id": "2756",
            "name": "利通区"
        },
        {
            "city_id": "322",
            "id": "2757",
            "name": "盐池县"
        },
        {
            "city_id": "322",
            "id": "2758",
            "name": "同心县"
        },
        {
            "city_id": "322",
            "id": "2759",
            "name": "青铜峡市"
        },
        {
            "city_id": "323",
            "id": "2760",
            "name": "原州区"
        },
        {
            "city_id": "323",
            "id": "2761",
            "name": "西吉县"
        },
        {
            "city_id": "323",
            "id": "2762",
            "name": "隆德县"
        },
        {
            "city_id": "323",
            "id": "2763",
            "name": "泾源县"
        },
        {
            "city_id": "323",
            "id": "2764",
            "name": "彭阳县"
        },
        {
            "city_id": "324",
            "id": "2765",
            "name": "沙坡头区"
        },
        {
            "city_id": "324",
            "id": "2766",
            "name": "中宁县"
        },
        {
            "city_id": "324",
            "id": "2767",
            "name": "海原县"
        },
        {
            "city_id": "325",
            "id": "2768",
            "name": "天山区"
        },
        {
            "city_id": "325",
            "id": "2769",
            "name": "沙依巴克区"
        },
        {
            "city_id": "325",
            "id": "2770",
            "name": "新市区"
        },
        {
            "city_id": "325",
            "id": "2771",
            "name": "水磨沟区"
        },
        {
            "city_id": "325",
            "id": "2772",
            "name": "头屯河区"
        },
        {
            "city_id": "325",
            "id": "2773",
            "name": "达坂城区"
        },
        {
            "city_id": "325",
            "id": "2774",
            "name": "东山区"
        },
        {
            "city_id": "325",
            "id": "2775",
            "name": "乌鲁木齐县"
        },
        {
            "city_id": "326",
            "id": "2776",
            "name": "独山子区"
        },
        {
            "city_id": "326",
            "id": "2777",
            "name": "克拉玛依区"
        },
        {
            "city_id": "326",
            "id": "2778",
            "name": "白碱滩区"
        },
        {
            "city_id": "326",
            "id": "2779",
            "name": "乌尔禾区"
        },
        {
            "city_id": "327",
            "id": "2780",
            "name": "吐鲁番市"
        },
        {
            "city_id": "327",
            "id": "2781",
            "name": "鄯善县"
        },
        {
            "city_id": "327",
            "id": "2782",
            "name": "托克逊县"
        },
        {
            "city_id": "328",
            "id": "2783",
            "name": "哈密市"
        },
        {
            "city_id": "328",
            "id": "2784",
            "name": "巴里坤哈萨克自治县"
        },
        {
            "city_id": "328",
            "id": "2785",
            "name": "伊吾县"
        },
        {
            "city_id": "329",
            "id": "2786",
            "name": "昌吉市"
        },
        {
            "city_id": "329",
            "id": "2787",
            "name": "阜康市"
        },
        {
            "city_id": "329",
            "id": "2788",
            "name": "米泉市"
        },
        {
            "city_id": "329",
            "id": "2789",
            "name": "呼图壁县"
        },
        {
            "city_id": "329",
            "id": "2790",
            "name": "玛纳斯县"
        },
        {
            "city_id": "329",
            "id": "2791",
            "name": "奇台县"
        },
        {
            "city_id": "329",
            "id": "2792",
            "name": "吉木萨尔县"
        },
        {
            "city_id": "329",
            "id": "2793",
            "name": "木垒哈萨克自治县"
        },
        {
            "city_id": "330",
            "id": "2794",
            "name": "博乐市"
        },
        {
            "city_id": "330",
            "id": "2795",
            "name": "精河县"
        },
        {
            "city_id": "330",
            "id": "2796",
            "name": "温泉县"
        },
        {
            "city_id": "331",
            "id": "2797",
            "name": "库尔勒市"
        },
        {
            "city_id": "331",
            "id": "2798",
            "name": "轮台县"
        },
        {
            "city_id": "331",
            "id": "2799",
            "name": "尉犁县"
        },
        {
            "city_id": "331",
            "id": "2800",
            "name": "若羌县"
        },
        {
            "city_id": "331",
            "id": "2801",
            "name": "且末县"
        },
        {
            "city_id": "331",
            "id": "2802",
            "name": "焉耆回族自治县"
        },
        {
            "city_id": "331",
            "id": "2803",
            "name": "和静县"
        },
        {
            "city_id": "331",
            "id": "2804",
            "name": "和硕县"
        },
        {
            "city_id": "331",
            "id": "2805",
            "name": "博湖县"
        },
        {
            "city_id": "332",
            "id": "2806",
            "name": "阿克苏市"
        },
        {
            "city_id": "332",
            "id": "2807",
            "name": "温宿县"
        },
        {
            "city_id": "332",
            "id": "2808",
            "name": "库车县"
        },
        {
            "city_id": "332",
            "id": "2809",
            "name": "沙雅县"
        },
        {
            "city_id": "332",
            "id": "2810",
            "name": "新和县"
        },
        {
            "city_id": "332",
            "id": "2811",
            "name": "拜城县"
        },
        {
            "city_id": "332",
            "id": "2812",
            "name": "乌什县"
        },
        {
            "city_id": "332",
            "id": "2813",
            "name": "阿瓦提县"
        },
        {
            "city_id": "332",
            "id": "2814",
            "name": "柯坪县"
        },
        {
            "city_id": "333",
            "id": "2815",
            "name": "阿图什市"
        },
        {
            "city_id": "333",
            "id": "2816",
            "name": "阿克陶县"
        },
        {
            "city_id": "333",
            "id": "2817",
            "name": "阿合奇县"
        },
        {
            "city_id": "333",
            "id": "2818",
            "name": "乌恰县"
        },
        {
            "city_id": "334",
            "id": "2819",
            "name": "喀什市"
        },
        {
            "city_id": "334",
            "id": "2820",
            "name": "疏附县"
        },
        {
            "city_id": "334",
            "id": "2821",
            "name": "疏勒县"
        },
        {
            "city_id": "334",
            "id": "2822",
            "name": "英吉沙县"
        },
        {
            "city_id": "334",
            "id": "2823",
            "name": "泽普县"
        },
        {
            "city_id": "334",
            "id": "2824",
            "name": "莎车县"
        },
        {
            "city_id": "334",
            "id": "2825",
            "name": "叶城县"
        },
        {
            "city_id": "334",
            "id": "2826",
            "name": "麦盖提县"
        },
        {
            "city_id": "334",
            "id": "2827",
            "name": "岳普湖县"
        },
        {
            "city_id": "334",
            "id": "2828",
            "name": "伽师县"
        },
        {
            "city_id": "334",
            "id": "2829",
            "name": "巴楚县"
        },
        {
            "city_id": "334",
            "id": "2830",
            "name": "塔什库尔干塔吉克自治县"
        },
        {
            "city_id": "335",
            "id": "2831",
            "name": "和田市"
        },
        {
            "city_id": "335",
            "id": "2832",
            "name": "和田县"
        },
        {
            "city_id": "335",
            "id": "2833",
            "name": "墨玉县"
        },
        {
            "city_id": "335",
            "id": "2834",
            "name": "皮山县"
        },
        {
            "city_id": "335",
            "id": "2835",
            "name": "洛浦县"
        },
        {
            "city_id": "335",
            "id": "2836",
            "name": "策勒县"
        },
        {
            "city_id": "335",
            "id": "2837",
            "name": "于田县"
        },
        {
            "city_id": "335",
            "id": "2838",
            "name": "民丰县"
        },
        {
            "city_id": "336",
            "id": "2839",
            "name": "伊宁市"
        },
        {
            "city_id": "336",
            "id": "2840",
            "name": "奎屯市"
        },
        {
            "city_id": "336",
            "id": "2841",
            "name": "伊宁县"
        },
        {
            "city_id": "336",
            "id": "2842",
            "name": "察布查尔锡伯自治县"
        },
        {
            "city_id": "336",
            "id": "2843",
            "name": "霍城县"
        },
        {
            "city_id": "336",
            "id": "2844",
            "name": "巩留县"
        },
        {
            "city_id": "336",
            "id": "2845",
            "name": "新源县"
        },
        {
            "city_id": "336",
            "id": "2846",
            "name": "昭苏县"
        },
        {
            "city_id": "336",
            "id": "2847",
            "name": "特克斯县"
        },
        {
            "city_id": "336",
            "id": "2848",
            "name": "尼勒克县"
        },
        {
            "city_id": "337",
            "id": "2849",
            "name": "塔城市"
        },
        {
            "city_id": "337",
            "id": "2850",
            "name": "乌苏市"
        },
        {
            "city_id": "337",
            "id": "2851",
            "name": "额敏县"
        },
        {
            "city_id": "337",
            "id": "2852",
            "name": "沙湾县"
        },
        {
            "city_id": "337",
            "id": "2853",
            "name": "托里县"
        },
        {
            "city_id": "337",
            "id": "2854",
            "name": "裕民县"
        },
        {
            "city_id": "337",
            "id": "2855",
            "name": "和布克赛尔蒙古自治县"
        },
        {
            "city_id": "338",
            "id": "2856",
            "name": "阿勒泰市"
        },
        {
            "city_id": "338",
            "id": "2857",
            "name": "布尔津县"
        },
        {
            "city_id": "338",
            "id": "2858",
            "name": "富蕴县"
        },
        {
            "city_id": "338",
            "id": "2859",
            "name": "福海县"
        },
        {
            "city_id": "338",
            "id": "2860",
            "name": "哈巴河县"
        },
        {
            "city_id": "338",
            "id": "2861",
            "name": "青河县"
        },
        {
            "city_id": "338",
            "id": "2862",
            "name": "吉木乃县"
        },
        {
            "city_id": "214",
            "id": "2863",
            "name": "中山市"
        },
        {
            "city_id": "75",
            "id": "2864",
            "name": "新吴区"
        },
        {
            "city_id": "75",
            "id": "2865",
            "name": "梁溪区"
        },
        {
            "city_id": "213",
            "id": "2866",
            "name": "东莞港"
        },
        {
            "city_id": "343",
            "id": "2867",
            "name": "香港特别行政区"
        },
        {
            "city_id": "344",
            "id": "2868",
            "name": "澳门特别行政区"
        },
        {
            "city_id": "345",
            "id": "2869",
            "name": "台湾省"
        },
        {
            "city_id": "78",
            "id": "2870",
            "name": "姑苏区"
        },
        {
            "city_id": "233",
            "id": "2871",
            "name": "天涯区"
        },
        {
            "city_id": "241",
            "id": "2872",
            "name": "利州区"
        },
        {
            "city_id": "78",
            "id": "2873",
            "name": "工业园区"
        },
        {
            "city_id": "213",
            "id": "2874",
            "name": "石龙镇"
        },
        {
            "city_id": "292",
            "id": "2875",
            "name": "华州区"
        },
        {
            "city_id": "30",
            "id": "2876",
            "name": "康巴什区"
        },
        {
            "city_id": "43",
            "id": "2877",
            "name": "经济技术开发区"
        },
        {
            "city_id": "214",
            "id": "2878",
            "name": "东区街道办事处"
        },
        {
            "city_id": "156",
            "id": "2879",
            "name": "开发区"
        },
        {
            "city_id": "180",
            "id": "2880",
            "name": "随县"
        },
        {
            "city_id": "271",
            "id": "2881",
            "name": "宁洱哈尼族彝族自治县"
        },
        {
            "city_id": "127",
            "id": "2882",
            "name": "经济技术开发区"
        },
        {
            "city_id": "141",
            "id": "2883",
            "name": "开发区"
        },
        {
            "city_id": "141",
            "id": "2884",
            "name": "高新区"
        },
        {
            "city_id": "135",
            "id": "2885",
            "name": "莱芜区"
        },
        {
            "city_id": "235",
            "id": "2886",
            "name": "郫都区"
        },
        {
            "city_id": "109",
            "id": "2887",
            "name": "经济开发区"
        },
        {
            "city_id": "101",
            "id": "2888",
            "name": "寿县"
        },
        {
            "city_id": "336",
            "id": "2889",
            "name": "霍尔果斯市"
        },
        {
            "city_id": "235",
            "id": "2890",
            "name": "高新区"
        },
        {
            "city_id": "247",
            "id": "2891",
            "name": "叙州区"
        },
        {
            "city_id": "127",
            "id": "2892",
            "name": "共青城市"
        },
        {
            "city_id": "98",
            "id": "2893",
            "name": "庐江县"
        },
        {
            "city_id": "104",
            "id": "2894",
            "name": "义安区"
        },
        {
            "city_id": "251",
            "id": "2895",
            "name": "恩阳区"
        },
        {
            "city_id": "213",
            "id": "2896",
            "name": "寮步镇"
        },
        {
            "city_id": "213",
            "id": "2897",
            "name": "石排镇"
        },
        {
            "city_id": "213",
            "id": "2898",
            "name": "大岭山镇"
        },
        {
            "city_id": "213",
            "id": "2899",
            "name": "松山湖管委会"
        },
        {
            "city_id": "213",
            "id": "2900",
            "name": "企石镇"
        },
        {
            "city_id": "213",
            "id": "2901",
            "name": "凤岗镇"
        },
        {
            "city_id": "213",
            "id": "2902",
            "name": "厚街镇"
        },
        {
            "city_id": "213",
            "id": "2903",
            "name": "常平镇"
        },
        {
            "city_id": "213",
            "id": "2904",
            "name": "清溪镇"
        },
        {
            "city_id": "213",
            "id": "2905",
            "name": "望牛墩镇"
        },
        {
            "city_id": "104",
            "id": "2906",
            "name": "枞阳县"
        },
        {
            "city_id": "322",
            "id": "2907",
            "name": "红寺堡区"
        },
        {
            "city_id": "241",
            "id": "2908",
            "name": "昭化区"
        },
        {
            "city_id": "171",
            "id": "2909",
            "name": "郧阳区"
        },
        {
            "city_id": "12",
            "id": "2910",
            "name": "开发区"
        },
        {
            "city_id": "240",
            "id": "2911",
            "name": "高新区"
        },
        {
            "city_id": "213",
            "id": "2912",
            "name": "塘厦镇"
        },
        {
            "city_id": "213",
            "id": "2913",
            "name": "横沥镇"
        },
        {
            "city_id": "213",
            "id": "2914",
            "name": "大朗镇"
        },
        {
            "city_id": "213",
            "id": "2915",
            "name": "黄江镇"
        },
        {
            "city_id": "213",
            "id": "2916",
            "name": "东莞生态园"
        },
        {
            "city_id": "256",
            "id": "2917",
            "name": "观山湖区"
        },
        {
            "city_id": "260",
            "id": "2918",
            "name": "碧江区"
        },
        {
            "city_id": "37",
            "id": "2919",
            "name": "沈北新区"
        },
        {
            "city_id": "37",
            "id": "2920",
            "name": "浑南区"
        },
        {
            "city_id": "291",
            "id": "2921",
            "name": "彬州市"
        },
        {
            "city_id": "140",
            "id": "2922",
            "name": "开发区"
        },
        {
            "city_id": "257",
            "id": "2923",
            "name": "盘州市"
        },
        {
            "city_id": "4",
            "id": "2924",
            "name": "曹妃甸区"
        },
        {
            "city_id": "325",
            "id": "2925",
            "name": "米东区"
        },
        {
            "city_id": "37",
            "id": "2926",
            "name": "经济技术开发区"
        },
        {
            "city_id": "102",
            "id": "2927",
            "name": "博望区"
        },
        {
            "city_id": "285",
            "id": "2928",
            "name": "双湖县"
        },
        {
            "city_id": "213",
            "id": "2929",
            "name": "长安镇"
        },
        {
            "city_id": "213",
            "id": "2930",
            "name": "东坑镇"
        },
        {
            "city_id": "213",
            "id": "2931",
            "name": "樟木头镇"
        },
        {
            "city_id": "213",
            "id": "2932",
            "name": "莞城街道办事处"
        },
        {
            "city_id": "213",
            "id": "2933",
            "name": "中堂镇"
        },
        {
            "city_id": "213",
            "id": "2934",
            "name": "南城街道办事处"
        },
        {
            "city_id": "135",
            "id": "2935",
            "name": "高新区"
        },
        {
            "city_id": "81",
            "id": "2936",
            "name": "清江浦区"
        },
        {
            "city_id": "8",
            "id": "2937",
            "name": "竞秀区"
        },
        {
            "city_id": "319",
            "id": "2938",
            "name": "茫崖市"
        },
        {
            "city_id": "15",
            "id": "2939",
            "name": "平城区"
        },
        {
            "city_id": "76",
            "id": "2940",
            "name": "工业园区"
        },
        {
            "city_id": "213",
            "id": "2941",
            "name": "谢岗镇"
        },
        {
            "city_id": "213",
            "id": "2942",
            "name": "虎门镇"
        },
        {
            "city_id": "213",
            "id": "2943",
            "name": "麻涌镇"
        },
        {
            "city_id": "213",
            "id": "2944",
            "name": "万江街道办事处"
        },
        {
            "city_id": "213",
            "id": "2945",
            "name": "洪梅镇"
        },
        {
            "city_id": "213",
            "id": "2946",
            "name": "东城街道办事处"
        },
        {
            "city_id": "213",
            "id": "2947",
            "name": "茶山镇"
        },
        {
            "city_id": "213",
            "id": "2948",
            "name": "石碣镇"
        },
        {
            "city_id": "213",
            "id": "2949",
            "name": "高埗镇"
        },
        {
            "city_id": "213",
            "id": "2950",
            "name": "道滘镇"
        },
        {
            "city_id": "214",
            "id": "2951",
            "name": "大涌镇"
        },
        {
            "city_id": "214",
            "id": "2952",
            "name": "南朗镇"
        },
        {
            "city_id": "214",
            "id": "2953",
            "name": "古镇镇"
        },
        {
            "city_id": "214",
            "id": "2954",
            "name": "坦洲镇"
        },
        {
            "city_id": "214",
            "id": "2955",
            "name": "西区街道办事处"
        },
        {
            "city_id": "214",
            "id": "2956",
            "name": "南区街道办事处"
        },
        {
            "city_id": "285",
            "id": "2957",
            "name": "色尼区"
        },
        {
            "city_id": "81",
            "id": "2958",
            "name": "淮安区"
        },
        {
            "city_id": "98",
            "id": "2959",
            "name": "巢湖市"
        },
        {
            "city_id": "199",
            "id": "2960",
            "name": "坪山区"
        },
        {
            "city_id": "102",
            "id": "2961",
            "name": "含山县"
        },
        {
            "city_id": "102",
            "id": "2962",
            "name": "和县"
        },
        {
            "city_id": "231",
            "id": "2963",
            "name": "江州区"
        },
        {
            "city_id": "72",
            "id": "2964",
            "name": "呼中区"
        },
        {
            "city_id": "72",
            "id": "2965",
            "name": "松岭区"
        },
        {
            "city_id": "72",
            "id": "2966",
            "name": "新林区"
        },
        {
            "city_id": "72",
            "id": "2967",
            "name": "加格达奇区"
        },
        {
            "city_id": "233",
            "id": "2968",
            "name": "崖州区"
        },
        {
            "city_id": "204",
            "id": "2969",
            "name": "经济技术开发区"
        },
        {
            "city_id": "262",
            "id": "2970",
            "name": "七星关区"
        },
        {
            "city_id": "240",
            "id": "2971",
            "name": "安州区"
        },
        {
            "city_id": "153",
            "id": "2972",
            "name": "禹王台区"
        },
        {
            "city_id": "98",
            "id": "2973",
            "name": "经济技术开发区"
        },
        {
            "city_id": "98",
            "id": "2974",
            "name": "高新技术开发区"
        },
        {
            "city_id": "199",
            "id": "2975",
            "name": "龙华区"
        },
        {
            "city_id": "327",
            "id": "2976",
            "name": "高昌区"
        },
        {
            "city_id": "81",
            "id": "2977",
            "name": "经济开发区"
        },
        {
            "city_id": "184",
            "id": "2978",
            "name": "渌口区"
        },
        {
            "city_id": "214",
            "id": "2979",
            "name": "横栏镇"
        },
        {
            "city_id": "214",
            "id": "2980",
            "name": "三角镇"
        },
        {
            "city_id": "214",
            "id": "2981",
            "name": "南头镇"
        },
        {
            "city_id": "214",
            "id": "2982",
            "name": "神湾镇"
        },
        {
            "city_id": "214",
            "id": "2983",
            "name": "东凤镇"
        },
        {
            "city_id": "214",
            "id": "2984",
            "name": "五桂山街道办事处"
        },
        {
            "city_id": "214",
            "id": "2985",
            "name": "黄圃镇"
        },
        {
            "city_id": "214",
            "id": "2986",
            "name": "小榄镇"
        },
        {
            "city_id": "214",
            "id": "2987",
            "name": "石岐区街道办事处"
        },
        {
            "city_id": "302",
            "id": "2988",
            "name": "麦积区"
        },
        {
            "city_id": "302",
            "id": "2989",
            "name": "秦州区"
        },
        {
            "city_id": "99",
            "id": "2990",
            "name": "三山区"
        },
        {
            "city_id": "99",
            "id": "2991",
            "name": "弋江区"
        },
        {
            "city_id": "142",
            "id": "2992",
            "name": "高新区"
        },
        {
            "city_id": "213",
            "id": "2993",
            "name": "桥头镇"
        },
        {
            "city_id": "233",
            "id": "2994",
            "name": "海棠区"
        },
        {
            "city_id": "233",
            "id": "2995",
            "name": "吉阳区"
        },
        {
            "city_id": "153",
            "id": "2996",
            "name": "祥符区"
        },
        {
            "city_id": "235",
            "id": "2997",
            "name": "简阳市"
        },
        {
            "city_id": "5",
            "id": "2998",
            "name": "经济技术开发区"
        },
        {
            "city_id": "39",
            "id": "2999",
            "name": "高新区"
        },
        {
            "city_id": "152",
            "id": "3000",
            "name": "经济技术开发区"
        },
        {
            "city_id": "152",
            "id": "3001",
            "name": "高新技术开发区"
        },
        {
            "city_id": "288",
            "id": "3002",
            "name": "鄠邑区"
        },
        {
            "city_id": "161",
            "id": "3003",
            "name": "建安区"
        },
        {
            "city_id": "306",
            "id": "3004",
            "name": "瓜州县"
        },
        {
            "city_id": "214",
            "id": "3005",
            "name": "民众镇"
        },
        {
            "city_id": "214",
            "id": "3006",
            "name": "阜沙镇"
        },
        {
            "city_id": "214",
            "id": "3007",
            "name": "东升镇"
        },
        {
            "city_id": "214",
            "id": "3008",
            "name": "板芙镇"
        },
        {
            "city_id": "214",
            "id": "3009",
            "name": "沙溪镇"
        },
        {
            "city_id": "214",
            "id": "3010",
            "name": "港口镇"
        },
        {
            "city_id": "214",
            "id": "3011",
            "name": "三乡镇"
        },
        {
            "city_id": "2",
            "id": "3012",
            "name": "蓟州区"
        },
        {
            "city_id": "248",
            "id": "3013",
            "name": "前锋区"
        },
        {
            "city_id": "56",
            "id": "3014",
            "name": "浑江区"
        },
        {
            "city_id": "78",
            "id": "3015",
            "name": "高新区"
        },
        {
            "city_id": "144",
            "id": "3016",
            "name": "经济技术开发区"
        },
        {
            "city_id": "127",
            "id": "3017",
            "name": "柴桑区"
        },
        {
            "city_id": "127",
            "id": "3018",
            "name": "濂溪区"
        },
        {
            "city_id": "99",
            "id": "3019",
            "name": "无为县"
        },
        {
            "city_id": "226",
            "id": "3020",
            "name": "福绵区"
        },
        {
            "city_id": "249",
            "id": "3021",
            "name": "达川区"
        },
        {
            "city_id": "26",
            "id": "3022",
            "name": "白云鄂博矿区"
        },
        {
            "city_id": "234",
            "id": "3023",
            "name": "开州区"
        },
        {
            "city_id": "8",
            "id": "3024",
            "name": "莲池区"
        },
        {
            "city_id": "221",
            "id": "3025",
            "name": "龙圩区"
        },
        {
            "city_id": "15",
            "id": "3026",
            "name": "云冈区"
        },
        {
            "city_id": "15",
            "id": "3027",
            "name": "云州区"
        },
        {
            "city_id": "154",
            "id": "3028",
            "name": "瀍河回族区"
        },
        {
            "city_id": "282",
            "id": "3029",
            "name": "卡若区"
        },
        {
            "city_id": "79",
            "id": "3030",
            "name": "高新区"
        },
        {
            "city_id": "199",
            "id": "3031",
            "name": "光明区"
        },
        {
            "city_id": "214",
            "id": "3032",
            "name": "火炬开发区街道办事处"
        },
        {
            "city_id": "83",
            "id": "3033",
            "name": "经济开发区"
        },
        {
            "city_id": "105",
            "id": "3034",
            "name": "宜秀区"
        },
        {
            "city_id": "90",
            "id": "3035",
            "name": "南湖区"
        },
        {
            "city_id": "163",
            "id": "3036",
            "name": "陕州区"
        },
        {
            "city_id": "278",
            "id": "3037",
            "name": "芒市"
        },
        {
            "city_id": "287",
            "id": "3038",
            "name": "巴宜区"
        },
        {
            "city_id": "51",
            "id": "3039",
            "name": "经济技术开发区"
        },
        {
            "city_id": "258",
            "id": "3040",
            "name": "播州区"
        },
        {
            "city_id": "197",
            "id": "3041",
            "name": "南沙区"
        },
        {
            "city_id": "328",
            "id": "3042",
            "name": "伊州区"
        },
        {
            "city_id": "124",
            "id": "3043",
            "name": "高新区"
        },
        {
            "city_id": "124",
            "id": "3044",
            "name": "经济技术开发区"
        },
        {
            "city_id": "284",
            "id": "3045",
            "name": "桑珠孜区"
        },
        {
            "city_id": "136",
            "id": "3046",
            "name": "开发区"
        },
        {
            "city_id": "4",
            "id": "3047",
            "name": "滦州市"
        },
        {
            "city_id": "172",
            "id": "3048",
            "name": "经济开发区"
        },
        {
            "city_id": "173",
            "id": "3049",
            "name": "襄州区"
        },
        {
            "city_id": "228",
            "id": "3050",
            "name": "平桂区"
        },
        {
            "city_id": "148",
            "id": "3051",
            "name": "陵城区"
        },
        {
            "city_id": "271",
            "id": "3052",
            "name": "思茅区"
        },
        {
            "city_id": "92",
            "id": "3053",
            "name": "柯桥区"
        },
        {
            "city_id": "111",
            "id": "3054",
            "name": "叶集区"
        },
        {
            "city_id": "291",
            "id": "3055",
            "name": "杨陵区"
        },
        {
            "city_id": "2",
            "id": "3056",
            "name": "滨海新区"
        },
        {
            "city_id": "31",
            "id": "3057",
            "name": "扎赉诺尔区"
        },
        {
            "city_id": "330",
            "id": "3058",
            "name": "阿拉山口市"
        },
        {
            "city_id": "293",
            "id": "3059",
            "name": "吴起县"
        },
        {
            "city_id": "17",
            "id": "3060",
            "name": "上党区"
        },
        {
            "city_id": "17",
            "id": "3061",
            "name": "潞州区"
        },
        {
            "city_id": "43",
            "id": "3062",
            "name": "北镇市"
        },
        {
            "city_id": "213",
            "id": "3063",
            "name": "沙田镇"
        },
        {
            "city_id": "147",
            "id": "3064",
            "name": "兰陵县"
        },
        {
            "city_id": "167",
            "id": "3065",
            "name": "经济开发区"
        },
        {
            "city_id": "220",
            "id": "3066",
            "name": "荔浦市"
        },
        {
            "city_id": "193",
            "id": "3067",
            "name": "零陵区"
        },
        {
            "city_id": "346",
            "id": "3068",
            "name": "中沙群岛的岛礁及其海域"
        },
        {
            "city_id": "346",
            "id": "3069",
            "name": "西沙群岛"
        },
        {
            "city_id": "346",
            "id": "3070",
            "name": "南沙群岛"
        },
        {
            "city_id": "299",
            "id": "3071",
            "name": "新城镇"
        },
        {
            "city_id": "299",
            "id": "3072",
            "name": "镜铁区"
        },
        {
            "city_id": "299",
            "id": "3073",
            "name": "长城区"
        },
        {
            "city_id": "299",
            "id": "3074",
            "name": "雄关区"
        },
        {
            "city_id": "299",
            "id": "3075",
            "name": "文殊镇"
        },
        {
            "city_id": "299",
            "id": "3076",
            "name": "峪泉镇"
        },
        {
            "city_id": "299",
            "id": "3077",
            "name": "市辖区"
        }
    ]
}`