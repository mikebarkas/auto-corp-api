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
	e.GET("/search", func(c echo.Context) error {
		return handleSearch(db, c)
	})

	e.Logger.Fatal(e.Start(":8080"))
}

type Auto struct {
	ID      int64  `json:"id"`
	Year    int64  `json:"year"`
	Make    string `json:"make"`
	Model   string `json:"model"`
	Color   string `json:"color"`
	Price   string `json:"price"`
	Mileage int64  `json:"mileage"`
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

func handleSearch(db *sql.DB, c echo.Context) error {
	var value string
	var search string

	// Only one param allowed
	// TODO: improve how the params are defined
	params := c.QueryParams()
	if val, ok := params["make"]; ok {
		value = val[0]
		search = "make"
	}
	if val, ok := params["price"]; ok {
		value = val[0]
		search = "price"
	}
	if val, ok := params["mileage"]; ok {
		value = val[0]
		search = "mileage"
	}
	// TODO: add param error checking and return

	rows, err := searchParam(db, c, search, value)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, struct {
		Status  string
		Data    []Auto
		Results int
	}{Status: "OK", Results: len(rows), Data: rows})

}

func searchParam(db *sql.DB, c echo.Context, search string, value string) ([]Auto, error) {

	var autos []Auto
	// Comparison changes if numbers are given
	compare := "="
	if search == "price" || search == "mileage" {
		compare = "<="
	}

	queryString := fmt.Sprintf("SELECT * FROM autos WHERE %s%s'%v'", search, compare, value)

	rows, err := db.Query(queryString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var auto Auto
		if err := rows.Scan(
			&auto.ID,
			&auto.Year,
			&auto.Make,
			&auto.Model,
			&auto.Color,
			&auto.Price,
			&auto.Mileage,
		); err != nil {
			return nil, fmt.Errorf("handleSearch: %v", err)
		}
		autos = append(autos, auto)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("handleSearch: %v", err)
	}
	return autos, nil
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
