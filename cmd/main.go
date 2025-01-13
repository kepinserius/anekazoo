 
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"anekazoo/internal/models"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/gorilla/mux"
)

var db *sql.DB

// Initialize the database connection
func initDB() {
	var err error
	connStr := "host=db user=kepinserius password=Kevinarjuda23. dbname=anekaZoo sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Cannot ping database:", err)
	}
}

// Create a new animal
func createAnimal(w http.ResponseWriter, r *http.Request) {
	var animal models.Animal
	if err := json.NewDecoder(r.Body).Decode(&animal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM animals WHERE id = $1)", animal.ID).Scan(&exists)
	if exists {
		http.Error(w, "Animal already exists", http.StatusConflict)
		return
	}

	_, err = db.Exec("INSERT INTO animals (id, name, class, legs) VALUES ($1, $2, $3, $4)",
		animal.ID, animal.Name, animal.Class, animal.Legs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(animal)
}

// Get all animals
func getAnimals(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, class, legs FROM animals")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var animals []models.Animal
	for rows.Next() {
		var animal models.Animal
		if err := rows.Scan(&animal.ID, &animal.Name, &animal.Class, &animal.Legs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		animals = append(animals, animal)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animals)
}

// Get a single animal by ID
func getAnimal(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var animal models.Animal
	err := db.QueryRow("SELECT id, name, class, legs FROM animals WHERE id = $1", id).Scan(&animal.ID, &animal.Name, &animal.Class, &animal.Legs)
	if err == sql.ErrNoRows {
		http.Error(w, "Hewan Tidak Ditemukan", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animal)
}

// Update an animal by ID
func updateAnimal(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var animal models.Animal
	if err := json.NewDecoder(r.Body).Decode(&animal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := db.Exec("UPDATE animals SET name = $1, class = $2, legs = $3 WHERE id = $4",
		animal.Name, animal.Class, animal.Legs, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		http.Error(w, "Hewan Tidak Ditemukan", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(animal)
}

// Delete an animal by ID
func deleteAnimal(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	res, err := db.Exec("DELETE FROM animals WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		http.Error(w, "Hewan Tidak Ditemukan", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Animal with ID %d deleted", id)
}

func main() {
	initDB()
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/v1/animals", createAnimal).Methods("POST")
	router.HandleFunc("/v1/animals", getAnimals).Methods("GET")
	router.HandleFunc("/v1/animals/{id}", getAnimal).Methods("GET")
	router.HandleFunc("/v1/animals/{id}", updateAnimal).Methods("PUT")
	router.HandleFunc("/v1/animals/{id}", deleteAnimal).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
