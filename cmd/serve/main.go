package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	playground "github.com/happybydefault/opentelemetry-playground"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

const name = "github.com/happybydefault/opentelemetry-playground/cmd/serve"
const version = "v0.1.0"

var environment = os.Getenv("PLAYGROUND_ENVIRONMENT")

func main() {
	var address string
	flag.StringVar(&address, "address", "0.0.0.0:", "Listener's TCP Address")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := serve(ctx, address); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func serve(ctx context.Context, address string) error {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint())

	r, err := resource.New(context.Background(),
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithTelemetrySDK(),
		// resource.WithProcess(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(name),
			semconv.ServiceVersionKey.String(version),
			semconv.DeploymentEnvironmentKey.String(environment),
		),
	)
	if err != nil {
		return err
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.ParentBased(trace.AlwaysSample())),
		trace.WithBatcher(exporter),
		trace.WithResource(r),
	)
	defer func() {
		if err := tp.Shutdown(context.TODO()); err != nil {
			log.Printf("could not shutdown tracer provider: %v", err)
		}
	}()

	p := playground.New(
		playground.WithInstrumentation(tp),
	)

	handler := otelhttp.NewHandler(
		p,
		"handle request",
		otelhttp.WithTracerProvider(tp),
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	)

	svr := http.Server{
		Addr:    address,
		Handler: handler,
	}

	l, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("could not listen: %w", err)
	}
	log.Println("listening at", l.Addr())

	errCh := make(chan error, 1)
	go func() {
		errCh <- svr.Serve(l)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Println("shutting down gracefully")
		return svr.Shutdown(context.TODO())
	}
}
