package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"os"
	metertracer "otel-go/pkg/metrics/exporter"
	exptracer "otel-go/pkg/tracer/exporter"
)

var (
	serviceName = os.Getenv("SERVICE_NAME")
)

func main() {
	cleanup := exptracer.InitTracer()
	defer cleanup(context.Background())

	provider := metertracer.InitMeter()
	defer func(provider *metricsdk.MeterProvider, ctx context.Context) {
		err := provider.Shutdown(ctx)
		if err != nil {
			panic(err)
		}
	}(provider, context.Background())

	meter := provider.Meter(serviceName)
	metertracer.GenerateMetrics(meter)

	r := gin.Default()
	r.Use(otelgin.Middleware(serviceName))
	r.GET("/", getHandler)
	r.GET("/error", errorHandler)

	if err := r.Run("8085"); err != nil {
		panic(err)
	}
}

func getService(c *gin.Context, payload string) error {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("controller", "books"))
	span.AddEvent("Get Service", trace.WithAttributes(attribute.Int("pid", 4328), attribute.String("handler.method", "GET")))

	return nil
}

func getHandler(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("controller", "books"))
	span.AddEvent("Get Handler", trace.WithAttributes(attribute.Int("pid", 4328), attribute.String("handler.method", "GET")))

	if err := getService(c, "payload"); err != nil {

		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "error hit get handler",
		})
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"message": "success hit get handler",
	})
}

func errorHandler(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("controller", "books"))
	span.AddEvent("Error Handler", trace.WithAttributes(attribute.Int("pid", 4328), attribute.String("handler.method", "GET")))

	c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"status":  "error",
		"message": "error hit error handler",
	})
}
