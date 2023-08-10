package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/Rosalita/k8s-service/foundation/web"
	"go.uber.org/zap"
)

// Logger uses the power of closures to inject a logger into a
// middleware function that can handle requests with logging.
// The logger will be a single allocation to the heap, meaning it
// exist until the app shuts down, allowing the app to maintain state.
func Logger(log *zap.SugaredLogger) web.Middleware {

	// Create a middleware function as a literal function.
	m := func(handler web.Handler) web.Handler {

		// Create a handler function as a literal function.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			v := web.GetValues(ctx)

			log.Infow("request started", "trace_id", v.TraceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr)

			err := handler(ctx, w, r)

			log.Infow("request completed", "trace_id", v.TraceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr, "status_code", v.StatusCode, "since", time.Since(v.Now))

			return err
		}
		return h
	}
	return m
}
