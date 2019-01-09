package vanilla

import (
	"fmt"
	"strings"

	"reflect"

	"github.com/kfchen81/beego"
	"os"
)

//RESOURCES 所有资源名的集合
var RESOURCES = make([]string, 0, 100)

var enableDevTestResource = (os.Getenv("ENABLE_DEV_TEST_RESOURCE") == "1")

//Router 添加路由
func Router(r RestResourceInterface) {
	//check whether is dev RESOURCE
	if r.IsForDevTest() && !enableDevTestResource {
		return
	}
	
	resource := r.Resource()
	RESOURCES = append(RESOURCES, resource)

	items := strings.Split(resource, ".")
	
	// html url
	if r.EnableHTMLResource() {
		url := fmt.Sprintf("/%s/", strings.Join(items, "/"))
		beego.Info(fmt.Sprintf("[resource]: %s -> %s", url, reflect.TypeOf(r)))
		beego.Router(url, r)
		return
	}
	
	//standard url
	{
		url := fmt.Sprintf("/%s/", strings.Join(items, "/"))
		beego.Info(fmt.Sprintf("[resource]: %s -> %s", url, reflect.TypeOf(r)))
		beego.Router(url, r)
	}

	// api url
	{
		lastIndex := len(items) - 1
		lastItem := items[lastIndex]
		items[lastIndex] = "api"

		itemSclie := items[:]
		itemSlice := append(itemSclie, lastItem)
		url := fmt.Sprintf("/%s/", strings.Join(itemSlice, "/"))
		beego.Info(fmt.Sprintf("[resource]: %s -> %s", url, reflect.TypeOf(r)))
		beego.Router(url, r)
	}
}

func init() {
	beego.Router("/console/console/", &ConsoleController{})
	beego.Router("/op/health/", &OpHealthController{})
	beego.Router("/", &IndexController{})
}