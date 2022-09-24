package main

import (
	"fmt"
	"net/http"
)

type pen struct {
	Name  string `json:"name"`
	Price string `json:"price"`
}

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	http.HandleFunc("/api/v1/pens", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("succes"))
	})

	port:="8000"
	fmt.Println("server run on port", port)
	http.ListenAndServe(":"+port,nil)
		
	}
		
	

