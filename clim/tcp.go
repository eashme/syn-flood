package clim

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sender/packet"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

const (
	FIN = 1  // 00 0001
	SYN = 2  // 00 0010
	RST = 4  // 00 0100
	PSH = 8  // 00 1000
	ACK = 16 // 01 0000
	URG = 32 // 10 0000
)

type TCPHeader struct {
	SrcPort   uint16
	DstPort   uint16
	SeqNum    uint32
	AckNum    uint32
	Offset    uint8
	Flag      uint8
	Window    uint16
	Checksum  uint16
	UrgentPtr uint16
}

type PsdHeader struct {
	SrcAddr   uint32
	DstAddr   uint32
	Zero      uint8
	ProtoType uint8
	TcpLength uint16
}

func inetAddr(host string) uint32 {
	var (
		segments []string = strings.Split(host, ".")
		ip       [4]uint64
		ret      uint64
	)
	for i := 0; i < 4; i++ {
		ip[i], _ = strconv.ParseUint(segments[i], 10, 64)
	}
	ret = ip[3]<<24 + ip[2]<<16 + ip[1]<<8 + ip[0]
	return uint32(ret)
}

func htons(port uint16) uint16 {
	//var (
	//	high uint16 = port >> 8
	//	ret  uint16 = port<<8 + high
	//)
	//return ret
	return port<<8 + port>>8
}

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += sum >> 16
	return uint16(^sum)
}

func NewPsdHeader(srcHost string, dstHost string) *PsdHeader {
	// 创建TCP外部包头内容
	return &PsdHeader{
		SrcAddr:   inetAddr(srcHost),
		DstAddr:   inetAddr(dstHost),
		Zero:      0,
		ProtoType: syscall.IPPROTO_TCP,
		TcpLength: 0,
	}
}

func NewTcpHeader(srcPort uint16, dstPort uint16) *TCPHeader {
	// 创建TCP内部包头内容
	return &TCPHeader{
		SrcPort:  srcPort,
		DstPort:  dstPort,
		SeqNum:   0,
		AckNum:   0,
		Offset:   uint8(uint16(unsafe.Sizeof(TCPHeader{}))/4) << 4,
		Flag:     SYN,
		Window:   60000,
		Checksum: 0,
	}
}

func NewTcpData(psdHeader *PsdHeader, tcpHeader *TCPHeader, data string) []byte {
	psdHeader.TcpLength = uint16(unsafe.Sizeof(TCPHeader{})) + uint16(len(data))
	var buffer bytes.Buffer
	/*buffer用来写入两种首部来求得校验和*/
	binary.Write(&buffer, binary.BigEndian, psdHeader)
	binary.Write(&buffer, binary.BigEndian, tcpHeader)
	tcpHeader.Checksum = CheckSum(buffer.Bytes())
	/*接下来清空buffer，填充实际要发送的部分*/
	buffer.Reset()
	binary.Write(&buffer, binary.BigEndian, tcpHeader)
	binary.Write(&buffer, binary.BigEndian, data)
	return buffer.Bytes()
}

func NewTCPfd() (int, error) {
	sockfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		return 0, err
	}
	// 设置tcp伪造包体级别
	err = syscall.SetsockoptInt(sockfd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
	if err != nil {
		return 0, err
	}
	return sockfd, nil
}

func SendPacket(fd int, srcIP, dstIP net.IP, dstPort int, body string) (err error) {
	// 设定目标地址
	addr := syscall.SockaddrInet4{Port: dstPort}
	copy(addr.Addr[:4], dstIP)
	err = syscall.Sendto(fd, NewTcpData(NewPsdHeader(fmt.Sprintf("%s", srcIP), fmt.Sprintf("%s", dstIP)), NewTcpHeader(uint16(dstPort), uint16(dstPort)), body), 0, &addr)
	if err != nil {
		return fmt.Errorf("Sendto error : %s ", err)
	}
	return nil
}

func Closefd(fd int) {
	syscall.Shutdown(fd, syscall.SHUT_RDWR)
}

func SendTCP(fd int, srcIP, dstIP uint32, dstPort int, body string, addr *syscall.SockaddrInet4) (err error) {
	var (
		tcpByte, ipv4Byte []byte
	)
	if tcpByte, err = packet.GetTcpHeader(srcIP, dstIP, uint16(dstPort), uint16(dstPort)); err != nil {
		return err
	}
	if ipv4Byte, err = packet.GetIPv4Header(srcIP, dstIP); err != nil {
		return err
	}
	buffs := make([]byte, 0)
	buffs = append(buffs, ipv4Byte...)
	buffs = append(buffs, tcpByte...)
	// 写入body
	buffs = append(buffs, body...)
	err = syscall.Sendto(fd, buffs, 0, addr)
	if err != nil {
		return fmt.Errorf("Sendto error : %s ", err)
	}
	return nil
}
