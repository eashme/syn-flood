package packet

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"syscall"
	"time"
	"unsafe"
)

const (
	tcpHeaderLen    = 20
	tcpMaxHeaderLen = 60
)

type PsdHeader struct {
	SrcAddr   [4]byte
	DstAddr   [4]byte
	Zero      uint8
	ProtoType uint8
	TcpLength uint16
}

// A tcp header
type tcpHeader struct {
	Src     uint16    //源端口
	Dst     uint16    //目的端口
	Seq     int    //序号
	Ack     int    //确认号
	Len     int    //头部长度
	Rsvd    int    //保留位
	Flag    int    //标志位
	Win     int    //窗口大小
	Sum     int    //校验和
	Urp     int    //紧急指针
	Options []byte // 选项, extension headers
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

// Marshal encode tcp header
func (h *tcpHeader) Marshal() ([]byte, error) {
	if h == nil {
		return nil, syscall.EINVAL
	}

	hdrlen := tcpHeaderLen + len(h.Options)
	b := make([]byte, hdrlen)
	//版本和头部长度
	binary.BigEndian.PutUint16(b[0:2], h.Src)
	binary.BigEndian.PutUint16(b[2:4], h.Dst)

	binary.BigEndian.PutUint32(b[4:8], uint32(h.Seq))
	binary.BigEndian.PutUint32(b[8:12], uint32(h.Ack))

	b[12] = uint8(hdrlen/4<<4 | 0)
	//TODO  Rsvd
	b[13] = uint8(h.Flag)
	binary.BigEndian.PutUint16(b[14:16], uint16(h.Win))
	binary.BigEndian.PutUint16(b[16:18], uint16(h.Sum))
	binary.BigEndian.PutUint16(b[18:20], uint16(h.Urp))

	if len(h.Options) > 0 {
		copy(b[tcpHeaderLen:], h.Options)
	}
	return b, nil
}

func GetTcpHeader(srcIp, dstIp uint32, srcPort, dstPort uint16) ([]byte, error) {
	rand.Seed(time.Now().UnixNano())
	h := &tcpHeader{
		Src:  srcPort,
		Dst:  dstPort,
		Seq:  rand.Intn(1<<32 - 1),
		Ack:  0,
		Flag: 0x02,
		Win:  60000,
		Urp:  0,
	}
	/*buffer用来写入两种首部来求得校验和*/
	buffs, _ := h.Marshal()
	h.Sum = int(CheckSum(append(GetPsdByte(srcIp, dstIp),buffs...)))
	return h.Marshal()
}

func GetTcpByte(srcPort, dstPort uint16) ([]byte, error) {
	rand.Seed(time.Now().UnixNano())
	h := &tcpHeader{
		Src:  srcPort,
		Dst:  dstPort,
		Seq:  rand.Intn(1<<32 - 1),
		Ack:  0,
		Flag: 0x02,
		Win:  60000,
		Urp:  0,
	}
	/*buffer用来写入两种首部来求得校验和*/
	return h.Marshal()
}

// 获取伪首部字节码
func GetPsdByte(srcIp, dstIp uint32) []byte {
	var psdheader PsdHeader
	/*填充TCP伪首部*/
	binary.BigEndian.PutUint32(psdheader.SrcAddr[:4], srcIp)
	binary.BigEndian.PutUint32(psdheader.DstAddr[:4], dstIp)
	psdheader.Zero = 0
	psdheader.ProtoType = syscall.IPPROTO_TCP
	psdheader.TcpLength = uint16(unsafe.Sizeof(tcpHeader{})) + uint16(0)
	var (
		buffer bytes.Buffer
	)
	binary.Write(&buffer, binary.BigEndian, psdheader)
	return buffer.Bytes()
}

// 改变包体内容
func ChangeTcpByte(psdHeader []byte, tcpHeader []byte, srcIP, dstIP uint32, srcPort, dstPort uint16) []byte {
	// 替换伪首部包体的源IP和目标IP
	binary.BigEndian.PutUint32(psdHeader[0:4], srcIP)
	binary.BigEndian.PutUint32(psdHeader[4:8], dstIP)
	// 替换TCP包体的源IP和目标IP
	binary.BigEndian.PutUint16(tcpHeader[0:2], srcPort)
	binary.BigEndian.PutUint16(tcpHeader[2:4], dstPort)
	// 重置check sum
	binary.BigEndian.PutUint16(tcpHeader[16:18], 0)
	// 计算并替换变更后的checkSum
	binary.BigEndian.PutUint16(tcpHeader[16:18], CheckSum(append(psdHeader, tcpHeader...)))
	return tcpHeader
}

