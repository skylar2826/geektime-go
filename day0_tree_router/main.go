package __tree_router

import (
	"fmt"

	"net/http"
)

func home(c *Context) {
	fmt.Fprintf(c.W, "Welcome to the home page!")
}

func user(c *Context) {
	fmt.Fprintf(c.W, "Welcome to the user page! userId: %d", c.PathParams)
}

type MyError struct {
}

func (m *MyError) Error() string {
	return "失败"
}

func RunTreeRouter() {
	var myError error = &MyError{}
	server := NewServer("test-server", HandleLog, HandleAccess)
	err := server.Route(http.MethodPost, "/signUp", SignUp)
	if err != nil {
		fmt.Println(myError.Error())
	}
	err = server.Route(http.MethodGet, "/home", home)
	if err != nil {
		fmt.Println(myError.Error())
	}
	err = server.Route(http.MethodGet, "/user/:id", user)
	if err != nil {
		fmt.Println(myError.Error())
	}
	err = server.Start("localhost:8080")
	if err != nil {
		fmt.Println(myError.Error())
	}
}
