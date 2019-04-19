package caller

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/easy-oj/common/logs"
	"github.com/easy-oj/common/proto/oss"
	"github.com/easy-oj/common/proto/queue"
	"github.com/easy-oj/common/settings"
)

var (
	OSSClient   oss.OSSServiceClient
	QueueClient queue.QueueServiceClient
)

func InitCaller() {
	OSSClient = oss.NewOSSServiceClient(dial("OSS", fmt.Sprintf("%s:%d", settings.OSS.Hosts[0], settings.OSS.Port)))
	QueueClient = queue.NewQueueServiceClient(dial("Queue", fmt.Sprintf("%s:%d", settings.Queue.Hosts[0], settings.Queue.Port)))
}

func dial(service, target string) *grpc.ClientConn {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	logs.Info("[Caller] dial %s service on %s", service, target)
	return conn
}
