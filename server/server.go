// Package server 导出 wiki-api 的可装配 HTTP 服务。
//
// 单仓部署走 cmd/server；聚合部署（starcat-api）import 本包并挂到网关。
// 业务实现仍在 internal/，本包只负责 env 装配、路由注册、定时任务与生命周期。
package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	kitenv "github.com/starcat-app/starcat-api-kit/env"
	"github.com/starcat-app/starcat-wiki-api/internal/handler"
	"github.com/starcat-app/starcat-wiki-api/internal/middleware"
	"github.com/starcat-app/starcat-wiki-api/internal/probe"
	"github.com/starcat-app/starcat-wiki-api/internal/scheduler"
	"github.com/starcat-app/starcat-wiki-api/internal/store"
	"github.com/starcat-app/starcat-wiki-api/internal/version"
)

const defaultPort = "5004"

const defaultStoreFile = "./wiki.db"

// Options 控制 wiki 服务装配。聚合网关可显式传入，单仓部署通常用 FromEnv。
type Options struct {
	Port                   string
	StoreFile              string
	APIKeys                []string
	ProbeUserAgent         string
	EnableCodewikiBatchRPC bool
	SkipListenLogEndpoints bool
}

// Service 是已装配的 wiki HTTP 服务。
type Service struct {
	opts         Options
	handler      http.Handler
	sqliteStore  *store.SQLiteStore
	sch          *scheduler.Scheduler
	probeHandler *handler.ProbeHandler
}

// Name 返回聚合网关识别用的稳定服务名。
func Name() string { return "wiki" }

// DefaultPort 返回单仓默认监听端口。
func DefaultPort() string { return defaultPort }

// FromEnv 从环境变量装配服务（与历史 cmd/server 行为一致）。
func FromEnv() (*Service, error) {
	apiKeys, err := kitenv.RequiredCSV("API_KEYS")
	if err != nil {
		return nil, fmt.Errorf("API_KEYS env is required")
	}
	opt := Options{
		Port:           kitenv.OrDefault("PORT", defaultPort),
		StoreFile:      kitenv.OrDefault("STORE_FILE", defaultStoreFile),
		APIKeys:        apiKeys,
		ProbeUserAgent: kitenv.OrDefault("PROBE_USER_AGENT", ""),
		// 历史行为只认字面 "true"，不用 ParseBool，避免 "1"/"TRUE" 语义漂移。
		EnableCodewikiBatchRPC: kitenv.OrDefault("ENABLE_CODEWIKI_BATCHEXECUTE", "") == "true",
	}
	return New(opt)
}

// New 按 Options 装配服务。
func New(opt Options) (*Service, error) {
	if strings.TrimSpace(opt.Port) == "" {
		opt.Port = defaultPort
	}
	if strings.TrimSpace(opt.StoreFile) == "" {
		opt.StoreFile = defaultStoreFile
	}
	if len(opt.APIKeys) == 0 {
		return nil, fmt.Errorf("APIKeys is required")
	}

	sqliteStore, err := store.NewSQLiteStore(opt.StoreFile)
	if err != nil {
		return nil, fmt.Errorf("initialize SQLite: %w", err)
	}

	baseReq := probe.NewBaseRequest(opt.ProbeUserAgent)
	probes := probe.DefaultRegistry(baseReq, opt.EnableCodewikiBatchRPC)
	probeHandler := handler.NewProbeHandler(sqliteStore, probes)
	sch := scheduler.New(sqliteStore, probeHandler.RetryErrors)

	authMW := middleware.NewBearerAuth(opt.APIKeys)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthzHandler)
	mux.Handle("GET /api/v1/ping", authMW.Wrap(handler.HandlePingV1(version.Service, version.Version)))
	mux.Handle("GET /api/v1/wikis", authMW.Wrap(http.HandlerFunc(probeHandler.HandleProbeV1)))
	mux.Handle("POST /api/v1/wikis/batch", authMW.Wrap(http.HandlerFunc(probeHandler.HandleProbeBatchV1)))
	mux.Handle("POST /internal/sync/probe", authMW.Wrap(handler.HandleAdminSyncProbe(sch)))
	mux.Handle("POST /internal/refresh/owner", authMW.Wrap(handler.HandleAdminRefreshOwner(sch)))

	if !opt.SkipListenLogEndpoints {
		log.Printf("starcat-wiki-api %s starting on port %s", version.Version, opt.Port)
		log.Printf("Endpoints:")
		log.Printf("  GET  /api/v1/ping                   - Connectivity probe for Starcat client")
		log.Printf("  GET  /api/v1/wikis?owner=X&repo=Y  - Single probe")
		log.Printf("  POST /api/v1/wikis/batch             - Batch probe (max 50, async)")
		log.Printf("  POST /internal/sync/probe             - Manual sync trigger")
		log.Printf("  POST /internal/refresh/owner          - Owner refresh")
		log.Printf("  GET  /healthz                         - Health check")
	}

	return &Service{
		opts:         opt,
		handler:      middleware.CORS(mux),
		sqliteStore:  sqliteStore,
		sch:          sch,
		probeHandler: probeHandler,
	}, nil
}

// Handler 返回已包 CORS 的根 handler，可供聚合网关挂载。
func (s *Service) Handler() http.Handler { return s.handler }

// Addr 返回建议监听地址（":port"）。
func (s *Service) Addr() string { return ":" + s.opts.Port }

// StartBackground 启动 cron 与 probing 恢复（与原 main 中 goroutine 一致）。
func (s *Service) StartBackground() {
	go s.sch.Start()
	go s.probeHandler.RecoverPendingProbes()
}

// Close 停止调度器并关闭 SQLite。
func (s *Service) Close() error {
	if s.sch != nil {
		s.sch.Stop()
	}
	if s.sqliteStore != nil {
		return s.sqliteStore.Close()
	}
	return nil
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
