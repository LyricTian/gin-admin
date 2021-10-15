package middleware

import (
	"log"
	"os"

	"github.com/LyricTian/gin-admin/v8/internal/app/config"

	// "github.com/LyricTian/gin-admin/v8/pkg/util/trace"
	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	// "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	environmentKey = "environment"
	IDKey          = "ID"
)

// TODO: Add other type of exporter
func newExporter() (sdktrace.SpanExporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.C.Trace.Endpoint)))
}

func initTracer() *sdktrace.TracerProvider {

	exporter, err := newExporter()

	if err != nil {
		log.Fatal(err)
	}
	pid := int64(os.Getpid())
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(config.C.Trace.Sample)),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.C.App.Name),
			attribute.String(environmentKey, config.C.App.Environment),
			attribute.Int64(IDKey, pid),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}
func TraceMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	initTracer()
	handler := otelgin.Middleware(config.C.App.Name)
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}
		handler(c)
	}
}
