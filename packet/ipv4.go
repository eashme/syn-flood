package packet

import (
	"encoding/binary"
	"syscall"
)

const (
	ipv4Version      = 4
	ipv4HeaderLen    = 20
	ipv4MaxHeaderLen = 60
)

// A ipv4 header
type ipv4Header struct {
	Version  int    // 协议版本 4bit
	Len      int    // 头部长度 4bit
	TOS      int    // 服务类   8bit
	TotalLen int    // 包长		16bit
	ID       int    // id		8bit
	Flags    int    // flags	3bit
	FragOff  int    // 分段偏移量 13bit
	TTL      int    // 生命周期 4bit
	Protocol int    // 上层服务协议4bit
	Checksum uint16 // 头部校验和16bit
	Src      uint32 // 源IP  	32bit
	Dst      uint32 // 目的IP  	32bit
	Options  []byte // 选项, extension headers
}

// Marshal encode ipv4 header
func (h *ipv4Header) Marshal() ([]byte, error) {
	if h == nil {
		return nil, syscall.EINVAL
	}

	hdrlen := ipv4HeaderLen + len(h.Options)
	b := make([]byte, hdrlen)

	//版本和头部长度
	b[0] = byte(ipv4Version<<4 | (hdrlen >> 2 & 0x0f))
	b[1] = byte(h.TOS)

	binary.BigEndian.PutUint16(b[2:4], uint16(h.TotalLen))
	binary.BigEndian.PutUint16(b[4:6], uint16(h.ID))

	flagsAndFragOff := (h.FragOff & 0x1fff) | int(h.Flags<<13)
	binary.BigEndian.PutUint16(b[6:8], uint16(flagsAndFragOff))

	b[8] = byte(h.TTL)
	b[9] = byte(h.Protocol)

	binary.BigEndian.PutUint16(b[10:12], uint16(h.Checksum))
	binary.BigEndian.PutUint32(b[12:16], h.Src)
	binary.BigEndian.PutUint32(b[16:20], h.Dst)
	if len(h.Options) > 0 {
		copy(b[ipv4HeaderLen:], h.Options)
	}
	return b, nil
}

func GetIPv4Header(srcIP, dstIP uint32) ([]byte, error) {
	h := ipv4Header{
		Version:  4,
		TTL:      64,  // 最好随机TTL  55 -  64
		Protocol: syscall.IPPROTO_TCP,
		Src:      srcIP,
		Dst:      dstIP,
	}
	return h.Marshal()
}

func ChangeIPV4Header(ipv4Header []byte, srcIP, dstIP uint32,ttl uint8) []byte{
	binary.BigEndian.PutUint32(ipv4Header[12:16], srcIP)
	binary.BigEndian.PutUint32(ipv4Header[16:20], dstIP)
	ipv4Header[8] = ttl // 重置ttl
	return ipv4Header
}
