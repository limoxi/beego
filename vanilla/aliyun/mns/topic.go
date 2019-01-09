package mns

import (
	"fmt"
)

var (
	DefaultNumOfMessages int32 = 16
	DefaultQPSLimit      int32 = 2000
)

const (
	PROXY_PREFIX = "MNS_PROXY_"
	GLOBAL_PROXY = "MNS_GLOBAL_PROXY"
)

type AliMNSTopic interface {
	Name() string
	SendMessage(message TopicMessageSendRequest) (resp MessageSendResponse, err error)
}

type MNSTopic struct {
	name       string
	client     MNSClient
	qpsLimit   int32
	decoder    MNSDecoder
}

func NewMNSTopic(name string, client MNSClient, qps ...int32) AliMNSTopic {
	if name == "" {
		panic("ali_mns: queue name could not be empty")
	}

	topic := new(MNSTopic)
	topic.client = client
	topic.name = name
	topic.qpsLimit = DefaultQPSLimit
	topic.decoder = NewAliMNSDecoder()

	if qps != nil && len(qps) == 1 && qps[0] > 0 {
		topic.qpsLimit = qps[0]
	}

	return topic
}

func (p *MNSTopic) Name() string {
	return p.name
}

func (p *MNSTopic) SendMessage(message TopicMessageSendRequest) (resp MessageSendResponse, err error) {
	//p.checkQPS()
	_, err = send(p.client, p.decoder, POST, nil, message, fmt.Sprintf("topics/%s/messages", p.name), &resp)
	return
}

func (p *MNSTopic) checkQPS() {
	//p.qpsMonitor.Pulse()
	//if p.qpsLimit > 0 {
	//	for p.qpsMonitor.QPS() > p.qpsLimit {
	//		time.Sleep(time.Millisecond * 10)
	//	}
	//}
}
