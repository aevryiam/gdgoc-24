package service

import (
	"Pelatihan-KMTETI-GoHTTP/src/db"
	"Pelatihan-KMTETI-GoHTTP/src/models"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var bookCollection *mongo.Collection

func init() {
	db.ConnectDB()
	bookCollection = db.Client.Database("go-http-server").Collection("books")
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	var books []models.BookResponse
	cursor, err := bookCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var book models.Book
		cursor.Decode(&book)
		books = append(books, models.BookResponse{
			Title:        book.Title,
			Author:       book.Author,
			Published_at: book.Published_at,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func GetBookDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var book models.Book
	err := bookCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := bookCollection.InsertOne(context.TODO(), book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := bookCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": book})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := bookCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
