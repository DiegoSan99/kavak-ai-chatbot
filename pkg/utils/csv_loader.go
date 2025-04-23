package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/tmc/langchaingo/schema"
)

func LoadCSVToDocuments(filePath string) ([]schema.Document, error) {

	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	var documents []schema.Document
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		metadata := make(map[string]any)
		for i, value := range record {
			if i < len(header) {
				if num, err := strconv.ParseFloat(value, 64); err == nil {
					metadata[header[i]] = num
				} else {
					metadata[header[i]] = value
				}
			}
		}

		carInfo := fmt.Sprintf("%s %s %s %s con precio de $%.0f",
			metadata["make"],
			metadata["model"],
			metadata["version"],
			metadata["year"],
			metadata["price"])

		documents = append(documents, schema.Document{
			PageContent: carInfo,
			Metadata:    metadata,
		})
	}

	log.Printf("Loaded %d car records from CSV\n", len(documents))
	return documents, nil
}
