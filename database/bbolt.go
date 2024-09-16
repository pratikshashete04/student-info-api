package database

import (
	"log"
	"os"
	"time"

	"go.etcd.io/bbolt"
)

// DB is the global bbolt database connection
var DB *bbolt.DB

// Initialize the database
func InitDatabase() {
	var err error

	DB, err = bbolt.Open("student-info.db", 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("students"))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database initialized")
}

// Close the database when application exits
func CloseDatabase() {
	if err := DB.Close(); err != nil {
		log.Fatal(err)
	}
}

// SaveStudent stores a student record in the students bucket
func SaveStudent(id string, student []byte) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("students"))
		return bucket.Put([]byte(id), student)
	})
}

// GetStudent retrieves a student record by ID
func GetStudent(id string) ([]byte, error) {
	var student []byte
	err := DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("students"))
		student = bucket.Get([]byte(id))
		if student == nil {
			return os.ErrNotExist
		}
		return nil
	})
	return student, err
}

// GetAllStudents retrieves all student records from the students bucket
func GetAllStudents() ([][]byte, error) {
	var students [][]byte
	err := DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("students"))
		return bucket.ForEach(func(k, v []byte) error {
			students = append(students, v)
			return nil
		})
	})
	return students, err
}
