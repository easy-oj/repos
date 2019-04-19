package initial

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/easy-oj/common/logs"
	"github.com/easy-oj/common/proto/repos"
	"github.com/easy-oj/common/settings"
	"github.com/easy-oj/repos/common/caller"
	"github.com/easy-oj/repos/common/database"
	"github.com/easy-oj/repos/service"
	"github.com/easy-oj/repos/service/http"
)

func Initialize() {
	caller.InitCaller()
	database.InitDatabase()

	address := fmt.Sprintf("0.0.0.0:%d", settings.Repos.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	repos.RegisterReposServiceServer(server, service.NewReposHandler())
	reflection.Register(server)
	go func() {
		if err := server.Serve(lis); err != nil {
			panic(err)
		}
	}()
	logs.Info("[Initialize] service served on %s", address)

	http.StartHTTPService()
}
