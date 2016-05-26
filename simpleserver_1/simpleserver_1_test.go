package main

import (
	"io/ioutil"
	"net/http"
	"testing"
)

// All you need to do is to have this function signature and "go" will automatically
// run this test for you. It also requires you to capitalize the first letter.
// so that it can separate the wheat from the chaff.
// Just to be sure we will do a negative and a positive test
func TestServerFail(t *testing.T) {

	// Being efficient and passing the whole function
	testRig(func() {

		// our server is really simple, so a Get call would do for now
		res, err := http.Get("http://localhost:3333")
		if err != nil {
			t.Fatal(err)
		}

		// We use a utility function from the standard Go package ioutil
		// to read the response's body as a byte[]
		msg1, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		// You shall not pass
		failMsg := "It's bigger on the inside"
		// we will convert byte[] to a string. simply string(myByteArray)
		if failMsg != string(msg1) {
			t.Fatal("Unexpected message: ", string(msg1))
		}

	})

	// What are you waiting for? go test
}

func TestServerPass(t *testing.T) {
	testRig(func() {
		res, err := http.Get("http://localhost:3333")
		if err != nil {
			t.Fatal(err)
		}

		msg1, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		//You shall, uh, pass.
		if msg != string(msg1) {
			t.Fatal("Unexpected message")
		}

	})
}
