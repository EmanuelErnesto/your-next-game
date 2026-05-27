package httpserver

import (
	"log/slog"
	"time"
	"your-next-game-backend/internal/platform/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoggerMiddleware intercepta todas as requisições para injetar Correlation ID e emitir a Canonical Log Line no final.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 1. Extrair ou gerar Correlation ID
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = "req-" + uuid.New().String()[:8]
		}
		c.Header("X-Correlation-ID", correlationID)

		// 2. Colocar Correlation ID e Accumulator no context.Context do request
		ctx := logger.WithCorrelationID(c.Request.Context(), correlationID)
		ctx, accumulator := logger.WithAccumulator(ctx)
		c.Request = c.Request.WithContext(ctx)

		// 3. Processar a requisição
		c.Next()

		// 4. Extrair dados pós-execução e emitir a Canonical Log Line
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Reunir contexto acumulado durante o processamento do request
		metadata := accumulator.Map()

		// Construir os atributos padrões da requisição
		attrs := []any{
			slog.String("correlation_id", correlationID),
			slog.String("http.method", c.Request.Method),
			slog.String("http.path", c.Request.URL.Path),
			slog.Int("http.status", statusCode),
			slog.Duration("http.latency_ms", latency/time.Millisecond),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
		}

		// Adicionar metadados acumulados
		for k, v := range metadata {
			attrs = append(attrs, slog.Any(k, v))
		}

		if errorMessage != "" {
			attrs = append(attrs, slog.String("error", errorMessage))
		}

		// Emitir log no nível apropriado
		logMessage := "Canonical Request Log"
		if statusCode >= 500 {
			slog.ErrorContext(ctx, logMessage, attrs...)
		} else if statusCode >= 400 {
			slog.WarnContext(ctx, logMessage, attrs...)
		} else {
			slog.InfoContext(ctx, logMessage, attrs...)
		}
	}
}
