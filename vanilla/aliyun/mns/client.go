package mns

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"net"
	"encoding/xml"
)

const (
	version = "2015-06-06"
)

const (
	DefaultTimeout int64 = 35
)

type Method string


func init() {
}

const (
	GET    Method = "GET"
	PUT           = "PUT"
	POST          = "POST"
	DELETE        = "DELETE"
)

type MNSClient interface {
	Send(method Method, headers map[string]string, message interface{}, resource string) (resp *http.Response, err error)
	SetProxy(url string)
}

type AliMNSClient struct {
	Timeout      int64
	url          string
	credential   Credential
	accessKeyId  string
	clientLocker sync.Mutex
	client       *http.Client
	proxyURL     string
}

func NewAliMNSClient(url, accessKeyId, accessKeySecret string) MNSClient {
	if url == "" {
		panic("ali-mns: message queue url is empty")
	}
	
	credential := NewAliMNSCredential(accessKeySecret)
	
	aliMNSClient := new(AliMNSClient)
	aliMNSClient.credential = credential
	aliMNSClient.accessKeyId = accessKeyId
	aliMNSClient.url = url
	
	timeoutInt := DefaultTimeout
	
	if aliMNSClient.Timeout > 0 {
		timeoutInt = aliMNSClient.Timeout
	}
	
	timeout := time.Second * time.Duration(timeoutInt)
	
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	aliMNSClient.client = &http.Client{
		Timeout:   timeout,
		Transport: netTransport,
	}
	
	return aliMNSClient
}

func (p *AliMNSClient) SetProxy(url string) {
	p.proxyURL = url
}

func (p *AliMNSClient) proxy(req *http.Request) (*url.URL, error) {
	if p.proxyURL != "" {
		return url.Parse(p.proxyURL)
	}
	return nil, nil
}

func (p *AliMNSClient) authorization(method Method, headers map[string]string, resource string) (authHeader string, err error) {
	if signature, e := p.credential.Signature(method, headers, resource); e != nil {
		return "", e
	} else {
		authHeader = fmt.Sprintf("MNS %s:%s", p.accessKeyId, signature)
	}
	
	return
}

func (p *AliMNSClient) Send(method Method, headers map[string]string, message interface{}, resource string) (resp *http.Response, err error) {
	var xmlContent []byte
	
	if message == nil {
		xmlContent = []byte{}
	} else {
		switch m := message.(type) {
		case []byte:
			{
				xmlContent = m
			}
		default:
			if bXml, e := xml.Marshal(message); e != nil {
				err = err
				return
			} else {
				xmlContent = bXml
			}
		}
	}
	
	xmlMD5 := md5.Sum(xmlContent)
	strMd5 := fmt.Sprintf("%x", xmlMD5)
	
	if headers == nil {
		headers = make(map[string]string)
	}
	
	headers[MQ_VERSION] = version
	headers[CONTENT_TYPE] = "text/xml"
	headers[CONTENT_MD5] = base64.StdEncoding.EncodeToString([]byte(strMd5))
	headers[DATE] = time.Now().UTC().Format(http.TimeFormat)
	
	if authHeader, e := p.authorization(method, headers, fmt.Sprintf("/%s", resource)); e != nil {
		err = err
		return
	} else {
		headers[AUTHORIZATION] = authHeader
	}
	
	url := p.url + "/" + resource
	
	//beego.Notice(string(xmlContent))
	postBodyReader := strings.NewReader(string(xmlContent))
	
	// 莫名的lock 加这个是为了啥 想不通。。 推拉模式 加lock 这是直接限流请求了
	// p.clientLocker.Lock()
	// defer p.clientLocker.Unlock()
	
	var req *http.Request
	if req, err = http.NewRequest(string(method), url, postBodyReader); err != nil {
		err = err
		return
	}
	
	for header, value := range headers {
		req.Header.Add(header, value)
	}
	
	if resp, err = p.client.Do(req); err != nil {
		err = err
		return
	}
	
	return
}

//func ParseError(resp ErrorMessageResponse, resource string) (err error) {
//	if errCodeTemplate, exist := errMapping[resp.Code]; exist {
//		err = errCodeTemplate.New(errors.Params{"resp": resp, "resource": resource})
//	} else {
//		err = ERR_MNS_UNKNOWN_CODE.New(errors.Params{"resp": resp, "resource": resource})
//	}
//	return
//}
