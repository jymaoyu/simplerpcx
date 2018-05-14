package main

import (
	"context"
	"flag"
	"log"

	"github.com/smallnest/rpcx/client"
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

func main() {
	flag.Parse()
	//使用p2p的方式连接服务器参数为“协议@主机地址：端口”，其中协议可以为tcp、http、unix、quic或kcp
	d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
	//NewXClient为建立一个xclient对象，xclient封装好的client对象，增加了服务发现和管理功能
	//NewXClient第一个参数必须为服务名称，第二个参数为失败模式，第三个参数是选择器，第四个参数是注册中心，最后的为其它选项

	/*失败模式。rpcx中有四种方式处理请求失败
	1、Failfast。在这种模式下，当rpcx请求失败时，便会终止请求

	2、Failover。在这种模式下，当rpcx请求失败时，会尝试请求其他节点，它会尝试重试知道服务返回正常响应。重试测试可在defaultOption定义

	3、Failtry。在这种模式下，当rpcx请求失败时，会继续重试当前节点，它会尝试重试知道服务返回正常响应。重试测试可在defaultOption定义

	4、Failbackup。在这种模式下，如果第一个请求在给定时间内没有返回，则rpcx会向另一个服务器发送另一个请求*/

	/*选择器。rpcx中有多个选择器用于多个相同服务节点间的选择
	1、RandomSelect。这个选择器将随机选取一个节点。

	2、Roundrobin。以循环方式选择节点。

	3、WeightedRoundRobin。使用了Nginx中的加权轮询算法。
	根据服务器的不同处理能力，给每个服务器分配不同的权值，使其能够接受相应权值数的服务请求。
	比如序列{a, a, a, a, a, b, c}中，前五个请求都会分配给服务器a，
	这就是一种不均匀的分配方法，更好的序列应该是：{a, a, b, a, c, a, a}。

	4、WeightedICMP。使用ping（ICMP）的结果来设置每个节点的权重。ping时间越短，节点的权重越高。

	5、ConsistentHash。一致性哈希算法为相同的servicePath、serviceMethod和参数配置相同的节点。

	6、Geography。通过地理位置选择最近的节点。

	7、SelectByUser Selector。用于自定义选择器。*/

	/*注册中心。服务注册中心用来实现服务发现，rpcx中有多种注册中心。
	1、peer2peer。最简单的方式，其实这种方式没有注册中心，客户端直接指定服务器地址，因为只有一个节点所以选择器没有意义。
	注册方式如：d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")。

	2、peer2peer。有多个服务却没有注册中心时使用，直接在客户端指定多个服务器地址。
	注册方式如：d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr1}, {Key: *addr2}})。

	3、zookeeper。Apache的一个命名注册项目，使用时服务器必须用ZooKeeperRegisterPlugin插件注册到zookeeper。
	而客户端则需要使用ZookeeperDiscovery来获取服务信息。需要设置服务器地址、zookeeper地址和basspath地址（用于防止命名冲突）。
	注册方式如：d := client.NewZookeeperDiscovery(*basePath, "Arith",[]string{*zkAddr}, nil)。

	4、etcd。由CoreOS 团队发起的开源项目，是一个键值存储仓库用于配置共享和服务发现。其在rpcx中的使用和zookeeper类似，都
	需要在服务器和客户端中配置插件，插件配置的方法也基本相似，不在累述。
	注册方式如:d := client.NewEtcdDiscovery(*basePath, "Arith",[]string{*etcdAddr}, nil)。

	5、Consul。由HashiCorp公司推出的开源工具，用于实现分布式系统的服务发现与配置。用法与zookeeper和etcd相似。
	注册方式如：d := client.NewConsulDiscovery(*basePath, "Arith",[]string{*consulAddr}, nil)。

	6、mDNS。通过组播dns协议使局域网内的主机实现互相发现和通信。同样需要在服务器和客户端中设置插件。
	注册方式如：d := client.NewMDNSDiscovery("Arith", 10*time.Second, 10*time.Second, "")

	注：以上四种插件在注册时都可设置度量（Metrics）和刷新租约（UpdateInterval）

	7、In process。用于测试
	*/
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := Args{
		Numa: 222,
		Numb: 333,
	}

	reply := &Reply{}

	/*call方法用于同步服务调用，调用方法为：
	Call(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
	参数为context元素，服务方法名，需要传入的参数，返回的参数。返回的是一个error类型
	而go方法则用于异步服务调用，方法为：
	Go(ctx context.Context, serviceMethod string, args interface{}, reply interface{}, done chan *Call) (*Call, error)
	参数为context元素，服务方法名，需要传入的参数，返回的参数，call指针。返回一个call指针和一个error类型*/
	err := xclient.Call(context.Background(), "Add", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Printf("%d * %d = %d\n"+"%s", args.Numa, args.Numb, reply.Sum, reply.D)

}
