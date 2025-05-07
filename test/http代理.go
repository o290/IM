package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
}

func (Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remote, _ := url.Parse("http://127.0.0.1:20023")
	//创建一个单主机的反向代理，将请求转发到 remote 所指向的目标地址
	reverseProxy := httputil.NewSingleHostReverseProxy(remote)
	reverseProxy.ServeHTTP(w, r)
}
func main() {
	addr := "127.0.0.1:8001"
	proxy := Proxy{}

	fmt.Printf("proxy server on%s\n", addr)
	http.ListenAndServe(addr, proxy)
}
