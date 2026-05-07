package grpcserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/usecases"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Server struct {
	api.UnimplementedEventServiceServer
	logger     Logger
	grpcServer *grpc.Server
	usecases   map[usecases.Type]usecases.Usecase
}

func NewServer(logger Logger, ucs map[usecases.Type]usecases.Usecase) *Server {
	s := &Server{
		logger:   logger,
		usecases: ucs,
	}

	s.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(s.loggingInterceptor),
	)
	api.RegisterEventServiceServer(s.grpcServer, s)
	reflection.Register(s.grpcServer)

	return s
}

func (s *Server) loggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)

	msg := fmt.Sprintf("method: %s, duration: %s", info.FullMethod, duration)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s, error: %v", msg, err))
	} else {
		s.logger.Info(msg)
	}

	return resp, err
}

func (s *Server) Start(host, port string) error {
	addr := net.JoinHostPort(host, port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.logger.Info("starting grpc server on " + addr)
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

func (s *Server) Stop() {
	s.logger.Info("stopping grpc server")
	s.grpcServer.GracefulStop()
}
