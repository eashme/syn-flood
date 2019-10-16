package packet

import (
	"encoding/binary"
	"net"
	"syscall"
)

const (
	udpHeaderLen    = 20 // udp包长度
	udpMaxHeaderLen = 60
)

type udpHeader struct {
	SrcIP    net.IP
	DstIP    net.IP
	Protocol uint8
	SrcPort  uint16
	DstPort  uint16
	Check    uint16
	Body     []byte
}

func (u *udpHeader) Marshal() []byte {
	var b = make([]byte, udpHeaderLen)
	// 源IP
	copy(b[0:4], u.SrcIP.To4())
	// 目标IP
	copy(b[4:8], u.DstIP.To4())
	// 0 和 协议
	copy(b[8:10], []byte{0, u.Protocol})
	// 包长度
	binary.BigEndian.PutUint16(b[10:12], uint16(len(u.Body)+8))
	// 源端口
	binary.BigEndian.PutUint16(b[12:14], u.SrcPort)
	// 目标端口
	binary.BigEndian.PutUint16(b[14:16], u.DstPort)
	// 包长度
	binary.BigEndian.PutUint16(b[16:18], uint16(len(u.Body)+8))
	// 校验和
	binary.BigEndian.PutUint16(b[18:20], 0)
	// 重新计算 校验和
	binary.BigEndian.PutUint16(b[18:20], CheckSum(append(b, u.Body...)))
	// 返回组合好的udp包内容
	return append(b, u.Body...)
}

// 获取组合完成的udp包体
func GetUdpPacket(srcIp, dstIp net.IP, srcPort, dstPort uint16, body string) []byte {
	h := udpHeader{
		SrcIP:srcIp,
		DstIP:dstIp,
		Protocol:syscall.IPPROTO_UDP,
		SrcPort: srcPort,
		DstPort:dstPort,
		Body:[]byte(body),
	}
	return h.Marshal()
}
