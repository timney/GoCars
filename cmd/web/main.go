package main

import (
	"carsdb/internal/handler"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", handler.HomeHandler)
	e.GET("/make/:make", handler.MakeHandler)
	e.GET("/make/:make/:model", handler.ModelHandler)
	e.Logger.Fatal(e.Start(":1323"))
}
