package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"net/http"
  "os/signal"
  "syscall"
	"flag"
)

var msg = "Dr. Who say's alonzi\n"
var anothermsg = "Winter is in Season 10\n"

type server struct {
	httpServer *http.Server
	listener   net.Listener
}

func (s *server) listenAndServe() error {

	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return err
	}
	s.listener = listener
	go s.httpServer.Serve(s.listener)
	fmt.Println("Server now listening")
	return nil

}

func (s *server) shutdown() error {

	if s.listener != nil {
		err := s.listener.Close()
		s.listener = nil
		if err != nil {
			return err
		}
	}
	fmt.Println("Shutting down server")
	return nil

}

func newServer(port string) *server {

	handle := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, msg)
	}
	anotherhandle := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, anothermsg)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handle)
	mux.HandleFunc("/another", anotherhandle)

	httpServer := &http.Server{Addr: ":" + port, Handler: mux}
	return &server{httpServer: httpServer}

}

func testRig(f func()) {
	server := newServer("3333")
	server.listenAndServe()
	defer server.shutdown()
	f()
}

func main() {
	// Let's get port as a commandline input
	var port string
  flag.StringVar(&port, "port", "3333", "./simpleserver_3 -port 3333")
	flag.Parse()


	// a channel to receive unix signals
  sigs := make(chan os.Signal, 1)
	// a channel to receive a finito confirmation on interrupt
  done := make(chan bool, 1)

	// signal.Notify is a method to create a channel which receives
	// SIGINT, SIGTERM unix signals.
  signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	moveAlong := func() {
		fmt.Println("Not the droid you lookin for...")
	}

	server := newServer(port)
	server.listenAndServe()
	defer moveAlong()

	// Another goroutine! This one blocks itself waiting for signal to be received
	// As soon as it receives a signal, it sends a 'true' to 'done' which unblocks too
	// and the program exits.
	go func() {
        sig := <-sigs
        fmt.Println()
        fmt.Println(sig)
				// this is graceful
				server.shutdown()
        done <- true
    }()


	// Ctrl-C sends a SIGINT signal to the program
	fmt.Println("Ctrl-C to interrupt...")
	<-done
	fmt.Println("Exiting...")

	// cd simpleserver_3
	// go build .
	// ./simpleserver_3 -port 4444
}
