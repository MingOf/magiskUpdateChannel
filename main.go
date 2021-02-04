package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type Conf struct {
	App App `json:"app"`
	Stub Stub `json:"stub"`
	Uninstaller Uninstaller `json:"uninstaller"`
	Magisk Magisk `json:"magisk"`
}
type App struct {
	Version     string `json:"version"`
	VersionCode string `json:"versionCode"`
	Link        string `json:"link"`
	Note        string `json:"note"`
}
type Stub struct {
	VersionCode string `json:"versionCode"`
	Link string `json:"link"`
}
type Uninstaller struct {
	Link string `json:"link"`
}
type Magisk struct{
	Version string `json:"version"`
	VersionCode string `json:"versionCode"`
	Link string `json:"link"`
	Note string `json:"note"`
	Md5 string `json:"md5"`
}

func main() {
	domain := flag.String("domain",getExternalIP(),"域名或IP")
	debug:=flag.Bool("debug",false,"详细日志")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	if *debug==true {
		gin.SetMode(gin.DebugMode)
	}
	fmt.Println(gin.Mode())

	//1. 获取 beta.json 配置
	fmt.Println("已指定域名:",*domain)

	go cron(getConfig,*domain)

	//4. 启动服务器
	r:=gin.Default()
	//下载功能
	r.GET("/magisk.apk", func(context *gin.Context) {
		// 获取要返回的文件数据流
		file,err:=os.OpenFile("./magisk.apk",os.O_RDONLY,400)
		if err != nil {
			fmt.Println("打开./magisk.apk失败")
		}
		content,err:=ioutil.ReadAll(file)
		context.Writer.WriteHeader(200)
		context.Header("Content-Disposition", "attachment; filename=magisk.apk")
		context.Header("Content-Type","application/text/plain")
		//c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
		context.Header("Accept-Length", "200")
		context.Data(200,"application/vnd.android.package-archive",content)
	})
	r.GET("/magisk.zip", func(context *gin.Context) {
		file,err:=os.OpenFile("./magisk.zip",os.O_RDONLY,400)
		if err != nil {
			fmt.Println("打开./magisk.zip失败")
		}
		content,err:=ioutil.ReadAll(file)
		context.Writer.WriteHeader(200)
		context.Header("Content-Disposition", "attachment; filename=magisk.apk")
		context.Header("Content-Type","application/text/plain")
		context.Header("Accept-Length", "200")
		context.Data(200,"application/application/zip",content)
	})
	r.GET("/beta.json", func(context *gin.Context) {
		file,err:=os.OpenFile("./beta.json",os.O_RDONLY,400)
		if err != nil {
			fmt.Println("打开./magisk.zip失败")
		}
		content,err:=ioutil.ReadAll(file)
		var data Conf
		err = json.Unmarshal(content, &data)
		if err!=nil {
			fmt.Println("err:",err)
			return
		}
		fmt.Println(data)
		context.JSON(200, data)
	})
	r.LoadHTMLFiles("./index.tmpl")
	r.GET("/", func(context *gin.Context) {
		context.HTML(200,"index.tmpl",gin.H{
			"host":"http://"+*domain+"/beta.json",
		})
	})
	r.Run()
}

func getAndSaveMagisk(ctx context.Context,link string, path string) {
	select {
	case <-ctx.Done():
		fmt.Println("download magisk timeout")
		return
	default:
		resp,err:=http.Get(link)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		out,err:=os.Create(path)
		if err!= nil {
			fmt.Printf("创建 %s 失败,err:%s\n",path,err)
			return
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			fmt.Printf("写入 %s 失败,err:%s\n",path,err)
			return
		}
		fmt.Println("下载",path,"成功")
	}
}
func getConfig(ctx context.Context,domain string) {
	select {
	case <-ctx.Done():
		fmt.Println("get config timeout")
		return
	default:
		resp,err:=http.Get("https://raw.githubusercontent.com/topjohnwu/magisk_files/master/beta.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		data,_:=ioutil.ReadAll(resp.Body)
		fmt.Println(string(data))

		var cfg = Conf{}
		var copyCfg = Conf{}
		err = json.Unmarshal(data, &cfg)
		err = json.Unmarshal(data, &copyCfg)

		//2. 复制配置的副本，修改副本，替换 link 为自己的 link. 并保存到本地
		copyCfg.App.Link = "http://"+domain+"/magisk.apk"
		copyCfg.Magisk.Link = "http://"+domain+"/magisk.zip"

		fmt.Println("获取远程配置成功")

		out,err:=os.Create("./beta.json")
		if err!= nil {
			fmt.Printf("创建本地配置文件失败,err:%s\n",err)
			return
		}
		copyString,err := json.Marshal(copyCfg)
		_, err = io.Copy(out, strings.NewReader(string(copyString)))
		if err!= nil {
			fmt.Printf("写入本地配置文件失败,err:%s\n",err)
			return
		}
		fmt.Println("更新本地配置成功")
		//3. 通过未修改的配置(cfg)中的 link 下载 magisk.apk 和 magisk.zip
		go getAndSaveMagisk(ctx ,cfg.App.Link,"./magisk.apk")
		go getAndSaveMagisk(ctx ,cfg.Magisk.Link,"./magisk.zip")
	}
}

func getExternalIP() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil && err!= io.EOF{
		fmt.Println("无法获取公网IP地址，err:",err)
		return ""
	}
	defer resp.Body.Close()
	ip,err:=ioutil.ReadAll(resp.Body)
	if err != nil &&err !=io.EOF {
	    fmt.Println("无法获取公网IP地址，err:",err)
	    return ""
	}
	return string(ip)
}

func cron(fn func(ctx context.Context,domain string),domain string) {
	for {
		ctx,cancel:=context.WithTimeout(context.Background(),30*time.Minute)
		go fn(ctx,domain)
		time.Sleep(24*time.Hour)
		cancel()
	}
}