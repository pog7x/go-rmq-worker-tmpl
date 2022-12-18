package server

import (
	"context"
	"encoding/json"
	"net/http"

	"net/http/pprof"

	"github.com/pog7x/go-rmq-worker-tmpl/internal/app/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type HTTPServer struct {
	logger *zap.Logger
	server *http.Server
}

func NewHTTPServer(logger *zap.Logger, cfg *config.Config) HTTPServer {
	return HTTPServer{
		logger: logger,

		server: &http.Server{
			Addr:         cfg.ServerListenAddr,
			ReadTimeout:  cfg.ServerReadTimeout,
			WriteTimeout: cfg.ServerWriteTimeout,
		},
	}
}

func (s HTTPServer) Start() error {
	s.registerHandlers()
	s.logger.Sugar().Infof("Diagnostic server is running on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s HTTPServer) registerHandlers() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	mux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		respBody, _ := json.Marshal(map[string]bool{"success": true})

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respBody)
	})
	s.server.Handler = mux
}

func (s HTTPServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
