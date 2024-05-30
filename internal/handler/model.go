package handler

import (
	"carsdb/internal/database"
	md "carsdb/internal/database/model"
	"carsdb/internal/views"
	"log"

	"github.com/labstack/echo/v4"
)

func ModelHandler(c echo.Context) error {
	makeName := c.Param("make")
	modelName := c.Param("model")
	dbName := database.GetLatestDBName()
	if dbName != nil {
		log.Println("DB name:", *dbName)
	}
	db := database.GetGormDB(dbName)

	var make md.Make
	db.Where("name = ?", makeName).First(&make)

	var model md.Model
	db.Where("name = ?", modelName).First(&model)

	var models []md.ModelResult
	db.Where("model_id = ?", model.ID).Find(&models)

	log.Printf("Models found %d %d", len(models), model.ID)

	return views.ModelPage(model, models).Render(c.Request().Context(), c.Response())
}
