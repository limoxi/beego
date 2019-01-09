package aliyun

import (
	"gskep/vanilla/aliyun/green"
	"gskep/vanilla/uuid"
	"github.com/bitly/go-simplejson"
	"github.com/kfchen81/beego"
	"encoding/json"
)

const MAX_RATE = 90.0

type GreenService struct {
}

const accessKeyId string = "LTAIO6NdtE5IWTIC"
const accessKeySecret string = "tmknbdsCUkC1212eSozF63keM9LcQc"

func (this *GreenService) shouldBlock(response string) bool {
	jsonObj := new(simplejson.Json)
	err := jsonObj.UnmarshalJSON([]byte(response))
	if err != nil {
		beego.Error(err)
		return false
	}
	beego.Debug(jsonObj)
	
	code, _ := jsonObj.Get("code").Int()
	if code != 200 {
		msg, _ := jsonObj.Get("msg").String()
		beego.Warn("[aliyun] green image check fail: ", msg)
		return false
	}
	
	var shouldBlock bool = false
	checkResults := jsonObj.Get("data").MustArray()
	for _, checkResult := range checkResults {
		checkResultMap := checkResult.(map[string]interface{})
		code, _ := checkResultMap["code"].(json.Number).Int64()
		if code == 200 {
			sceneResults := checkResultMap["results"].([]interface{})
			for _, sceneResult := range sceneResults {
				sceneResultMap := sceneResult.(map[string]interface{})
				suggestion := sceneResultMap["suggestion"].(string)
				if suggestion == "block" {
					shouldBlock = true
					break
				}
				
				rate, _ := sceneResultMap["rate"].(json.Number).Float64()
				if suggestion == "view" && rate > MAX_RATE {
					shouldBlock = true
					break
				}
			}
		}
	}
	
	if shouldBlock {
		return true
	} else {
		return false
	}
}

func (this *GreenService) ShouldBlockImage(url string) bool {
	profile := green.Profile{AccessKeyId:accessKeyId, AccessKeySecret:accessKeySecret}
	path := "/green/image/scan";
	
	clientInfo := green.ClinetInfo{Ip:"127.0.0.1"}
	bizType := "Green"
	scenes := []string{"porn"}
	
	task := green.Task{DataId:uuid.Rand().Hex(), Url:url}
	tasks := []green.Task{task}
	
	bizData := green.BizData{ bizType, scenes, tasks}
	
	var client green.IAliYunClient = green.DefaultClient{Profile:profile}
	
	resp := client.GetResponse(path, clientInfo, bizData)

	return this.shouldBlock(resp)
}

func (this *GreenService) ShouldBlockImages(urls []string) bool {
	profile := green.Profile{AccessKeyId:accessKeyId, AccessKeySecret:accessKeySecret}
	path := "/green/image/scan";
	
	clientInfo := green.ClinetInfo{Ip:"127.0.0.1"}
	bizType := "Green"
	scenes := []string{"porn"}
	
	tasks := make([]green.Task, 0)
	
	for _, url := range urls {
		task := green.Task{DataId:uuid.Rand().Hex(), Url:url}
		tasks = append(tasks, task)
	}
	
	bizData := green.BizData{ bizType, scenes, tasks}
	
	var client green.IAliYunClient = green.DefaultClient{Profile:profile}
	
	resp := client.GetResponse(path, clientInfo, bizData)
	
	return this.shouldBlock(resp)
}

func (this *GreenService) ShouldBlockText(content string) bool {
	profile := green.Profile{AccessKeyId:accessKeyId, AccessKeySecret:accessKeySecret}
	path := "/green/text/scan";
	
	clientInfo := green.ClinetInfo{Ip:"127.0.0.1"}
	bizType := "Green"
	scenes := []string{"antispam"}
	
	task := green.Task{DataId:uuid.Rand().Hex(), Content:content}
	tasks := []green.Task{task}
	
	bizData := green.BizData{ bizType, scenes, tasks}
	
	var client green.IAliYunClient = green.DefaultClient{Profile:profile}
	
	resp := client.GetResponse(path, clientInfo, bizData)
	
	return this.shouldBlock(resp)
}

func NewGreenService() *GreenService {
	return new(GreenService)
}
