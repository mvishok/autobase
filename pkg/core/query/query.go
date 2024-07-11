package query

func Run(rows [][]string, query map[string][]string) []map[string]string {
	if len(rows) < 2 {
		return nil
	}

	headers := rows[0]
	var results []map[string]string

	for _, row := range rows[1:] {
		match := true
		record := make(map[string]string)
		for i, value := range row {
			record[headers[i]] = value
			for key, val := range query {
				if headers[i] == key && value != val[0] {
					match = false
					break
				}
			}
		}
		if match {
			results = append(results, record)
		}
	}

	return results
}
