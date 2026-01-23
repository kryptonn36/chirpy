package main

import "net/http"

type Server struct{
	Addr string
	Handler http.Handler
}

func main(){
	mux := http.NewServeMux()
	mux.Handle("/",http.FileServer(http.Dir(".")))
	server := Server{
		Addr    : ":8080",
		Handler : mux ,
	}
	http.ListenAndServe(server.Addr,server.Handler)
}