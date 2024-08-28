package __template_and_file

import (
	"fmt"
)

func Home(c *Context) {
	fmt.Fprintf(c.W, "Welcome to the home page! pathParams: %v", c.PathParams)

	//c, span := tracer.Start(ctx.R.Context(), "第一个")
	//defer span.End()
	//
	//c, second := tracer.Start(c, "第二个")
	//time.Sleep(100 * time.Microsecond)
	//
	//c, third1 := tracer.Start(c, "第三个的第一个")
	//time.Sleep(100 * time.Microsecond)
	//third1.End()
	//c, third2 := tracer.Start(c, "第三个的第二个")
	//time.Sleep(100 * time.Microsecond)
	//third2.End()
	//second.End()

}

func User(c *Context) {
	fmt.Fprintf(c.W, "Welcome to the user page! userId: %d", c.PathParams)
}
