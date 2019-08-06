package restws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kfchen81/beego"
	beecontext "github.com/kfchen81/beego/context"
	"github.com/kfchen81/beego/vanilla"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type WsResponse struct {
	*vanilla.Response
	Rid string `json:"rid"`
}

type fakeResponseWriter struct{}

func (f *fakeResponseWriter) Header() http.Header {
	return http.Header{}
}
func (f *fakeResponseWriter) Write(b []byte) (int, error) {
	return 0, nil
}
func (f *fakeResponseWriter) WriteHeader(n int) {}

func handleRequest(restReq RestRequest, bCtx context.Context) (resp WsResponse) {
	ctx := beecontext.NewContext()
	defer func() {
		err := recover()
		if err != nil {
			resp = RecoverPanic(err, ctx, restReq)
		}
	}()
	restReq.Method = strings.ToUpper(restReq.Method)

	cr := beego.BeeApp.Handlers

	if !strings.HasPrefix(restReq.Path, "/") {
		restReq.Path = fmt.Sprintf("/%s", restReq.Path)
	}
	if !strings.HasSuffix(restReq.Path, "/") {
		restReq.Path = fmt.Sprintf("%s/", restReq.Path)
	}
	req := &http.Request{
		URL:    &url.URL{Scheme: "http", Host: "localhost", Path: restReq.Path},
		Method: restReq.Method,
	}
	ctx.Reset(&fakeResponseWriter{}, req)
	ctx.Input.SetData("bContext", bCtx)

	formData := url.Values{}
	data := make(map[string]interface{}, 0)
	json.Unmarshal([]byte(restReq.Params), &data)
	for k, v := range data {
		formData.Set(k, fmt.Sprint(v))
	}
	req.Form = formData

	cron, findRouter := cr.FindRouter(ctx)
	if !findRouter {
		resp = WsResponse{
			Response: &vanilla.Response{
				Code: 404,
				Data: vanilla.Map{
					"endpoint": restReq.Path,
				},
				ErrMsg:  "",
				ErrCode: "restws:404",
			},
			Rid: restReq.Rid,
		}
		return
		//exception("404", context)
	}
	execController := cron.Init()
	execController.Init(ctx, "restws.request", req.Method, execController)

	//call prepare function
	execController.Prepare()
	switch req.Method {
	case "GET":
		execController.Get()
	case "POST":
		execController.Post()
	case "PUT":
		execController.Put()
	case "DELETE":
		execController.Delete()
	}
	execController.Finish()
	vcData := reflect.ValueOf(execController).Elem().FieldByName("Data")
	respData := vcData.Interface().(map[interface{}]interface{})["json"]
	resp = WsResponse{respData.(*vanilla.Response), restReq.Rid}
	return
}
