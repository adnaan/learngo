package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

var host = "http://localhost:3333"

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
		if failMsg == string(msg1) {
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

func testHTTPRequest(verb string, resource string, body string) (*http.Response, error) {
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	r, _ := http.NewRequest(verb, fmt.Sprintf("%s%s", host, resource), strings.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	return client.Do(r)
}

func TestJSON(t *testing.T) {
	testRig(func() {
		response, err := testHTTPRequest("POST", "/", `{"name":"Hodor"}`)
		if err != nil {
			t.Fatalf("Request failed %v", err)
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Error reading response body %v", err)
		}

		fmt.Println("test response", string(body))

	})

}
