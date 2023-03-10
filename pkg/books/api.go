package books

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
)

func BooksHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		books := GetAllBooks()
		json.NewEncoder(w).Encode(books)
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var book Book
		err = json.Unmarshal(body, &book)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		errs, err := validateBook(book)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if errs != "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "%s", errs)
			return
		}
		id, created := CreateBook(book)
		if created {
			w.Header().Add("Location", "/api/books/"+id)
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "Book with given id already exists")
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported request method"))
	}
}

func BookHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := strings.ToLower(params["id"])
	switch method := r.Method; method {
	case http.MethodGet:
		book, found := GetBook(id)
		if found {
			json.NewEncoder(w).Encode(book)
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Book with given id not present")
		}
	case http.MethodPut:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var book Book
		err = json.Unmarshal(body, &book)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		errs, err := validateBook(book)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if errs != "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "%s", errs)
			return
		}
		newid, updated := UpdateBook(id, book)
		if updated {
			w.Header().Add("Location", "/api/books/"+newid)
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Book with given id not present")
		}
	case http.MethodDelete:
		deleted := DeleteBook(id)
		if deleted {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Book with given id not present")
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported request method"))
	}
}
