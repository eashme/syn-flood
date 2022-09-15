# syn-flood
参考https://github.com/jiangeZh/SYN_Flood这位大佬的C语言版本实现了一个 go版本的syn-flood,

> 1.程序只支持linux版本编译运行

> 2.单机单线程测试发包效率为 单线程 60M/s 左右

> 3.使用config文件夹下的cfg.yaml文件对程序进行控制

> 4.依然可以继续优化至逼近单机极限

####

>>目前版本思路为,提前构造好包,每个线程开启循环, 根据config文件生成的ip去修改包体对应位置应该修改的字节,以达到发出不同包的目的

## TODO LIST
    > 性能瓶颈在生成随机数的地方,但是随机数的最终目标是服从均匀分布的去生成数值,所以我们可以在有限的数值中,预先生成全部可能,然后一起发送,可以达到同样的效果,不再生成随机数,性能也可以得到巨大提升。


# 后记
 > 1. IP_HDRINCL(原始数据包) 高版本golang在windows下编译已经不在支持,我写这个程序是在(go 1.11环境下), win7(据说高版本的windows也不支持)
 > 2. 之前加深网络方面知识做的玩具,如果有想试试这个程序的欢迎提issues 