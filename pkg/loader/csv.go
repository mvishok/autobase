package loader

import (
	"encoding/csv"
	"os"
	"statix/pkg/log"
)

func ReadCSV(filePath string) ([][]string, error) {
	log.Info("Reading CSV file: " + filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}
