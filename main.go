package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/json", handleJson)

	e.Logger.Fatal(e.Start(":8080"))
}

func handleJson(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		Status string
		Hello  string
	}{Status: "OK", Hello: "Hello from Docker"})

}
