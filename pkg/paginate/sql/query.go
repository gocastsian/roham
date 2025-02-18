package pagesql

import (
	"fmt"
	"strings"

	"roham/pkg/paginate"
)

var (
	DefaultSortColumn = "id"
)

// WriteQuery generates a SQL query for paginated results and a count query based on the provided filters, sorting, and pagination parameters.
// It constructs a SELECT query with conditions for filtering, ordering, and pagination (LIMIT and OFFSET), and a COUNT query for the total number of records.
func WriteQuery(table string, fields []string, filters map[paginate.FilterParameter]paginate.Filter, sortColumn string, descending bool, limit, offset uint64) (query string, countQuery string, args []interface{}) {

	selectFields := "*"
	if len(fields) > 0 {
		selectFields = strings.Join(fields, ", ")
	}

	// Base query for pagination
	query = fmt.Sprintf("SELECT %s FROM %s", selectFields, table)

	// Base query for total count
	countQuery = "SELECT COUNT(*) FROM " + table

	// Add filters to the query
	if len(filters) == 0 {
		// Use DefaultSortColumn if sortColumn is empty
		if sortColumn == "" {
			sortColumn = DefaultSortColumn
		}
		orderClause := fmt.Sprintf(" ORDER BY %s %s", sortColumn, orderDirection(descending))

		// Set the pagination arguments
		args = []interface{}{limit, offset}
		return query + orderClause + " LIMIT $1 OFFSET $2;", countQuery + ";", args
	}

	query += " WHERE "
	conditions := []string{}
	paramIndex := len(args) + 1 // Start parameter numbering after pagination arguments

	for p, f := range filters {
		switch f.Operator {
		case paginate.FilterOperatorEqual:
			conditions = append(conditions, fmt.Sprintf("%s = $%d", p, paramIndex))
			args = append(args, f.Values[0])
			paramIndex++
		case paginate.FilterOpratorNotEqual:
			conditions = append(conditions, fmt.Sprintf("%s != $%d", p, paramIndex))
			args = append(args, f.Values[0])
			paramIndex++
		case paginate.FilterOperatorGreater:
			conditions = append(conditions, fmt.Sprintf("%s > $%d", p, paramIndex))
			args = append(args, f.Values[0])
			paramIndex++
		case paginate.FilterOperatorGreaterEqual:
			conditions = append(conditions, fmt.Sprintf("%s >= $%d", p, paramIndex))
			args = append(args, f.Values[0])
			paramIndex++
		case paginate.FilterOperatorLess:
			conditions = append(conditions, fmt.Sprintf("%s < $%d", p, paramIndex))
			args = append(args, f.Values[0])
			paramIndex++
		case paginate.FilterOperatorLessEqual:
			conditions = append(conditions, fmt.Sprintf("%s <= $%d", p, paramIndex))
			args = append(args, f.Values[0])
			paramIndex++
		case paginate.FilterOperatorIN:
			placeholders := []string{}
			for _, value := range f.Values {
				placeholders = append(placeholders, fmt.Sprintf("$%d", paramIndex))
				args = append(args, value)
				paramIndex++
			}
			conditions = append(conditions, fmt.Sprintf("%s IN (%s)", p, strings.Join(placeholders, ", ")))
		case paginate.FilterOperatorNotIn:
			placeholders := []string{}
			for _, value := range f.Values {
				placeholders = append(placeholders, fmt.Sprintf("$%d", paramIndex))
				args = append(args, value)
				paramIndex++
			}
			conditions = append(conditions, fmt.Sprintf("%s NOT IN (%s)", p, strings.Join(placeholders, ", ")))
		case paginate.FilterOperatorBetween:
			conditions = append(conditions, fmt.Sprintf("%s BETWEEN $%d AND $%d", p, paramIndex, paramIndex+1))
			args = append(args, f.Values[0], f.Values[1])
			paramIndex += 2
		}
	}

	// Add the conditions to both the query and countQuery
	query += strings.Join(conditions, " AND ")

	// Use "ID" as the default column if sortColumn is empty
	if sortColumn == "" {
		sortColumn = DefaultSortColumn
	}

	// Add sorting and pagination
	query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, orderDirection(descending))
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d;", paramIndex, paramIndex+1)
	countQuery += " WHERE " + strings.Join(conditions, " AND ") + ";"

	// Add pagination arguments
	args = append(args, limit, offset)

	return query, countQuery, args
}

func orderDirection(descending bool) string {
	if descending {
		return "DESC"
	}
	return "ASC"
}
