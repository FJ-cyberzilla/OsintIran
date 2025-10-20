// pkg/observability/tracing.go
func NewTracer(serviceName string) (trace.TracerProvider, error) {
    return jaeger.New(jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint(os.Getenv("JAEGER_ENDPOINT")),
    ))
}
