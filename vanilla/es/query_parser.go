package es

import (
	"github.com/deckarep/golang-set"
	"github.com/kfchen81/beego/vanilla"
	"strings"
)

type Query struct{
	data map[string]interface{}
}

func (this *Query) SetData(data map[string]interface{}){
	this.data = data
}

func (this *Query) RemoveKey(key string){
	delete(this.data, key)
}

func (this Query) Source() (interface{}, error){
	return this.data, nil
}

// 将rest filter语法转换成es query语法、sort语法
type QueryParser struct{
	vanilla.ServiceBase

	query *Query
	mustArray []map[string]interface{}
	mustNotArray []map[string]interface{}

}


func NewQueryParser() *QueryParser{
	parser := new(QueryParser)
	parser.mustArray = make([]map[string]interface{}, 0)
	parser.mustNotArray = make([]map[string]interface{}, 0)
	parser.query = new(Query)
	parser.query.SetData(map[string]interface{}{
		"match_all": map[string]interface{}{},
		"bool": map[string]interface{}{
			"filter": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": &parser.mustArray,
					"must_not": &parser.mustNotArray,
				},
			},
		},
	})
	return parser
}


/**
处理 [range, icontains, gte, lte, range, notin, not, in, gt, lt]
	形如：{
		'bid': '1242464682341',
		'createdAt__range': ['2018-07-08 12:12:22', '2018-07-09 13:12:22'],
		'finalMoney__gte': 50,
		'corp_name__contains': u'小',
		'status__in': [1, 3],
	} ===>>>

	{
		'bool': {
			'must': [{
				'term': {
					'bid': '1242464682341'
				}
			}, {
				'range': {
					'created_at': {
						'gte': '2018-07-08 12:12:22',
						'lte': '2018-07-09 13:12:22'
					}
				}
			}, {
				'range': {
					'final_money': {
						'gte': 50
					}
				}
			}, {
				'match': {
					'corp.name': u'小'
				}
			}, {
				'terms': {
					'status': [1, 3]
				}
			}]
		}
	}
 */
func(this *QueryParser) Parse(filters map[string]interface{}) *Query{
	if filters == nil{
		this.query.RemoveKey("bool")
	}else{
		this.query.RemoveKey("match_all")
	}

	var nestPath2Query map[string]map[string]interface{}
	for k, v := range filters{
		realKey := strings.Split(k, "__")[0]
		isNested := len(strings.Split(realKey, ".")) > 1

		if len(strings.Split(k, ">")) > 1{
			k = strings.Replace(k, ">", ".", 1)
		}

		mustNode, mustNotNode := this.makeQuery(k, v)

		var nestedPath string
		nestMustArray := make([]interface{}, 0)
		nestMustNotArray := make([]interface{}, 0)

		if isNested{
			nestedPath = strings.Split(realKey, ".")[0]
			nestPath2Query = map[string]map[string]interface{}{
				nestedPath: {
					"must": &nestMustArray,
					"must_not": &nestMustNotArray,
				},
			}
		}

		if mustNode != nil{
			if !isNested{
				this.mustArray = append(this.mustArray, mustNode)
			}else{
				nestMustArray = append(nestMustArray, mustNode)
			}
		}

		if mustNotNode != nil{
			if !isNested{
				this.mustNotArray = append(this.mustNotArray, mustNotNode)
			}else{
				nestMustNotArray = append(nestMustNotArray, mustNotNode)
			}
		}
	}

	for nestPath, nestQuery := range nestPath2Query{
		if len(nestQuery["must"].([]interface{})) > 0{
			this.mustArray = append(this.mustArray, map[string]interface{}{
				"nested": map[string]interface{}{
					"path": nestPath,
					"query": map[string]interface{}{
						"bool": map[string]interface{}{
							"must": nestQuery["must"],
						},
					},
				},
			})
		}
		if len(nestQuery["must_not"].([]interface{})) > 0{
			this.mustNotArray = append(this.mustNotArray, map[string]interface{}{
				"nested": map[string]interface{}{
					"path": nestPath,
					"query": map[string]interface{}{
						"bool": map[string]interface{}{
							"must_not": nestQuery["must_not"],
						},
					},
				},
			})
		}
	}

	return this.query
}

func(this *QueryParser) makeQuery(k string, v interface{})(map[string]interface{}, map[string]interface{}){
	splits := strings.Split(k, "__")
	var mustNode map[string]interface{}
	var mustNotNode map[string]interface{}

	if len(splits) == 1{
		mustNode = 	map[string]interface{}{
			"term": map[string]interface{}{
				k: v,
			},
		}
	}else if len(splits) == 2{
		key := splits[0]
		op := splits[1]

		if mapset.NewSetFromSlice([]interface{}{"icontain", "contains", "contain"}).Contains(op){
			mustNode = 	map[string]interface{}{
				"match_phrase": map[string]interface{}{
					key: v,
				},
			}
		}else if op == "range"{
			gte := v.([]interface{})[0]
			lte := v.([]interface{})[1]
			mustNode = 	map[string]interface{}{
				"range": map[string]map[string]interface{}{
					key: {
						"gte": gte,
						"lte": lte,
					},
				},
			}
		}else if op == "in"{
			mustNode = 	map[string]interface{}{
				"terms": map[string]interface{}{
					key: v.([]interface{}),
				},
			}
		}else if op == "lt" || op == "gt"{
			mustNode = 	map[string]interface{}{
				"range": map[string]map[string]interface{}{
					key: {
						op: v,
					},
				},
			}
		}else if op == "notin"{
			mustNotNode = map[string]interface{}{
				"terms": map[string]interface{}{
					key: v,
				},
			}
		}else if op == "not"{
			mustNotNode = map[string]interface{}{
				"term": map[string]interface{}{
					key: v,
				},
			}
		}
	}

	return mustNode, mustNotNode
}