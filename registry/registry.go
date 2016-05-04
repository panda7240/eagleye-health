package registry

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
	"os"
	"os/exec"
	"encoding/json"
	"github.com/siye1982/eagleye-health/config"

	"strconv"
)

const(
	// eagleye环境变量
	EAGLEYE_HOST = "EAGLEYE_HOST"

	// 日期格式
	DATA_FORMAT = "2006-01-02 15:04:05"

)

//throughput per min
var tpm uint64 = 0

//total throughput
var tt uint64 = 0

// boot time
var bootTime string = time.Now().Format(DATA_FORMAT)


type Heartbeat struct {
	//`json:"body,omitempty"` 如果为空置则忽略字段
	//`json:"-"`  直接忽略字段
	Pid int `json:"-"`
	Tt uint64 `json:"tt"`
	Tpm uint64 `json:"tpm"`
	Host string `json:"-"`
	Config string `json:"config"`
	Group string `json:"-"`
	Btime string `json:"btime"`
}


/**
组装心跳数据
 */
func assemble_health_info() Heartbeat{
	var heartbeat Heartbeat
	heartbeat.Btime = bootTime
	heartbeat.Config = config.HeartbeatConfig
	heartbeat.Group = "packetbeat"
	heartbeat.Host = getHost()
	heartbeat.Tpm = getTpm()
	heartbeat.Tt = getTt()
	heartbeat.Pid = os.Getpid()

	_, err := json.Marshal(heartbeat)
	if err != nil {
		fmt.Println("json err:", err)
	}
	return heartbeat

}

/**
注册信息到etcd
 */
func Regist(heartbeat Heartbeat) {
	key := heartbeat.Host + "," + strconv.Itoa(heartbeat.Pid)
	value, err := json.Marshal(heartbeat)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Println("upload heartbeat : " + key + "   " + string(value))
	config.SetHeartbeatDataToEtcd(key, string(value))
}



/**
每分钟生产力计数器
 */
func TpmCounter(){
	// 使用AddUint64函数为计数器进行自增操作，向其传递计数器的内存地址作为第一个参数
	atomic.AddUint64(&tpm, 1)
	// 允许其他协程进行处理
	runtime.Gosched()
}

/**
获取每分钟生产力
 */
func getTpm() uint64{
	throughput := atomic.LoadUint64(&tpm)
	return throughput
}



/**
生产力总量计数器
 */
func TtCounter(){
	// 使用AddUint64函数为计数器进行自增操作，向其传递计数器的内存地址作为第一个参数
	atomic.AddUint64(&tt, 1)
	// 允许其他协程进行处理
	runtime.Gosched()
}

/**
获取总生产力
 */
func getTt() uint64{
	// 当其他协程正在更新的时候，为了安全使用计数器，我们通过 LoadUint64 释出一份当前值的拷贝到 opsFinal 中
	// 和上面一样，我们需要给这个函数传递计数器的内存地址
	total_throughput := atomic.LoadUint64(&tt)
	return total_throughput
}

/**
获取host,优先获取环境变量 EAGLEYE_HOST, 如果没有获取本机ip地址
 */
func getHost() string{
	// 获取环境变量
	eagleye_host := os.Getenv(EAGLEYE_HOST)
	if eagleye_host == "" {//如果为空, 则去本机ip地址
		cmd := exec.Command("/bin/sh", "-c",`/sbin/ifconfig -a|grep inet|grep -v 127.0.0.1|grep -v inet6|awk '{print $2}'|tr -d "addr:"`)
		result, err := cmd.Output()
		if err != nil {
			panic(err.Error())
		}
		return string(result[:])
	}
	return eagleye_host
}


func init() {



	//go func() {
	//	for {
	//
	//		time.Sleep(60 * time.Second)
	//		tpm = 0
	//	}
	//}()

	go func() {
		for {

			Regist(assemble_health_info())
			//60秒对tpm进行一次清零
			tpm = 0
			// 每隔60秒上传一次心跳信息
			time.Sleep(60 * time.Second)

		}
	}()
}


/**
启动入口
 */
func Start(){

}



