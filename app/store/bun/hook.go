package bun

import (
	"context"
	"github.com/uptrace/bun"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type OtelQueryHook struct {
	Tracer trace.Tracer
	DBName string
}

func NewOtelQueryHook(tracerName, dbName string) *OtelQueryHook {
	return &OtelQueryHook{
		Tracer: otel.Tracer(tracerName),
		DBName: dbName,
	}
}

func (h *OtelQueryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	spanName := "SQL: " + event.Query
	ctx, span := h.Tracer.Start(ctx, spanName, trace.WithAttributes(
		attribute.String("db.system", "mysql"),
		attribute.String("db.statement", event.Query),
		attribute.String("db.name", h.DBName),
	))
	return trace.ContextWithSpan(ctx, span)
}

func (h *OtelQueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	span := trace.SpanFromContext(ctx)
	if event.Err != nil {
		span.RecordError(event.Err)
	}
	span.End()
}
