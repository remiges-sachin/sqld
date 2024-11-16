package sqld

// Query represents the user's request for data
type Query struct {
	Select []string               `json:"select"` // Fields to return
	From   string                 `json:"from"`   // Table name
	Where  map[string]interface{} `json:"where"`  // Simple key-value filters
}

// QueryResult represents the generic response structure
type QueryResult struct {
	Data  []map[string]interface{} `json:"data"`
	Error string                   `json:"error,omitempty"`
}
