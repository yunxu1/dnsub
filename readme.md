# dnsub

![GitHub release (latest by date)](https://img.shields.io/github/v/release/yunxu1/dnsub)     ![GitHub All Releases](https://img.shields.io/github/downloads/yunxu1/dnsub/total?style=plastic)     ![GitHub stars](https://img.shields.io/github/stars/yunxu1/dnsub?style=plastic)     	



​		本工具通过字典枚举的方式用于扫描探测子域名,意在帮助用户梳理子域名资产使用,dnsub使用`go`语言高并发扫描，并可展示 **子域名**、**IP** 、 **CNAME**、**域名信息**，处理了枚举中常见的泛解析问题，支持加载多个字典，枚举探测更多深度的子域名信息，帮助用户快速掌握域名资产，扫描速度快效率高且跨平台，希望这款工具能帮助你有更多的收获 : )

 [点击下载](https://github.com/yunxu1/dnsub/releases "Releases")

##### 参数: 

```shell
Usage of ./dnsub:
  -d string
    	target domain (子域名目标)
  -depth int
    	enumerating subdomain depth using a param[f2] file content (default 2) （子域名爆破深度）
  -dns string
    	dns server address (default "9.9.9.9") （指定dns服务器）
  -f string
    	load subdomain filepath. eg: dnsubnames.txt (default "dict/dnsubnames.txt") (子域名字典)
  -f2 string
    	load subdomain filepath. eg: dnsub_next.txt (default "dict/dnsub_next.txt") (子域名字典,这个参数主要用于爆破2级及以后深度所使用的字典)
  -o string
    	output result to csv,set file path.(子域名扫描结果输出到csv文件,指定路径及文件名)
  -t int
    	thread pool numbers (default 20) (设置扫描的线程池数)
  -timeout int
    	dns question timeout,unit is second (default 5) (dns请求超时时间,默认5秒)
  -v int
    	show verbose level (default 1) (展示扫描信息等级,默认为1,如果大于1则会展示域名banner信息)
  -debug
    	enable debug output log info (打印debug信息)
  -h	help （帮助）
```



简单扫描:

```shell
dnsub -d example.com
```

##### 常规扫描:

```
dnsub -d example.com -t 200 -v 2 -o out.csv 

# -d 指定目标
# -t 指定线程数
# -o 指定输出文件
#	-v 2 扫描到子域名后展示banner信息
```

![example](https://raw.githubusercontent.com/yunxu1/dnsub/master/img/s.png)


##### 指定双字典扫描 && 设置扫描深度

本功能根据枚举深度加载不同的字典,推荐使用一些小字典多深度爆破下级域名

```shell
dnsub -d example.com -t 200 -f subnames.txt -f2 next.txt -depth 3 -o out.csv -v 2

#通过指定 -f 参数加载一个一级子域名枚举字典;
#通过指定 -f2 参数加载一个2级以上的子域名枚举字典
#通过指定 -depth 参数设置子域名枚举深度为3,即 a.b.c.example ,该参数默认值是2
```

如果你希望所有深度子域名枚举都使用一个字典可以指定 

```	shell
dnsub -f sub.txt -f2 sub.txt
```


##### 版本更新记录:
+ 优化泛解析识别;
+ 缩短打印行;
+ 优化域名访问探测,增强稳定性,忽略请求超时的域名;
- 下版本计划加入API接口查询功能.