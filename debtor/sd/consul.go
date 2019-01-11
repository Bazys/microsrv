package consulsd

import (
	"math/rand"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
)

// ConsulRegister method.
func ConsulRegister(consulAddress string,
	consulPort string,
	advertiseAddress string,
	advertisePort string,
	grpcPort int) sd.Registrar {

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	// Service discovery domain. In this example we use Consul.
	var client consul.Client
	{
		consulConfig := api.DefaultConfig()
		consulConfig.Address = consulAddress + ":" + consulPort
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		client = consul.NewClient(consulClient)
	}

	check := api.AgentServiceCheck{
		HTTP:     "http://" + advertiseAddress + ":" + advertisePort + "/health",
		Interval: "10s",
		Timeout:  "1s",
		Notes:    "Basic health checks",
	}

	asr := api.AgentServiceRegistration{
		ID:      "debtor",
		Name:    "debtor",
		Address: advertiseAddress,
		Port:    grpcPort,
		Tags:    []string{"debtor", "timetable"},
		Check:   &check,
	}
	return consul.NewRegistrar(client, &asr, logger)

}
