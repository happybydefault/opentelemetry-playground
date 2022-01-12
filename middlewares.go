package playground

import (
	"net/http"
)

func (p *Playground) withInstrumentation(name string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx, span := p.tracer.Start(
			r.Context(),
			name,
		)
		defer span.End()

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return fn
}
