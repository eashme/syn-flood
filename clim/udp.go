package clim

import (
	"net"
	"sender/packet"
	"syscall"
)

// 创建UDP连接句柄
func NewUDPfd() (int, error) {
	// 创建udp发包连接句柄fd
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_UDP)
	if err != nil {
		return 0, err
	}
	// 设置发包修改包体权限
	if err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1); err != nil {
		return 0, err
	}
	return fd, nil
}

// 发送udp 包体
func SendUDP(fd int, srcIP, dstIP net.IP, srcPort, dstPort uint16, body string) (err error) {
	// 构造包体
	b := packet.GetUdpPacket(srcIP, dstIP, srcPort, dstPort, body)

	// 设定目标地址
	addr := syscall.SockaddrInet4{Port: int(dstPort)}
	copy(addr.Addr[:4], dstIP)
	return syscall.Sendto(fd, b, 0, &addr)
}

