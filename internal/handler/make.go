package handler

import (
	"carsdb/internal/database"
	md "carsdb/internal/database/model"
	"carsdb/internal/views"
	"log"

	"github.com/labstack/echo/v4"
)

func MakeHandler(c echo.Context) error {
	makeName := c.Param("make")
	dbName := database.GetLatestDBName()
	if dbName != nil {
		log.Println("DB name:", *dbName)
	}
	db := database.GetGormDB(dbName)

	var make md.Make
	db.Where("name = ?", makeName).First(&make)

	var models []md.Model
	db.Where("make_id = ?", make.ID).Find(&models)

	return views.MakePage(make, models).Render(c.Request().Context(), c.Response())
}
