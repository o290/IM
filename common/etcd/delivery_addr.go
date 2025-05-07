package etcd

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/netx"
	"server/core"
	"strings"
)

// DeliveryAddr 上送服务地址,将服务地址存储到etcd中
// etcdAddr etcd服务地址 serviceName 目标服务名称 addr目标服务地址
func DeliveryAddr(etcdAddr string, serviceName string, addr string) {
	//1.分割目标服务地址
	list := strings.Split(addr, ":")
	if len(list) != 2 {
		logx.Errorf("地址错误%s", addr)
		return
	}

	//2.如果目标服务地址是0.0.0.0，则获取本地IP地址，并将服务地址替换
	if list[0] == "0.0.0.0" {
		//服务注册需要提供真实的IP地址，InternalIp会过滤出服务条件的IP地址
		ip := netx.InternalIp()
		//替换 0.0.0.0:20022=>27.17.135.114:20022
		addr = strings.ReplaceAll(addr, "0.0.0.0", ip)
	}

	//3.初始化etcd客户端
	client := core.InitEtcd(etcdAddr)

	//4.以服务名称为键，将addr存储到etcd中
	_, err := client.Put(context.Background(), serviceName, addr)
	if err != nil {
		logx.Errorf("地址上送失败%s", err.Error())
		return
	}
	logx.Infof("地址上送成功%s %s", serviceName, addr)
}

// GetServiceAddr 从 Etcd 服务中获取指定服务名称对应的服务地址
// etcdAddr Etcd服务的地址
// serviceName 要获取地址的服务名称
func GetServiceAddr(etcdAddr string, serviceName string) (addr string) {
	//初始化etcd客户端
	client := core.InitEtcd(etcdAddr)
	//从 Etcd 获取服务地址
	res, err := client.Get(context.Background(), serviceName)
	//处理获取到的结果
	if err == nil && len(res.Kvs) > 0 {
		return string(res.Kvs[0].Value)
	}
	return ""
}
