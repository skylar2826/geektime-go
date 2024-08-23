package __shutdown

import "fmt"

func Home(c *Context) {
	fmt.Fprintf(c.W, "Welcome to the home page!")
}

func User(c *Context) {
	fmt.Fprintf(c.W, "Welcome to the user page! userId: %d", c.PathParams)
}
