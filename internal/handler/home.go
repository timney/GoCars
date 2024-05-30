package handler

import (
	"carsdb/internal/database"
	md "carsdb/internal/database/model"
	"carsdb/internal/views"
	"log"

	"github.com/labstack/echo/v4"
)

func HomeHandler(c echo.Context) error {
	dbName := database.GetLatestDBName()
	if dbName != nil {
		log.Println("DB name:", *dbName)
	}
	db := database.GetGormDB(dbName)
	var makes []md.Make
	db.Find(&makes)

	// db.Table("make").Select("make.*").Joins("inner join model on make.id = model.make_id inner join model_result on model.id = model_result.model_id").Scan(&makes)
	db.Raw("select distinct make.* from make inner join model on make.id = model.make_id inner join model_result on model.id = model_result.model_id").Scan(&makes)

	return views.HomePage(makes).Render(c.Request().Context(), c.Response())
}
