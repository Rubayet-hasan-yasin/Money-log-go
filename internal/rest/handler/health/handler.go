package healthhandler

import (
	"database/sql"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	db        *sql.DB
	log       *zap.Logger
	startTime time.Time
}

func NewHandler(db *sql.DB, log *zap.Logger) *Handler {
	return &Handler{
		db:        db,
		log:       log,
		startTime: time.Now(),
	}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"message":   "MoneyLog API is running",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    time.Since(h.startTime).Seconds(),
	})
}

func (h *Handler) DetailedHealth(c *gin.Context) {
	dbStatus := "connected"
	dbMessage := "Database is operational"
	statusCode := http.StatusOK

	if err := h.db.PingContext(c.Request.Context()); err != nil {
		dbStatus = "disconnected"
		dbMessage = "Database connection failed: " + err.Error()
		statusCode = http.StatusServiceUnavailable
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	c.JSON(statusCode, gin.H{
		"status": func() string {
			if statusCode == http.StatusOK {
				return "ok"
			}
			return "degraded"
		}(),
		"message":   "MoneyLog API Health Check",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    time.Since(h.startTime).Seconds(),
		"services": gin.H{
			"database": gin.H{
				"status":  dbStatus,
				"message": dbMessage,
			},
			"api": gin.H{
				"status":  "ok",
				"message": "API is operational",
			},
		},
		"system": gin.H{
			"go_version": runtime.Version(),
			"goroutines": runtime.NumGoroutine(),
			"memory": gin.H{
				"alloc_mb":       float64(memStats.Alloc) / 1024 / 1024,
				"total_alloc_mb": float64(memStats.TotalAlloc) / 1024 / 1024,
				"sys_mb":         float64(memStats.Sys) / 1024 / 1024,
			},
		},
	})
}
