package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
)

var port = "3333"
var msg = "Dr. Who say's alonzi"
var anothermsg = "Winter is in Season 10"

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

  // This tiny function was kinda alone out there so lets just treat it as a local var.
  handle := func (w http.ResponseWriter, r *http.Request) {
  	io.WriteString(w, msg)
  }
  // and lets setup another route
  anotherhandle := func (w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, anothermsg)
  }
  // You guessed it. ServeMux is an HTTP request multiplexer.
  mux:=http.NewServeMux()
  // Functions are first class citizens in Go. And like any other good citizen
  // they can be passed as arguments.
  mux.HandleFunc("/",handle)
  mux.HandleFunc("/another",anotherhandle)

	httpServer := &http.Server{Addr: ":" + port, Handler:mux }
	return &server{httpServer: httpServer}

}

func testRig(f func()) {

	server := newServer(port)
	server.listenAndServe()
	defer server.shutdown()
	f()
}

func main() {

  // Here comes the chan type, the Jackie Chan of concurrency(sorry). A "Channel"
  // is a typed pipe through which you can send and receive values across goroutines.
  // This particular channel is of type struct{}, (an empty struct occupies zero bytes
  // of storage. It's totally cool to have chan int, chan bool chan MyStruct etc.)
  //since we are using it for a mundane task.
  ch := make(chan struct{})

  moveAlong := func(){
    fmt.Println("Not the droid you lookin for...")
  }

	server := newServer(port)
	server.listenAndServe()
	defer server.shutdown()
	defer moveAlong()

  // This channel will wait to receive a value(and in our case it will wait for
  // eternity). While it's waiting further execution of the main goroutine will remain
  // blocked hence serving our purpose. To unblock this channel someone needs to
  // do : ch <- someVal or close(ch). More of this coming up.
  // Now you can do : go run simpleserver_2.go  and expect it to work.
  // Goto http://localhost:3333/another to check whether it did.
  <-ch
}
