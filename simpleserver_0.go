// Package main is the entry point of an executable Go program
// Are there non-executable Go programs? Well Yes, they are called...wait-for-it
// ...Packages
package main

// The import block. Each import is a package
// Notice how the methods we use from a package have their
//first letter capitalized.
import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// Declaring it outside a function body hence
// the following is valid
var port = "3333"
// while this is not
// port:=3333

// Handling the request and replying back
// Umm Guy, you have your argument types reversed...
// Yeah this is how we roll. Is this easier to read?
func handle(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, fmt.Sprintf("Hello from port %q \n", port))
}

// Entry point of the package main
func main() {
  //But hey we can do that weird (:=) declaration in here
  // Also if it looks like a string, prints like a string, it's probably a string
  duck := "QuackyQuack"
  // Handle a route
	http.HandleFunc("/", handle)
  // Oh look format specifiers %q, %s, %d
	log.Printf("%q on localhost:%q\n",duck,port)
  // And serve it and pipe the error value directly into a log
	log.Fatal(http.ListenAndServe(":"+port, nil))
  // How do we test this? Let's bring out the curls and the postman's?
  // Ah no, lets TDD this program.
}
