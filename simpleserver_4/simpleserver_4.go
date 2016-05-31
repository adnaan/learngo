package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var msg = "Dr. Who say's alonzi\n"
var anothermsg = "Winter is in Season 10\n"

// Response The thing inside `` are called struct tags. Struct tags are used to attach
// meta data to a struct field. Under the hood, packages can use reflection to extract
// this meta data. For e.g. here, encoding/json package uses it to encode/decode the struct
// from and to a json. Also notice the first letter of the type and fields is capitalized.
// The front letter capitlization is Go's way of indicating that the type(and embeddded fields)
// are public.
type Response struct {
	Msg        string `json:"msg,omitempty"`
	Anothermsg string `json:"another_msg,omitempty"`
}

// Request ...
type Request struct {
	Name string `json:"name"`
}

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
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("error reading request body", err)
			return
		}

		// The very useful unmarshal method decodes the incoming json onto the request
		// type
		var request Request
		err = json.Unmarshal(body, &request)
		if err != nil {
			log.Println("json unmarshal error", err)
			return
		}
		log.Println("request:", request)
		// Marshal encodes the struct type into a json byte[]
		response := &Response{Msg: "Hello " + request.Name}
		b, err := json.Marshal(response)
		if err != nil {
			log.Println("json marshal error", err)
			return
		}
		log.Println("response:", string(b))
		io.WriteString(w, string(b))
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
	time.Sleep(500)
	defer server.shutdown()
	f()
}

func main() {
	var port string
	flag.StringVar(&port, "port", "3333", "./simpleserver_4 -port 3333")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	moveAlong := func() {
		fmt.Println("Not the droid you lookin for...")
	}

	server := newServer(port)
	server.listenAndServe()
	defer moveAlong()

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		server.shutdown()
		done <- true
	}()

	fmt.Println("Ctrl-C to interrupt...")
	<-done
	fmt.Println("Exiting...")

	// Take a look at simpleserver_4_test.go
	// go test
	// btw this file is gettig too messy. let's do something about that.
}
