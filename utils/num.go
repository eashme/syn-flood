package utils

import (
	cRand "crypto/rand"
	"errors"
	"math/big"
	mRand "math/rand"
	"net"
	"os"
	"time"
)

var randSeed int64 = 10517284972

func MathRandNum(max int64) int64 {
	mRand.Seed(time.Now().Unix())
	return mRand.Int63n(max)
}

func RandNum(max int64) int64 {
	if max <= 0 {
		return 0
	}
	index, err := cRand.Int(cRand.Reader, big.NewInt(max))
	if err != nil {
		return MathRandNum(max)
	}
	return index.Int64()
}

func String2IPV4(ip string) (net.IP, error) {
	ipv4 := net.ParseIP(ip).To4()
	if ipv4 == nil {
		return nil, errors.New("ip is not valid")
	}
	return ipv4, nil
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// 根据上一个数值生成伪随机数,计算速度要更快
func FakeRandNum(lastValue int64,max int64) int64{
	if max <= 0{
		return 0
	}
	lastValue = lastValue ^ randSeed
	// c++ 伪随机数生成计算法
	return (lastValue ^ (lastValue >> 11) ^ (lastValue ^ (lastValue << 4) >> 8) + 1231) % max
}

func UInt64FakeRandNum(lastValue uint64,max uint64) uint64{
	if max <= 0{
		return 0
	}
	lastValue = lastValue ^ uint64(randSeed)
	// c++ 伪随机数生成计算法
	return (lastValue ^ (lastValue >> 11) ^ (lastValue ^ (lastValue << 4) >> 8) + uint64(randSeed)) % max
}

// 每隔一段时间更换种子,保证一段时间内
func ChangeRandSeed(){
	randSeed = mRand.Int63()
}