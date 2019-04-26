package vanilla

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
	
	"github.com/kfchen81/beego"
	"github.com/kfchen81/beego/orm"
	"github.com/bitly/go-simplejson"
	"os"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go"
)

var _PLATFORM_SECRET string
var _USER_LOGIN_SECRET string

type ResourceResponse struct {
	RespData *simplejson.Json
}

func (this *ResourceResponse) IsSuccess() bool {
	code, _ := this.RespData.Get("code").Int()
	return code == 200
}

func (this *ResourceResponse) Data() *simplejson.Json {
	return this.RespData.Get("data")
}

/*RestResource 扩展beego.Controller, 作为rest中各个资源的基类
 */
type Resource struct {
	Ctx context.Context
	CustomJWTToken string
}

func (this *Resource) request(method string, service string, resource string, data Map) (respData *ResourceResponse, err error) {
	var jwtToken string
	if this.CustomJWTToken != "" {
		jwtToken = this.CustomJWTToken
	} else {
		var ok bool
		if jwtToken, ok = this.Ctx.Value("jwt").(string); ok {
		
		} else {
			jwtToken = ""
		}
	}
	
	usePeanutPure := os.Getenv("USE_PEANUT_PURE")
	if usePeanutPure == "1" && service == "peanut" {
		service = "peanut_pure"
	}
	
	apiServerHost := beego.AppConfig.String("api::API_SERVER_HOST")
	//创建client
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	//构建url.Values
	params := url.Values{"_v": {"1"}}

	//处理resource
	pos := strings.LastIndexByte(resource, '.')
	resource = fmt.Sprintf("%s/%s", resource[:pos], resource[pos+1:])

	//构建request
	//bytes, _ := json.Marshal(ids)
	apiUrl := fmt.Sprintf("http://%s/%s/%s/", apiServerHost, service, resource)
	var req *http.Request
	if method == "GET" {
		for k, v := range data {
			value := ""
			switch t := v.(type) {
			case int:
				value = fmt.Sprintf("%d", v)
			case bool:
				value = fmt.Sprintf("%t", v)
			case string:
				value = v.(string)
			case float64:
				value = fmt.Sprintf("%f", v)
			default:
				beego.Warn("unknown type: ", t)
			}
			params.Set(k, value)
		}
		apiUrl += "?" + params.Encode()
		beego.Warn("apiUrl: ", apiUrl)
		//strings.NewReader(values.Encode())

		req, err = http.NewRequest("GET", apiUrl, nil)
	} else {
		if method == "PUT" {
			params.Set("_method", "put")
		} else if method == "DELETE" {
			params.Set("_method", "delete")
		}
		apiUrl += "?" + params.Encode()
		beego.Warn("apiUrl: ", apiUrl)

		values := url.Values{}
		for k, v := range data {
			value := ""
			switch t := v.(type) {
			case int:
				value = fmt.Sprintf("%d", v)
			case bool:
				value = fmt.Sprintf("%t", v)
			case string:
				value = v.(string)
			case float64:
				value = fmt.Sprintf("%f", v)
			default:
				beego.Warn("unknown type: ", t)
			}
			values.Set(k, value)
		}

		req, err = http.NewRequest("POST", apiUrl, strings.NewReader(values.Encode()))
	}
	//req, err := http.NewRequest("GET", apiUrl, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	if method != "GET" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("AUTHORIZATION", jwtToken)
	
	//inject open tracing
	span := opentracing.SpanFromContext(this.Ctx)
	if span != nil {
		ext.SpanKindRPCClient.Set(span)
		ext.HTTPUrl.Set(span, apiUrl)
		ext.HTTPMethod.Set(span, method)
		span.Tracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header),
		)
	}

	//执行request，获得response
	resp, err := netClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	//获取response的内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	jsonObj := new(simplejson.Json)
	err = jsonObj.UnmarshalJSON(body)
	if err != nil {
		return nil, err
	}

	resourceResp := new(ResourceResponse)
	resourceResp.RespData = jsonObj
	//fmt.Println(string(body))

	if resourceResp.IsSuccess() {
		return resourceResp, nil
	} else {
		beego.Error(jsonObj)
		errCode := jsonObj.Get("errCode")
		if errCode == nil {
			return resourceResp, errors.New("remote_service_error")
		} else {
			return resourceResp, errors.New(errCode.MustString())
		}
	}
}

