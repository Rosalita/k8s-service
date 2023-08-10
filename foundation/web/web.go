// Package web contains a small web framework extension.
package web

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
)

// A Handler is a type that handles a http request within our own little mini framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data.logic on this App struct.
type App struct {
	*httptreemux.ContextMux // embedded pointer, grants App the methods of ContextMux
	shutdown                chan os.Signal
	mw                      []Middleware
}

// NewApp creates an App value that handles a set of routes for the application.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw:         mw,
	}
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mix. This method overrides the the Handle
// method on *httptreemux.ContextMux.
// This prevents *httptreemux.ContextMux's Handle method from being
// promoted up to App.
func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {
	// Bind route level middleware first.
	handler = wrapMiddleware(mw, handler)
	// Then bind application level middleware.
	handler = wrapMiddleware(a.mw, handler)

	// A literal function called h is defined to be the outer function.
	// The handler is called by h. This allows for code to be added
	// before and after h is called.
	h := func(w http.ResponseWriter, r *http.Request) {

		// Generate a unique ID and add to context.
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now(),
		}
		ctx := context.WithValue(r.Context(), key, &v)

		if err := handler(ctx, w, r); err != nil {
			// to do: handle error
			return
		}
	}

	a.ContextMux.Handle(method, path, h)
}
