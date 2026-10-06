package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

const limiterTTL = 10 * time.Minute

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type ipRateLimiter struct {
	mu          sync.Mutex
	clients     map[string]*clientLimiter
	limit       rate.Limit
	burst       int
	ttl         time.Duration
	lastCleanup time.Time
}

type rateLimitResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Seconds int    `json:"seconds"`
}

func newIPRateLimiter(limit rate.Limit, burst int, ttl time.Duration) *ipRateLimiter {
	return &ipRateLimiter{
		clients:     make(map[string]*clientLimiter),
		limit:       limit,
		burst:       burst,
		ttl:         ttl,
		lastCleanup: time.Now(),
	}
}

func (l *ipRateLimiter) getLimiter(ip string) *rate.Limiter {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	if now.Sub(l.lastCleanup) > l.ttl {
		for key, client := range l.clients {
			if now.Sub(client.lastSeen) > l.ttl {
				delete(l.clients, key)
			}
		}
		l.lastCleanup = now
	}

	client, exists := l.clients[ip]
	if !exists {
		limiter := rate.NewLimiter(l.limit, l.burst)
		client = &clientLimiter{limiter: limiter, lastSeen: now}
		l.clients[ip] = client
	} else {
		client.lastSeen = now
	}

	return client.limiter
}

func RateLimiter(requests int, window time.Duration) gin.HandlerFunc {
	if requests <= 0 || window <= 0 {
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "rate limiter misconfigured",
				"data":    nil,
			})
		}
	}

	limiter := newIPRateLimiter(rate.Every(window/time.Duration(requests)), requests, limiterTTL)
	return func(c *gin.Context) {
		ip := realip.FromRequest(c.Request)
		if ip == "" {
			ip = c.ClientIP()
		}
		if ip == "" {
			ip = "unknown"
		}

		l := limiter.getLimiter(ip)
		if !l.Allow() {
			r := l.Reserve()
			delay := r.Delay()
			r.Cancel()

			retryAfter := int(delay.Seconds())
			if retryAfter <= 0 {
				retryAfter = 1
			}

			c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "rate limit exceeded",
				"data": &rateLimitResponse{
					Success: false,
					Message: "Too many requests. Please try again later.",
					Seconds: retryAfter,
				},
			})
			return
		}

		c.Next()
	}
}
