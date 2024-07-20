package query

import (
	"net/url"
	"strings"
	"syncengin/pkg/loader"
	"syncengin/pkg/log"
)

var headers []string

func Run(csvPath string, query url.Values, access_level string) []map[string]string {
	rows, err := loader.ReadCSV(csvPath)
	if err != nil {
		log.Error("An error occurred while reading the CSV file: " + err.Error())
		return nil
	}
	if len(rows) < 2 {
		return nil
	}

	results := []map[string]string{}

	//if empty query, return all rows
	if len(query) == 0 {
		for i, row := range rows {
			if i == 0 {
				continue
			}
			result := make(map[string]string)
			for i, value := range row {
				result[rows[0][i]] = value
			}
			results = append(results, result)
		}
		return results
	}

	//if select query, return only selected columns
	if query.Get("select") != "" {

		cols := strings.Split(query.Get("select"), ",")

		for i, row := range rows {
			if i == 0 {
				headers = row
				continue
			}
			if whereClause(row, query) {
				result := make(map[string]string)
				for i, value := range row {
					if contains(cols, rows[0][i]) {
						result[rows[0][i]] = value
					}
				}
				results = append(results, result)
			}
		}
	}

	return results
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// whereClause returns true if the row satisfies the where clause
func whereClause(row []string, query url.Values) bool {
	for key, value := range query {
		if key == "where" {
			for _, v := range value {
				//syntax: &where=key:op:val,key:op:val...
				conditions := strings.Split(v, ",")

				for _, condition := range conditions {
					//syntax: key:op:val
					parts := strings.Split(condition, ":")

					if len(parts) != 3 {
						return false
					}

					key := parts[0]
					op := parts[1]
					val := parts[2]

					index := getHeaderIndex(key)
					if index == -1 {
						return false
					}

					switch op {
					case "eq":
						if row[index] != val {
							return false
						}
					default:
						return false
					}
				}
			}
		}
	}
	return true
}

func getHeaderIndex(key string) int {
	for i, h := range headers {
		if h == key {
			return i
		}
	}
	return -1
}
