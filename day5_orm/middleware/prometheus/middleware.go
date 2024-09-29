package prometheus

import (
	"context"
	"geektime-go/day5_orm/my_orm_mysql"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type MiddlewareBuilder struct {
	NameSpace string
	Subsystem string
	Name      string
	Help      string
}

func (m *MiddlewareBuilder) Build() my_orm_mysql.Middleware {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      m.Name,
		Subsystem: m.Subsystem,
		Namespace: m.NameSpace,
		Help:      m.Help,
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.90:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, []string{"type", "table"})

	prometheus.MustRegister(vector)

	return func(next my_orm_mysql.Handler) my_orm_mysql.Handler {
		return func(ctx context.Context, queryCtx *my_orm_mysql.QueryContext) *my_orm_mysql.QueryResult {
			startTime := time.Now()
			defer func() {
				// 执行时间
				vector.WithLabelValues(queryCtx.Type, queryCtx.Model.TableName).Observe(float64(time.Since(startTime).Milliseconds()))
			}()
			return next(ctx, queryCtx)
		}
	}
}
