package carsdb

import (
	"database/sql"
	"log"

	database "carsdb/internal/database"

	_ "github.com/mattn/go-sqlite3"
)

type ModelSourceMapping struct {
	ID      int
	ModelID sql.NullInt32
	Make    string
	Model   string
}

func GetModelSourceMapping() []ModelSourceMapping {
	db := database.GetDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, model_id, make, model FROM model_source_mapping")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	msmArr := make([]ModelSourceMapping, 0)

	for rows.Next() {
		msm := ModelSourceMapping{}
		err := rows.Scan(&msm.ID, &msm.ModelID, &msm.Make, &msm.Model)
		if err != nil {
			log.Fatal(err)
		}
		msmArr = append(msmArr, msm)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return msmArr
}
