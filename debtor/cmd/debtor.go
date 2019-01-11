package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"

	"microsrv/config"

	"github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
	debtorendpoint "microsrv/debtor/endpoint"
	"microsrv/debtor/sd"
	"microsrv/debtor/service"
	"microsrv/debtor/transport"
	"microsrv/pb"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Parameters{}
	pwd, _ := os.Getwd()
	ex, _ := os.Executable()
	ex = filepath.Base(ex)
	iniFile := ""
	fs := flag.NewFlagSet("debtor", flag.ExitOnError)
	configFile := fs.String("cfg", "", "Location of config file")
	if *configFile == "" {
		iniFile = filepath.Join(pwd, strings.TrimSuffix(ex, filepath.Ext(ex))+".ini")
	} else {
		iniFile = *configFile
	}
	err := cfg.Read(iniFile)
	if err != nil {
		fmt.Println(err)
	}
	cfg.Save(iniFile)
	var (
		debugPort  = fs.String("debug.port", fmt.Sprintf(":%d", cfg.Service.DebugPort), "Debug and metrics listen address")
		grpcPort   = fs.String("grpc.port", fmt.Sprintf("%d", cfg.Service.GrpcPort), "gRPC listen address")
		httpAddr   = fs.String("http.addr", fmt.Sprintf("%s", cfg.Service.HTTPAddr), "HTTP Listen Address")
		httpPort   = fs.String("http.port", fmt.Sprintf("%d", cfg.Service.HTTPPort), "HTTP Listen Port")
		consulAddr = fs.String("consul.addr", fmt.Sprintf("%s", cfg.Service.ConsulAddr), "Consul Address")
		consulPort = fs.String("consul.port", fmt.Sprintf("%d", cfg.Service.ConsulPort), "Consul Port")
		db         = fs.String("db.database", fmt.Sprintf("%s", cfg.DB.DB), "Database name")
		user       = fs.String("db.user", fmt.Sprintf("%s", cfg.DB.DbUser), "Database user")
		password   = fs.String("db.password", fmt.Sprintf("%s", cfg.DB.DbPassword), "Database password")
	)
	dsn := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", *user, *password, *db)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	fs.Parse(os.Args[1:])

	iGrpcPort, _ := strconv.Atoi(*grpcPort)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	var service debtorservice.Service
	{
		service = debtorservice.NewDB(dsn)
		service = debtorservice.LoggingMiddleware(logger)(service)
	}

	var (
		endpoints   = debtorendpoint.MakeServerEndpoints(service)
		httpHandler = transport.NewHTTPHandler(endpoints, logger)
		grpcServer  = transport.NewGRPCServer(endpoints, logger)
		registar    = consulsd.ConsulRegister(*consulAddr, *consulPort, *httpAddr, *httpPort, iGrpcPort)
	)
	var g group.Group
	{
		// The debug listener mounts the http.DefaultServeMux, and serves up
		// stuff like the Go debug and profiling routes, and so on.
		debugListener, err := net.Listen("tcp", *debugPort)
		if err != nil {
			logger.Log("transport", "debug/HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "debug/HTTP", "addr", *debugPort)
			return http.Serve(debugListener, http.DefaultServeMux)
		}, func(error) {
			debugListener.Close()
		})
	}
	{
		// The service discovery registration.
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", *httpAddr, "port", *httpPort)
			registar.Register()
			return http.ListenAndServe(":"+*httpPort, httpHandler)
		}, func(error) {
			registar.Deregister()
		})
		defer registar.Deregister()
	}
	{
		// The gRPC listener mounts the Go kit gRPC server we created.
		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%s", *grpcPort))
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "gRPC", "addr", *grpcPort)
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
			pb.RegisterDebtorSvcServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())

}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}
