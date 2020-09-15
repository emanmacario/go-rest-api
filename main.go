package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Init books variable as a slice books struct
var books []Book

// Book struct
type Book struct {
	ID     string  `json:"id"`
	ISBN   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author struct
type Author struct {
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
}

// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Get single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Get params
	params := mux.Vars(r)
	// Loop through books and find book with id
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	// Return empty book if id not found
	json.NewEncoder(w).Encode(&Book{})
}

// Create a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(1000000)) // Mock ID - not safe (do not use for production)
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

// Updates a book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var updatedBook Book
	_ = json.NewDecoder(r.Body).Decode(&updatedBook)
	for index, item := range books {
		if item.ID == params["id"] {
			updatedBook.ID = item.ID
			books = append(books[:index], books[index+1:]...)
			books = append(books, updatedBook)
		}
	}
	json.NewEncoder(w).Encode(books)
}

// Deletes a book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...) // Appending another slice requires ... syntax
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	// Initialise router
	r := mux.NewRouter()

	// Mock data
	// TODO: Implement database
	books = append(books, Book{ID: "1", ISBN: "42069", Title: "Hairy Potter", Author: &Author{Firstname: "Jane", Lastname: "Doe"}})
	books = append(books, Book{ID: "2", ISBN: "42168", Title: "Hairy Potter 2", Author: &Author{Firstname: "Joe", Lastname: "Smith"}})

	// Route handlers / Endpoints
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// book := Book{}
	// bp := &book
	// fmt.Println(book)
	// fmt.Printf("typeof book: %T\n", book)
	// fmt.Printf("typeof bp  : %T\n", bp)

	log.Fatal(http.ListenAndServe(":8000", r))
}
