package __shutdown

import "fmt"

func HandleLog(next Filter) Filter {
	return func(c *Context) {
		fmt.Println("打印日志1")
		next(c)
	}
}

func HandleAccess(next Filter) Filter {
	return func(c *Context) {
		fmt.Println("处理跨域2")
		next(c)
	}
}

func init() {
	RegisterFilterBuilder("handleLog", HandleLog)
	RegisterFilterBuilder("handleAccess", HandleAccess)
}
