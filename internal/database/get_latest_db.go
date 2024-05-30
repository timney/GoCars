package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const filePrefix = "cars_*"

func GetLatestDBName() *string {
	files, err := os.ReadDir(".")
	if err != nil {
		log.Println("Error reading directory:", err)
		return nil
	}

	dbFiles := []DbFile{}

	for _, file := range files {
		fileName := file.Name()

		// Parsing the filename (example)
		fileExtension := filepath.Ext(fileName)
		matched, _ := filepath.Match(filePrefix, fileName)

		if fileExtension == ".db" && matched {
			fmt.Println("File:", fileName)
			parsedTime, err := parseFilenameDate(fileName)
			if err != nil {
				log.Fatalln("Error parsing time:", err)
			}
			dbFiles = append(dbFiles, DbFile{Filename: fileName, Date: parsedTime})
		}
	}

	sort.Slice(dbFiles, func(i, j int) bool {
		return dbFiles[j].Date.Before(dbFiles[i].Date)
	})

	if len(dbFiles) > 0 {
		return &dbFiles[0].Filename
	}

	return nil
}

func parseFilenameDate(filename string) (time.Time, error) {
	parts := strings.Split(filename, "_")
	log.Println(parts[2])
	timeStr := strings.Split(parts[2], ".")
	log.Println("timepart", timeStr)

	dateTimeStr := parts[1] + "_" + timeStr[0] // Combine date and time

	layout := "2006-01-02_15:04:05" // Matching the format in the filename
	parsedTime, err := time.Parse(layout, dateTimeStr)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

type DbFile struct {
	Filename string
	Date     time.Time
}
