package config

import(
	"log"
	"time"
	"golang.org/x/net/context"
	"github.com/coreos/etcd/client"
	//"strings"
	"strings"
	"fmt"
)

const EAGLEHE_HEALTH_PATH = "/eagleye/health"

var EtcdHosts string

var GroupName string

var HeartbeatConfig string

var kapi client.KeysAPI

func InitEtcdClient() {
	//flag.StringVar(&EtcdHosts, "etcdHosts", "", "Please input etct hosts, eg: http://xxx.xxx.xxx.xxx:2379,http://xxx.xxx.xxx.xxx:2379")
	//flag.StringVar(&GroupName, "groupName", "", "Please input group name, eg: packetbeats")
	//
	//flag.Parse()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Init Etcd Client . Recovering, but please report this: [", r, "]")
		}
	}()

	log.Println("init config")

	if EtcdHosts==""{
		log.Println("EtcdHosts is nil, please try it again, eg: http://xxx.xxx.xxx.xxx:2379,http://xxx.xxx.xxx.xxx:2379")
		return
	}

	if GroupName==""{
		log.Println("GroupName is nil, please try is again, eg: packetbeats")
		return
	}

	//for _, host := range strings.Split(EtcdHosts, ","){
	//	log.Println(host)
	//}

	cfg := client.Config{
		Endpoints: strings.Split(EtcdHosts, ","),
		//Endpoints: []string{"http://192.168.10.235:2379"},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	api := client.NewKeysAPI(c)
	kapi = api

	//创建该group对应的目录
	createGroupDir()

}




func createGroupDir(){

	// 如果存在不会重复覆盖, 直接返回key already exists
	setOptions := &client.SetOptions{
		//TTL: time.Second * 60,
		Dir: true,
		PrevExist: client.PrevNoExist,
	}



	resp, err := kapi.Set(context.Background(), EAGLEHE_HEALTH_PATH + "/"+GroupName, "", setOptions)
	if err != nil {
		//log.Printf(string(err))
	} else {
		// print common key info
		log.Printf("Set is done. Metadata is %q\n", resp)
	}

	log.Printf("Etcd hosts : [%s]", EtcdHosts)
}


/**

 */
func SetHeartbeatDataToEtcd(key string, value string){

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Upload packetbeat health to etcd . Recovering, but please report this: [", r, "]")
		}
	}()

	setOptions := &client.SetOptions{
		TTL: time.Second * 65,
		Dir: false,
		PrevExist: client.PrevIgnore,
	}

	keyPath :=  EAGLEHE_HEALTH_PATH + "/" + GroupName + "/" + key

	_, err := kapi.Set(context.Background(), keyPath, value, setOptions)
	if err != nil {
		fmt.Println("Upload packetbeat health throw exception : [", err, "]")
		//log.Fatal(err)
	} else {
		// print common key info
		//log.Printf("Set is done. Metadata is %q\n", resp)
	}

}






