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
	"time"

	"github.com/juju/fslock"
	cmc "github.com/zytzjx/anthenacmc/cmcserverinfo"
	"github.com/zytzjx/anthenacmc/datacentre"
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
	Log.Log.Info("version:20.12.14.0; author:Jeffery zhang")
	logfile := flag.String("logfile", "", "upload to cmc server log zip file")
	jsonfile := flag.String("jsonfile", "", "transcation to cmc server data")
	isserver := flag.Bool("start-service", false, "upload failed list")
	flag.Parse()

	var configInstall cmc.ConfigInstall //map[string]interface{}
	if err := configInstall.LoadFile("serialconfig.json"); err != nil {
		Log.Log.Error(err)
		os.Exit(10)
	}
	datacentre.IsEmptySaveSerialConfig(configInstall)
	staticurl := configInstall.Results[0].Staticfileserver
	serviceserver := configInstall.Results[0].Webserviceserver

	if *isserver {
		lockSer := fslock.New(".service.lock")
		lockErrser := lockSer.TryLock()
		if lockErrser != nil {
			Log.Log.Info("Service has run")
			return
		}
		defer lockSer.Unlock()
		for {
			lock := fslock.New(".uploadlist.lock")
			lockErr := lock.TryLock()
			if lockErr != nil {
				time.Sleep(2 * time.Minute)
				continue
			}
			// release the lock
			lock.Unlock()

			reportcmc.SendLocalFiletoCMC(serviceserver, staticurl)
			time.Sleep(2 * time.Minute)
		}
	}

	lock := fslock.New(".uploadlist.lock")
	lockErr := lock.TryLock()
	// release the lock
	defer lock.Unlock()
	if lockErr != nil {
		time.Sleep(200 * time.Microsecond)
	}

	var uuid string
	var productid string
	if *jsonfile == "" {
		// postFile("https://httpbin.org/post", "reportbase.UUID", "reportbase.Productid", *logfile)
		reportbase, _, err := reportcmc.ReportCMC(*logfile)
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

		uuid = reportbase.UUID
		productid = reportbase.Productid
	} else {
		jsonFile, err := os.Open(*jsonfile)
		// if we os.Open returns an error then handle it
		if err != nil {
			Log.Log.Error(err)
			os.Exit(5)
		}
		// defer the closing of our jsonFile so that we can parse it later on
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		var items map[string]interface{}
		if err := json.Unmarshal(byteValue, &items); err != nil {
			Log.Log.Error(err)
			os.Exit(6)
		}
		reportcmc.Transcation(serviceserver, items)
		uuid = fmt.Sprintf("%v", items["uuid"])
		productid = fmt.Sprintf("%v", items["productid"])
	}

	if _, err := os.Stat(*logfile); os.IsNotExist(err) {
		Log.Log.Fatalf("update load zip,but logfile not exist %s", *logfile)
		os.Exit(3)
	}
	u, _ := url.Parse(staticurl)
	u.Path = path.Join(u.Path, "uploadlog/")
	postFile(u.String()+"/", uuid, productid, *logfile)
	os.Exit(0)
}
