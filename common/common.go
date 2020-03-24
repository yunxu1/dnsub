package common

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

var lock=sync.Mutex{}

func  GetRandomString(l int) string {
	str := "abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//输出内容in到csv文件
func OutPutCsv(filename string,in []string){

	lock.Lock()
	ofs,err:=os.OpenFile(filename,os.O_RDWR|os.O_CREATE|os.O_SYNC,0666)

	if err!=nil{
		fmt.Printf("output result fail:%v\n",err)
		return
	}

	ofs.WriteString("\xEF\xBB\xBF") //UTF-8
	defer ofs.Close()

	ofs.Seek(0,io.SeekEnd)

	w:=csv.NewWriter(ofs)

	w.Comma=','
	w.UseCRLF=true

	err=w.Write(in)

	if err!=nil{
		fmt.Printf("output result fail:%v\n",err)
	}
	w.Flush()
	lock.Unlock()
}

//判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}


func CSVFileNameRepair(pathfile string) (string){

	dir:=path.Dir(pathfile)
	filename:=path.Base(pathfile)
	fileSuffix:=path.Ext(pathfile)
	filenameOnly := strings.TrimSuffix(filename, fileSuffix)
	csvfilepath:=[]string{dir,filenameOnly+".csv"}

	return path.Join(csvfilepath...)


}