// internal/pkg/tracing/iran_tracer.go
package tracing

type IranTracerConfig struct {
    SampleRate    float64
    LocalStorage  bool          // Store traces locally if external unavailable
    Compression   bool          // Compress spans for slow networks
    BatchSize     int           // Smaller batches for unstable networks
}

func NewIranTracer(serviceName string) (trace.TracerProvider, error) {
    config := IranTracerConfig{
        SampleRate:  0.1, // Lower sample rate for Iran
        LocalStorage: true,
        Compression: true,
        BatchSize:   50,
    }
    
    exporter, err := newResilientExporter(config)
    if err != nil {
        // Fallback to local storage
        return newLocalTracer(serviceName), nil
    }
    
    return trace.NewTracerProvider(
        trace.WithBatcher(exporter,
            trace.WithBatchTimeout(5*time.Second),
            trace.WithMaxExportBatchSize(config.BatchSize),
        ),
        trace.WithResource(resource.NewWithAttributes(
            semantic.SchemaURL,
            semantic.ServiceNameKey.String(serviceName),
            semantic.ServiceVersionKey.String("1.0.0"),
            semantic.DeploymentEnvironmentKey.String("iran-production"),
        )),
    ), nil
}
