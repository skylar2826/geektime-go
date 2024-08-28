package __template_and_file

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func init() {
	RegisterFilterBuilder("handleLog", HandleLog)
	RegisterFilterBuilder("handleAccess", HandleAccess)
}

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

type FileUploader struct {
	fileField string
	// 用户传路径的原因：考虑文件重名的问题； 不重名用uuid
	//uuid.New().String()
	DstPathFunc func(header *multipart.FileHeader) string
}

func (f *FileUploader) Handle() HandlerFunc {
	return func(c *Context) {
		/*
			1. 读文件内容
			2. 计算出目标存储路径
			3. 保存文件
			4. 返回响应
		*/

		file, fileHeader, err := c.R.FormFile(f.fileField)
		if err != nil {
			c.BadRequest(err)
			return
		}

		dst := f.DstPathFunc(fileHeader)

		/*
			可以尝试把沿途的dir都建好，不然OpenFile找不到文件夹会报错
		*/
		//os.MkdirAll()

		// 如果已经存在就给他O_TRUNC 清空掉，不存在就O_CREATE 创建
		dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o666)
		if err != nil {
			c.SystemError(err)
			return
		}

		// buffer 指每次传多大块的数据，复用，影响性能
		_, err = io.CopyBuffer(dstFile, file, nil)
		if err != nil {
			c.SystemError(err)
			return
		}
		c.RequestOk("上传成功")
	}
}

/*
	New(options) + 直接方法 // options可以动态配置载入初始化哪些需要的能力

	func (h *XXX) handle() HanlderFunc {
	    // 可以在这里初始化设置默认值
		return func(c *Context) {}
	}
*/

type FileDownloader struct {
	Dir string
}

func (f *FileDownloader) Handle() HandlerFunc {
	return func(c *Context) {
		val, err := c.QueryValue("file")
		if err != nil {
			c.BadRequest(err)
			return
		}
		val = filepath.Clean(val)
		dst := filepath.Join(f.Dir, val)

		//// 做校验，防止攻击者使用相对路径越权访问服务 比如file="../../1.sh"
		//dst, err := filepath.Abs(dst)
		//if strings.Contains(dst, f.Dir) {
		//	// todo
		//}

		fn := filepath.Base(dst)
		header := c.W.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")

		http.ServeFile(c.W, c.R, dst)
	}
}
