package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-memdb"
)

type Note struct {
	Id    string
	Title string
	Text  string
}

var db *memdb.MemDB

func GetNotesHandler(w http.ResponseWriter, r *http.Request) {
	txn := db.Txn(false)
	defer txn.Abort()
	it, err := txn.Get("notes", "id")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var notes []*Note
	for obj := it.Next(); obj != nil; obj = it.Next() {
		note := obj.(*Note)
		notes = append(notes, note)
	}

	jsonResponse, jsonError := json.Marshal(notes)
	if jsonError != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(string(jsonResponse))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func CreateNotesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	txn := db.Txn(true)

	note.Id = uuid.New().String()
	if err := txn.Insert("notes", &note); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	txn.Commit()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	txn := db.Txn(false)

	it, err := txn.Get("notes", "id", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(it)
	obj := it.Next()
	if obj == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var note Note
	err = json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	note.Id = id

	txn = db.Txn(true)

	if err := txn.Insert("notes", &note); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	txn.Commit()
	json.NewEncoder(w).Encode(&note)

}

func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	txn := db.Txn(true)
	var note Note
	note.Id = id
	_, err := txn.DeleteAll("notes", "id", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	txn.Commit()
}

func initDB() {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"notes": &memdb.TableSchema{
				Name: "notes",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.UUIDFieldIndex{Field: "Id"},
					},
				},
			},
		},
	}

	conn, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}
	db = conn

}

func main() {
	initDB()
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/notes", GetNotesHandler).Methods("GET")
	r.HandleFunc("/notes", CreateNotesHandler).Methods("POST")
	r.HandleFunc("/notes/{id}", UpdateNoteHandler).Methods("PUT")
	r.HandleFunc("/notes/{id}", DeleteNoteHandler).Methods("DELETE")

	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
