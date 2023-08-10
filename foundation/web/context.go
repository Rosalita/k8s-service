package web

import (
	"context"
	"time"
)

// It is important for ctxKey to be unexported.
type ctxKey int

const key ctxKey = 1

// Values represent state for each request
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// GetValues returns the values from the context.
// If values don't exist, construct them.
func GetValues(ctx context.Context) *Values {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return &Values{
			TraceID: "0000000-0000-0000-0000-000000000000",
			Now:     time.Now(),
		}
	}
	return v
}

// GetTraceID returns the trace ID from the context.
// If there is no trace ID, return a default trace ID.
func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return "0000000-0000-0000-0000-000000000000"
	}
	return v.TraceID
}

// GetTime returns the time from the context.
// If there is no time in the context, it returns the time now.
func GetTime(ctx context.Context) time.Time {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return time.Now()
	}
	return v.Now
}

// SetStatusCode sets the status code back into the context.
func SetStatusCode(ctx context.Context, statusCode int) {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return
	}
	v.StatusCode = statusCode
}
