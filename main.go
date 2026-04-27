package main

import (
    "log"
    "note-app/controllers"
    "note-app/models"
    "note-app/repository"
    "os"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/sirupsen/logrus"       // add this
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "note_app_requests_total",
            Help: "Total number of HTTP requests by method, path, and status code",
        },
        []string{"method", "path", "status"},
    )
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "note_app_request_duration_seconds",
            Help:    "Histogram of response latencies for HTTP requests",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )
    logger = logrus.New()   // add this
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)

    // Setup logrus
    logger.SetFormatter(&logrus.JSONFormatter{})
    logger.SetOutput(os.Stdout)

    // Setup Loki hook
    lokiURL := os.Getenv("LOKI_URL")
    if lokiURL == "" {
        lokiURL = "http://loki:3100/loki/api/v1/push"
    }
    logger.AddHook(NewLokiHook(lokiURL, map[string]string{
        "app": "note-app",
        "env": os.Getenv("ENV"),
    }))
}

func prometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.FullPath()
        if path == "" {
            path = "unknown"
        }
        c.Next()
        status := strconv.Itoa(c.Writer.Status())
        duration := time.Since(start).Seconds()

        httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
        httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)

        // Log each request to Loki via logrus
        logger.WithFields(logrus.Fields{
            "method":   c.Request.Method,
            "path":     path,
            "status":   status,
            "duration": duration,
        }).Info("request handled")
    }
}

func main() {
    dsn := os.Getenv("DB_DSN")
    var db *gorm.DB
    var err error
    for i := 0; i < 5; i++ {
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {
            break
        }
        logger.Warnf("Retrying DB connection... (%d/5): %v", i+1, err)
        time.Sleep(2 * time.Second)
    }
    if err != nil {
        logger.Fatalf("Failed to connect to database: %v", err)
    }

    db.AutoMigrate(&models.Note{})

    noteRepo := &repository.NoteRepository{DB: db}
    noteCtrl := &controllers.NoteController{Repo: noteRepo}

    r := gin.Default()
    r.LoadHTMLGlob("templates/*")
    r.Use(prometheusMiddleware())
    r.GET("/metrics", gin.WrapH(promhttp.Handler()))
    r.GET("/", noteCtrl.Index)
    r.POST("/notes", noteCtrl.Create)
    r.GET("/notes/edit/:id", noteCtrl.Edit)
    r.POST("/notes/update/:id", noteCtrl.Update)
    r.POST("/notes/delete/:id", noteCtrl.Delete)

    logger.Info("Note App running on :8080")
    log.Println("Metrics available at /metrics")
    r.Run(":8080")
}
