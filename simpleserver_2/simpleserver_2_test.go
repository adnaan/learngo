package main

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestServerFail(t *testing.T) {

	testRig(func() {

		res, err := http.Get("http://localhost:3333/another")
		if err != nil {
			t.Fatal(err)
		}

		msg1, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		failMsg := "It's bigger on the inside"
		if failMsg != string(msg1) {
			t.Fatal("Unexpected message: ", string(msg1))
		}

	})
}

func TestServerPass(t *testing.T) {
	testRig(func() {
		res, err := http.Get("http://localhost:3333/another")
		if err != nil {
			t.Fatal(err)
		}

		msg1, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		//You shall, uh, pass.
		if anothermsg != string(msg1) {
			t.Fatal("Unexpected message")
		}

	})
}
