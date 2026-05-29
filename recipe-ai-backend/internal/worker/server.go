package worker

import (
	"time"
	"recipe-ai-backend/internal/pkg/config"

	"github.com/hibiken/asynq"
)

// Server Worker服务器
type Server struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

// NewServer 创建Worker服务器
func NewServer(redisOpt asynq.RedisConnOpt, cfg config.WorkerConfig) *Server {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: cfg.Concurrency,
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				return time.Duration(n+1) * 10 * time.Second
			},
		},
	)

	return &Server{
		server: server,
		mux:    asynq.NewServeMux(),
	}
}

// Run 启动Worker
func (s *Server) Run() error {
	return s.server.Run(s.mux)
}

// Shutdown 停止Worker
func (s *Server) Shutdown() {
	s.server.Shutdown()
}

// GetMux 获取ServeMux
func (s *Server) GetMux() *asynq.ServeMux {
	return s.mux
}
