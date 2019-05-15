package es

import (
	"context"
	"fmt"
	"github.com/kfchen81/beego"
	"github.com/kfchen81/beego/vanilla"
	"github.com/olivere/elastic"
	"strings"
)

type ESClient struct{
	vanilla.ServiceBase

	client *elastic.Client

	indexName string
	docType string
}

func (this *ESClient) Use(indexName string) *ESClient{
	this.indexName = indexName
	return this
}

func (this *ESClient) Select(docType string) *ESClient{
	this.docType = docType
	return this
}

func (this *ESClient) prepareUpdateData(updateService *elastic.UpdateByQueryService, data map[string]interface{}, filters map[string]interface{}){
	// make query
	parser := NewQueryParser()
	query := parser.Parse(filters)
	updateService.Query(query)
	// make script
	script := make([]string, 0)
	for k, _ := range data{
		script = append(script, fmt.Sprintf("ctx._source.%s=params.%s", k, k))
	}
	scriptStr := strings.Join(script, ";")
	eScript := elastic.NewScript(scriptStr)
	eScript.Lang("painless").Params(data)
	updateService.Script(eScript)
	// make params
	updateService.Conflicts("proceed").Slices(50).WaitForCompletion(false).Size(-1)
}

func (this *ESClient) Update(data map[string]interface{}, filters map[string]interface{}){
	updateService := this.client.UpdateByQuery(this.indexName).Type(this.docType)
	this.prepareUpdateData(updateService, data, filters)
	_, err := updateService.Do(this.Ctx)
	if err != nil{
		beego.Error(err)
		panic(vanilla.NewSystemError("es_update:failed", "更新索引数据失败"))
	}
}

func (this *ESClient) Push(id string, data interface{}) {
	// Add a document
	indexResult, err := this.client.Index().
		Index(this.indexName).
		Type(this.docType).
		Id(id).
		BodyJson(&data).
		Do(this.Ctx)
	if err != nil {
		errMsg := fmt.Errorf("es_push doc(id:%d) to index %s: %v", id, this.indexName, err)
		beego.Error(errMsg)
		panic(vanilla.NewSystemError("es_push:failed", errMsg.Error()))
	}
	if indexResult == nil {
		errMsg := fmt.Errorf("es_push doc(id:%d) to index %s: result is %v",
			id, this.indexName, indexResult)
		beego.Error(errMsg)
		panic(vanilla.NewSystemError("es_push:failed", errMsg.Error()))
	}
}

func NewESClient(ctx context.Context) *ESClient{

	host := beego.AppConfig.String("es::ES_SEARCH_HOST")
	port := beego.AppConfig.String("es::ES_SEARCH_PORT")
	user := beego.AppConfig.String("es::ES_AUTH_USER")
	pwd := beego.AppConfig.String("es::ES_AUTH_SECRET")

	beego.Info(host, port, user, pwd)

	optionFunc := func (c *elastic.Client) error{
		var err error
		if user != ""{
			err = elastic.SetBasicAuth(user, pwd)(c)
		}
		err = elastic.SetURL("http://"+host+":"+port)(c)
		return err
	}

	client := new(ESClient)
	client.Ctx = ctx
	c, err := elastic.NewSimpleClient(optionFunc)
	if err != nil{
		beego.Error(err)
		panic(vanilla.NewSystemError("es:link_failed", "连接es服务失败"))
	}
	client.client = c
	return client
}

