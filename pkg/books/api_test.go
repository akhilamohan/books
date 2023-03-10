package books

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddBookHandler(t *testing.T) {
	EmptyBooks()
	testBook := make(map[string]Book)
	//emptyBook := make(map[string]Book)
	testBook["id-123"] = Book{Title: "title1", Author: "author1", PublicationDate: "1991-10-04", ID: "id-123"}

	tt := []struct {
		name       string
		method     string
		body       string
		expBooks   map[string]Book
		want       string
		statusCode int
	}{
		{
			name:       "Add book",
			method:     http.MethodPost,
			body:       `{"title": "title1", "author": "author1", "publication-date": "1991-10-04", "id": "id-123"}`,
			expBooks:   testBook,
			want:       ``,
			statusCode: http.StatusCreated,
		},
		{
			name:       "Add book with empty entry",
			method:     http.MethodPost,
			body:       `{"title": "", "author": "author1", "publication-date": "1991-10-04", "id": "id-124"}`,
			expBooks:   testBook,
			want:       `Title is a required field`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Add book with invalid date format",
			method:     http.MethodPost,
			body:       `{"title": "abc", "author": "author1", "publication-date": "04/101991", "id": "id-124"}`,
			expBooks:   testBook,
			want:       `Invalid date format`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Add book with same id",
			method:     http.MethodPost,
			body:       `{"title":"abc","author":"author1","publication-date":"1991-10-04","id":"id-123"}`,
			expBooks:   testBook,
			want:       `Book with given id already exists`,
			statusCode: http.StatusConflict,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//defer cleanup()
			request := httptest.NewRequest(tc.method, "/api/books", strings.NewReader(tc.body))
			responseRecorder := httptest.NewRecorder()

			handler := http.HandlerFunc(BooksHandler)
			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}

			assert.Equal(t, tc.expBooks, books, "both should be equal")
		})
	}
}

func TestGetBooksHandler(t *testing.T) {
	EmptyBooks()
	testBook := make(map[string]Book)
	emptyBook := make(map[string]Book)
	testBook["id-123"] = Book{Title: "title1", Author: "author1", PublicationDate: "1991-10-04", ID: "id-123"}

	tt := []struct {
		name       string
		method     string
		body       string
		expBooks   map[string]Book
		want       string
		statusCode int
	}{
		{
			name:       "Empty books",
			method:     http.MethodGet,
			expBooks:   emptyBook,
			want:       `null`,
			statusCode: http.StatusOK,
		},
		{
			name:       "after adding book",
			method:     http.MethodGet,
			expBooks:   testBook,
			want:       `[{"title":"title1","author":"author1","publication-date":"1991-10-04","id":"id-123"}]`,
			statusCode: http.StatusOK,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//defer cleanup()
			if tc.name == "after adding book" {
				CreateBook(Book{Title: "title1", Author: "author1", PublicationDate: "1991-10-04", ID: "id-123"})
			}
			request := httptest.NewRequest(tc.method, "/api/books", nil)
			responseRecorder := httptest.NewRecorder()

			handler := http.HandlerFunc(BooksHandler)
			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSuffix(strings.TrimSpace(responseRecorder.Body.String()), "\n") != strings.TrimSpace(tc.want) {
				t.Errorf("Want '%s', got '%s'", strings.TrimSpace(tc.want), strings.TrimSuffix(strings.TrimSpace(responseRecorder.Body.String()), "\n"))
			}

			assert.Equal(t, tc.expBooks, books, "both should be equal")
		})
	}
}

func TestGetBookHandler(t *testing.T) {
	EmptyBooks()
	CreateBook(Book{Title: "title1", Author: "author1", PublicationDate: "1991-10-04", ID: "id-123"})
	testBook := make(map[string]Book)
	testBook["id-123"] = Book{Title: "title1", Author: "author1", PublicationDate: "1991-10-04", ID: "id-123"}

	tt := []struct {
		name       string
		method     string
		expBooks   map[string]Book
		want       string
		statusCode int
		id         string
	}{
		{
			name:       "Id exists",
			method:     http.MethodGet,
			expBooks:   testBook,
			want:       `{"title":"title1","author":"author1","publication-date":"1991-10-04","id":"id-123"}`,
			id:         "id-123",
			statusCode: http.StatusOK,
		},
		{
			name:       "Id not present",
			method:     http.MethodGet,
			expBooks:   testBook,
			want:       `Book with given id not present`,
			id:         "id-124",
			statusCode: http.StatusNotFound,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			requrl := "/api/book/" + tc.id
			request := httptest.NewRequest(tc.method, requrl, nil)
			param := map[string]string{"id": tc.id}
			request = mux.SetURLVars(request, param)
			responseRecorder := httptest.NewRecorder()

			handler := http.HandlerFunc(BookHandler)
			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSuffix(strings.TrimSpace(responseRecorder.Body.String()), "\n") != strings.TrimSpace(tc.want) {
				t.Errorf("Want '%s', got '%s'", strings.TrimSpace(tc.want), strings.TrimSuffix(strings.TrimSpace(responseRecorder.Body.String()), "\n"))
			}

			assert.Equal(t, tc.expBooks, books, "both should be equal")
		})
	}
}

