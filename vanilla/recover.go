package vanilla

import (
	raven "github.com/getsentry/raven-go"
	"github.com/kfchen81/beego/context"
	"github.com/kfchen81/beego/metrics"

	"bytes"
	"fmt"
	"runtime"
	"strings"

	"github.com/kfchen81/beego"
	"github.com/kfchen81/beego/logs"
	"github.com/kfchen81/beego/orm"
)

func isEnableSentry() bool {
	return beego.AppConfig.DefaultBool("sentry::ENABLE_SENTRY", false)
}

// CapturePanicToSentry will collect error info then send to sentry
func CapturePanicToSentry(ctx *context.Context, err error) {
	if !isEnableSentry() {
		return
	}
	rvalStr := fmt.Sprint(err)
	skipFramesCount := beego.AppConfig.DefaultInt("sentry::SKIP_FRAMES_COUNT", 3)
	contextLineCount := beego.AppConfig.DefaultInt("sentry::CONTEXT_LINE_COUNT", 5)
	appRootPath := beego.AppConfig.String("appname")
	inAppPaths := []string{appRootPath}

	var sStacktrace *raven.Stacktrace
	var sError, ok = err.(error)
	if ok {
		// if sError, ok := err.(error); ok {
		sStacktrace = raven.GetOrNewStacktrace(sError, skipFramesCount, contextLineCount, inAppPaths)
	} else {
		sStacktrace = raven.NewStacktrace(skipFramesCount, contextLineCount, inAppPaths)
	}
	sException := raven.NewException(sError, sStacktrace)
	var packet *raven.Packet
	packet = raven.NewPacket(rvalStr, sException, raven.NewHttp(ctx.Request))
	raven.Capture(packet, nil)
}

func RecoverPanic(ctx *context.Context) {
	if err := recover(); err != nil {
		//rollback tx
		o := ctx.Input.Data()["sessionOrm"]
		if o != nil {
			o.(orm.Ormer).Rollback()
			beego.Warn("[ORM] rollback transaction")
		}

		//finish span
		span := ctx.Input.GetData("span")
		if span != nil {
			beego.Info("[Tracing] finish span in recoverPanic")
			span.(opentracing.Span).Finish()
		}

		//释放锁
		if mutex, ok := ctx.Input.Data()["sessionRestMutex"]; ok {
			if mutex != nil {
				beego.Debug("[lock] release resource lock @2")
				mutex.(*redsync.Mutex).Unlock()
			}
		}

		//记录panic counter
		//1. 非BusinessError需要记录
		//2. IsPanicError为true的BusinessError需要记录
		if be, ok := err.(*BusinessError); ok {
			if be.IsPanicError() {
				metrics.GetPanicCounter().Inc()
			}
		} else {
			metrics.GetPanicCounter().Inc()
		}

		if err == beego.ErrAbort {
			return
		}
		if !beego.BConfig.RecoverPanic {
			panic(err)
		}

		errMsg := ""
		if be, ok := err.(*BusinessError); ok {
			errMsg = fmt.Sprintf("%s:%s", be.ErrCode, be.ErrMsg)
		} else {
			errMsg = fmt.Sprintf("%s", err)
		}

		//log error info
		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf("[Unprocessed_Exception] %s\n", errMsg))
		buffer.WriteString(fmt.Sprintf("Request URL: %s\n", ctx.Input.URL()))
		for i := 1; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			buffer.WriteString(fmt.Sprintf("%s:%d\n", file, line))
		}
		if beego.BConfig.RunMode == "dev" {
			logs.Critical(buffer.String())
		} else {
			logs.Critical(strings.Replace(buffer.String(), "\n", ";  ", -1))
		}

		//return error response
		var resp Map
		if be, ok := err.(*BusinessError); ok {
			resp = Map{
				"code":        500,
				"data":        nil,
				"errCode":     be.ErrCode,
				"errMsg":      be.ErrMsg,
				"innerErrMsg": "",
			}
		} else {
			endpoint := ctx.Request.RequestURI
			pos := strings.Index(endpoint, "?")
			if pos != -1 {
				endpoint = endpoint[:pos]
			}
			resp = Map{
				"code": 531,
				"data": Map{
					"endpoint": endpoint,
				},
				"errCode":     "system:exception",
				"errMsg":      fmt.Sprintf("%s", err),
				"innerErrMsg": "",
			}
		}
		ctx.Output.JSON(resp, true, true)
	}
}

func init() {
	if isEnableSentry() {
		raven.SetDSN(beego.AppConfig.String("sentry::SENTRY_DSN"))
	}
}
