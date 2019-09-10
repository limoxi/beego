package vanilla

import (
	"context"
	"strings"
)

var (
	REQUEST_MODE_PROD = "PROD"
	REQUEST_MODE_TEST = "TEST"
	REQUEST_HEADER_FORMAT = "Request-Mode"
)

type requestMode struct {
	define string
}

func (this *requestMode) String() string{
	return strings.ToUpper(this.define)
}

func (this *requestMode) IsTest() bool{
	return strings.HasSuffix(this.String(), "TEST")
}

func (this *requestMode) IsProd() bool{
	return strings.HasSuffix(this.String(), "PROD")
}

// GetRequestModeFromCtx 获取请求模式
// 默认prod
func GetRequestModeFromCtx(ctx context.Context) *requestMode{
	mode := new(requestMode)
	mode.define = REQUEST_MODE_PROD
	modeIf := ctx.Value("REQUEST_MODE")
	if modeIf == nil{
		return mode
	}
	mode.define = modeIf.(string)
	return mode
}