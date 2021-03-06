package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Conf struct {
	App         App         `json:"app"`
	Stub        Stub        `json:"stub"`
	Uninstaller Uninstaller `json:"uninstaller"`
	Magisk      Magisk      `json:"magisk"`
}
type App struct {
	Version     string `json:"version"`
	VersionCode string `json:"versionCode"`
	Link        string `json:"link"`
	Note        string `json:"note"`
}
type Stub struct {
	VersionCode string `json:"versionCode"`
	Link        string `json:"link"`
}
type Uninstaller struct {
	Link string `json:"link"`
}
type Magisk struct {
	Version     string `json:"version"`
	VersionCode string `json:"versionCode"`
	Link        string `json:"link"`
	Note        string `json:"note"`
	Md5         string `json:"md5"`
}

var domain = flag.String("d", "", "域名或IP")
var port = flag.String("p", "80", "端口")
var listenPort = flag.String("listenPort", "80", "监听端口")
var listenAddress = flag.String("listenAddress", "0.0.0.0", "监听地址")
var isSsl = flag.Bool("s", false, "使用ssl")
var debug = flag.Bool("debug", false, "详细日志")

func main() {

	var host string
	var realListen string
	var s = make(map[bool]string, 2)

	flag.Parse()

	s[true] = "https"
	s[false] = "http"

	gin.SetMode(gin.ReleaseMode)
	if *debug == true {
		gin.SetMode(gin.DebugMode)
	}
	//fmt.Println(gin.Mode())

	if *domain == "" {
		*domain = getExternalIP()
	}
	if !strings.HasPrefix(*domain, "http://") && !strings.HasPrefix(*domain, "https://") {
		host = fmt.Sprintf("%s://%s:%s", s[*isSsl], *domain, *port)
		if *port == "80" {
			host = fmt.Sprintf("%s://%s", s[*isSsl], *domain)
		}
	} else {
		host = fmt.Sprintf("%s:%s", *domain, *port)
		if *port == "80" {
			host = *domain
		}
	}

	fmt.Println("已指定域名:", *domain)
	fmt.Println("已指定端口:", *port)
	fmt.Println("生成自定义更新通道访问路径:", host+"/beta.json")

	realListen = *listenAddress + ":" + *listenPort

	//1. 获取 beta.json 配置
	go cron(getConfig, host)
	//4. 启动服务器
	r := gin.Default()
	//下载功能
	r.StaticFile("/magisk.apk", "./magisk.apk")
	r.StaticFile("/magisk.zip", "./magisk.zip")
	r.StaticFile("/beta.json", "./beta.json")
	r.LoadHTMLFiles("./index.tmpl")
	r.GET("/", func(context *gin.Context) {
		context.HTML(200, "index.tmpl", gin.H{
			"host": host + "/beta.json",
		})
	})
	if err := r.Run(realListen); err != nil {
		fmt.Println(err)
	}
	fmt.Println("已监听:", realListen)
	//r.GET("/magisk.apk", func(context *gin.Context) {
	//	// 获取要返回的文件数据流
	//	file, err := os.OpenFile("./magisk.apk", os.O_RDONLY, 400)
	//	if err != nil {
	//		fmt.Println("打开./magisk.apk失败")
	//	}
	//	content, err := ioutil.ReadAll(file)
	//	context.Writer.WriteHeader(200)
	//	context.Header("Content-Disposition", "attachment; filename=magisk.apk")
	//	context.Header("Content-Type", "application/text/plain")
	//	//c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
	//	context.Header("Accept-Length", "200")
	//	context.Data(200, "application/vnd.android.package-archive", content)
	//})
	//r.GET("/magisk.zip", func(context *gin.Context) {
	//	file, err := os.OpenFile("./magisk.zip", os.O_RDONLY, 400)
	//	if err != nil {
	//		fmt.Println("打开./magisk.zip失败")
	//	}
	//	content, err := ioutil.ReadAll(file)
	//	context.Writer.WriteHeader(200)
	//	context.Header("Content-Disposition", "attachment; filename=magisk.zip")
	//	context.Header("Content-Type", "application/text/plain")
	//	context.Header("Accept-Length", "200")
	//	context.Data(200, "application/application/zip", content)
	//})
	//r.GET("/beta.json", func(context *gin.Context) {
	//	file, err := os.OpenFile("./beta.json", os.O_RDONLY, 400)
	//	if err != nil {
	//		fmt.Println("打开./beta.json失败")
	//	}
	//	content, err := ioutil.ReadAll(file)
	//	var data Conf
	//	err = json.Unmarshal(content, &data)
	//	if err != nil {
	//		fmt.Println("err:", err)
	//		return
	//	}
	//	fmt.Println(data)
	//	context.JSON(200, data)
	//})
}

func getAndSaveMagisk(ctx context.Context, link string, savePath string) {
	select {
	case <-ctx.Done():
		fmt.Println("download ", savePath, " timeout")
		return
	default:
		resp, err := http.Get(link)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		out, err := os.Create(savePath)
		if err != nil {
			fmt.Printf("创建 %s 失败,err:%s\n", savePath, err)
			return
		}
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			fmt.Printf("写入 %s 失败,err:%s\n", savePath, err)
			return
		}
		if err = out.Close(); err != nil {
			log.Printf("关闭%s失败,err:%v", savePath, err)
		}
		fmt.Println("下载", savePath, "成功")
	}
}
func getConfig(ctx context.Context, host string) {
	select {
	case <-ctx.Done():
		fmt.Println("get config timeout")
		return
	default:
		resp, err := http.Get("https://raw.githubusercontent.com/topjohnwu/magisk_files/master/beta.json")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		data, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(data))

		var cfg = Conf{}
		var copyCfg = Conf{}
		err = json.Unmarshal(data, &cfg)
		err = json.Unmarshal(data, &copyCfg)

		//2. 复制配置的副本，修改副本，替换 link 为自己的 link. 并保存到本地
		copyCfg.App.Link = fmt.Sprintf("%s/%s", host, "magisk.apk")
		copyCfg.Magisk.Link = fmt.Sprintf("%s/%s", host, "magisk.zip")

		fmt.Println("获取远程配置成功")

		out, err := os.Create("./beta.json")
		if err != nil {
			fmt.Printf("创建本地配置文件失败,err:%s\n", err)
			return
		}
		copyString, err := json.Marshal(copyCfg)
		_, err = io.Copy(out, strings.NewReader(string(copyString)))
		if err != nil {
			fmt.Printf("写入本地配置文件失败,err:%s\n", err)
			return
		}
		if err = out.Close(); err != nil {
			log.Println("关闭beta.json失败", err)
		}
		fmt.Println("更新本地配置成功")
		//3. 通过未修改的配置(cfg)中的 link 下载 magisk.apk 和 magisk.zip
		fmt.Println(cfg.Magisk.Link)
		go getAndSaveMagisk(ctx, cfg.App.Link, "./magisk.apk")
		go getAndSaveMagisk(ctx, cfg.Magisk.Link, "./magisk.zip")
	}
}

func getExternalIP() string {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get("http://myexternalip.com/raw")
	if err != nil && err != io.EOF {
		fmt.Println("无法自动获取公网IP地址，其使用[-d]手动指定域名或者IP")
		panic(err)
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		fmt.Println("无法自动获取公网IP地址，其使用[-d]手动指定域名或者IP")
		panic(err)
	}
	return string(ip)
}

func cron(fn func(ctx context.Context, host string), host string) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		go fn(ctx, host)
		time.Sleep(24 * time.Hour)
		cancel()
	}
}
