package utils

import (
	"fmt"
	"io/ioutil"
	"net"
	"reflect"
	"strconv"
	"strings"
)

func RandomPort(ports [2]int) int {
	if ports[0] == ports[1] {
		return ports[0]
	}
	start, end := ports[0], ports[1]
	if start > end {
		start, end = end, start
	}
	len := end - start + 1
	return start + int(RandNum(int64(len)))
}

func RandomIP(ips [2]net.IP) (net.IP, error) {
	if ips[0] == nil && ips[1] == nil {
		return nil, nil
	}
	if ips[0] == nil {
		return ips[1], nil
	}
	if ips[1] == nil {
		return ips[0], nil
	}
	if reflect.DeepEqual(ips[0], ips[1]) {
		return ips[0], nil
	}
	start := IP2Uint32(ips[0])
	end := IP2Uint32(ips[1])
	if end < start {
		start, end = end, start
	}

	len := end - start + 1
	rand := RandNum(int64(len))
	got := uint32(rand) + start
	return Uint322IP(got), nil
}

func GetIPRange(start string, end string) <-chan string {
	s, _ := strconv.ParseInt(strings.SplitN(start, ".", 4)[3], 10, 64)
	e, _ := strconv.ParseInt(strings.SplitN(end, ".", 4)[3], 10, 64)
	prefix := start[:strings.LastIndex(start, ".")]
	if prefix != end[:strings.LastIndex(start, ".")] {
		panic("起始ip和截止ip网段不同错误")
	}
	sc := make(chan string)
	go func() {
		for i := s; i < e+1; i ++ {
			sc <- fmt.Sprintf("%s.%d", string(prefix), i)
		}
		close(sc)
	}()
	return sc
}

func RandomGetIP(start string, end string) string {
	s, _ := strconv.ParseInt(strings.SplitN(start, ".", 4)[3], 10, 64)
	e, _ := strconv.ParseInt(strings.SplitN(end, ".", 4)[3], 10, 64)
	prefix := start[:strings.LastIndex(start, ".")]
	if prefix != end[:strings.LastIndex(start, ".")] {
		panic("起始ip和截止ip网段不同错误")
	}
	if s > e {
		e, s = s, e
	}
	return fmt.Sprintf("%s.%d", prefix, RandNum(e-s + 1)+s)
}

func Uint322IP(ip uint32) net.IP {
	res := make([]byte, 0, 4)
	res = append(res, uint8(ip>>24))
	res = append(res, uint8(ip>>16))
	res = append(res, uint8(ip>>8))
	res = append(res, uint8(ip))
	return net.IP(res)
}

func IP2Uint32(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3])
}

func Uint32TO4Byte(ip uint32) [4]byte{
	var b [4]byte
	b[0] = uint8(ip>>24)
	b[1] = uint8(ip>>16)
	b[2] = uint8(ip>>8)
	b[3] = uint8(ip)
	return b
}

func GetIPs(path string) []uint32 {
	var (
		ips []uint32
		b   []byte
		ip net.IP
		err error
	)
	if Exists(path) {
		if b, err = ioutil.ReadFile(path); err != nil {
			return nil
		}
		d := strings.Split(string(b), "\n")
		for i := 0; i < len(d); i ++ {
			ip = net.ParseIP(strings.TrimSpace(d[i])).To4()
			if ip == nil{
				continue
			}
			ips = append(ips,IP2Uint32(ip))
		}
	}
	return ips
}
