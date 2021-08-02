package main

import (
	services "github.com/awesome-demo-app/todolist-api/core/services"
	handlers "github.com/awesome-demo-app/todolist-api/handlers"
	repositories "github.com/awesome-demo-app/todolist-api/repositories"

	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func bootstrapTracingWithOpenTelemetry() (*sdktrace.TracerProvider, context.Context) {
	// We only export traces to stdout for now
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		log.Fatalf("failed to initialize stdouttrace export pipeline: %v", err)
	}

	ctx := context.Background()
	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(traceExporter)
	traceProvider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(batchSpanProcessor))

	// Set the trace provider as global so we don't have to pass it everywhere
	otel.SetTracerProvider(traceProvider)

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.Baggage{},
		propagation.TraceContext{},
	)
	// ??
	otel.SetTextMapPropagator(propagator)

	return traceProvider, ctx
}

func main() {
	// Tracing
	traceProvider, ctx := bootstrapTracingWithOpenTelemetry()
	defer func() { traceProvider.Shutdown(ctx) }()

	// Hexagonal Architecture components
	todoRepository := repositories.NewSQLiteDB(
		"file::memory:?cache=shared&mode=rwc&_journal_mode=WAL",
	)
	todoService := services.New(todoRepository)
	todoHandler := handlers.NewHTTPHandler(todoService)

	// Main loop
	todoHandler.HandleRequests()
}
