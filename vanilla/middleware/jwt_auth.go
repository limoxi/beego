package middleware

import (
	"github.com/kfchen81/beego/context"

	go_context "context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/kfchen81/beego/vanilla"

	"github.com/bitly/go-simplejson"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/kfchen81/beego"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/kfchen81/beego/orm"
)

var SALT string = "030e2cf548cf9da683e340371d1a74ee"
var SKIP_JWT_CHECK_URLS []string = make([]string, 0)

var JWTAuthFilter = func(ctx *context.Context) {
	uri := ctx.Request.RequestURI
	
	if uri == "/" {
		if gBContextFactory != nil {
			bCtx := gBContextFactory.NewContext(go_context.Background(), ctx.Request, 0, "", nil) //bCtx is for "business context"
			o := orm.NewOrm()
			bCtx = go_context.WithValue(bCtx, "orm", o)
			ctx.Input.SetData("bContext", bCtx)
		}
		return
	}
	
	for _, skipUrl := range SKIP_JWT_CHECK_URLS {
		if strings.Contains(uri, skipUrl) {
			beego.Debug("[jwt_middleware] skip jwt check", "url", skipUrl)
			if gBContextFactory != nil {
				bCtx := gBContextFactory.NewContext(go_context.Background(), ctx.Request, 0, "", nil) //bCtx is for "business context"
				o := orm.NewOrm()
				bCtx = go_context.WithValue(bCtx, "orm", o)
				ctx.Input.SetData("bContext", bCtx)
			}
			return
		}
	}

	jwtToken := ctx.Input.Header("AUTHORIZATION")
	
	if jwtToken == "" {
		jwtToken = ctx.Input.Query("_jwt")
	}
	
	if jwtToken != "" {
		items := strings.Split(jwtToken, ".")
		if len(items) != 3 {
			//jwt token 格式不对
			response := vanilla.MakeErrorResponse(500, "jwt:invalid_jwt_token", fmt.Sprintf("无效的jwt token 1 - [%s]", jwtToken))
			ctx.Output.JSON(response, true, false)
			return
		}

		headerB64Code, payloadB64Code, expectedSignature := items[0], items[1], items[2]
		message := fmt.Sprintf("%s.%s", headerB64Code, payloadB64Code)

		h := hmac.New(sha256.New, []byte(SALT))
		h.Write([]byte(message))
		actualSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

		if expectedSignature != actualSignature {
			//jwt token的signature不匹配
			response := vanilla.MakeErrorResponse(500, "jwt:invalid_jwt_token", fmt.Sprintf("无效的jwt token 2 - [%s]", jwtToken))
			ctx.Output.JSON(response, true, false)
			return
		}

		decodeBytes, err := base64.StdEncoding.DecodeString(payloadB64Code)
		if err != nil {
			log.Fatalln(err)
		}
		js, err := simplejson.NewJson([]byte(decodeBytes))

		if err != nil {
			response := vanilla.MakeErrorResponse(500, "jwt:invalid_jwt_token", fmt.Sprintf("无效的jwt token 3 - [%s]", jwtToken))
			ctx.Output.JSON(response, true, false)
			return
		}

		//beego.Warn("**********>>>>>>>>>>**********>>>>>>>>>>")
		//beego.Warn(string(decodeBytes))
		//beego.Warn("**********>>>>>>>>>>**********>>>>>>>>>>")
		jwtType, err := js.Get("type").Int()
		if err != nil {
			log.Fatalln(err)
			response := vanilla.MakeErrorResponse(500, "jwt:invalid_jwt_token", fmt.Sprintf("无效的jwt token 4.1 - [%s]", jwtToken))
			ctx.Output.JSON(response, true, false)
			return
		}
		
		var userId int
		var err2 error
		if jwtType == 1 {
			userId, err2 = js.Get("user_id").Int()
		} else if jwtType == 2 {
			userId, err2 = js.Get("uid").Int()
		} else if jwtType == 3 {
			userId, err2 = js.Get("user").Get("uid").Int()
		} else {
			err2 = errors.New(fmt.Sprintf("invalid jwt type: %d", jwtType))
		}
		if err2 != nil {
			beego.Error(err2)
			response := vanilla.MakeErrorResponse(500, "jwt:invalid_jwt_token", fmt.Sprintf("无效的jwt token 4.2 - [%s]", jwtToken))
			ctx.Output.JSON(response, true, false)
			return
		}
		
		bCtx := gBContextFactory.NewContext(go_context.Background(), ctx.Request, userId, jwtToken, js) //bCtx is for "business context"
		
		//enhance business context
		{
			//add tracing span
			spanCtx, _ := vanilla.Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(ctx.Request.Header))
			uri := ctx.Request.URL.Path
			operationName := fmt.Sprintf("%s %s", ctx.Request.Method, uri)
			span := vanilla.Tracer.StartSpan(operationName, ext.RPCServerOption(spanCtx))
			bCtx = opentracing.ContextWithSpan(bCtx, span)
			
			//add orm
			bCtx = go_context.WithValue(bCtx, "jwt", jwtToken)
			o := orm.NewOrmWithSpan(span)
			bCtx = go_context.WithValue(bCtx, "orm", o)
		}
		
		ctx.Input.SetData("bContext", bCtx)
		ctx.Input.SetData("span", opentracing.SpanFromContext(bCtx))
	} else {
		response := vanilla.MakeErrorResponse(500, "jwt:invalid_jwt_token", fmt.Sprintf("无效的jwt token 5 - [%s]", jwtToken))
		ctx.Output.JSON(response, true, false)
		return
	}

}

func init() {
	skipUrls := beego.AppConfig.String("SKIP_JWT_CHECK_URLS")
	if skipUrls == "" {
		beego.Info("SKIP_JWT_CHECK_URLS is empty")
	} else {
		SKIP_JWT_CHECK_URLS = strings.Split(skipUrls, ";")
	}
	
	beego.Info("SKIP_JWT_CHECK_URLS: ", SKIP_JWT_CHECK_URLS)
}