package playground

import "go.opentelemetry.io/otel/trace"

type Option func(p *Playground)

func WithInstrumentation(tp trace.TracerProvider) Option {
	return func(p *Playground) {
		p.tracer = tp.Tracer(name,
			trace.WithInstrumentationVersion(version),
		)
	}
}
