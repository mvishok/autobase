package query

import (
	"net/url"
	"sort"
	"strconv"
	"strings"
	"syncengin/pkg/log"
)

func Run(rows [][]string, query url.Values, access_level string) []map[string]string {
	if len(rows) < 2 {
		return nil
	}

	headers := rows[0]
	results := []map[string]string{}

	for _, row := range rows[1:] {
		record := make(map[string]string)
		for i, value := range row {
			record[headers[i]] = value
		}
		if rowMatchesQuery(record, query) {
			results = append(results, record)
		}
	}

	//if it is select query
	if query.Get("select") != "" {
		if access_level == "read" || access_level == "write" {
			results = selectFields(results, query.Get("select"))
		}
	} else {
		//return empty
		log.Error("Unkown query type")
		return []map[string]string{}
	}

	results = orderBy(results, query.Get("order"))

	return results
}

func selectFields(rows []map[string]string, fields string) []map[string]string {
	if fields == "" {
		return rows
	}

	selectedFields := strings.Split(fields, ",")
	var selected []map[string]string

	for _, row := range rows {
		selectedRow := make(map[string]string)
		for _, field := range selectedFields {
			if value, ok := row[field]; ok {
				selectedRow[field] = value
			}
		}
		selected = append(selected, selectedRow)
	}

	return selected
}

func rowMatchesQuery(record map[string]string, query url.Values) bool {
	where := query.Get("where")
	if where == "" {
		return true
	}

	conditions := strings.Split(where, ",")
	for _, condition := range conditions {
		parts := strings.SplitN(condition, " ", 3)
		if len(parts) != 3 {
			return false
		}
		field, operator, value := parts[0], parts[1], parts[2]

		recordValue := record[field]

		switch operator {
		case "eq":
			if recordValue != value {
				return false
			}
		case "neq":
			if recordValue == value {
				return false
			}
		case "gt":
			recordFloat, err := strconv.ParseFloat(recordValue, 64)
			queryFloat, err2 := strconv.ParseFloat(value, 64)
			if err != nil || err2 != nil || recordFloat <= queryFloat {
				return false
			}

		case "gte":
			recordFloat, err := strconv.ParseFloat(recordValue, 64)
			queryFloat, err2 := strconv.ParseFloat(value, 64)
			if err != nil || err2 != nil || recordFloat < queryFloat {
				return false
			}

		case "lt":
			recordFloat, err := strconv.ParseFloat(recordValue, 64)
			queryFloat, err2 := strconv.ParseFloat(value, 64)
			if err != nil || err2 != nil || recordFloat >= queryFloat {
				return false
			}

		case "lte":
			recordFloat, err := strconv.ParseFloat(recordValue, 64)
			queryFloat, err2 := strconv.ParseFloat(value, 64)
			if err != nil || err2 != nil || recordFloat > queryFloat {
				return false
			}

		case "contains":
			if !strings.Contains(recordValue, value) {
				return false
			}

		default:
			return false
		}
	}
	return true
}

func orderBy(rows []map[string]string, order string) []map[string]string {
	if order == "" {
		return rows
	}

	orderFields := strings.Split(order, ",")
	sort.SliceStable(rows, func(i, j int) bool {
		for _, field := range orderFields {
			parts := strings.Split(field, " ")
			fieldName := parts[0]
			desc := false
			if len(parts) > 1 && parts[1] == "desc" {
				desc = true
			}

			if rows[i][fieldName] == rows[j][fieldName] {
				continue
			}

			if desc {
				return rows[i][fieldName] > rows[j][fieldName]
			}
			return rows[i][fieldName] < rows[j][fieldName]
		}
		return false
	})

	return rows
}
