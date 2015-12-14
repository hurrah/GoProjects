package main

import (
	"io"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hey You, How's it going?")
}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":9000", nil)
}
