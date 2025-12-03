package main

import (
	"context"
	"flag"
	// "fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

// Generic gRPC server main that:
// - listens on a configurable address
// - supports optional TLS
// - installs basic unary interceptors (logging + panic recovery)
// - registers gRPC health service and reflection
// - performs graceful shutdown with timeout
//
// To expose your own services, register them inside registerServices(s).

func main() {
	var (
		addr           = flag.String("addr", ":50051", "gRPC listen address")
		certFile       = flag.String("tls-cert", "", "TLS certificate file (optional)")
		keyFile        = flag.String("tls-key", "", "TLS key file (optional)")
		shutdownTimeout = flag.Duration("shutdown-timeout", 10*time.Second, "graceful shutdown timeout")
	)
	flag.Parse()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", *addr, err)
	}

	var opts []grpc.ServerOption
	// Interceptors
	opts = append(opts, grpc.ChainUnaryInterceptor(loggingUnaryInterceptor, recoveryUnaryInterceptor))

	// TLS if specified
	if *certFile != "" && *keyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("failed to load TLS credentials: %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(opts...)

	// Register health server
	healthSrv := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthSrv)
	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// Register reflection for clients like grpcurl
	reflection.Register(grpcServer)

	// Place to register your own services:
	registerServices(grpcServer)

	// Serve in goroutine
	serverErrCh := make(chan error, 1)
	go func() {
		log.Printf("starting gRPC server on %s (tls=%v)", *addr, *certFile != "" && *keyFile != "")
		if err := grpcServer.Serve(lis); err != nil {
			serverErrCh <- err
		}
		close(serverErrCh)
	}()

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("received signal %v, initiating graceful shutdown", sig)
	case err := <-serverErrCh:
		if err != nil {
			log.Printf("server stopped with error: %v", err)
		} else {
			log.Printf("server stopped")
		}
	}

	// Graceful stop with timeout
	done := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Printf("graceful shutdown completed")
	case <-time.After(*shutdownTimeout):
		log.Printf("graceful shutdown timed out after %s, forcing stop", shutdownTimeout.String())
		grpcServer.Stop()
	}

	// set health to NOT_SERVING before exit
	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	log.Printf("server exited")
}

// registerServices is a placeholder where you should register your gRPC services.
// Example (requires generated pb code):
//   pb.RegisterYourServiceServer(s, &yourServiceImpl{})
func registerServices(s *grpc.Server) {
	// TODO: Register your service implementations here.
	// e.g. myservice.RegisterMyServiceServer(s, myServiceInstance)
	_ = s // keep compiler happy if empty
}

// loggingUnaryInterceptor logs basic info about each unary RPC.
func loggingUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()
	p, _ := peer.FromContext(ctx)
	resp, err = handler(ctx, req)
	duration := time.Since(start)

	clientAddr := "unknown"
	if p != nil {
		clientAddr = p.Addr.String()
	}

	st := status.Convert(err)
	log.Printf("method=%s client=%s duration=%s code=%s msg=%q",
		info.FullMethod, clientAddr, duration, st.Code(), st.Message())

	return resp, err
}

// recoveryUnaryInterceptor recovers from panics in handlers and returns an INTERNAL error.
func recoveryUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			// You might want to capture stack trace here.
			err = status.Errorf(codes.Internal, "panic: %v", r)
		}
	}()
	return handler(ctx, req)
}