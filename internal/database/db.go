package database

import (
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func GetGormDB(databaseName *string) *gorm.DB {
	var dbname string
	if databaseName == nil {
		dbname = "cars.sqlite3"
	} else {
		dbname = *databaseName
	}

	db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func GenerateStructs() {
	gormdb := GetGormDB(nil)

	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/database/generated",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(gormdb) // reuse your gorm db

	g.ApplyBasic(
		// Generate struct `User` based on table `users`
		g.GenerateModel("model_source_mapping"),
		g.GenerateModel("make"),
		g.GenerateModel("model"),
		g.GenerateModel("model_result"),
		g.GenerateModel("job_run"),
		g.GenerateModel("job_source"),
	)
	g.Execute()
}

// func GetDB() *sql.DB {
// 	db, err := sql.Open("sqlite3", "cars.sqlite3")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return db
// }