func TestDeleteBookHandler(t *testing.T) {
	EmptyBooks()
	CreateBook(Book{Title: "title1", Author: "author1", PublicationDate: "1991-10-04", ID: "id-123"})
	testBook := make(map[string]Book)
	emptyBook := make(map[string]Book)
	testBook["id-123"] = Book{Title: "title1", Author: "author1", PublicationDate: "1991-10-04", ID: "id-123"}

	tt := []struct {
		name       string
		method     string
		expBooks   map[string]Book
		want       string
		statusCode int
		id         string
	}{
		{
			name:       "Id not present",
			method:     http.MethodDelete,
			expBooks:   testBook,
			want:       `Book with given id not present`,
			id:         "id-124",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Id present",
			method:     http.MethodDelete,
			expBooks:   emptyBook,
			want:       ``,
			id:         "id-123",
			statusCode: http.StatusOK,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			requrl := "/api/book/" + tc.id
			request := httptest.NewRequest(tc.method, requrl, nil)
			param := map[string]string{"id": tc.id}
			request = mux.SetURLVars(request, param)
			responseRecorder := httptest.NewRecorder()

			handler := http.HandlerFunc(BookHandler)
			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSuffix(strings.TrimSpace(responseRecorder.Body.String()), "\n") != strings.TrimSpace(tc.want) {
				t.Errorf("Want '%s', got '%s'", strings.TrimSpace(tc.want), strings.TrimSuffix(strings.TrimSpace(responseRecorder.Body.String()), "\n"))
			}

			assert.Equal(t, tc.expBooks, books, "both should be equal")
		})
	}
}

func TestUpdateBookHandler(t *testing.T) {
	EmptyBooks()
	CreateBook(Book{Title: "title1", Author: "author1", PublicationDate: "1991-10-04", ID: "id-123"})
	testBook := make(map[string]Book)
	updatedBook := make(map[string]Book)
	testBook["id-123"] = Book{Title: "title1", Author: "author1", PublicationDate: "1991-10-04", ID: "id-123"}
	updatedBook["id-123"] = Book{Title: "title1", Author: "author2", PublicationDate: "1991-10-04", ID: "id-123"}

	tt := []struct {
		name       string
		method     string
		expBooks   map[string]Book
		want       string
		body       string
		statusCode int
		id         string
	}{
		{
			name:       "Id not present",
			method:     http.MethodPut,
			expBooks:   testBook,
			want:       `Book with given id not present`,
			body:       `{"title": "title1", "author": "author1", "publication-date": "1991-10-04", "id": "id-123"}`,
			id:         "id-124",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Id present",
			method:     http.MethodPut,
			expBooks:   updatedBook,
			body:       `{"title": "title1", "author": "author2", "publication-date": "1991-10-04", "id": "id-123"}`,
			want:       ``,
			id:         "id-123",
			statusCode: http.StatusOK,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			requrl := "/api/book/" + tc.id
			request := httptest.NewRequest(tc.method, requrl, strings.NewReader(tc.body))
			param := map[string]string{"id": tc.id}
			request = mux.SetURLVars(request, param)
			responseRecorder := httptest.NewRecorder()

			handler := http.HandlerFunc(BookHandler)
			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSuffix(strings.TrimSpace(responseRecorder.Body.String()), "\n") != strings.TrimSpace(tc.want) {
				t.Errorf("Want '%s', got '%s'", strings.TrimSpace(tc.want), strings.TrimSuffix(strings.TrimSpace(responseRecorder.Body.String()), "\n"))
			}

			assert.Equal(t, tc.expBooks, books, "both should be equal")
		})
	}
}
