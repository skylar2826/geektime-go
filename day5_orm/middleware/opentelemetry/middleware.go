package opentelemetry

import (
	"context"
	"fmt"
	"geektime-go/day5_orm/my_orm_mysql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var instrumentationName = "geektime-go/day5_orm/middleware/opentelemetry"

type MiddlewareBuilder struct {
	tracer trace.Tracer
}

func (m MiddlewareBuilder) Build() my_orm_mysql.Middleware {
	if m.tracer == nil {
		m.tracer = otel.GetTracerProvider().Tracer(instrumentationName)
	}
	return func(next my_orm_mysql.Handler) my_orm_mysql.Handler {
		return func(ctx context.Context, queryCtx *my_orm_mysql.QueryContext) *my_orm_mysql.QueryResult {
			// spanName = "select-test_model"
			spanName := fmt.Sprintf("%s-%s", queryCtx.Type, queryCtx.Model.TableName)
			spanCtx, span := m.tracer.Start(ctx, spanName)
			defer span.End()

			q, _ := queryCtx.Builder.Build()
			if q != nil {
				span.SetAttributes(attribute.String("sql", q.SQL))
			}
			span.SetAttributes(attribute.String("table", queryCtx.Model.TableName))
			span.SetAttributes(attribute.String("component", "orm"))
			res := next(spanCtx, queryCtx)
			if res.Err != nil {
				span.RecordError(res.Err)
			}
			return res
		}
	}
}
