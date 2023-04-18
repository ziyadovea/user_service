package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	delivery_grpc "github.com/ziyadovea/task_manager/users/internal/app/delivery/grpc"
	"github.com/ziyadovea/task_manager/users/internal/app/delivery/grpc/interceptors"
	"github.com/ziyadovea/task_manager/users/internal/app/repository/postgresql"
	"github.com/ziyadovea/task_manager/users/internal/app/usecase"
	"github.com/ziyadovea/task_manager/users/internal/authentication/jwt"
	"github.com/ziyadovea/task_manager/users/internal/config"
	"github.com/ziyadovea/task_manager/users/proto/v1/pb"
)

func Run(cfg config.Config) {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// init repo layer
	db, err := sqlx.Connect("pgx", cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	repo := postgresql.NewUserRepository(db)

	// init JWT authenticator
	auth := jwt.NewAuthenticator(
		[]byte(os.Getenv(config.AccessTokenSecretEnvKey)),
		[]byte(os.Getenv(config.RefreshTokenSecretEnvKey)),
		cfg.AccessTokenExpirationDuration,
		cfg.RefreshTokenExpirationDuration,
	)

	// init usecase layer
	uc := usecase.NewUserUsecase(repo, auth)

	// init delivery layer
	userGRPCService := delivery_grpc.NewUserService(uc)

	// start the gRPC server
	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.Logging(),
			interceptors.Auth(auth),
		),
	)
	pb.RegisterUserServiceServer(gRPCServer, userGRPCService)
	reflection.Register(gRPCServer)

	gRPCListener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// start the gRPC gateway server
	gatewayMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err = pb.RegisterUserServiceHandlerFromEndpoint(ctx, gatewayMux, ":"+cfg.GRPCPort, opts); err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}
	gatewayServer := &http.Server{
		Addr:    ":" + cfg.RestPort,
		Handler: gatewayMux,
	}

	// start the gRPC server goroutine
	go func() {
		log.Printf("starting gRPC server on %s", gRPCListener.Addr())
		if err := gRPCServer.Serve(gRPCListener); err != nil {
			log.Fatal(err)
		}
	}()

	// start the gRPC gateway server goroutine
	go func() {
		log.Printf("starting gRPC gateway server on %s", gatewayServer.Addr)
		if err := gatewayServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// listen for the interrupt signal
	<-ctx.Done()

	// restore default behavior on the interrupt signal and notify user of shutdown
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// perform application shutdown with a maximum timeout of 5 seconds
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shut down gracefully
	gRPCServer.GracefulStop()
	if err := gatewayServer.Shutdown(timeoutCtx); err != nil {
		log.Fatalf("gateway server shutdown errror: %v", err)
	}
}
