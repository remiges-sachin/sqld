package sqld

import (
	"encoding/json"
)

// Example of how to use the package
func Example() {
	queryJSON := `{
		"select": ["id", "name"],
		"from": "users",
		"where": {
			"status": "active"
		}
	}`

	var q Query
	if err := json.Unmarshal([]byte(queryJSON), &q); err != nil {
		// handle error
	}

	builder, err := BuildQuery(q)
	if err != nil {
		// handle error
	}

	sql, args, err := builder.ToSql()
	// use sql and args with database
}
