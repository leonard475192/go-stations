package main

import (
	"fmt"
	"net/http"
)

func healthz(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello World")
}
