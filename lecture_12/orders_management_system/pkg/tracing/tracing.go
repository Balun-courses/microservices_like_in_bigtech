package tracing

import (
	"context"
	"os"

	"github.com/moguchev/microservices_courcse/orders_management_system/pkg/closer"
	"github.com/uber/jaeger-client-go/config"
)

func Init(serviceName string) error {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: os.Getenv("JAEGER_HOST"),
		},
	}

	close, err := cfg.InitGlobalTracer(serviceName)
	if err != nil {
		return err
	}

	closer.Add(func(ctx context.Context) error {
		return close.Close()
	})

	return nil
}
