package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/pkg/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	httpServer *http.Server
	logger     Logger
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

func NewServer(logger Logger, grpcAddr string, host, port string) *Server {
	s := &Server{
		logger: logger,
	}

	ctx := context.Background()
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := api.RegisterEventServiceHandlerFromEndpoint(ctx, gwmux, grpcAddr, opts)
	if err != nil {
		logger.Error("failed to register grpc gateway: " + err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.Handle("/", gwmux)
	mux.HandleFunc("GET /swagger", s.serveSwaggerUI)
	mux.HandleFunc("GET /swagger.json", s.serveSwaggerJSON)
	mux.HandleFunc("GET /proto", s.serveProto)

	handler := loggingMiddleware(logger, mux)

	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf("%s:%s", host, port),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return s
}

func (s *Server) serveSwaggerUI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, swaggerUIHTML)
}

func (s *Server) serveSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("pkg", "api", "EventService.swagger.json"))
}

func (s *Server) serveProto(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("api", "EventService.proto"))
}

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@3/swagger-ui.css">
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@3/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@3/swagger-ui-standalone-preset.js"></script>
    <script>
    window.onload = function() {
      const ui = SwaggerUIBundle({
        url: "/swagger.json",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      })
      window.ui = ui
    }
    </script>
</body>
</html>`

func (s *Server) Start(_ context.Context) error {
	s.logger.Info("starting http gateway server on " + s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http gateway server failed: %w", err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("stopping http gateway server")
	return s.httpServer.Shutdown(ctx)
}
