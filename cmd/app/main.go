package main

import (
	database "carsdb/internal/database"
	md "carsdb/internal/database/model"
	"carsdb/internal/dealers"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	fmt.Println("CarsDB")

	dbname, err := setupDb()
	if err != nil {
		log.Fatal(err)
	}

	db := database.GetGormDB(dbname)

	var msm []md.ModelSourceMapping
	db.Where("model_id not null").Find(&msm)
	fmt.Println("Got model source mapping: ", len(msm))

	ch := make(chan []md.ModelResult)
	go dealers.ScrapeCazoo(ch, msm)
	go dealers.ScrapeArnoldClark(ch, msm)
	go dealers.ScrapeCinch(ch, msm)

	carModels := []md.ModelResult{}
	for i := 0; i < 3; i++ {
		res := <-ch
		carModels = append(carModels, res...)
	}
	close(ch)

	log.Printf("Got %d car models", len(carModels))

	batchSize := 100
	for i := 0; i < len(carModels); i += batchSize {
		end := i + batchSize
		if end > len(carModels) {
			end = len(carModels)
		}
		batch := carModels[i:end]
		db.Create(&batch)
	}

}

func setupDb() (*string, error) {
	// copy database and rename with date
	currentDate := time.Now().Format("2006-01-02_15:04:05")
	originalFile := "cars.sqlite3"
	newFile := "cars_" + currentDate + ".db"

	// 1. Open the original file for reading
	source, err := os.Open(originalFile)
	if err != nil {
		fmt.Println("Error opening original file:", err)
		return nil, err
	}
	defer source.Close()

	// 2. Create the new file (destination)
	destination, err := os.Create(newFile)
	if err != nil {
		fmt.Println("Error creating new file:", err)
		return nil, err
	}
	defer destination.Close()

	// 3. Copy the contents
	_, err = io.Copy(destination, source)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return nil, err
	}

	// 4. Optionally, check for a successful copy
	err = destination.Sync()
	if err != nil {
		fmt.Println("Error syncing file:", err)
		return nil, err
	}

	// 5. Close files (optional, but good practice)
	source.Close()
	destination.Close()
	fmt.Println("New database created with name: ", newFile)

	return &newFile, nil
}
