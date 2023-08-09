package mid

import (
	"context"
	"net/http"

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

			log.Infow("request started", "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

			err := handler(ctx, w, r)

			log.Infow("request completed", "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

			return err
		}
		return h
	}
	return m
}
