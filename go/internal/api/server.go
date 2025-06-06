// Run the API to serve the quiz JSON

package api

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/herrewig/tomedome/go/internal/dota"
	"github.com/sirupsen/logrus"
)

// We allow JS to call the API from these origins:
// - https://dota.tomedome.io
// - http://localhost:8080
func getCorsOrigin(reqHost string) string {
	exceptions := []string{
		"https://dota.tomedome.io",
	}
	for _, exception := range exceptions {
		if reqHost == exception {
			return exception
		}
	}
	return "http://localhost:8080"
}

// Sets up all the middleware for the API
// The middleware order matters:
//
// Outer layer: validate the route -- return 404 for any invalid routes
// Next: ensure no params -- return 400 for any calls with params
// Next: ensure Client IP is correctly parsed out of X-Forwarded-For header
//
//	for Google Cloud L7 LBs (second-to-last IP)
//
// Inner layer: rate limit the calls with in-memory db
func newHandler(log *logrus.Entry, enableRateLimiting bool, routes []string, mux http.Handler) http.Handler {
	// Nothing to do if we're not rate limiting
	if !enableRateLimiting {
		return mux
	}
	handler := newLimiterMiddleware(mux)
	handler = newClientIpMiddleware(handler)
	handler = newParamValidationMiddleware(handler)
	handler = newRouteValidationMiddleware(routes, handler)
	return newOuterMiddleware(log, handler)
}

// Do this for all calls to the API
// For now, all it does is reject any calls with params
// by returning a 400 Bad Request
func newParamValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := getLogger(r.Context())
		if len(r.URL.Query()) > 0 {
			log.Warn("call has params. dropping")
			// Don't be too specific in the error message so bad people
			// can't figure stuff out
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Before all the param validation and rate-limiting stuff, return a 404 if it's
// not a valid route.
func newRouteValidationMiddleware(routes []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := getLogger(r.Context())

		var found bool = false
		for _, route := range routes {
			if r.URL.Path == route {
				found = true
				next.ServeHTTP(w, r)
			}
		}
		if !found {
			log.Warn("invalid route, returning 404")
			http.Error(w, "not found", http.StatusNotFound)
		}
	})
}

// Returns *http.Server with all the routes and handlers, and middleware
func newServer(log *logrus.Entry, enableRateLimiting bool, host string, dotes dota.DotaService) *http.Server {

	// Assign handlers to routes
	routeHandlers := map[string]http.Handler{
		"/api/v1/healthz": newHealthzHandler(log, dotes),
		"/api/v1/quiz":    newQuizHandler(log, dotes),
	}

	// Create a new ServeMux and assign handlers to routes
	mux := http.NewServeMux()
	for path, handler := range routeHandlers {
		mux.Handle(path, handler)
	}

	// Create a server with a middleware-wrapped mux
	handler := newHandler(log, enableRateLimiting, mapKeys(routeHandlers), mux)
	return &http.Server{
		Addr:    host,
		Handler: handler,
	}
}

// RunServer manages the lifecycle of HTTP server that serves the Dota quiz API
// Callers can shutdown the server by canceling the context
func RunServer(ctx context.Context, log *logrus.Entry, enableRateLimiting bool, host string, dotes dota.DotaService) {
	server := newServer(log, enableRateLimiting, host, dotes)

	go func() {
		log.WithField("host", host).Info("starting server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithField("error", err).Fatal("ListenAndServe failed")
		}
	}()

	// Block on context cancellation
	<-ctx.Done()
	log.Info("context canceled. shutting down server")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.WithField("error", err).Fatal("server shutdown failed", err)
	}
	log.Info("server exited")
}

// Returns slice of keys from a map of http.Handlers
func mapKeys(m map[string]http.Handler) []string {
	keys := []string{}
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}

// Retrieve logger from context
func getLogger(ctx context.Context) *logrus.Entry {
	logger, ok := ctx.Value("log").(*logrus.Entry)
	if !ok {
		// There's absolutely no reason this should ever happen. Fail hard
		// if it does
		panic("failed to get logger from context")
	}
	return logger
}

// Add logger to context and pass it all the way through the middleware chain
func newOuterMiddleware(log *logrus.Entry, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newLogger := log.WithFields(logrus.Fields{
			"reqId":           uuid.New().String(),
			"method":          r.Method,
			"x-forwarded-for": r.Header.Get("X-Forwarded-For"),
			"remoteAddr":      r.RemoteAddr,
			"userAgent":       r.UserAgent(),
			"referer":         r.Referer(),
			"requestUri":      r.URL.RequestURI(),
			"proto":           r.Proto,
			"contentType":     r.Header.Get("Content-Type"),
		})
		newLogger.Info("request received")
		ctx := context.WithValue(r.Context(), "log", newLogger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
