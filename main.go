package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/eriklupander/dtls"

	"github.com/eriklupander/tradfri-go/grpc_server"
	pb "github.com/eriklupander/tradfri-go/grpc_server/golang"
	"github.com/eriklupander/tradfri-go/router"
	"github.com/eriklupander/tradfri-go/tradfri"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFlags = pflag.NewFlagSet("config", pflag.ExitOnError)
var commandFlags = pflag.NewFlagSet("commands", pflag.ExitOnError)

func init() {

	configFlags.String("gateway_ip", "", "IP to your gateway. No protocol or port here!")
	configFlags.String("gateway_address", "", "Address to your gateway. Including port here!")
	configFlags.String("psk", "", "Pre-shared key on bottom of Gateway")
	configFlags.String("client_id", "", "Your client id, make something up or use the NNN-NNN-NNN on the bottom of your Gateway")
	configFlags.String("loglevel", "info", "Log level. Allowed values: fatal, error, warn, info, debug, trace")

	commandFlags.Bool("server", false, "Start in server mode?")
	commandFlags.Bool("authenticate", false, "Perform PSK exchange?")
	commandFlags.String("get", "", "URL to GET")
	commandFlags.String("put", "", "URL to PUT")
	commandFlags.String("payload", "", "Payload for PUT")
	commandFlags.String("listen_host", "", "Host to listen on. Default empty allows connections from anywhere. Use \"127.0.0.1\" to only allow local connections.")
	commandFlags.Int("port", 8080, "Port of the REST server. Set to 0 to disable REST server.")
	commandFlags.Int("grpc_port", 8081, "Port of the gRPC server. Set to 0 to disable gRPC server.")

	commandFlags.AddFlagSet(configFlags)
	_ = commandFlags.Parse(os.Args[1:])

	_ = viper.BindPFlags(configFlags)
	viper.AutomaticEnv()
	viper.AddConfigPath(".") // e.g. reads ./config.json or config.yaml
	err := viper.ReadInConfig()
	if err != nil {
		slog.Info(err.Error())
		slog.Info("You probably have to run --authenticate first")
	}
	viper.RegisterAlias("pre_shared_key", "psk")
}

func main() {
	// configure logging
	levelStr := viper.GetString("loglevel")
	if levelStr == "" {
		levelStr = "info"
	}
	fmt.Printf("Using loglevel: %v\n", levelStr)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetLogLoggerLevel(parseLevel(levelStr))
	slog.SetDefault(logger)
	dtls.SetLogFunc(func(ts time.Time, level string, peer string, msg string) {
		switch level {
		case "error":
			slog.Error(msg, slog.String("level", level), slog.String("peer", peer))
		case "warn":
			slog.Warn(msg, slog.String("level", level), slog.String("peer", peer))
		case "info":
			slog.Info(msg, slog.String("level", level), slog.String("peer", peer))
		case "debug":
			slog.Debug(msg, slog.String("level", level), slog.String("peer", peer))
		}
	})
	dtls.SetLogLevel(resolveDTLSLogLevel(levelStr))

	gatewayAddress := viper.GetString("gateway_address")
	if gatewayAddress == "" {
		gatewayAddress = viper.GetString("gateway_ip") + ":5684"
	}
	psk := viper.GetString("psk")
	clientID := viper.GetString("client_id")
	serverMode, _ := commandFlags.GetBool("server")
	authenticate, _ := commandFlags.GetBool("authenticate")
	get, getErr := commandFlags.GetString("get")
	put, putErr := commandFlags.GetString("put")
	payload, _ := commandFlags.GetString("payload")
	listenHost, _ := commandFlags.GetString("listen_host")
	port, _ := commandFlags.GetInt("port")
	grpcPort, _ := commandFlags.GetInt("grpc_port")

	// Handle the special authenticate use-case
	if authenticate {
		performTokenExchange(gatewayAddress, clientID, psk)
		return
	}

	checkRequiredConfig(gatewayAddress, clientID, psk)

	// Check running mode...
	if serverMode {
		slog.Info("Running in server mode")

		tc := tradfri.NewTradfriClient(gatewayAddress, clientID, psk)
		wg := sync.WaitGroup{}
		// REST
		if port > 0 {
			wg.Add(1)
			slog.Info(fmt.Sprintf("REST: %s:%d", listenHost, port))
			go func() {
				defer wg.Done()
				router.SetupChi(tc, fmt.Sprintf("%s:%d", listenHost, port))
			}()
		}
		// gRPC
		if grpcPort > 0 {
			wg.Add(1)
			slog.Info(fmt.Sprintf("gRPC: %s:%d", listenHost, grpcPort))
			go func() {
				defer wg.Done()
				go registerGrpcServer(tc, fmt.Sprintf("%s:%d", listenHost, grpcPort))
			}()
		}

		wg.Wait()
	} else {
		// client mode
		if getErr == nil && get != "" {
			resp, _ := tradfri.NewTradfriClient(gatewayAddress, clientID, psk).Get(get)
			slog.Info(string(resp.Payload))
		} else if putErr == nil && put != "" {
			resp, _ := tradfri.NewTradfriClient(gatewayAddress, clientID, psk).Put(put, payload)
			slog.Info(string(resp.Payload))
		} else {
			slog.Info("No client operation was specified, supported one(s) are: get, put, authenticate")
		}
	}

}

