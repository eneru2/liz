package main

import (
	"net/http"
	// "embed"
	"github.com/go-chi/chi/v5"
	handlers "liz/.liz/files"
)

func main() {
	r := chi.NewRouter()
	// r.Handle("/*", public())
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("static"))))
	handlers.Routes(r)
	
	
	http.ListenAndServe(":4300", r)
}

// //go:embed static/*
// var public_folder embed.FS

// func public() http.Handler {
// 	return http.FileServerFS(public_folder)
// }