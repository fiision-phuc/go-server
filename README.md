GO-SERVER
=========
go-server is a simple and lightweight server written in Go and the implementation is best suit for micro services.
------------------------------------------------------------------------------------------------------------------

### Getting Started

After installing Go and setting up your [GOPATH](http://golang.org/doc/code.html#GOPATH), install the go-server package (**go 1.5** or greater is required):
~~~
go get github.com/phuc0302/go-server
~~~

Then create your first `.go` file. Let call it `server.go`
~~~ go
package main

import (
	"github.com/phuc0302/go-server"
	"github.com/phuc0302/go-server/util"
)

func main() {
    // 1. Initialize server's environment. Either in sandbox mode or production mode
	server.Initialize(true)

	// 2. Bind a handler to HTTP 1.1 GET /
	server.BindGet("/", func(c *server.RequestContext) {
		c.OutputText(util.Status200(), "Hello world!")
	})
	
	// 3. Start HTTP server
	server.Run()
}
~~~

Then run your server:
~~~
go run server.go
~~~

You will now have a go-server running on `http://localhost:8080`.

### Table of Contents
* [go-server](#go-server)
  * [Handler](#handler)
  * [Routing](#routing)
  * [Request Context](#request-context)

### go-server
to initialize server's environment, either sandbox or production, the server allows us to define 2 different configuration files. Depend on sandboxMode is true or false, the server will load `server.debug.cfg` or `server.release.cfg`.
~~~ go
  sandboxMode := true
  server.Initialize(sandboxMode)
~~~

#### Handler
There are 2 types of handlers:
- **HandleGroupFunc:** _a type alias for group func callback handler._
- **HandleContextFunc:** a type alias for request context func callback handler.

This func will append `/api/v1` before `/sample`. Thus, the full path will be: `/api/v1/sample`.
~~~ go
// HandleGroupFunc
server.GroupRoute("/api/v1", func() {
	server.BindGet("/sample", func(c *server.RequestContext) {
		c.OutputText(util.Status200(), "Hello World!")
	})
})
~~~

~~~ go
// HandleContextFunc
server.BindGet("/sample", func(c *server.RequestContext) {
	c.OutputText(util.Status200(), "Hello World!")
})
~~~

Incase if you want to intercept your HandleContextFunc, we can do this:
~~~ go
handler := func(c *server.RequestContext) {
	c.OutputText(util.Status200(), "Hello world!")
}

adapter := func() {
	return func(f server.HandleContextFunc) server.HandleContextFunc {
		return func(c *server.RequestContext) {
			defer fmt.Println("After...")
			fmt.Println("Before...")
			f(c)
		}
	}
}
// 2. Bind a handler to HTTP 1.1 GET /
server.BindGet("/", server.Adapt(handler, adapter))
~~~

#### Routing
In go-server, a route is a node. Each node contains a single URL-matching pattern and one or more paired `HTTP method - HandleContextFunc`. Routes are matched in the order they are defined, the first route that matches the request is invoked.
~~~ go
server.BindGet("/", func(c *server.RequestContext) {
    // Read
})

server.BindPatch("/", func(c *server.RequestContext) {
    // Update
})

server.BindPost("/", func(c *server.RequestContext) {
    // Create
})

server.BindPut("/", func(c *server.RequestContext) {
    // Replace
})

server.BindDelete("/", func(c *server.RequestContext) {
    // Delete
})
~~~

Route patterns may include named parameters.
~~~ go
server.BindGet("/user/{userName}", func(c *server.RequestContext) {
    c.OutputText(util.Status200(), fmt.Sprintf("Hello %s!", c.PathParams["userName"]))
})
~~~

Route groups can be added too using the `GroupRoute` func.
~~~ go
server.GroupRoute("/items", func() {
    server.BindGet("", GetItems)
    server.BindPost("", NewItem)
    server.BindGet("/{itemID}", GetItem)
    server.BindPut("/{itemID}", UpdateItem)
    server.BindDelete("/{itemID}", DeleteItem)
})
~~~

#### Request Context
Request Context represent request scope when server received a request from client. The context will be created by server and send to handler.