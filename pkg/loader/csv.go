package loader

import (
	"encoding/csv"
	"os"
	"strings"
	"autobase/pkg/log"
)

func ReadCSV(filePath string) ([][]string, error) {
	log.Info("Reading CSV file: " + filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	allRows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(allRows) == 0 {
		log.Info("No rows found in CSV file: " + filePath)
		return [][]string{}, nil
	}

	headers := allRows[0]
	var validIndices []int
	for i, header := range headers {
		if !strings.HasPrefix(header, "#") {
			validIndices = append(validIndices, i)
		}
	}

	filteredRows := make([][]string, len(allRows))
	for i, row := range allRows {
		var filteredRow []string
		for _, idx := range validIndices {
			filteredRow = append(filteredRow, row[idx])
		}
		filteredRows[i] = filteredRow
	}

	return filteredRows, nil
}

func UpdateCSV(filePath string, rows [][]string) error {
	log.Info("Updating CSV file: " + filePath)

	//check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Error("File does not exist: " + filePath)
		return err
	}

	//open the file
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Error("An error occurred while opening the file: " + err.Error())
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	err = writer.WriteAll(rows)
	if err != nil {
		log.Error("An error occurred while writing to the file: " + err.Error())
		return err
	}

	return nil
}
