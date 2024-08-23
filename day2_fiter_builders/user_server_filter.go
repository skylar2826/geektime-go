package __filter_builder

import (
	"encoding/json"
	"fmt"
)

type accessLog struct {
	Host       string `json:"host,omitempty"`
	Route      string `json:"route,omitempty"`
	HTTPMethod string `json:"http_method,omitempty"`
	Path       string `json:"path,omitempty"`
}

func HandleLog(next Filter) Filter {
	return func(c *Context) {
		defer func() {
			l := accessLog{
				Host:       c.R.URL.Host,
				Route:      c.MatchRoute, // 匹配到的完整路由
				HTTPMethod: c.R.Method,
				Path:       c.R.URL.Path,
			}
			//log.Fatalln(l)
			data, err := json.Marshal(l)
			if err != nil {
				panic("marshal access log error")
			}
			fmt.Printf("accesslog: %v\n", string(data))
		}()

		next(c)
	}
}

func HandleAccess(next Filter) Filter {
	return func(c *Context) {
		fmt.Println("处理跨域...")
		next(c)
	}
}

//
//var instrumentationName = "test-tracer"
//var tracer = otel.GetTracerProvider().Tracer(instrumentationName)
//
//func HandleTracer(next Filter) Filter {
//	return func(c *Context) {
//		reqCtx := c.R.Context()
//		reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(c.R.Header))
//
//		reqCtx, span := tracer.Start(reqCtx, "unknown") // spanName一般为router
//
//		span.SetAttributes(attribute.String("http.method", c.R.Method))
//		span.SetAttributes(attribute.String("http.path", c.R.URL.Path))
//		span.SetAttributes(attribute.String("http.host", c.R.URL.Host))
//		span.SetAttributes(attribute.String("http.path", c.R.URL.Path))
//		span.SetAttributes(attribute.String("http.scheme", c.R.URL.Scheme))
//
//		c.R = c.R.WithContext(reqCtx)
//
//		next(c)
//
//		// 防止next中发生panic，所以放在defer中
//		defer func() {
//			span.SetName(c.MatchRoute)
//			span.SetAttributes(attribute.Int("http.status", c.RespStatusCode))
//			span.End()
//		}()
//	}
//}
//
//func HandlePrometheus(next Filter) Filter {
//	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
//		Name:      "test",
//		Subsystem: "test2",
//		Namespace: "test3",
//		Objectives: map[float64]float64{
//			0.5:   0.01,
//			0.75:  0.01,
//			0.90:  0.01,
//			0.99:  0.001,
//			0.999: 0.0001,
//		},
//	}, []string{"pattern", "method", "status"})
//	return func(c *Context) {
//		startTime := time.Now()
//		defer func() {
//			duration := time.Now().Sub(startTime).Milliseconds()
//			pattern := c.MatchRoute
//			if pattern == "" {
//
//				pattern = "unknown"
//			}
//			vector.WithLabelValues(pattern, c.R.Method, strconv.Itoa(c.RespStatusCode)).Observe(float64(duration))
//		}()
//		next(c)
//	}
//}

func init() {
	RegisterFilterBuilder("handleLog", HandleLog)
	RegisterFilterBuilder("handleAccess", HandleAccess)
}
