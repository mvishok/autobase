package loader

import (
	"encoding/csv"
	"os"
	"syncengin/pkg/log"
	"strings"
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
