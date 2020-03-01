package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func getSnippets(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) { // http.Handler interface
	return func(w http.ResponseWriter, r *http.Request) {
		var snippets []Snippet
		if err := db.Find(&snippets).Error; err != nil {
			fmt.Errorf("error retrieving snippets:", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		response, err := json.Marshal(snippets)
		if err != nil {
			fmt.Errorf("marshal error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func getSnippet(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		snippetID := vars["id"]

		var snippet Snippet
		if err := db.First(&snippet, snippetID).Error; err != nil {
			fmt.Errorf("error retrieving snippet:", err)
			http.Error(w, fmt.Sprintf("snippet with ID: %s not found", snippetID), http.StatusNotFound)
			return
		}

		response, err := json.Marshal(snippet)
		if err != nil {
			fmt.Errorf("marshal error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func createSnippet(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var snippet Snippet

		err := json.NewDecoder(r.Body).Decode(&snippet)
		if err != nil {
			fmt.Errorf("Decode error: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = db.Create(&snippet).Error; err != nil {
			fmt.Errorf("Error creating snippet: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := json.Marshal(snippet)
		if err != nil {
			fmt.Errorf("marshal error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	}
}

func updateSnippet(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var snippet Snippet
		var updated Snippet

		vars := mux.Vars(r)
		snippetID := vars["id"]

		if err := db.First(&snippet, snippetID).Error; err != nil {
			fmt.Errorf("error retrieving snippet:", err)
			http.Error(w, fmt.Sprintf("Snippet with ID: %s not found", snippetID), http.StatusNotFound)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&updated)
		if err != nil {
			fmt.Errorf("decode error: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := db.Model(&snippet).Updates(updated).Error; err != nil {
			fmt.Errorf("Error updating snippet: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := json.Marshal(snippet)
		if err != nil {
			fmt.Errorf("marshal error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func deleteSnippet(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		snippetID := vars["id"]

		var snippet Snippet
		if err := db.First(&snippet, snippetID).Error; err != nil {
			fmt.Errorf("error retrieving snippet:", err)
			http.Error(w, fmt.Sprintf("Snippet with ID: %s not found", snippetID), http.StatusNotFound)
			return
		}

		if err := db.Delete(&snippet).Error; err != nil {
			fmt.Errorf("error deleting snippet: ", err)
			http.Error(w, fmt.Sprint("Error deleting snippet: ", err), http.StatusBadRequest)
			return
		}

		response, err := json.Marshal(snippet)
		if err != nil {
			fmt.Errorf("marshal error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func main() {
	dbConn, err := gorm.Open("sqlite3", "./snippetbin.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/snippets", getSnippets(dbConn)).Methods("GET") // GET /snippets
	router.HandleFunc("/snippets", createSnippet(dbConn)).Methods("POST")
	router.HandleFunc("/snippets/{id}", getSnippet(dbConn)).Methods("GET")
	router.HandleFunc("/snippets/{id}", updateSnippet(dbConn)).Methods("PATCH")
	router.HandleFunc("/snippets/{id}", deleteSnippet(dbConn)).Methods("DELETE")

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", router))
}
