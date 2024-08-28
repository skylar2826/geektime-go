package __template_and_file

import (
	"log"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
	"text/template"
)

func RunTemplateAndFile() {
	tpl, err := template.ParseGlob("testdata/tpls/*.gohtml")
	if err != nil {
		panic("模板解析失败")
	}
	engine := &GoTemplateEngine{
		T: tpl,
	}

	server := NewServer("test", ServerWithTemplateEngine(engine))
	//server := NewServer("router1", HandleLog, HandleAccess)
	server.Route(http.MethodGet, "/home", func(c *Context) {
		err := c.Render("home.gohtml", nil)
		if err != nil {
			log.Println("渲染home.gohtml返回前端出错")
		}
	})

	server.Route(http.MethodGet, "/upload", func(c *Context) {
		err := c.Render("upload.gohtml", nil)
		if err != nil {
			log.Println("渲染upload.gohtml返回到前端出错")
		}
	})

	fileUploader := &FileUploader{
		fileField: "myfile",
		DstPathFunc: func(header *multipart.FileHeader) string {
			// 传的是相对路径
			// filepath能解决不同操作系统间差异
			return filepath.Join("testdata", "file", header.Filename)
		},
	}
	server.Route(http.MethodPost, "/upload", fileUploader.Handle())

	fileDownloader := &FileDownloader{
		path.Join("testdata", "file"),
	}
	server.Route(http.MethodGet, "/download", fileDownloader.Handle())

	srh := NewStaticResourceHandler2(path.Join("testdata", "file"), StaticWithExtContentType(map[string]string{
		"gif": "image/gif",
	}), StaticWithFileCache(20*1024*1024, 100))
	//if err != nil {
	//	panic(fmt.Sprint("静态资源服务创建失败: ", err))
	//}
	server.Route(http.MethodGet, "/static/:file", srh.Handle)

	server.Start("127.0.0.1:8000")
}
