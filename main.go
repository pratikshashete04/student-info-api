package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"student-info-api/database"

	"github.com/gorilla/mux"
)

type Student struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Grade string `json:"grade"`
}

// Initialize database and router
func init() {
	database.InitDatabase()
}

// Close the database when the app exits
func main() {
	defer database.CloseDatabase()

	router := mux.NewRouter()
	router.HandleFunc("/students", GetStudents).Methods("GET")
	router.HandleFunc("/students/{id}", GetStudentByID).Methods("GET")
	router.HandleFunc("/students", CreateStudent).Methods("POST")

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Get all students
func GetStudents(w http.ResponseWriter, r *http.Request) {
	studentsData, err := database.GetAllStudents()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var students []Student
	for _, studentData := range studentsData {
		var student Student
		json.Unmarshal(studentData, &student)
		students = append(students, student)
	}

	json.NewEncoder(w).Encode(students)
}

// Get student by ID
func GetStudentByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	studentData, err := database.GetStudent(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Student not found"})
		return
	}

	var student Student
	json.Unmarshal(studentData, &student)
	json.NewEncoder(w).Encode(student)
}

// Create a new student
func CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student Student
	json.NewDecoder(r.Body).Decode(&student)

	// Generate a random ID for the student
	student.ID = fmt.Sprintf("%d", rand.Int())

	// Marshal student struct to JSON
	studentData, _ := json.Marshal(student)

	// Save to bbolt DB
	err := database.SaveStudent(student.ID, studentData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(student)
}
