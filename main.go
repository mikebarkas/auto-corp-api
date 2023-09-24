package main

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initDb: %s", err)
	}
	defer db.Close()

	e.GET("/json", handleJson)
	e.GET("/count", func(c echo.Context) error {
		return handleCount(db, c)
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func handleCount(db *sql.DB, c echo.Context) error {
	count, err := countRecords(db)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, fmt.Sprintf("countRecords err: %s", err.Error()))
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("There are %d cars for sale.", count))

}

func handleJson(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		Status string
		Hello  string
	}{Status: "OK", Hello: "Hello from Docker"})

}

func initDB() (*sql.DB, error) {
	connString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
	)
	var (
		db  *sql.DB
		err error
	)

	db, err = sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func countRecords(db *sql.DB) (int, error) {
	rows, err := db.Query("SELECT COUNT(*) FROM autos")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			rows.Close()
		}
	}
	return count, nil
}
