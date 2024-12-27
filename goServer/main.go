package main

import (
	"database/sql"
	"log"
	"os"

	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type CarRecord struct {
	ID          string `json:"id"`
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

	if _, err := os.Stat("./goServer/cars.db"); os.IsNotExist(err) {
		log.Println("База данных не найдена, создается новая база данных.")
	}

	db, err = sql.Open("sqlite3", "./goServer/cars.db")
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

	rows, err := db.Query("SELECT id, title, description, carNumber, carBrand, clientName, clientPhone FROM cars")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	carRecords = []CarRecord{}

	for rows.Next() {
		var car CarRecord
		if err := rows.Scan(&car.ID, &car.Title, &car.Description, &car.CarNumber, &car.CarBrand, &car.ClientName, &car.ClientPhone); err != nil {
			panic(err)
		}
		carRecords = append(carRecords, car)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}
}

func getCars(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, description, carNumber, carBrand, clientName, clientPhone FROM cars")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var carRecords []CarRecord
	for rows.Next() {
		var car CarRecord
		if err := rows.Scan(&car.ID, &car.Title, &car.Description, &car.CarNumber, &car.CarBrand, &car.ClientName, &car.ClientPhone); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		carRecords = append(carRecords, car)
	}

	c.JSON(http.StatusOK, carRecords)
}

func addCar(c *gin.Context) {
	var newCar CarRecord
	if err := c.ShouldBindJSON(&newCar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sqlStatement := `
        INSERT INTO cars (title, description, carNumber, carBrand, clientName, clientPhone)
        VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(sqlStatement, newCar.Title, newCar.Description, newCar.CarNumber, newCar.CarBrand, newCar.ClientName, newCar.ClientPhone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newCar)
}

func deleteCar(c *gin.Context) {
	carNumber := c.Param("id")

	sqlStatement := `DELETE FROM cars WHERE carNumber = ?`
	result, err := db.Exec(sqlStatement, carNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car deleted"})
}

func updateCar(c *gin.Context) {
	carNumber := c.Param("id")
	var updatedCar CarRecord

	if err := c.ShouldBindJSON(&updatedCar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sqlStatement := `
        UPDATE cars 
        SET title = ?, description = ?, carBrand = ?, clientName = ?, clientPhone = ?
        WHERE carNumber = ?`
	_, err := db.Exec(sqlStatement, updatedCar.Title, updatedCar.Description, updatedCar.CarBrand, updatedCar.ClientName, updatedCar.ClientPhone, carNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCar)
}

func main() {
	initDB()
	defer db.Close()

	loadCarRecords()

	router := gin.Default()
	router.Static("/static", "./static")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	router.GET("/api/cars", getCars)
	router.POST("/api/cars", addCar)
	router.DELETE("/api/cars/:id", deleteCar)
	router.PUT("/api/cars/:id", updateCar)

	router.Run(":8080")
}
