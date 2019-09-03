package base

import (
	log "github.com/sirupsen/logrus"
	"imooc.com/resk/infra"
	"net"
	"net/rpc"
	"reflect"
)

var rpcServer *rpc.Server

func RpcServer() *rpc.Server {
	Check(rpcServer)
	return rpcServer
}

func RpcRegister(ri interface{}) {
	typ := reflect.TypeOf(ri)
	log.Infof("goRPC Register: %s", typ.String())
	RpcServer().Register(ri)
}

type GoRPCStarter struct {
	infra.BaseStarter
	server *rpc.Server
}

func (s *GoRPCStarter) Init(ctx infra.StarterContext) {
	log.Info("GoRPCStarter Init()")
	s.server = rpc.NewServer()
	rpcServer = s.server
}

func (s *GoRPCStarter) Start(ctx infra.StarterContext) {
	port := ctx.Props().Section("app").Key("rpc.port").MustString("18082")
	// 监听网络端口
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Panic(err)
	}
	log.Info("tcp port listened for rpc:", port)
	// 处理网络连接和请求
	go s.server.Accept(listener)
}
