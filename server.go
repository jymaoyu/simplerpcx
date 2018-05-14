package main

import (
	"context"
	"flag"

	"github.com/smallnest/rpcx/server"
)

var (
	addr = flag.String("addr", "localhost:8972", "server address")
)

type Args struct {
	Numa int
	Numb int
}

type Reply struct {
	Sum int
	D   string
}
type Arith struct{}

// 需要调用的方法
func (t *Arith) Add(ctx context.Context, args Args, reply *Reply) error {
	reply.Sum = args.Numa + args.Numb
	reply.D = "hello rpcx"
	return nil
}

func main() {
	flag.Parse()

	s := server.NewServer()
	//使用Register方法注册的服务，其服务名称为类型名
	//s.Register(new(Arith), "")
	//使用RegisterName方法注册的服务，其服务名称则为方法名
	s.RegisterName("Arith", new(Arith), "")
	//通过Serve或ServeHttp方法启动监听，
	err := s.Serve("tcp", *addr)
	if err != nil {
		panic(err)
	}
	/*在注册中心时，P2P和P2M方式里，服务器的配置都不需要改动。而在使用插件注册时，则需要在服务器上配置插件。
		ZooKeeper中，服务器插件ZooKeeperRegisterPlugin。主要有5个参数：
		ServiceAddress: 本机的监听地址， 这个对外暴露的监听地址， 格式为tcp@ipaddress:port
		ZooKeeperServers: Zookeeper集群的地址
		BasePath: 服务前缀。 如果有多个项目同时使用zookeeper，避免命名冲突，可以设置这个参数，为当前的服务设置命名空间
		Metrics: 用来更新服务的TPS
		UpdateInterval: 服务的刷新间隔， 如果在一定间隔内(当前设为2 * UpdateInterval)没有刷新,服务就会从Zookeeper中删除
		注意，插件必须要在注册服务前配置添加，否则无效
		配置的方法如下：
		func addRegistryPlugin(s *server.Server) {
	    r := &serverplugin.ZooKeeperRegisterPlugin{
	        ServiceAddress:   "tcp@" + *addr,
	        ZooKeeperServers: []string{*zkAddr},
	        BasePath:         *basePath,
	        Metrics:          metrics.NewRegistry(),
	        UpdateInterval:   time.Minute,
	    }
	    err := r.Start()
	    if err != nil {
	        log.Fatal(err)
	    }
	    s.Plugins.Add(r)
	}
		而etcd插件和consul插件与zookeeper插件配置基本类似，不在累述
	*/
}
