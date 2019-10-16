package app

import (
	"fmt"
	"io/ioutil"
	"net"
	"sender/utils"
)

func GetBody() string {
	var (
		b    []byte
		err  error
		body = Cfg.Packet.Body
	)
	if utils.Exists(Cfg.Packet.BodyFile) {
		if b, err = ioutil.ReadFile(Cfg.Packet.BodyFile); err == nil {
			body = string(b)
		}
	}
	return body
}

// 根据输入的Ip计算IP地址波动范围
func GetIPRange(start string, end string) int64 {
	sip := net.ParseIP(start)
	eip := net.ParseIP(end)
	if sip == nil || eip == nil {
		panic(fmt.Errorf("IP 地址解析失败"))
	}
	s := utils.IP2Uint32(sip)
	e := utils.IP2Uint32(eip)
	if e > s {
		return int64(e - s)
	}
	return int64(s - e)
}

type getIP func(count int) uint32

type getPort func(count int) uint16

// 构造获取IP的函数, 返回值会超高频率执行,所以构造好的函数需要非常精简
func constructGetIP(ips []uint32,startIP , endIP string) getIP {
	il := len(ips)
	if il > 0 { // 获取到了,则从该列表中取即可
		// 对传入的ip列表进行从小到大排序
		for i := 0;i < il; i ++ {
			for j := il - 1;j > i; j -- {
				if ips[i] > ips[j]{
					ips[i] ,ips[j] = ips[j] ,ips[i]
				}
			}
		}
		return func(count int) uint32 {
			return ips[count % il]
		}
	}
	IPStart := utils.IP2Uint32(net.ParseIP(startIP))
	ipRange := int(GetIPRange(startIP, endIP))
	return func(count int) uint32 {
		return uint32(count % ipRange) + IPStart
	}
}

// 构造获取端口的函数, 返回值会超高频率执行,所以构造好的函数需要非常精简
func constructGetPort(startPort,endPort int) getPort{
	pRange := int(endPort - startPort)
	sp := uint16(startPort)
	return func(count int) uint16 {
		return uint16(count % pRange) + sp
	}
}