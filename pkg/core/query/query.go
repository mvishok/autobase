package query

import (
	"autobase/pkg/loader"
	"autobase/pkg/log"
	"net/url"
	"strings"
)

var headers []string
var indexMap = make(map[string][]string)
var primaryKey int = -1

type Response struct {
	Status string              `json:"status"`
	Count  int                 `json:"count,omitempty"`
	Data   []map[string]string `json:"data,omitempty"`
	Error  string              `json:"error,omitempty"`
}

func Run(csvPath string, query url.Values, access_level string) Response {
	rows, err := loader.ReadCSV(csvPath)
	if err != nil {
		log.Error("An error occurred while reading the CSV file: " + err.Error())
		return Response{
			Status: "error",
			Error:  "Internal server error",
		}
	}
	if len(rows) < 1 {
		log.Error("No columns defined in CSV file")
		return Response{
			Status: "error",
			Error:  "Internal server error",
		}
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
		return Response{
			Status: "success",
			Count:  len(results) - 1,
			Data:   results,
		}
	}

	//primary key detection
	for i := range rows[0] {
		if strings.HasPrefix(rows[0][i], "$") {
			rows[0][i] = strings.TrimPrefix(rows[0][i], "$")
			primaryKey = i
			break
		}
	}

	//if primary key is defined, update indexMap
	if primaryKey != -1 {
		for i := range rows {
			if i == 0 {
				continue
			}
			//append the row to the indexMap with the primary key as the key
			indexMap[rows[i][primaryKey]] = rows[i]
		}
	}

	headers = rows[0]

	//if select query, return only selected columns
	if query.Get("select") != "" {
		if access_level != "read" && access_level != "" {
			return Response{
				Status: "ClientError",
				Error:  "Unauthorized",
			}
		}

		cols := strings.Split(query.Get("select"), ",")

		for _, col := range cols {
			if col == "*" {
				cols = append(rows[0], cols...)
				break
			}
		}

		//if where clause is defined and contains only primary key, use whereClause with that row from indexMap
		if query.Get("where") != "" {
			for key, value := range query {
				if key == "where" {
					value := strings.Split(value[0], ":")
					if len(value) > 2 && primaryKey != -1 && value[0] == rows[0][primaryKey] {
						row, ok := indexMap[value[2]]
						if ok {
							if whereClause(row, query) {
								result := make(map[string]string)
								for i, value := range row {
									if contains(cols, rows[0][i]) {
										result[rows[0][i]] = value
									}
								}
								return Response{
									Status: "success",
									Count:  1,
									Data:   []map[string]string{result},
								}
							}
						}
					}
				}
			}
		}

		for i, row := range rows {
			if i == 0 {
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

		return Response{
			Status: "success",
			Count:  len(results),
			Data:   results,
		}
	}

	//if update query, update the rows
	if query.Get("update") != "" {
		if access_level != "write" && access_level != "" {
			return Response{
				Status: "ClientError",
				Error:  "Unauthorized",
			}
		}
		cols := strings.Split(query.Get("update"), ",")
		for i, row := range rows {
			if i == 0 {
				continue
			}
			if whereClause(row, query) {

				//syntax: ?update=key:newValue,key2:newValue2...
				for _, col := range cols {
					parts := strings.Split(col, ":")
					if len(parts) != 2 {
						log.Error("Invalid update syntax: " + col)
						return Response{
							Status: "ClientError",
							Error:  "Invalid update syntax",
						}
					}
					key := parts[0]
					newValue := parts[1]
					index := getHeaderIndex(key)
					if index == -1 {
						return Response{
							Status: "ClientError",
							Error:  "Invalid column name",
						}
					}
					rows[i][index] = newValue

					result := make(map[string]string)
					for i, value := range row {
						result[rows[0][i]] = value
					}
					results = append(results, result)
				}
			}
		}
		loader.UpdateCSV(csvPath, rows)
		return Response{
			Status: "success",
			Count:  len(rows) - len(results),
			Data:   results,
		}
	}

	//if insert query, insert a new row
	if query.Get("insert") != "" {
		if access_level != "write" && access_level != "" {
			return Response{
				Status: "ClientError",
				Error:  "Unauthorized",
			}
		}
		newRow := strings.Split(query.Get("insert"), ",")
		if len(newRow) != len(rows[0]) {
			log.Error("Invalid number of columns")
			return Response{
				Status: "ClientError",
				Error:  "Invalid number of columns",
			}
		}
		rows = append(rows, newRow)
		loader.UpdateCSV(csvPath, rows)
		return Response{
			Status: "success",
			Count:  1,
		}
	}

	//if delete query, delete the rows
	if query.Get("delete") == "true" {
		if access_level != "write" && access_level != "" {
			return Response{
				Status: "ClientError",
				Error:  "Unauthorized",
			}
		}

		var newRows [][]string
		for i, row := range rows {
			if i == 0 {
				newRows = append(newRows, row) // Always keep the header
				continue
			}
			if !whereClause(row, query) {
				newRows = append(newRows, row)
			}
		}
		loader.UpdateCSV(csvPath, newRows)
		return Response{
			Status: "success",
			Count:  len(rows) - len(newRows),
		}
	}

	//if none of the above, return error
	return Response{
		Status: "ClientError",
		Error:  "Invalid query",
	}
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
						log.Warning("Invalid column name: " + key)
						return false
					}

					switch op {
					case "eq":
						if row[index] != val {
							return false
						}
					case "ne":
						if row[index] == val {
							return false
						}
					case "gt":
						if row[index] <= val {
							return false
						}
					case "lt":
						if row[index] >= val {
							return false
						}
					case "ge":
						if row[index] < val {
							return false
						}
					case "le":
						if row[index] > val {
							return false
						}
					default:
						log.Warning("Invalid operator: " + op)
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
