package utils

import (
	"fmt"
	"github.com/awnumar/fastrand"
	"math/rand"
	"net"
	"sender/packet"
	"testing"
	"time"
)

func TestGetIPRange(t *testing.T) {
	GetIPRange("192.168.2.1", "192.168.2.255")
}

func TestExists(t *testing.T) {
	timer := time.NewTimer(time.Second)
	for {
		select {
		case <-timer.C:
			return
		default:
			break
		}
		fmt.Println("ok")
	}
}

func TestRandomIP(t *testing.T) {
	dst, err := RandomIP([2]net.IP{
		net.ParseIP("192.168.1.23").To4(),
		net.ParseIP("192.168.2.33").To4(),
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dst)
}

func TestInt2TCPPort(t *testing.T) {
	for i := 0; i < 200; i ++ {
		fmt.Println(net.ParseIP(RandomGetIP("172.16.33.1", "172.16.33.255")).To4())
	}
}

func TestIP2Uint32(t *testing.T) {
	for ip := range GetIPRange("172.16.33.2", "172.16.33.2") {
		p := net.ParseIP(ip).To4()
		fmt.Println(fmt.Sprintf("%s", p))
		var addr [4]byte
		for i := 0; i < len(p); i ++ {
			addr[i] = p[i]
		}
		fmt.Println(addr)
	}
}

func TestUint322IP(t *testing.T) {
	s := IP2Uint32(net.ParseIP("192.168.3.183"))
	e := IP2Uint32(net.ParseIP("192.169.7.152"))
	fmt.Println(e - s)
}

func TestMathRandNum(t *testing.T) {
	var (
		a = 500
		b = 8000
	)
	count := 0
	for i := 0; i < b; i ++ {
		if i%(b/a+1) == 0 {
			fmt.Println(i)
			count ++
		}
	}
	fmt.Println(count)
}

func TestRandNum(t *testing.T) {
	var (
		srcIP, dstIP     uint32 = 12345678, 12345678
		srcPort, dstPort uint16 = 3306, 3306
	)
	tcpB, err := packet.GetTcpByte(srcPort, dstPort)
	if err != nil {
		panic(err)
	}
	psdB := packet.GetPsdByte(srcIP, dstIP)
	b, err := packet.GetTcpHeader(srcIP, dstIP, srcPort, dstPort)
	if err != nil {
		panic(err)
	}
	a := packet.ChangeTcpByte(psdB, tcpB, 888853, 777752, uint16(srcPort), uint16(dstPort))
	fmt.Println(a, b)
}

func TestFakeRandNumSpeed(t *testing.T) {
	var (
		v   int64 = 0
		max int64 = 65535
	)
	start := time.Now().UnixNano()
	for i := 0; i < 10000000; i ++ {
		RandNum(max)
	}
	fmt.Println("D-Value: ", time.Now().UnixNano()-start)

	start = time.Now().UnixNano()
	for i := 0; i < 10000000; i ++ {
		v = FakeRandNum(v, max)
	}
	fmt.Println("D-Value: ", time.Now().UnixNano()-start)
}


func TestFakeRandNumSpeed2(t *testing.T) {
	var ips = []uint32{1736792064, 1736792065, 1736792066, 1736792067, 1736792068, 1736792069, 1736792070, 1736792071, 1736792072, 1736792073 ,1736792074, 1736792075, 1736792076, 1736792077, 1736792078, 1736792079, 1736792080}
	il := len(ips)
	a := func(last uint32) uint32{
		return ips[FakeRandNum(int64(last),int64(il))]
	}

	b := func(v uint32) uint32{
		return uint32(FakeRandNum(int64(v),16) + 1736792064)
	}

	var v uint32
	start := time.Now().UnixNano()
	for i :=0 ; i < 100000000; i ++ {
		v = a(v)
	}
	fmt.Println("D-value: ", time.Now().UnixNano() - start)

	start = time.Now().UnixNano()
	for i :=0 ; i < 100000000; i ++ {
		v = b(v)
	}
	fmt.Println("D-value: ", time.Now().UnixNano() - start)

	start = time.Now().UnixNano()
	for i := 0; i < 100000000; i ++ {
		rand.Intn(100)
	}
	fmt.Println("D-value: ", time.Now().UnixNano() - start)

	start = time.Now().UnixNano()
	for i := 0; i < 100000000; i ++ {
		fastrand.Intn(100)
	}
	fmt.Println("D-value: ", time.Now().UnixNano() - start)
}

func getRandIP(lastIP uint32) uint32 {
	return uint32(FakeRandNum(int64(lastIP), 255)) + IP2Uint32(net.ParseIP("172.16.3.0"))
}

func TestFakeRandNum(t *testing.T) {
	var (
		v uint32
		lastv uint32
	)
	for i := 0 ;i < 100 ; i ++ {
		lastv = v
		v = getRandIP(lastv)
		fmt.Println(lastv - v)
	}
}