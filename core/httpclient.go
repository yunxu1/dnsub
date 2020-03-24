package core

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)


type HttpBanner struct {
	Schema string
	Title string
	Banner string
	XPowerBy string
}

func (banner *HttpBanner) String() string{

	fmtString:=""

	if banner.Schema!=""{
		fmtString+="|"+banner.Schema
	}
	if banner.Title!=""{
		fmtString+="|"+banner.Title
	}

	if banner.Banner!=""{
		fmtString+="|"+banner.Banner
	}

	if banner.XPowerBy!=""{
		fmtString+="|"+banner.XPowerBy
	}

	return strings.TrimLeft(fmtString,"|")
}


type HttpClient struct {
	client *http.Client
}

func NewHttpClient(timeout int) *HttpClient{
	transport:= &http.Transport{Dial: (&net.Dialer{
	
		//解决连接过多出现err:too many open file.
		// https://colobu.com/2016/07/01/the-complete-guide-to-golang-net-http-timeouts/
		// http://craigwickesser.com/2015/01/golang-http-to-many-open-files/
		Timeout:   time.Duration(timeout) * time.Second,
		Deadline:  time.Now().Add(time.Duration(timeout) * time.Second),
		KeepAlive: time.Duration(timeout) * time.Second,
	}).Dial,
		TLSHandshakeTimeout:time.Duration(timeout)* time.Second,
		//忽略证书校验
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},}

	return &HttpClient{client:&http.Client{
		Transport:transport,
		Timeout:time.Duration(timeout)*time.Second,
	},}
}



func (client *HttpClient) newRequest(method,target string) *http.Request{

	req,err:=http.NewRequest(http.MethodGet,target,nil)

	if err!=nil{
		log.Printf("error newrequest:%s\n",err)
		return nil
	}

	req.Header.Set("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.70 Safari/537.36")
	req.Header.Set("Connection","close")

	return req
}

func (httpclient *HttpClient) Verify(domain string)  (httpbanner HttpBanner,error error){
	//fmt.Println(">>>>>>:"+domain)
	var schemas = []string{"https","http"}
	for _,schema :=range schemas{


		target:=schema+"://"+domain

		httpbanner,error=httpclient.reqGet(schema,target)

		//如果一次执行成功,就返回首次访问的数据
		//比如探测https协议的网页存在,就不再继续探测http的
		if error==nil{
			break
		}
	}

	return
}


func (httpclient *HttpClient) reqGet(schema,target string)  (httpbanner HttpBanner,error error){

	const timeout  = 2*time.Second
	deadline := time.Now().Add(timeout)

	for retry:=1;time.Now().Before(deadline);retry++{
		error=nil //初始化error

		req:=httpclient.newRequest(http.MethodGet,target)

		resp,err:=httpclient.client.Do(req)

		if err!=nil{
			log.Printf("Server not respond (%v);the %d times retry....", err, retry)
			error=err
			time.Sleep(2*time.Second)
			continue
		}

		defer resp.Body.Close()

		banner:=resp.Header.Get("Server")
		//contentlength := resp.Header.Get("Content-Length")
		x_power_by:=resp.Header.Get("X-Powered-By")
		//contenttype:=""
		//pair := strings.SplitN(resp.Header.Get("Content-Type"), ";", 2)
		//if len(pair) == 2 {
		//	contenttype= pair[0]
		//}

		body,err:=ioutil.ReadAll(resp.Body)
		if err!=nil{
			//出现响应读取异常直接break跳出结束
			//log.Printf("Server respond reader (%v);the %d times retry....", err, retry)
			error=err
			break

		}


		re:=regexp.MustCompile("<title>([\\s\\S]*?)</title>")

		title:=re.FindSubmatch(body)


		httpbanner.Schema=schema
		httpbanner.Banner=banner
		httpbanner.XPowerBy=x_power_by
		if len(title)>0{
			httpbanner.Title=strings.TrimSpace(string(title[len(title)-1]))
		}
		break
	}
	return httpbanner,error

}