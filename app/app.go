package app

import (
	"context"
	"log"
	"math/rand"
	"sender/clim"
	"sender/packet"
	"sender/utils"
	"sync"
	"syscall"
	"time"
)

func Run() {
	var (
		// 读取包体文件
		body                   = GetBody()
		fd                     int
		err                    error
		getSrcIP, getDstIP     getIP
		getSrcPort, getDstPort getPort
	)
	srcIPS := utils.GetIPs(Cfg.Src.File)
	dstIPS := utils.GetIPs(Cfg.Dst.File)
	// 构造获取IP 和端口的函数
	getSrcIP = constructGetIP(srcIPS, Cfg.Src.Start, Cfg.Src.End)
	getDstIP = constructGetIP(dstIPS, Cfg.Dst.Start, Cfg.Dst.End)

	getSrcPort = constructGetPort(Cfg.Src.Port.Start, Cfg.Src.Port.End)
	getDstPort = constructGetPort(Cfg.Dst.Port.Start, Cfg.Dst.Port.End)
	// 默认设置最大连接数为500,系统内置,不可修改
	// 创建fd数量 最大不超过1024
	wg := new(sync.WaitGroup)
	wg.Add(Cfg.ThreadCount)
	for i := 0; i < Cfg.ThreadCount; i ++ {
		if i%(Cfg.ThreadCount/500+1) == 0 { // 根据线程数决定每多少个线程共享一个连接
			if Cfg.Packet.Protocol {
				fd, err = clim.NewTCPfd()
			} else {
				fd, err = clim.NewUDPfd()
			}
			if err != nil {
				log.Println("create fd failed ")
				continue
			}
		}
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			if Cfg.Packet.Protocol {
				runSender(fd, getSrcIP, getDstIP, getSrcPort, getDstPort, body)
			} else {
				UdpSender(fd, getSrcIP, getDstIP, getSrcPort, getDstPort, body)
			}
		}(wg)
	}
	wg.Wait()
}

func runSender(fd int, getSrcIP getIP, getDstIP getIP, getSrcPort getPort, getDstPort getPort, body string) {
	// 构造目标地址对象
	var (
		addr         syscall.SockaddrInet4
		dstIP, srcIP uint32
		dstPort      uint16
		srcPort      uint16
		err          error
		count        int
	)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(Cfg.TimeOut)*time.Second)
	// 构造初步包体
	// 1. psd header
	psd := packet.GetPsdByte(getSrcIP(0), getDstIP(0))
	// 2. tcp header
	tcp, _ := packet.GetTcpByte(getSrcPort(0), getDstPort(0))
	// 3. ipv4 header
	ipv4, _ := packet.GetIPv4Header(getSrcIP(0), getDstIP(0))

	for {
		count ++
		select {
		case <-ctx.Done(): // 发送时间结束
			return
		default:
			dstIP = getDstIP(count)
			srcIP = getSrcIP(count * 2)
			dstPort = getDstPort(count)
			srcPort = getSrcPort(count * 2)
			addr.Port = int(dstPort)
			addr.Addr = utils.Uint32TO4Byte(dstIP)
			if err = syscall.Sendto(fd, append(append(packet.ChangeIPV4Header(ipv4, srcIP, dstIP, uint8(rand.Intn(15)) + 50), packet.ChangeTcpByte(psd, tcp, srcIP, dstIP, srcPort, dstPort)...), body...), 0, &addr); err != nil {
				// 在发包过程中 可能会出现 send Failed : operation not permitted 异常,
				// 因为发包过快linux的连接跟踪表满了(syn是建立连接的信号包,会被操作系统记录本次建立连接请求)
				// 导致linux认为你不能再继续进行发包操作了,所以报本次发包权限错误
				// 关闭linux防火墙 可解决该异常,当然调低发包速度,也不会有该异常产生
				// 但实际上,该异常导致的发包失败,影响小到可以忽略掉
				log.Println("send Failed :", err)
			}
		}
	}
}

func UdpSender(fd int, getSrcIP getIP, getDstIP getIP, getSrcPort getPort, getDstPort getPort, body string) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(Cfg.TimeOut)*time.Second)
	for {
		select {
		case <-ctx.Done(): // 发送时间结束
			return
		default:
			// todo UDP发包
		}
	}
}

/*
	entry
      |
	[main] -- router
		|		|
	server	server
	|			|
  router	router
	|	  |
deferent service
		|
	culture
*/

/*
10.10
//. 提高发包效率
1. ip全部使用uint32 进制数字进行处理,加快速度
2. 发包每次源ip都要不同
3. 缓存组包数据,每次发包,只变更源ip,端口和校验和

10.11
// todo 优化点
// 1. 节省内存分配时间，改变包体内容的时候不再去使用 tcp header 的副本进行置换,直接使用tcp header进行置换
// 2. 生成随机数,占用了较多CPU资源(随机数用公式加减计算出来,不用真正的随机数,降低cpu资源消耗)
// 3. 循环时是否去除select ?  select 会有一个判断时间,发送时间控制是否可以去掉？

需求
//  1. 读取文件中的ip进行发包
//  2. 发包方式, 轮询发包,每个ip逐个进行发包(队列),且兼用之前的范围随机选择发包方式,通过配置文件进行配置
*/
