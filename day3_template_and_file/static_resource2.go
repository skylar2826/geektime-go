package __template_and_file

import (
	lru "github.com/hashicorp/golang-lru"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

type StaticResourceHandler2 struct {
	dir               string
	cache             *lru.Cache
	maxSize           int
	extContentTypeMap map[string]string
}

type StaticResourceHandlerOption2 func(*StaticResourceHandler2)

func StaticWithFileCache(maxSize int, cnt int) StaticResourceHandlerOption2 {
	return func(handler *StaticResourceHandler2) {
		cache, err := lru.New(cnt)
		if err != nil {
			log.Println("文件缓存创建失败：", err)
		}
		handler.maxSize = maxSize
		handler.cache = cache
	}
}

func StaticWithExtContentType(extContentTypeMap map[string]string) StaticResourceHandlerOption2 {
	return func(handler *StaticResourceHandler2) {
		for ext, contentType := range extContentTypeMap {
			handler.extContentTypeMap[ext] = contentType
		}
	}
}

func NewStaticResourceHandler2(dir string, opts ...StaticResourceHandlerOption2) *StaticResourceHandler2 {
	srh := &StaticResourceHandler2{dir: dir, extContentTypeMap: map[string]string{
		"jpeg": "image/jpeg",
		"jpe":  "image/jpeg",
		"jpg":  "image/jpeg",
		"png":  "image/png",
		"pdf":  "image/pdf",
	}}
	for _, opt := range opts {
		opt(srh)
	}
	return srh
}

type fileCacheItem2 struct {
	fileName    string
	fileSize    int
	data        []byte
	contentType string
}

func (s *StaticResourceHandler2) readFileFormData(fileName string) (*fileCacheItem2, bool) {
	if s.cache != nil {
		if item, ok := s.cache.Get(fileName); ok {
			log.Printf("静态资源访问缓存")
			return item.(*fileCacheItem2), true
		}
	}
	return nil, false
}

func (s *StaticResourceHandler2) writeFileFormData(item *fileCacheItem2) {
	if s.cache != nil && item.fileSize <= s.maxSize {
		s.cache.Add(item.fileName, item)
	}
}

func (s *StaticResourceHandler2) writeResp(item *fileCacheItem2, c *Context) {
	header := c.W.Header()
	header.Set("Content-Type", item.contentType)
	header.Set("Content-Length", strconv.Itoa(item.fileSize))
	c.RespStatusCode = http.StatusOK
	c.RespData = item.data
}

func (s *StaticResourceHandler2) Handle(c *Context) {
	/*
		1. 获取完整路径，读文件；
			1. 文件路径校验，避免越权访问
			2. 设置缓存
				优先读取缓存，设置maxSize
		2. 计算值：content-type\content-length
		3. 返回文件
	*/
	fileName, err := c.PathValue("file")
	if err != nil {
		c.BadRequest(err)
		return
	}
	item, ok := s.readFileFormData(fileName)
	if ok {
		s.writeResp(item, c)
		return
	}

	dst := path.Join(s.dir, fileName)
	data, err := os.ReadFile(dst)
	if err != nil {
		c.SystemError(err)
		return
	}

	ext := path.Ext(dst)

	item = &fileCacheItem2{
		fileName:    fileName,
		fileSize:    len(data),
		data:        data,
		contentType: s.extContentTypeMap[ext[1:]],
	}
	s.writeFileFormData(item)
	s.writeResp(item, c)
}
