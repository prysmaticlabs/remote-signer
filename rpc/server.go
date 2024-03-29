package rpc

import (
	"context"
	"fmt"
	"net"

	validatorpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/validator-client"
	"github.com/prysmaticlabs/remote-signer/keyvault"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var log logrus.FieldLogger

func init() {
	log = logrus.WithField("prefix", "rpc")
}

// Config options for the gRPC server.
type Config struct {
	Host     string
	Port     string
	CertFlag string
	KeyFlag  string
	KeyVault keyvault.Store
}

// Server defining a gRPC server for the remote signer API.
type Server struct {
	ctx             context.Context
	cancel          context.CancelFunc
	host            string
	port            string
	listener        net.Listener
	withCert        string
	withKey         string
	credentialError error
	grpcServer      *grpc.Server
	keyVault        keyvault.Store
}

// NewServer instantiates a new gRPC server.
func NewServer(ctx context.Context, cfg *Config) *Server {
	ctx, cancel := context.WithCancel(ctx)
	return &Server{
		ctx:      ctx,
		cancel:   cancel,
		host:     cfg.Host,
		port:     cfg.Port,
		withCert: cfg.CertFlag,
		withKey:  cfg.KeyFlag,
		keyVault: cfg.KeyVault,
	}
}

// Start the gRPC server.
func (s *Server) Start() {
	// Setup the gRPC server options and TLS configuration.
	address := fmt.Sprintf("%s:%s", s.host, s.port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("Could not listen to port in Start() %s: %v", address, err)
	}
	s.listener = lis

	opts := make([]grpc.ServerOption, 0)
	if s.withCert != "" && s.withKey != "" {
		creds, err := credentials.NewServerTLSFromFile(s.withCert, s.withKey)
		if err != nil {
			log.Errorf("Could not load TLS keys: %s", err)
			s.credentialError = err
		}
		opts = append(opts, grpc.Creds(creds))
	} else {
		log.Fatal("Cannot use an insecure gRPC connection. Provide a certificate and key to connect securely")
	}
	log.WithFields(logrus.Fields{
		"crt-path": s.withCert,
		"key-path": s.withKey,
	}).Info("Loaded TLS certificates")
	s.grpcServer = grpc.NewServer(opts...)

	// Instantiate a remote signer server.
	remoteSigner := NewRemoteSigner(s.ctx, s.keyVault)

	// Register services available for the gRPC server.
	validatorpb.RegisterRemoteSignerServer(s.grpcServer, remoteSigner)
	reflection.Register(s.grpcServer)

	go func() {
		if s.listener != nil {
			if err := s.grpcServer.Serve(s.listener); err != nil {
				log.Errorf("Could not serve: %v", err)
			}
		}
	}()
	log.WithField("address", address).Info("gRPC server listening on address")
}

// Stop the gRPC server.
func (s *Server) Stop() error {
	s.cancel()
	if s.listener != nil {
		s.grpcServer.GracefulStop()
		log.Debug("Initiated graceful stop of server")
	}
	return nil
}

// Status returns nil or credentialError.
func (s *Server) Status() error {
	return s.credentialError
}