func (this *Resource) Get(service string, resource string, data Map) (resp *ResourceResponse, err error) {
	return this.request("GET", service, resource, data)
}

func (this *Resource) Put(service string, resource string, data Map) (resp *ResourceResponse, err error) {
	return this.request("PUT", service, resource, data)
}

func (this *Resource) Post(service string, resource string, data Map) (resp *ResourceResponse, err error) {
	return this.request("POST", service, resource, data)
}

func (this *Resource) Delete(service string, resource string, data Map) (resp *ResourceResponse, err error) {
	return this.request("DELETE", service, resource, data)
}

func (this *Resource) LoginAs(username string) *Resource {
	if _PLATFORM_SECRET == "" {
		beego.Error("_PLATFORM_SECRET is '', Please set _PLATFORM_SECRET in your *.conf file")
		return nil
	}
	
	resp, err := this.Put("gskep", "login.logined_corp_user", Map{
		"username": username,
		"password": _PLATFORM_SECRET,
	})
	if err != nil {
		beego.Error(err)
		return nil
	}
	
	respData := resp.Data()
	this.CustomJWTToken, _ = respData.Get("sid").String()
	return this
}

func (this *Resource) LoginAsUser(unionid string) *Resource {
	if _USER_LOGIN_SECRET == "" {
		beego.Error("_USER_LOGIN_SECRET is '', Please set _USER_LOGIN_SECRET in your *.conf file")
		return nil
	}
	
	resp, err := this.Put("gskep", "login.logined_h5_user", Map{
		"unionid": unionid,
		"secret": _USER_LOGIN_SECRET,
	})
	if err != nil {
		beego.Error(err)
		return nil
	}
	
	respData := resp.Data()
	this.CustomJWTToken, _ = respData.Get("sid").String()
	return this
}

func (this *Resource) LoginAsManager() *Resource {
	return this.LoginAs("manager")
}

func CronLogin(o orm.Ormer) (*Resource, error) {
	apiServerHost := beego.AppConfig.String("api::API_SERVER_HOST")
	apiUrl := fmt.Sprintf("http://%s/skep/account/logined_corp_user", apiServerHost)
	params := url.Values{"_v": {"1"}}
	params.Set("_method", "put")
	apiUrl += "?" + params.Encode()
	
	values := url.Values{}
	
	values.Set("username", "manager")
	values.Set("password", "dc120c3e372d9ba9998a52c9d8edcdcb")
	
	resp, err := http.PostForm(apiUrl, values)
	
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	jsonObj := new(simplejson.Json)
	err = jsonObj.UnmarshalJSON(body)
	
	jsonData, err := jsonObj.Map()
	
	if err != nil {
		return nil, err
	}
	
	dataMap := jsonData["data"].(map[string]interface{})
	jwt := dataMap["sid"].(string)
	
	ctx := context.Background()
	span := opentracing.StartSpan("CRONTAB")
	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx = context.WithValue(ctx, "jwt", jwt)
	ctx = context.WithValue(ctx, "orm", o)
	resource := NewResource(ctx)
	resource.CustomJWTToken = jwt
	
	return resource, nil
}

func NewResource(ctx context.Context) *Resource {
	resource := new(Resource)
	resource.Ctx = ctx
	return resource
}

func ToJsonString(obj interface{}) string {
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}

func init() {
	_PLATFORM_SECRET = beego.AppConfig.String("system::PLATFORM_SECRET")
	_USER_LOGIN_SECRET = beego.AppConfig.String("system::USER_LOGIN_SECRET")
	beego.Info("[init] use _PLATFORM_SECRET: " + _PLATFORM_SECRET)
	beego.Info("[init] use _USER_LOGIN_SECRET: " + _USER_LOGIN_SECRET)
}
