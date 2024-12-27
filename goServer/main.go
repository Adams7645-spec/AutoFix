package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"

	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

type CarRecord struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CarNumber   string `json:"carNumber"`
	CarBrand    string `json:"carBrand"`
	ClientName  string `json:"clientName"`
	ClientPhone string `json:"clientPhone"`
}

var (
	carRecords []CarRecord
	mu         sync.Mutex
	db         *sql.DB
)

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./cars.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS cars (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        description TEXT,
        carNumber TEXT,
        carBrand TEXT,
        clientName TEXT,
        clientPhone TEXT
    );`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatal(err)
	}
}

func loadCarRecords() {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.Open("data.json")
	if err != nil {
		if os.IsNotExist(err) {
			carRecords = []CarRecord{}
			return
		}
		panic(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &carRecords); err != nil {
		panic(err)
	}
}

func saveCarRecords() {
	data, err := json.Marshal(carRecords)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("data.json", data, 0644); err != nil {
		panic(err)
	}
}

func getCars(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(http.StatusOK, carRecords)
}

func addCar(c *gin.Context) {
	var newCar CarRecord
	if err := c.ShouldBindJSON(&newCar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	carRecords = append(carRecords, newCar)
	mu.Unlock()

	c.JSON(http.StatusCreated, newCar)
}

func deleteCar(c *gin.Context) {
	carNumber := c.Param("id")

	mu.Lock()
	defer mu.Unlock()

	for i, car := range carRecords {
		if car.CarNumber == carNumber {
			carRecords = append(carRecords[:i], carRecords[i+1:]...)
			saveCarRecords()
			c.JSON(http.StatusOK, gin.H{"message": "Car deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
}

func main() {
	loadCarRecords()

	initDB()
	defer db.Close()

	router := gin.Default()
	router.Static("/static", "./static")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	router.GET("/api/cars", getCars)
	router.POST("/api/cars", addCar)
	router.DELETE("/api/cars/:id", deleteCar)
	//router.PUT("/api/cars/:id", updateCar)

	router.Run(":8080")
}
