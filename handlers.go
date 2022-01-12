package playground

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/codes"
)

func (p *Playground) helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	err := p.work(r.Context(), 150*time.Millisecond)
	if err != nil {
		writeError(w, http.StatusNotFound)
		return
	}

	fmt.Fprintln(w, "Hello world!")
}

func (p *Playground) work(ctx context.Context, duration time.Duration) error {
	_, span := p.tracer.Start(
		ctx,
		"work",
	)
	defer span.End()

	time.Sleep(duration)

	fail := rand.Intn(2) == 0
	if fail {
		err := errors.New("failed work")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func writeError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
