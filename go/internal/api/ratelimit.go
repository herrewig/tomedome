package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/ulule/limiter/v3"
	mhttp "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// For rate limiting, we need to know the real client IP address.
// We use a GCP L7 load balancer. GCP docs explain the
// second-to-last IP in the X-Forwarded-For header is the real
// client IP. The limiter package doesn't know how to do this
// natively so we have to write a custom middleware for it
func newClientIpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := getLogger(r.Context())
		header := r.Header.Get("X-Forwarded-For")

		if header == "" {
			// Just ignore calls to /healthz because they usually
			// come from the GCP L7 LB. They'll still get rate
			// limited via request.RemoteAddr
			if r.URL.Path != "/api/v1/healthz" {
				log.Warn("X-Forwarded-For header not found")
			}
		} else {
			ips := strings.Split(header, ",")
			if len(ips) < 2 {
				log.Warn("X-Forwarded-For header has less than 2 IPs")
			} else {
				// Use the second-to-last IP
				r.RemoteAddr = strings.TrimSpace(ips[len(ips)-2])
			}
		}
		next.ServeHTTP(w, r)
	})
}

func newLimiterMiddleware(next http.Handler) http.Handler {
	// Limit to 10 requests per IP per minute
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  10,
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	middleware := mhttp.NewMiddleware(instance)
	return middleware.Handler(next)
}