func parseLevel(str string) slog.Level {
	switch strings.ToLower(str) {
	case "debug":
		return slog.LevelDebug
	case "error":
		return slog.LevelError
	case "warn":
		return slog.LevelWarn
	case "info":
		fallthrough
	default:
		return slog.LevelInfo
	}
}

func checkRequiredConfig(gatewayAddress, clientID, psk string) {
	if gatewayAddress == "" {
		fail("Unable to resolve gatewayAddress from command-line flag or config.json file")
	}
	if clientID == "" {
		fail("Unable to resolve clientID from command-line flag or config.json file")
	}
	if psk == "" {
		fail("Unable to resolve psk (pre shared key) from command-line flag or config.json file")
	}
}

func performTokenExchange(gatewayAddress, clientID, psk string) {
	if len(clientID) < 1 || len(psk) < 10 {
		fail("Both clientID and psk args must be specified when performing key exchange")
	}

	done := make(chan bool)
	defer func() { done <- true }()
	go func() {
		select {
		case <-time.After(time.Second * 5):
			slog.Info("(Please note that the key exchange may appear to be stuck at \"Connecting to peer at\" if the PSK from the bottom of your Gateway is not entered correctly.)")
		case <-done:
		}
	}()

	// Note that we hard-code "Client_identity" here before creating the DTLS client,
	// required when performing token exchange
	dtlsClient := tradfri.NewTradfriClient(gatewayAddress, "Client_identity", psk)

	authToken, err := dtlsClient.AuthExchange(clientID)
	if err != nil {
		fail(err.Error())
	}
	viper.Set("client_id", clientID)
	viper.Set("gateway_address", gatewayAddress)
	viper.Set("psk", authToken.Token)
	err = viper.WriteConfigAs("config.json")
	if err != nil {
		fail(err.Error())
	}
	slog.Info("Your configuration including the new PSK and clientID has been written to config.json, keep this file safe!")
}

func registerGrpcServer(tc *tradfri.Client, listenAddress string) {
	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.StandardLogger())),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_logrus.StreamServerInterceptor(logrus.NewEntry(logrus.StandardLogger())),
		),
	)
	pb.RegisterTradfriServiceServer(s, grpc_server.New(tc))
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		slog.Info(fmt.Sprintf("failed to listen on grpc %s: %v", listenAddress, err.Error()))
		return
	}
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		slog.Info(err.Error())
	}
}

func fail(msg string) {
	slog.Info(msg)
	os.Exit(1)
}

// resolveDTLSLogLevel maps our logrus levels to the ones supported by the DTLS library.
func resolveDTLSLogLevel(level string) string {
	switch level {
	case "fatal":
		fallthrough
	case "error":
		return dtls.LogLevelError
	case "warn":
		return dtls.LogLevelWarn
	case "info":
		return dtls.LogLevelInfo
	case "debug":
		fallthrough
	case "trace":
		return dtls.LogLevelDebug
	}
	return "info"
}
