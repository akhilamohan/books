package books

import (
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	"time"
)

var books map[string]Book

type Book struct {
	Title           string `json:"title" validate:"required"`
	Author          string `json:"author" validate:"required"`
	PublicationDate string `json:"publication-date" validate:"required"`
	ID              string `json:"id" validate:"required"`
}

func init() {
	books = make(map[string]Book)
}

// Clear books
func EmptyBooks() {
	for k := range books {
		delete(books, k)
	}
}

// return all books in list
func GetAllBooks() []Book {
	var allbooks []Book
	for _, book := range books {
		allbooks = append(allbooks, book)
	}
	return allbooks
}

// Add a book to list
func CreateBook(book Book) (string, bool) {
	_, ok := books[book.ID]
	if ok {
		return "", false
	}
	books[book.ID] = book
	return book.ID, true
}

// Returns book of given id
func GetBook(id string) (Book, bool) {
	book, ok := books[id]
	if !ok {
		return Book{}, false
	}
	return book, true
}

// Update book with given id
func UpdateBook(id string, book Book) (string, bool) {
	_, ok := books[id]
	if !ok {
		return "", false
	}
	delete(books, id)
	books[book.ID] = book
	return book.ID, true
}

// Delete book with given id
func DeleteBook(id string) bool {
	_, ok := books[id]
	if !ok {
		return false
	}
	delete(books, id)
	return true
}

func validateBook(book Book) (string, error) {
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	if !found {
		return "", fmt.Errorf("translator not found")
	}

	v := validator.New()

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		return "", err
	}

	err := v.Struct(book)
	errs := translateError(err, trans)
	strerrs := ""
	for _, e := range errs {
		strerrs += e.Error()
	}
	_, err = time.Parse("2006-01-02", book.PublicationDate)
	if err != nil {
		fmt.Println(err)
		strerrs += "Invalid date format"
	}
	return strerrs, nil
}

func translateError(err error, trans ut.Translator) (errs []error) {
	if err == nil {
		return nil
	}
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(trans))
		errs = append(errs, translatedErr)
	}
	return errs
}
