package main

import (
	"fmt"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()

	_, err1 := fmt.Fprintln(w, "hello world")
	panic(err1)
}

func main() {
	http.HandleFunc("/", IndexHandler)
	err := http.ListenAndServe("127.0.0.1:8000", nil)
	if err != nil {
		panic(err)
	}
}
