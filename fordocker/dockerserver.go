package main

import (
	"net/http"
	"fmt"
)

func indexHanlder(w http.ResponseWriter, _ *http.Request){
	fmt.Fprint(w, "Hello docker...")
}

func main() {
	http.HandleFunc("/", indexHanlder)

	http.ListenAndServe(":8080", nil)
}
