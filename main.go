
func main() {

	//监听信号
	signal.Notify(c, os.Interrupt, os.Kill)

	//收到退出信号结束程序
	go func() {
		select {
		case s:=<-c:
			fmt.Printf("\n%v,bye!\n",s)

			os.Exit(0)
		}
	}()

	flag.Parse()

	if ver{
		fmt.Println(VERSION)
	}

	if h{
		flag.Usage()
	}

	if !debug{
		log.SetOutput(ioutil.Discard)
	}

	domainItems:=strings.Split(d,",")
	//当输出文件存在时,提示用户删除文件,或者保留继续追加
	if o!=""{

		outputfile=common.CSVFileNameRepair(o)
	}

	OutputFileHandler(outputfile)

	
	fmt.Println("\nDone!")
	c<-os.Kill
	wait:=make(chan bool)
	<-wait
}