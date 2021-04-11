package main

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/prysmaticlabs/remote-signer/keyvault"
	"github.com/prysmaticlabs/remote-signer/keyvault/deterministic"
	"github.com/prysmaticlabs/remote-signer/keyvault/mnemonic"
	"github.com/prysmaticlabs/remote-signer/rpc"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "main")

var (
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
	keyVaultFlag = flag.String(
		"keyvault",
		"deterministic",
		"Type of keyvault. Examples: "+
			"deterministic (default) | mnemonic | s3 (unimplemented) | hashicorp (unimplemented)",
	)
	numDeterministicKeysFlag = flag.Int(
		"num-deterministic-keys",
		1,
		"Number of deterministic keys to generate for a deterministic keyvault (demonstrative purposes)",
	)
	numMnemonicKeysFlag = flag.Int(
		"num-mnemonic-keys",
		1,
		"Number of keys to generate from a mnemonic phrase",
	)
	startIndexForMnemonicFlag = flag.Int(
		"start-index",
		0,
		"Start index for mnemonic keys generation",
	)
	mnemonicFileFlag = flag.String(
		"mnemonic-file",
		"",
		"Path to the mnemonic file, containing the mnemonic phrase",
	)
	mnemonicPasswordFlag = flag.String(
		"mnemonic-password",
		"",
		"Password of the mnemonic phrase",
	)
)

func main() {
	flag.Parse()
	ctx := context.Background()
	grpcServerHost := *grpcServerHostFlag
	grpcServerPort := *grpcServerPortFlag
	tlsCertPath := *tlsCertPathFlag
	tlsKeyPath := *tlsKeyPathFlag
	keyVaultKind := *keyVaultFlag
	numDeterministicKeys := *numDeterministicKeysFlag
	numMnemonicKeys := *numMnemonicKeysFlag
	startIndexForMnemonic := *startIndexForMnemonicFlag
	mnemonicFile := *mnemonicFileFlag
	mnemonicPassword := *mnemonicPasswordFlag

	if tlsCertPath == "" || tlsKeyPath == "" {
		log.Fatal("Expected --tls-crt-path and --tls-key-path flags for secure connections")
	}

	// Initialize keyvault kind as specified by user.
	var vault keyvault.Store
	var err error
	switch keyVaultKind {
	case "deterministic":
		log.Warn(
			"You are using a deterministic keyvault (only for reference purposes) " +
				"DO NOT USE in production",
		)
		vault, err = deterministic.NewStore(numDeterministicKeys)
	case "mnemonic":
		log.Warn(
			"Using a mnemonic to recover keys from",
		)

		if mnemonicFile == "" {
			log.Fatal("You must provide a --mnemonic-file")
		}

		content, err := ioutil.ReadFile(mnemonicFile)
		if err != nil {
			logrus.Info("File reading error", err)
			return
		}
		mnemonicPhrase := string(content)
		logrus.Info("Contents of file: ", mnemonicPhrase)

		vault, err = mnemonic.NewStore(
			mnemonicPhrase,
			mnemonicPassword,
			startIndexForMnemonic,
			numMnemonicKeys)
	default:
		log.Fatalf("Keyvault kind %s not yet supported", keyVaultKind)
	}
	if err != nil {
		log.Fatalf("Could not initialize keyvault: %v", err)
	}

	// Initialize new gRPC server.
	srv := rpc.NewServer(ctx, &rpc.Config{
		Host:     grpcServerHost,
		Port:     grpcServerPort,
		CertFlag: tlsCertPath,
		KeyFlag:  tlsKeyPath,
		KeyVault: vault,
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
