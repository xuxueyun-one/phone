package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"log"
	"github.com/tidwall/gjson"
)

const  (
	GetImsUrl = "http://218.205.115.239:8080/fy/phone/v2/ims/freeAuth/public/authed/info"
	GetPrefixUrl = "http://218.205.115.239:8080/fy/phone/v2/ims/freeAuth/public/rule"
)

func init(){
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}


func main() {
	// 同步通讯录
	http.HandleFunc("/syn", AscRecord)
	// 获取ims
	http.HandleFunc("/ims",GetIms)
	// 获取规则
	http.HandleFunc("/prefix",GetPrefix)

	// 音响打电话  传入 "给马宁打电话"
	http.HandleFunc("/dial",Dial)


	log.Fatal(http.ListenAndServe(":9090",nil))
}

type Record struct {
	Msisdn 		string		`json:"msisdn"`// ims号码
	Persons 	[]*Person	`json:"persons"`// 联系人
}
type Person struct {
	NickName 	string	`json:"nickname"`// 联系人昵称
	Number		string	`json:"number"`// 联系人电话
	Status		interface{}	`json:"status"`// 联系人状态
}

/*type ImsData struct {
	ImsNum		string		// ims号码
	ImsAccount 	string		// sbc账号
	Password	string		// sbc密码
	Sbc			string		// sbc地址
	Port 		string		// sbc端口
	Domain		string		// 域名
}*/


func AscRecord(writer http.ResponseWriter, request *http.Request) {

	// 获取信息
	request.ParseForm()
	log.Println("同步通讯录参数为：",request.Form)
	if len(request.Form["msisdn"]) == 0 || len(request.Form["persons"]) == 0{
		writer.Write([]byte("1010"))
		return
	}
	//msisdn := request.Form["msisdn"][0]
	persons := request.Form["persons"][0]

	// 进行同步

	log.Println("同步通讯录完成")

	parray := []Person{}
	err := json.Unmarshal([]byte(persons),&parray)
	if err != nil {
		log.Println("解析json字符串错误")
	}

	var rs = struct {
		Persons interface{}	`json:"persons"`
		Recode	int64		`json:"recode"`
	}{parray,1}

	rsb,err := json.Marshal(rs)
	if err != nil {
		log.Println("解析json字符串错误")
	}

	// 返回结果
	_,err = writer.Write(rsb)
	if err != nil{
		log.Println("同步通讯录时返回数据错误")
		return
	}


}

func GetIms(writer http.ResponseWriter, request *http.Request)  {

	// 	获取请求参数
	request.ParseForm()
	if len(request.Form["deviceId"]) == 0 {
		log.Println("音响请求参数不全，需要参数deviceId")
		return
	}
	deviceId := request.Form["deviceId"][0]
	log.Println("请求参数：",request.Form)

	imsData,err := getImsData(deviceId)
	if err != nil {
		log.Println("查询imsData时报错：【",err,"]")
		return
	}
	log.Println("请求结果：",imsData)
	writer.Write([]byte(imsData))

}

func GetPrefix(writer http.ResponseWriter, request *http.Request)  {

	// 	获取请求参数
	request.ParseForm()
	if len(request.Form["imsNum"]) == 0 || len(request.Form["msisdn"]) == 0{
		log.Println("音响请求参数不全，需要参数imsNum和msisdn")
	}
	log.Println("请求参数：",request.Form)
	imsNum := request.Form["imsNum"][0]
	msisdn := request.Form["msisdn"][0]

	// 发送请求
	prefix,err := getPrefix(imsNum,msisdn)
	if err != nil {
		log.Println("查询prefix时报错：【",err,"]")
		return
	}
	log.Println("请求结果：",prefix)
	writer.Write([]byte(prefix))
}


func Dial(writer http.ResponseWriter, request *http.Request)  {

	// 进过讯飞解析知道为打电话
	// 获取要拨打的电话  name = "马宁" 查库的to
	to := "17611571680"

	//获取deviceId
	request.ParseForm()
	if len(request.Form["deviceId"]) == 0 {
		log.Println("参数不全，需要参数deviceId")
		writer.Write([]byte("1010"))
		return
	}
	deviceId := request.Form["deviceId"][0]

	// 获取imsNum imsData
	imsDta,err := getImsData(deviceId)
	if err != nil {
		log.Println("获取imsData失败:[,",err,"]")
		writer.Write([]byte("1010"))
		return
	}
	imsNum := gjson.Get(string(imsDta),"imsNum").String()

	// 获取prefix
	prefix,err  := getPrefix(imsNum,to)
	if err != nil {
		log.Println("获取prefix错误:[,",err,"]")
		writer.Write([]byte("1010"))
		return
	}
	to = prefix + to


	// 完成进行返回结构体的封装
	var result struct{
		Answer	string	`json:"answer"`	// 回答
		Data	interface{}	`json:"data"`// imsData
	}

	result.Answer = "好的,这就为您拨打"
	result.Data = map[string]interface{}{"to":to,"imsData":imsDta}
	rb,err := json.Marshal(result)
	if err != nil {
		log.Println("json化result时错误:[,",err,"]")
		writer.Write([]byte("1010"))
		return
	}
	_,err = writer.Write(rb)
	if err != nil {
		log.Println("返回数据时错误")
		writer.Write([]byte("1010"))
		return
	}

}



func getImsData(deviceId string) (string,error) {

	/*
	http://218.205.115.239:8080/fy/phone/v2/ims/freeAuth/public/authed/info? deviceId=xxxxxx&boxType=1
	*/


	// 获取ims
	getImsUrl,err := url.Parse(GetImsUrl)
	if err != nil{
		log.Println("解析url错误",err)
		return "",err
	}

	values := url.Values{}
	values.Set("deviceId",deviceId)
	values.Set("boxType","3")
	getImsUrl.RawQuery = values.Encode()
	resp,err := http.Get(getImsUrl.String())
	if err != nil {
		log.Println(err)
		return "",err
	}
	b,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "",err
	}
	return string(b),nil
}

func getPrefix(imsNum,msisdn string) (string,error) {

	/*
	GET http://218.205.115.239:8080/fy/phone/v2/ims/freeAuth/public/rule?imsNum=0531-111&msisdn=18867101111
	*/

	getImsUrl,err := url.Parse(GetPrefixUrl)
	if err != nil{
		log.Println("解析url错误",err)
		return "",err
	}

	values := url.Values{}
	values.Set("imsNum",imsNum)
	values.Set("msisdn",msisdn)
	getImsUrl.RawQuery = values.Encode()
	resp,err := http.Get(getImsUrl.String())
	if err != nil {
		log.Println(err)
		return "",err

	}
	b,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "",err
	}

	prefix := gjson.Get(string(b),"prefix").String()
	if prefix == "-1"{
		return "",nil
	}
	return prefix,nil
}
