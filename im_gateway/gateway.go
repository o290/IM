package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"server/common/etcd"
	"strings"
)

// BaseResponse 通用响应结构体
type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// Proxy 实现了 http.Handler 接口的 ServeHTTP 方法，用于处理 HTTP 请求
type Proxy struct {
}

func (Proxy) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//1.路径匹配请求前缀，并提取服务名字
	///api/(.*?)/ 以api为开头，后面跟着任意字符并以/为结尾
	//Compile将正则表达式字符串编译成一个 Regexp 对象
	regex, _ := regexp.Compile(`/api/(.*?)/`)
	//用于在给定的字符串中查找与正则表达式匹配的子字符串，并返回所有捕获组的匹配结果
	//切片的第一个元素是整个正则表达式匹配到的字符串，后续元素是各个捕获组匹配到的内容
	addrList := regex.FindStringSubmatch(req.URL.Path)
	if len(addrList) != 2 {
		FileResponse("err", res)
		return
	}
	//提取出服务名称 如请求路径是 127.0.0.1:8080/api/user/login，则提取出user
	service := addrList[1]

	//2.服务地址查找，从etcd中获取请求服务相关的地址
	//GetServiceAddr从 Etcd 中获取与服务名称对应的服务地址
	addr := etcd.GetServiceAddr(config.Etcd, service+"_api")
	if addr == "" {
		logx.Errorf("%s 不匹配的服务", service)
		FileResponse("err", res)
		return
	}
	//logx.Infof(addr)
	//3.请求认证服务地址并构建URL
	//分割带到ip地址和端口号切片
	//remote 客户端地址和端口号 没啥用
	remoteAddr := strings.Split(req.RemoteAddr, ":")
	//authAddr auth服务上送的地址与端口号
	authAddr := etcd.GetServiceAddr(config.Etcd, "auth_api")
	//authUrl 根据authAddr构造出认证api地址
	authUrl := fmt.Sprintf("http://%s/api/auth/authentication", authAddr)
	//logx.Info("11111", remoteAddr, authAddr, authUrl)

	//4.打印日志
	//请求转发，并根据原始请求和请求方法创建一个新的http请求
	//proxyUrl 认证地址 没啥用
	proxyUrl := fmt.Sprintf("http://%s%s", addr, req.URL.String())
	//logx.Info(proxyUrl)
	logx.Infof("%s %s", remoteAddr[0], proxyUrl)

	//5.请求认证
	if !auth(authUrl, res, req) {
		return
	}

	//6.创建并执行反向代理
	//将地址解析为 url.URL 类型，用于创建反向代理
	remote, _ := url.Parse(fmt.Sprintf("http://%s", addr))
	//创建一个单主机的反向代理，这个反向代理被配置为将接收到的请求转发到 remote 所指向的目标地址
	reverseProxy := httputil.NewSingleHostReverseProxy(remote)
	//将原始请求 req 转发到目标服务，并将目标服务的响应通过响应写入器 res 返回给客户端
	reverseProxy.ServeHTTP(res, req)
}

// FileResponse 生成并返回错误响应
func FileResponse(msg string, res http.ResponseWriter) {
	response := BaseResponse{Code: 7, Msg: msg}
	byteData, _ := json.Marshal(response)
	//写入http响应
	res.Write(byteData)
}

// auth 请求认证 将原始请求的 Header 复制到认证请求中，并添加 token（如果存在）和请求路径到认证请求 Header
func auth(authAddr string, res http.ResponseWriter, req *http.Request) (ok bool) {
	//1.创建认证请求
	authReq, _ := http.NewRequest("POST", authAddr, nil)
	authReq.Header = req.Header

	//2.添加Token和ValidPath到认证请求，因为认证功能是需要这两个参数的
	//判断是否能够在query中拿到token，如果拿到就添加到token中
	token := req.URL.Query().Get("token")
	if token != "" {
		authReq.Header.Set("Token", token)
	}
	//这个主要判断是不是在白名单中
	authReq.Header.Set("ValidPath", req.URL.Path)

	//3.发送认证请求并处理响应
	//http.DefaultClient.Do用于发送请求，即认证api
	authRes, err := http.DefaultClient.Do(authReq)
	if err != nil {
		logx.Error(err)
		FileResponse("认证服务错误", res)
		return
	}
	//4.解析认证响应
	type Response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data *struct {
			UserID uint `json:"userID"`
			Role   int  `json:"role"`
		} `json:"data"`
	}
	var authResponse Response
	byteData, _ := io.ReadAll(authRes.Body)
	authErr := json.Unmarshal(byteData, &authResponse)
	if authErr != nil {
		logx.Error(authErr)
		FileResponse("认证服务错误", res)
	}
	//5.解析认证结果
	if authResponse.Code != 0 {
		logx.Error(err)
		res.Write(byteData)
		return
	}
	//6.添加用户信息到原始请求头
	if authResponse.Data != nil {
		//在请求头添加用户id和用户角色信息
		req.Header.Set("User-ID", fmt.Sprintf("%d", authResponse.Data.UserID))
		req.Header.Set("Role", fmt.Sprintf("%d", authResponse.Data.Role))
	}
	//7.返回认证结果
	return true
}

// flag可以定义命令，解析命令参数
var configFile = flag.String("f", "settings.yaml", "the config file")

type Config struct {
	Addr string
	Etcd string
	Log  logx.LogConf
}

var config Config

func main() {
	//解析命令
	flag.Parse()
	//加载配置 将配置文件加载到结构体中
	//*configFile 是通过 flag 包解析得到的配置文件路径
	//&config 是指向 Config 结构体实例的指针
	conf.MustLoad(*configFile, &config)
	//初始化日志
	logx.SetUp(config.Log)
	fmt.Printf("gateway running %s\n", config.Addr)
	//Proxy 结构体实现了 http.Handler 接口的 ServeHTTP 方法，用于处理 HTTP 请求
	proxy := Proxy{}
	// 绑定服务 启动一个 HTTP 服务器
	http.ListenAndServe(config.Addr, proxy)
}
