package beego

import (
	go_context "context"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/kfchen81/beego/context"
	"github.com/kfchen81/beego/logs"
	"runtime/debug"
	
	"os"
)

func isEnableSentry() bool {
	return AppConfig.DefaultBool("sentry::ENABLE_SENTRY", false)
}

// CapturePanicToSentry will collect error info then send to sentry
func CaptureErrorToSentry(ctx *context.Context, err error) {
	if !isEnableSentry() {
		beegoMode := os.Getenv("BEEGO_RUNMODE")
		if beegoMode == "prod" {
			Warn("Sentry is not enabled under prod mode, Please enable it!!!!")
		}
		return
	}
	
	rvalStr := fmt.Sprint(err)
	skipFramesCount := AppConfig.DefaultInt("sentry::SKIP_FRAMES_COUNT", 3)
	contextLineCount := AppConfig.DefaultInt("sentry::CONTEXT_LINE_COUNT", 5)
	appRootPath := AppConfig.String("appname")
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
	if ctx == nil {
		packet = raven.NewPacket(rvalStr, sException)
	} else {
		packet = raven.NewPacket(rvalStr, sException, raven.NewHttp(ctx.Request))
	}
	
	raven.Capture(packet, nil)
}

func CaptureTaskErrorToSentry(ctx go_context.Context, errMsg string, taskErrMsg string) {
	if !isEnableSentry() {
		return
	}
	tags := map[string]string{
		"task_name": ctx.Value("taskName").(string),
		"service_name": AppConfig.String("appname"),
	}
	var packet *raven.Packet
	
	packet = raven.NewPacket(errMsg)
	stack := string(debug.Stack())
	packet.Extra = map[string]interface{}{
		"errMsg": taskErrMsg,
		"stacktrace": stack,
	}
	raven.Capture(packet, tags)
	// local log
	if BConfig.RunMode == "dev"{
		logs.Critical(stack)
	}
}

func PushErrorToSentry(errMsg string) {
	if !isEnableSentry() {
		return
	}
	tags := map[string]string{
		"service_name": AppConfig.String("appname"),
	}
	var packet *raven.Packet
	
	packet = raven.NewPacket(errMsg)
	stack := string(debug.Stack())
	packet.Extra = map[string]interface{}{
		"stacktrace": stack,
	}
	raven.Capture(packet, tags)
}

func init() {
	if isEnableSentry() {
		raven.SetDSN(AppConfig.String("sentry::SENTRY_DSN"))
		Info(fmt.Sprintf("[sentry] enable:%t, dsn:%s ", isEnableSentry(), AppConfig.String("sentry::SENTRY_DSN")))
	}
}
