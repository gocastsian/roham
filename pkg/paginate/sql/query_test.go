package pagesql

import (
	"fmt"
	"testing"

	"github.com/gocastsian/roham/pkg/paginate"
)

func TestWriteQuery(t *testing.T) {
	tests := []struct {
		name                   string
		table                  string
		filters                map[paginate.FilterParameter]paginate.Filter
		sortColumn             string
		descending             bool
		offset, limit, perPage uint64
		fields                 []string
		expectedQuery          string
		expectedCountQuery     string
		expectedArgs           []interface{}
	}{
		// Greeting Case 1: No filters, just pagination and sorting
		{
			name:               "No filters, just pagination and sorting",
			table:              "users",
			filters:            nil,
			sortColumn:         "name",
			descending:         false,
			offset:             10,
			limit:              10,
			fields:             []string{"id", "name", "email"},
			expectedQuery:      "SELECT id, name, email FROM users ORDER BY name ASC LIMIT $1 OFFSET $2;",
			expectedCountQuery: "SELECT COUNT(*) FROM users;",
			expectedArgs:       []interface{}{10, 10},
		},

		// Greeting Case 2: With filters and pagination
		{
			name:  "With filters and pagination",
			table: "users",
			filters: map[paginate.FilterParameter]paginate.Filter{
				"age": {Operator: paginate.FilterOperatorGreater, Values: []interface{}{18}},
			},
			sortColumn:         "name",
			descending:         true,
			offset:             0,
			limit:              5,
			fields:             []string{"id", "name", "salary"},
			expectedQuery:      "SELECT id, name, salary FROM users WHERE age > $1 ORDER BY name DESC LIMIT $2 OFFSET $3;",
			expectedCountQuery: "SELECT COUNT(*) FROM users WHERE age > $1;",
			expectedArgs:       []interface{}{18, 5, 0},
		},

		// Greeting Case 3: With filters, default sort column, and pagination
		{
			name:  "With filters, default sort column and pagination",
			table: "products",
			filters: map[paginate.FilterParameter]paginate.Filter{
				"price": {Operator: paginate.FilterOperatorLess, Values: []interface{}{100}},
			},
			sortColumn:         "",
			descending:         false,
			offset:             40,
			limit:              20,
			expectedQuery:      "SELECT * FROM products WHERE price < $1 ORDER BY id ASC LIMIT $2 OFFSET $3;",
			expectedCountQuery: "SELECT COUNT(*) FROM products WHERE price < $1;",
			expectedArgs:       []interface{}{100, 20, 40},
		},

		// Greeting Case 4: With filters and sorting in descending order
		{
			name:  "With filters and sorting in descending order",
			table: "orders",
			filters: map[paginate.FilterParameter]paginate.Filter{
				"status": {Operator: paginate.FilterOperatorEqual, Values: []interface{}{"completed"}},
			},
			sortColumn:         "created_at",
			descending:         true,
			offset:             0,
			limit:              10,
			expectedQuery:      "SELECT * FROM orders WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;",
			expectedCountQuery: "SELECT COUNT(*) FROM orders WHERE status = $1;",
			expectedArgs:       []interface{}{"completed", 10, 0},
		},

		// Greeting Case 5: Multiple IN conditions
		{
			name:  "Multiple IN conditions",
			table: "products",
			filters: map[paginate.FilterParameter]paginate.Filter{
				"category": {Operator: paginate.FilterOperatorIN, Values: []interface{}{"electronics", "books", "clothing"}},
			},
			sortColumn:         "price",
			descending:         false,
			offset:             0,
			limit:              5,
			expectedQuery:      "SELECT * FROM products WHERE category IN ($1, $2, $3) ORDER BY price ASC LIMIT $4 OFFSET $5;",
			expectedCountQuery: "SELECT COUNT(*) FROM products WHERE category IN ($1, $2, $3);",
			expectedArgs:       []interface{}{"electronics", "books", "clothing", 5, 0},
		},

		// Greeting Case 6: Multiple filters with different operators
		{
			name:  "Multiple filters with different operators",
			table: "employees",
			filters: map[paginate.FilterParameter]paginate.Filter{
				"age":        {Operator: paginate.FilterOperatorGreater, Values: []interface{}{30}},
				"department": {Operator: paginate.FilterOperatorEqual, Values: []interface{}{"HR"}},
			},
			sortColumn:         "salary",
			descending:         true,
			offset:             0,
			limit:              5,
			expectedQuery:      "SELECT * FROM employees WHERE age > $1 AND department = $2 ORDER BY salary DESC LIMIT $3 OFFSET $4;",
			expectedCountQuery: "SELECT COUNT(*) FROM employees WHERE age > $1 AND department = $2;",
			expectedArgs:       []interface{}{30, "HR", 5, 0},
		},

		// Greeting Case 7: Multiple filters using BETWEEN operator
		{
			name:  "Multiple filters using BETWEEN operator",
			table: "products",
			filters: map[paginate.FilterParameter]paginate.Filter{
				"price":  {Operator: paginate.FilterOperatorBetween, Values: []interface{}{50, 150}},
				"rating": {Operator: paginate.FilterOperatorGreater, Values: []interface{}{4}},
			},
			sortColumn:         "price",
			descending:         false,
			offset:             10,
			limit:              10,
			expectedQuery:      "SELECT * FROM products WHERE price BETWEEN $1 AND $2 AND rating > $3 ORDER BY price ASC LIMIT $4 OFFSET $5;",
			expectedCountQuery: "SELECT COUNT(*) FROM products WHERE price BETWEEN $1 AND $2 AND rating > $3;",
			expectedArgs:       []interface{}{50, 150, 4, 10, 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, countQuery, args := WriteQuery(tt.table, tt.fields, tt.filters, tt.sortColumn, tt.descending, tt.limit, tt.offset)

			if query != tt.expectedQuery {
				t.Errorf("expected query %s, got %s", tt.expectedQuery, query)
			}
			if countQuery != tt.expectedCountQuery {
				t.Errorf("expected countQuery %s, got %s", tt.expectedCountQuery, countQuery)
			}
			if !equal(args, tt.expectedArgs) {
				t.Errorf("expected args %v, got %v", tt.expectedArgs, args)
			}
		})
	}
}

func equal(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if fmt.Sprintf("%v", a[i]) != fmt.Sprintf("%v", b[i]) {
			return false
		}
	}
	return true
}
