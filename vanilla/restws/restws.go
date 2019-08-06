package restws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kfchen81/beego"
	"github.com/kfchen81/beego/vanilla"
	"net/http"
	"strings"
	"time"
)

type RestWS struct {
	vanilla.RestResource
}

func (this *RestWS) Resource() string {
	return "restws"
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

func (this *RestWS) Get() {
	ws, err := upgrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)

	if err != nil {
		beego.Error(err)
		return
	}

	restwsGauge.Inc()

	bCtx := this.GetBusinessContext()
	bCtx = context.WithValue(bCtx, "isWsMode", true)
	reader(ws, bCtx)
}

func reader(ws *websocket.Conn, bCtx context.Context) {
	defer func() {
		ws.Close()
		restwsGauge.Dec()
	}()
	for {
		req := new(RestRequest)
		ws.SetReadDeadline(time.Now().Add(readWait))
		err := ws.ReadJSON(req)
		if err != nil {
			beego.Error("Read Error:", err)
			break
		}

		go handle(ws, bCtx, req)
	}
}

func handle(ws *websocket.Conn, bCtx context.Context, req *RestRequest) {
	defer func() {
		err := recover()
		if err != nil {
			beego.Error("Handle Error:", err)
		}
	}()
	log(req)
	resp := handleRequest(*req, bCtx)
	content, err := json.Marshal(resp)
	if err != nil {
		beego.Error("Encode websocket data error:", err)
		return
	}
	ws.SetWriteDeadline(time.Now().Add(writeWait))
	if err := ws.WriteMessage(websocket.TextMessage, content); err != nil {
		beego.Error("Write Error:", err)
		ws.Close()
		return
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
