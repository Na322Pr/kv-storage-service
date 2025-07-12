package main

import (
	"context"
	"fmt"
	kv_storage_service "github.com/Na322Pr/kv-storage-service/internal/app/kv-storage-service"
	"github.com/Na322Pr/kv-storage-service/internal/config"
	"github.com/Na322Pr/kv-storage-service/internal/service"
	"github.com/Na322Pr/kv-storage-service/internal/storage"
	desc "github.com/Na322Pr/kv-storage-service/pkg/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.MustLoad()
	grpcAddress := cfg.GetGRPCAddress()

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	//logger := zap.NewProduction()
	//defer logger.Sync()

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	logger.Info("Service config info",
		zap.Int("id", cfg.Node.ID),
		zap.Bool("isSeed", cfg.Node.IsSeed),
		zap.String("seedNodes", strings.Join(cfg.Node.SeedNodes, ",")),
		zap.String("grpcAddress", grpcAddress),
	)

	registry := prometheus.NewRegistry()

	// Expose the metrics endpoint
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	go http.ListenAndServe(":2112", nil)

	nodeService := service.NewNodeService(cfg.Node.ID, grpcAddress, cfg.Node.IsSeed, logger)

	keyValueStorage := storage.NewKeyValueInMemoryStorage()
	storageService := service.NewStorageService(keyValueStorage)

	storeApp := kv_storage_service.NewImplementation(nodeService, storageService)

	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	desc.RegisterKeyValueStorageServer(grpcServer, storeApp)

	logger.Info(fmt.Sprintf("Starting grpc server on %s...", grpcAddress))
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to listen and server grpc server: %v", err)
		}
	}()

	logger.Info("Starting gossiping...")
	if err := nodeService.Run(ctx, cfg.SeedNodes); err != nil {
		log.Fatalf("Failed to start node: %v", err)
	}

	<-stop
	fmt.Println("\nShutting down servers...")
	os.Exit(0)
}

// Ноды начинают процесс голосования, когда лидер нода выбывает
// (n *Node) CleanUpStalePeers
//
// Нода собирает голоса с остальных нол
// Лидер, сам рассылает информацию о себе
// Ноды сами определяют лидера
// Есть независимый сервис, который это все определяет.

// Должен быть какой-то сервис, который будет расспределять нагрузку между нодами.
// Соответственно через него, можно будет создать/удалить ноды.
