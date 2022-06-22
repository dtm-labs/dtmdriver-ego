package main

import (
	"context"

	"github.com/dtm-labs/dtmdriver-clients/gozero/trans/pb"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego-component/eetcd"
	"github.com/gotomicro/ego-component/eetcd/registry"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server"
	"github.com/gotomicro/ego/server/egrpc"
)

// export EGO_DEBUG=true EGO_NAME=hello && go run main.go --config=config.toml
func main() {
	if err := ego.New().
		Registry(registry.Load("registry").Build(registry.WithClientEtcd(eetcd.Load("etcd").Build()))).
		Serve(func() server.Server {
			server := egrpc.Load("server.grpc").Build()
			pb.RegisterTransSvcServer(server.Server, &Trans{server: server})
			return server
		}()).Run(); err != nil {
		elog.Panic("startup", elog.Any("err", err))
	}
}

type Trans struct {
	server *egrpc.Component
}

func (s *Trans) TransIn(ctx context.Context, in *pb.AdjustInfo) (*pb.Response, error) {
	return &pb.Response{}, nil
}

func (s *Trans) TransOut(ctx context.Context, in *pb.AdjustInfo) (*pb.Response, error) {
	return &pb.Response{}, nil
}
