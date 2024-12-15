# Type System in SQLD

SQLD provides a type system that ensures safety for dynamic queries. This document explains how SQLD's type system works and how to use it.

## Core Concepts

### 1. Query Types

SQLD supports two types of queries:

#### Structured Queries
```go
type Query struct {
    Select   []string          // Fields to select
    Where    map[string]string // WHERE conditions
    OrderBy  []OrderBy        // ORDER BY clauses
    Limit    *int             // LIMIT clause
    Offset   *int             // OFFSET clause
}
```

#### Raw Queries
```go
// Raw SQL with named parameters
query := "SELECT * FROM users WHERE status = :status AND age > :min_age"
```

### 2. Type Parameters

Both query types use two type parameters:

```go
// For structured queries
func Execute[P any, R any](ctx context.Context, db interface{}, q Query) ([]R, error)

// For raw queries
func ExecuteRaw[P any, R any](ctx context.Context, db interface{}, query string, params map[string]interface{}) ([]R, error)
```

Where:
- `P`: Parameter struct type (defines valid field names and types)
- `R`: Result struct type (defines how results are mapped)

## Usage Examples

### 1. Structured Queries

```go
// Define parameter and result types
type UserParams struct {
    Status string `db:"status"`
    Age    int    `db:"age"`
}

type UserResult struct {
    ID     int    `db:"id" json:"id"`
    Name   string `db:"name" json:"name"`
    Status string `db:"status" json:"status"`
    Age    int    `db:"age" json:"age"`
}

// Build and execute query
query := Query{
    Select: []string{"id", "name", "status", "age"},
    Where: map[string]string{
        "status": "active",
        "age": "21",
    },
}

results, err := Execute[UserParams, UserResult](ctx, db, query)
```

### 2. Raw Queries

```go
// Using the same types as above
params := map[string]interface{}{
    "status": "active",
    "min_age": 21,
}

results, err := ExecuteRaw[UserParams, UserResult](
    ctx, 
    db, 
    "SELECT * FROM users WHERE status = :status AND age > :min_age",
    params,
)
```

## Type Safety Features

### 1. Parameter Validation

- Validates field names against the parameter struct
- Ensures parameter types match struct field types
- Prevents SQL injection by validating field names
- Handles null values through pointer types

```go
type SearchParams struct {
    Status *string `db:"status"` // Optional parameter
    MinAge *int    `db:"min_age"` // Optional parameter
}
```

### 2. Result Mapping

- Maps database columns to struct fields using `db` tags
- Supports automatic JSON serialization via `json` tags
- Handles null values through pointer types
- Ignores unknown columns

```go
type UserResult struct {
    ID        int        `db:"id" json:"id"`
    Name      string     `db:"name" json:"name"`
    DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
```

## Error Handling

Common error cases:

1. Invalid field names in Query
```
field "unknown_field" not found in type UserParams
```

2. Type mismatches
```
field "age": expected int, got string
```

3. Missing required parameters
```
required parameter "status" not provided
```

## Best Practices

1. **Define Clear Types**
```go
// Good: Clear parameter definition
type UserParams struct {
    Status string `db:"status"`
    Age    int    `db:"age"`
}

// Bad: Using interface{} or generic maps
map[string]interface{}
```

2. **Use Optional Fields**
```go
// Good: Clear which fields are optional
type SearchParams struct {
    Status *string `db:"status"`
    MinAge *int    `db:"min_age"`
}
```

3. **Consistent Naming**
```go
// Good: DB fields match struct tags
type UserResult struct {
    ID        int    `db:"user_id" json:"id"`
    FirstName string `db:"first_name" json:"firstName"`
}
```

## Performance Notes

- Type information is cached after first use
- No reflection during query execution
- Zero allocation for type metadata
- Efficient parameter validation
