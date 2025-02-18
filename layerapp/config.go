package layerapp

import (
	"roham/layerapp/service/layer"
	"roham/pkg/grpc"
	httpserver "roham/pkg/http_server"
	"roham/pkg/logger"
	"roham/pkg/postgresql"
	"time"
)

type Config struct {
	LayerSvcCfg layer.Config
	Server      httpserver.Config `koanf:"server"`
	PostgresDB  postgresql.Config `koanf:"postgres_db"`
	Logger      logger.Config     `koanf:"logger"`
	GrpcClient  grpc.Client       `koanf:"grpc_client"`
	//OutboxScheduler      outbox.Config     `koanf:"outbox_scheduler"`
	//RabbitMQ             rabbitmq.Config   `koanf:"rabbitmq"`
	PathOfMigration string `koanf:"path_of_migration"`
	//GRPCServer           grpc.Config       `koanf:"grpc_server"`
	TotalShutdownTimeout time.Duration `koanf:"total_shutdown_timeout"`
}
