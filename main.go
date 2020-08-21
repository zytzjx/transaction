package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"

	Log "github.com/zytzjx/anthenacmc/loggersys"
	"github.com/zytzjx/anthenacmc/reportcmc"
)

// func postFile1(url, uuid, productid string, filePath string) error {

// }

func postFile(url, uuid, productid string, filePath string) error {

	//打开文件句柄操作
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer file.Close()

	//创建一个模拟的form中的一个选项,这个form项现在是空的
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作, 设置文件的上传参数叫uploadfile, 文件名是filename,
	//相当于现在还没选择文件, form项里选择文件的选项
	fileWriter, err := bodyWriter.CreateFormFile("fileobj", uuid+".zip")
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	//iocopy 这里相当于选择了文件,将文件放到form中
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return err
	}

	//获取上传文件的类型,multipart/form-data; boundary=...
	contentType := bodyWriter.FormDataContentType()

	//这里就是上传的其他参数设置,可以使用 bodyWriter.WriteField(key, val) 方法
	//也可以自己在重新使用  multipart.NewWriter 重新建立一项,这个再server 会有例子
	params := map[string]string{
		"uuid":      uuid,
		"productid": productid,
	}
	//这种设置值得仿佛 和下面再从新创建一个的一样
	for key, val := range params {
		_ = bodyWriter.WriteField(key, val)
	}
	//这个很关键,必须这样写关闭,不能使用defer关闭,不然会导致错误
	bodyWriter.Close()

	//发送post请求到服务端
	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(respbody))
	// ioutil.WriteFile("fileresult.txt", respbody, 0644)

	if resp.StatusCode == http.StatusOK {
		rsjs := make(map[string]interface{})
		if json.Unmarshal(respbody, &rsjs) != nil {
			if errstr, ok := rsjs["error"]; ok {
				return fmt.Errorf("http error, %s", errstr)
			}
		}
		return nil
	}
	return fmt.Errorf("http error, %s", resp.Status)

}

func main() {
	Log.NewLogger("reportcmc")
	Log.Log.Info("version:20.08.21.0; author:Jeffery zhang")
	logfile := flag.String("logfile", "", "upload to cmc server log zip file")
	flag.Parse()
	// postFile("https://httpbin.org/post", "reportbase.UUID", "reportbase.Productid", *logfile)
	reportbase, surl, err := reportcmc.ReportCMC()
	if err != nil {
		os.Exit(1)
	}

	if *logfile == "" {
		Log.Log.Info("need not update load zip")
		os.Exit(0)
	}

	if reportbase == nil {
		Log.Log.Fatal("update load zip,but missing parameter")
		os.Exit(2)
	}
	if _, err := os.Stat(*logfile); os.IsNotExist(err) {
		Log.Log.Fatalf("update load zip,but logfile not exist %s", *logfile)
		os.Exit(3)
	}
	u, err := url.Parse(surl)
	u.Path = path.Join(u.Path, "uploadlog/")
	postFile(u.String()+"/", reportbase.UUID, reportbase.Productid, *logfile)
	os.Exit(0)
}
