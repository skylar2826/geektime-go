package __template_and_file

import (
	"fmt"
	"strings"
)

func HandleUser(next Filter) Filter {
	return func(c *Context) {
		fmt.Printf("命中router-filter %v, to do something\n", c.R.URL.Path)
		next(c)
	}
}

func HandlerStaticPng(next Filter) Filter {
	return func(c *Context) {
		if strings.Contains(c.R.URL.Path, "1.png") {
			fmt.Printf("命中router-filter %v, to do something\n", c.R.URL.Path)
		}
		next(c)
	}
}

func HandlerX(next Filter) Filter {
	return func(c *Context) {
		fmt.Printf("this is x...\n")
		next(c)
	}
}

func HandlerA(next Filter) Filter {
	return func(c *Context) {
		fmt.Printf("this is a...\n")
		next(c)
	}
}

func HandlerAB(next Filter) Filter {
	return func(c *Context) {
		fmt.Printf("this is ab...\n")
		next(c)
	}
}
func HandlerAX(next Filter) Filter {
	return func(c *Context) {
		fmt.Printf("this is ax...\n")
		next(c)
	}
}
func HandlerAD(next Filter) Filter {
	return func(c *Context) {
		fmt.Printf("this is ad...\n")
		next(c)
	}
}
func HandlerABD(next Filter) Filter {
	return func(c *Context) {
		fmt.Printf("this is abd...\n")
		next(c)
	}
}
func HandlerABDE(next Filter) Filter {
	return func(c *Context) {
		fmt.Printf("this is abde...\n")
		next(c)
	}
}

func HandlerABDX(next Filter) Filter {
	return func(c *Context) {
		fmt.Printf("this is abdx...\n")
		next(c)
	}
}

func init() {
	RegisterFilterBuilder("/user", HandleUser)
	RegisterFilterBuilder("/1.png", HandlerStaticPng)
}
