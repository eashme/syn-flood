package app

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var Cfg Config

func init() {
	rep, err := ioutil.ReadFile("./config/cfg.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(rep, &Cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println(Cfg)
}

type Config struct {
	// 源IP 端口
	Src IP `yaml:"src"`
	// 目标IP 目标端口
	Dst IP `yaml:"dst"`
	// 线程数
	ThreadCount int `yaml:"thread_count"`
	// 连接数
	Connection int `yaml:"connection"`
	// 包体内容
	Packet Packet `yaml:"packet"`
	// 每个伪造源IP 需要发多久包
	TimeOut int `yaml:"timeout"`
}

type IP struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
	Port  Port   `yaml:"port"`
	File  string `yaml:"file"`
}

type Port struct {
	Start int `yaml:"start"`
	End   int `yaml:"end"`
}

type Packet struct {
	Protocol bool `yaml:"protocol"`    // tcp / udp
	Flag     struct {
		Syn bool `yaml:"syn"`
		Urg bool `yaml:"urg"`
		Ack bool `yaml:"ack"`
		Psh bool `yaml:"psh"`
		Rst bool `yaml:"rst"`
		Fin bool `yaml:"fin"`
	} `yaml:"flag"`                    // 标志位 URG ACK PSH RST SYN FIN
	Urgent   int    `yaml:"urgent"`    // 紧急标志
	BodyFile string `yaml:"body_file"` // 包体内容文件路径
	Body     string `yaml:"body"`      // 包体内容
}
