package test

import (
	"testing"
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"log"
)
const  (
	//BaseUrl = "http://localhost:9090"
	BaseUrl = "http://47.92.93.152:9090"
	DailUrl = BaseUrl + "/dial"
	SynUrl = BaseUrl + "/syn"
	GetImsUrl = BaseUrl + "/ims"
	GetPrefixUrl = BaseUrl + "/prefix"

	//GetImsUrl = BaseUrl + "/fy/phone/v2/ims/freeAuth/public/authed/info"
	//GetPrefixUrl = BaseUrl + "http://218.205.115.239:8080/fy/phone/v2/ims/freeAuth/public/rule"
)

func init()  {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

type Record struct {
	Msisdn 		string		`json:"msisdn"`// ims号码
	Persons 	[]*Person	`json:"persons"`// 联系人
}
type Person struct {
	NickName 	string	`json:"nickname"`// 联系人昵称
	Number		string	`json:"number"`// 联系人电话
	Status		string	`json:"status"`// 联系人状态
}


func TestGetsyn(t *testing.T)  {

	p1 := Person{"小马","17611571680","11"}
	p2 := Person{"小子超","15613451678","2"}
	ps := []*Person{&p1,&p2}

	rb,err := json.Marshal(ps)
	if err != nil {
		fmt.Print("解析json错误",err)
		return
	}

	values := url.Values{}
	values.Set("msisdn","1111")
	fmt.Println(string(rb))
	values.Set("persons",string(rb))
	resp,err := http.PostForm(SynUrl,values)
	if err != nil {
		log.Println(err)
		return
	}
	b,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(b))

}

func TestGetIms(t *testing.T) {
	getImsUrl,err := url.Parse(GetImsUrl)
	if err != nil{
		log.Println("解析url错误",err)
		return
	}

	values := url.Values{}
	values.Set("deviceId","rwt123456")
	values.Set("boxType","3")
	getImsUrl.RawQuery = values.Encode()
	resp,err := http.Get(getImsUrl.String())
	if err != nil {
		log.Println(err)
		return
	}
	b,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(b))
}


func TestGetPrefix(t *testing.T) {
	getImsUrl,err := url.Parse(GetPrefixUrl)
	if err != nil{
		log.Println("解析url错误",err)
		return
	}

	values := url.Values{}
	values.Set("imsNum","0531-58021024")
	values.Set("msisdn","17611571680")
	getImsUrl.RawQuery = values.Encode()
	resp,err := http.Get(getImsUrl.String())
	if err != nil {
		log.Println(err)
		return
	}
	b,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(b))
}


func TestDial(t *testing.T) {
	getImsUrl,err := url.Parse(DailUrl)
	if err != nil{
		log.Println("解析url错误",err)
		return
	}

	values := url.Values{}
	values.Set("deviceId","rwt123456")
	getImsUrl.RawQuery = values.Encode()
	resp,err := http.Get(getImsUrl.String())
	if err != nil {
		log.Println(err)
		return
	}
	b,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(b))
}