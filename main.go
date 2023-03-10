package main

import (
	"github.com/akhilamohan/books/pkg/books"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/books", books.BooksHandler).Methods("POST")
	router.HandleFunc("/api/books", books.BooksHandler).Methods("GET")
	router.HandleFunc("/api/book/{id}", books.BookHandler).Methods("GET")
	router.HandleFunc("/api/book/{id}", books.BookHandler).Methods("PUT")
	router.HandleFunc("/api/book/{id}", books.BookHandler).Methods("DELETE")

	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatal(err)
	}
}
