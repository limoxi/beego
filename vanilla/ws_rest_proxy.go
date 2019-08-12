package vanilla

import (
	"encoding/json"
	"fmt"
	"github.com/kfchen81/beego/metrics"
	"github.com/gorilla/websocket"
	"github.com/kfchen81/beego"
	beeContext "github.com/kfchen81/beego/context"
	"net/http"
	"strings"
	"time"
)

type RestProxy struct {
	RestResource
}

func (this *RestProxy) Resource() string {
	return "ws.rest_proxy"
}

const (
	// Time allowed to write data to the client.
	writeWait = 10 * time.Second

	// Time allowed to read data to the client.
	readWait = 90 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = readWait

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 8) / 10
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type RestRequest struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	Params string `json:"params"`
	Rid    string `json:"rid"`
}

func (this *RestProxy) Get() {
	ws, err := upgrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)

	if err != nil {
		beego.Error(err)
		return
	}

	metrics.GetRestwsGauge().Inc()

	respChan := make(chan WsResponse)
	defer close(respChan)
	go write(ws, respChan)
	reader(ws, this.Ctx, respChan)

}

func reader(ws *websocket.Conn, ctx *beeContext.Context, respChan chan<- WsResponse) {
	defer func() {
		ws.Close()
		metrics.GetRestwsGauge().Dec()
	}()
	for {
		req := new(RestRequest)
		ws.SetReadDeadline(time.Now().Add(readWait))
		err := ws.ReadJSON(req)
		if err != nil {
			beego.Info("Read Error:", err)
			break
		}

		go handle(ctx, req, respChan)
	}
}

func handle(ctx *beeContext.Context, req *RestRequest, respChan chan<- WsResponse) {
	defer func() {
		err := recover()
		if err != nil {
			beego.Error("Handle Error:", err)
		}
	}()
	log(req)
	resp := handleRequest(*req, ctx)
	respChan <- resp
}

func write(ws *websocket.Conn, respChan <-chan WsResponse) {
	defer func() {
		err := recover()
		if err != nil {
			beego.Error("Write Recover:", err)
		}
	}()
	for resp := range respChan{
		content, err := json.Marshal(resp)
		if err != nil {
			beego.Error("Encode websocket data error:", err)
			return
		}
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.TextMessage, content); err != nil {
			beego.Info("Write Error:", err)
			ws.Close()
			return
		}
	}
}

func log(req *RestRequest) {
	now := time.Now().Format("2006-01-02 15:04:05")
	if !strings.HasPrefix(req.Path, "/") {
		req.Path = fmt.Sprintf("/%s", req.Path)
	}
	beego.Info(fmt.Sprintf("[%s] Method:%s Path:%s Params:%s Rid:%s",
		now, req.Method, req.Path, req.Params, req.Rid))
}

