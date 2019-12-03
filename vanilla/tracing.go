package vanilla

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"io"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go"
	"github.com/kfchen81/beego"
	"time"
)

var Tracer opentracing.Tracer
var Closer io.Closer

func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	tracingMode := beego.AppConfig.DefaultString("tracing::MODE", "dev")
	var cfg *config.Configuration
	
	if tracingMode == "dev" {
		cfg = &config.Configuration{
			Sampler: &config.SamplerConfig{
				Type:  "const",
				Param: 1,
			},
			Reporter: &config.ReporterConfig{
				LogSpans: true,
				BufferFlushInterval: 1 * time.Second,
			},
		}
	} else if tracingMode == "disable" {
		cfg = &config.Configuration{
			Disabled: true,
		}
		beego.Warn("Open Tracing is disabled!!")
	} else {
		cfg = &config.Configuration{
			Sampler: &config.SamplerConfig{
				Type:  "probabilistic",
				Param: 0.2,
			},
			Reporter: &config.ReporterConfig{
				LogSpans: false,
				BufferFlushInterval: 5 * time.Second,
			},
		}
	}
	
	tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

func init() {
	serviceName := beego.AppConfig.DefaultString("appname", beego.BConfig.AppName)
	Tracer, Closer = initJaeger(serviceName)
	beego.Debug("[tracing] Tracer ", Tracer)
}