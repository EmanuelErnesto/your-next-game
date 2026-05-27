package logger

import (
	"context"
	"log"
	"log/slog"
	"os"
	"sync"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type contextKey string

const (
	correlationIDKey contextKey = "correlation_id"
	loggerKey        contextKey = "logger"
	accumulatorKey   contextKey = "log_accumulator"
)

// Accumulator armazena de forma thread-safe metadados adicionados durante o request.
type Accumulator struct {
	mu   sync.RWMutex
	data map[string]any
}

func NewAccumulator() *Accumulator {
	return &Accumulator{
		data: make(map[string]any),
	}
}

func (a *Accumulator) Add(key string, value any) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.data[key] = value
}

func (a *Accumulator) AddAll(attrs map[string]any) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for k, v := range attrs {
		a.data[k] = v
	}
}

func (a *Accumulator) Map() map[string]any {
	a.mu.RLock()
	defer a.mu.RUnlock()
	// Retorna uma cópia
	copyMap := make(map[string]any, len(a.data))
	for k, v := range a.data {
		copyMap[k] = v
	}
	return copyMap
}

// InitLogger configura o structured logger global integrado com o OpenTelemetry Collector.
func InitLogger(env string) {
	ctx := context.Background()

	// 1. Criar exportador OTLP HTTP
	exporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint("localhost:4318"),
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		log.Printf("Failed to create OTel log exporter: %v. Falling back to stdout.", err)
		setupFallbackLogger(env)
		return
	}

	// 2. Provedor de logs OTel SDK
	res, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("your-next-game-backend"),
			semconv.DeploymentEnvironmentKey.String(env),
		),
	)

	processor := sdklog.NewBatchProcessor(exporter)
	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(processor),
		sdklog.WithResource(res),
	)

	// Registrar globalmente
	global.SetLoggerProvider(loggerProvider)

	// 3. Ponte slog para OTel
	otelLogger := otelslog.NewLogger("your-next-game-backend")
	slog.SetDefault(otelLogger)
}

func setupFallbackLogger(env string) {
	var handler slog.Handler
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}
	slog.SetDefault(slog.New(handler))
}

// Conexão de Contexto com Correlation ID
func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationIDKey, correlationID)
}

func GetCorrelationID(ctx context.Context) string {
	if val, ok := ctx.Value(correlationIDKey).(string); ok {
		return val
	}
	return ""
}

// Conexão de Contexto com Accumulator
func WithAccumulator(ctx context.Context) (context.Context, *Accumulator) {
	acc := NewAccumulator()
	return context.WithValue(ctx, accumulatorKey, acc), acc
}

func GetAccumulator(ctx context.Context) *Accumulator {
	if acc, ok := ctx.Value(accumulatorKey).(*Accumulator); ok {
		return acc
	}
	return nil
}

// AddContext adiciona um par chave-valor ao acumulador de logs no contexto atual.
func AddContext(ctx context.Context, key string, value any) {
	if acc := GetAccumulator(ctx); acc != nil {
		acc.Add(key, value)
	}
}

// ContextLogger retorna um logger já configurado com correlation_id do contexto
func ContextLogger(ctx context.Context) *slog.Logger {
	id := GetCorrelationID(ctx)
	if id != "" {
		return slog.Default().With(slog.String("correlation_id", id))
	}
	return slog.Default()
}
