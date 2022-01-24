package playground

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	name    = "github.com/happybydefault/opentelemetry-playground"
	version = "v0.1.0"
)

type Playground struct {
	mux    *chi.Mux
	tracer trace.Tracer
}

func New(options ...Option) *Playground {
	p := Playground{
		mux: chi.NewMux(),
		tracer: otel.GetTracerProvider().Tracer(name,
			trace.WithInstrumentationVersion(version),
		),
	}

	for _, option := range options {
		option(&p)
	}

	p.registerRoutes()

	return &p
}

func (p *Playground) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.mux.ServeHTTP(w, r)
}
