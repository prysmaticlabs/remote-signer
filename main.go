package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/prysmaticlabs/remote-signer/rpc"
	"github.com/sirupsen/logrus"
)

var (
	//allowedOriginsFlag = flag.String(
	//	"gateway-allowed-origins",
	//	"*",
	//	"comma-separated, allowed origins for gateway cors",
	//)
	grpcServerHostFlag = flag.String(
		"grpc-server-host",
		"127.0.0.1",
		"host address for the grpc server",
	)
	grpcServerPortFlag = flag.String(
		"grpc-server-port",
		"4000",
		"port for the grpc server",
	)
	tlsCertPathFlag = flag.String(
		"tls-crt-path",
		"",
		"/path/to/server.crt for secure TLS connections",
	)
	tlsKeyPathFlag = flag.String(
		"tls-key-path",
		"",
		"/path/to/server.key for secure TLS connections",
	)
)

func main() {
	flag.Parse()
	ctx := context.Background()

	//grpcGatewayAddress := *grpcGatewayAddressFlag
	//grpcServerAddress := *grpcServerAddressFlag
	//lis, err := net.Listen("tcp", grpcGatewayAddress)
	//if err != nil {
	//	logrus.Errorf("Could not listen to port in Start() %s: %v", grpcServerAddress, err)
	//}

	//opts := []grpc.ServerOption{}
	//grpcServer := grpc.NewServer(opts...)
	//pb.RegisterAPIServer(grpcServer, &server{})
	//
	//go func() {
	//	logrus.WithField("address", grpcGatewayAddress).Info("gRPC server listening on port")
	//	if err := grpcServer.Serve(lis); err != nil {
	//		logrus.Errorf("Could not serve gRPC: %v", err)
	//	}
	//}()

	//allowedOrigins := []string{*allowedOriginsFlag}
	//if strings.Contains(*allowedOriginsFlag, ",") {
	//	allowedOrigins = strings.Split(*allowedOriginsFlag, ",")
	//}
	//gatewaySrv := gateway.New(
	//	ctx,
	//	grpcGatewayAddress,
	//	grpcServerAddress,
	//	nil, /*optional http mux */
	//	allowedOrigins,
	//)
	//gatewaySrv.Start()
	grpcServerHost := *grpcServerHostFlag
	grpcServerPort := *grpcServerPortFlag
	tlsCertPath := *tlsCertPathFlag
	tlsKeyPath := *tlsKeyPathFlag

	if tlsCertPath == "" || tlsKeyPath == "" {
		log.Fatal("Expected --tls-crt-path and --tls-key-path flags for secure connections")
	}

	srv := rpc.NewServer(ctx, &rpc.Config{
		Host:     grpcServerHost,
		Port:     grpcServerPort,
		CertFlag: tlsCertPath,
		KeyFlag:  tlsKeyPath,
	})
	srv.Start()

	// Listen for any process interrupts.
	stop := make(chan struct{})
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigc)
		<-sigc
		logrus.Info("Got interrupt, shutting down...")
		if err := srv.Stop(); err != nil {
			log.Fatal(err)
		}
		stop <- struct{}{}
	}()

	// Wait for stop channel to be closed.
	<-stop
}
