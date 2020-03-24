package core

import (
	"dnsub/common"
	"fmt"
	"github.com/chenhg5/collection"
	"github.com/miekg/dns"
	"log"
	"strings"
	"time"
)

var(
	//泛解析cname黑名单和IP地址黑名单
	blackCnameList []string
	blackAddrList[]string
)



type DnsHandler struct {
	Client *dns.Client
	DnsService string
}


func NewDnsClient(dnsService string,timeout int) *DnsHandler{
	return &DnsHandler{
		Client:&dns.Client{
			Timeout:        time.Duration(timeout) * time.Second,
			DialTimeout:    time.Duration(timeout) * time.Second,
			ReadTimeout:    time.Duration(timeout) * time.Second,
			WriteTimeout:   time.Duration(timeout) * time.Second,

		},
		DnsService:dnsService,
	}

}

//检查DNS泛解析
func (c *DnsHandler) DnsAnalysis(domain string){
	for i:=0;i<3;i++{

		//随机虚假域名访问,获取泛解析IP和CNAME
		pre:= common.GetRandomString(8)
		dnsdomain:=pre+"."+domain
		//fmt.Println(dnsdomain)
		log.Printf("随机域名:%s",dnsdomain)
		cname,ipa,err:=c.DnsResolve(dnsdomain)

		if err!=nil{
			fmt.Println(err)
		}

		for _,v :=range cname{
			blackCnameList=append(blackCnameList,v)
			//fmt.Println("dnsall cname:"+v)
		}

		for _,v:=range ipa{
			blackAddrList=append(blackAddrList,v)
			//fmt.Println("dnsall ip:"+v)
		}
	}
	blackCnameList=collection.Collect(blackCnameList).Unique().ToStringArray()
	blackAddrList=collection.Collect(blackAddrList).Unique().ToStringArray()

	if len(blackAddrList)>0 || len(blackCnameList)>0{
		log.Printf("域名:%s,检测到泛解析,忽略下列解析.\n",domain)
		for _,v :=range blackCnameList{

			log.Printf("[*] ignore cname:"+v)
		}

		for _,v:=range blackAddrList{

			log.Printf("[*] ignore address:"+v)
		}
	}
}

func (c *DnsHandler) DnsResolve(src string) (cname,ipa []string, err error){
	var lastErr error
	timeout:=c.Client.Timeout+(time.Second*1)//deadline等于超时时间+1秒
	deadline :=time.Now().Add(timeout)

	for i := 1; time.Now().Before(deadline); i++ {
		m := dns.Msg{}
		// 最终都会指向一个ip 也就是typeA, 这样就可以返回所有层的cname.
		m.SetQuestion(dns.Fqdn(src), dns.TypeA)
		r, _, err := c.Client.Exchange(&m, c.DnsService+":53")
		if err != nil {
			lastErr = err
			log.Printf("DNS请求:%s,异常:%+v,retry:%d\n",dns.Fqdn(src),err,i)
			time.Sleep(900 * time.Millisecond)
			continue
		}



		cname = []string{}
		ipa=[]string{}
		for _, ans := range r.Answer {

			record, isType := ans.(*dns.CNAME)
			if isType {
				cname_str:=strings.TrimSuffix(record.Target,".")

				if !collection.Collect(blackCnameList).Contains(cname_str){

					cname = append(cname,cname_str )
				}

			}

			record1, isType1:= ans.(*dns.A)
			if isType1 {

				if !collection.Collect(blackAddrList).Contains(record1.A.String()){
					ipa = append(ipa, record1.A.String())
				}
			}
		}
		lastErr = nil
		break
	}
	err = lastErr
	return
}



