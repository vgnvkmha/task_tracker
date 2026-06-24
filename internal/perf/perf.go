package perf

import (
	"context"
	"log"
	"time"
)

type endpointKeyType struct{}

var endpointKey = endpointKeyType{}

var trackedEndpoints = map[string]struct{}{
	"/user/list_active":             {},
	"/team/list":                    {},
	"/board/list":                   {},
	"/task/get_by_board/:board_id":  {},
	"/task/get_by_board/:board_id/": {},
}

func IsTrackedEndpoint(endpoint string) bool {
	_, ok := trackedEndpoints[endpoint]
	return ok
}

func WithEndpoint(ctx context.Context, endpoint string) context.Context {
	if !IsTrackedEndpoint(endpoint) {
		return ctx
	}
	return context.WithValue(ctx, endpointKey, endpoint)
}

func Endpoint(ctx context.Context) (string, bool) {
	endpoint, ok := ctx.Value(endpointKey).(string)
	if !ok || endpoint == "" {
		return "", false
	}
	return endpoint, true
}

func LogStep(ctx context.Context, step string) {
	endpoint, ok := Endpoint(ctx)
	if !ok {
		return
	}
	log.Printf("[PERF] endpoint=%s step=%s", endpoint, step)
}

func Track(ctx context.Context, step string) func() {
	endpoint, ok := Endpoint(ctx)
	if !ok {
		return func() {}
	}

	start := time.Now()
	return func() {
		log.Printf("[PERF] endpoint=%s step=%s duration=%s", endpoint, step, time.Since(start))
	}
}
