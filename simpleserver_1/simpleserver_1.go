package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
)

var port = "3333"
var msg = "Dr. Who say's alonzi"

// So we want to tests to drive our writing. For which we need the ability to
// start and stop the humble http server through a test function

// To do that we encapsulate things we need in a struct. structs are Go's construct
// for a collection of fields. It's also Go's way to define custom types
type server struct {
	httpServer *http.Server
	listener   net.Listener
}

// This method is defined for the type server and receives a pointer value.
// Oh so is (s Server) also valid? Yes it is, but then you would receive a copy
// of the value.
func (s *server) listenAndServe() error {
	// Multi-Value returns! It doesn't return you know, an "OBJECT" or even a struct
	// Doing this in Go is idiomatic(i.e the Go's highway of good design)
	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return err
	}

	//let's store this for posterity
	s.listener = listener

	// goroutine! Yes that's all you gotta do for concurrency. When this function call completes,
	// goroutine exits silenty, kinda like unix shell's "&".
	// A goroutine is simply a function executing concurrently with other goroutines
	// in the same address space. It's not a thread or a process. More about this coming up.
	// Let me try to tldr; this:
	// The following function call is "stepped out" of the main execution stack onto it's own
	// execution stack so that it can run concurrently.
	go s.httpServer.Serve(s.listener)

	fmt.Println("Server now listening")

	return nil

}

// Shutdown the server.
func (s *server) shutdown() error {

	if s.listener != nil {
		// Then stop the server.
		err := s.listener.Close()
		s.listener = nil
		if err != nil {
			return err
		}
	}

	fmt.Println("Shutting down server")
	return nil
}

// We have seen this before. http.HandlerFunc is an adapter which allows us to use this humble
// function to be used as a HTTP handler. See below.
func handle(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, msg)
}

// We create a server and return a pointer to a standard package struct type called Server
// A girl needs a pointer because it holds the memory address of a variable.
// but pointer in Go is not a many faced God like in C (no pointer arithmetic)
func newServer(port string) *server {

	httpServer := &http.Server{Addr: ":" + port, Handler: http.HandlerFunc(handle)}
	// Initializing a struct is easy. Simply use the name of the field as a key to the
	// value you want to pass. Also notice how we only did a partial initiliaztion.
	return &server{httpServer: httpServer}

}

// A lot happening here! This is our test rig.
// In Go we can pass function as an argument.
// Here we will pass our test logic.
func testRig(f func()) {
	// create a new server
	server := newServer(port)
	//start the server
	server.listenAndServe()
	// and "defer" it to shutdown when this function returns
	defer server.shutdown()
	// calling the function containing our test logic
	f()
	// yes it's totally cool in Go for a function to not return anything
	// For to return nothing, there is
}

func moveAlong() {
	fmt.Println("Not the droid you lookin for...")
}

func main() {

	// Hi, you are here. I hope you read this top to bottom.

	server := newServer(port)
	//start the server
	server.listenAndServe()
	// and "defer" it to shutdown when this function returns
	defer server.shutdown()
	// Notice which defer executes first
	defer moveAlong()

	// Now, go run simpleserver_1.go while I wait...
	// I wouldn't mind waiting for: go build && ./simpleserver_1 either.
	//...
	//...
	// What! It simply exited! Well, remember the "go" invocation we did earlier
	// and how it started a new goroutine? In the previous example our
	// humble server would block the current execution, hanging around, waiting for
	// that curl command(you didn't use that did you), but now since the main "goroutine"
	// has no work it simply,uh, goes away.
	// But let's find the solution to this on another day, for now let's just go onto
	// our test: simpleserver_1_test.go
}
