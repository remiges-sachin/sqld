# Query Building in SQLD

## Overview

SQLD provides two complementary approaches to building queries, each designed for different use cases:

1. **Structured Query System**: A high-level abstraction with runtime type validation and generic result types
2. **Raw Query System**: A flexible system for type-safe raw SQL queries with named parameters

## Structured Query System

### Overview
The Structured Query System provides a type-safe way to build and execute dynamic SQL queries using Go's type system.

### Request Structure
The `QueryRequest` struct is the main input for structured queries:

```go
type QueryRequest struct {
    // Required: Fields to retrieve
    Select []string

    // Optional: Filter conditions as key-value pairs
    Where map[string]interface{}

    // Optional: Sorting criteria
    OrderBy []OrderByClause

    // Optional: Page-based pagination (takes precedence over Limit/Offset)
    Pagination *PaginationRequest

    // Optional: Direct limit/offset controls
    Limit  *int
    Offset *int
}
```

#### Components:

1. **Select Fields** (Required)
   - List of field names to retrieve
   - Must match JSON field names from your model's struct tags
   - Cannot be empty
   - Each field is validated against model metadata
   ```go
   Select: []string{"id", "first_name", "email"}
   ```

2. **Where Conditions** (Optional)
   - Key-value pairs for filtering
   - Keys must match JSON field names
   - Values are type-checked against model field types
   - Each field is validated against model metadata
   ```go
   Where: map[string]interface{}{
       "is_active": true,
       "department": "Engineering",
   }
   ```

3. **Order By** (Optional)
   - Specify sorting criteria for multiple fields
   - Each clause contains field name and sort direction
   - Field names are validated against model metadata
   ```go
   OrderBy: []OrderByClause{
       {Field: "first_name", Desc: false},  // Ascending
       {Field: "hire_date", Desc: true},    // Descending
   }
   ```

4. **Pagination** (Optional)
   Two options available:
   
   a. Page-based pagination (recommended):
   - Takes precedence over direct limit/offset
   - Page numbers start at 1 (not 0)
     - Page 1 returns the first page of results
     - Page 2 returns the second page of results
     - And so on...
   - Default page size: 10 (DefaultPageSize)
   - Maximum page size: 100 (MaxPageSize)
   ```go
   // Get first page
   Pagination: &PaginationRequest{
       Page: 1,     // First page of results
       PageSize: 10, // 10 items per page
   }

   // Get second page
   Pagination: &PaginationRequest{
       Page: 2,     // Second page of results
       PageSize: 10, // 10 items per page
   }
   ```
   
   b. Direct limit/offset:
   - Only used if Pagination is not provided
   - Both values must be non-negative
   ```go
   // Get first 10 results
   Limit: &limit,   // e.g., limit := 10
   Offset: &offset, // e.g., offset := 0

   // Get next 10 results
   Limit: &limit,   // e.g., limit := 10
   Offset: &offset, // e.g., offset := 10
   ```

### Validation Rules
The system performs the following validations:
1. Select fields cannot be empty
2. All field names (in Select, Where, OrderBy) must exist in the model
3. Where clause values must match field types
4. Pagination:
   - Page numbers start at 1
   - Page size is capped at 100 (MaxPageSize)
   - Defaults to 10 items per page (DefaultPageSize)
5. Limit and Offset must be non-negative

### Key Features
1. Type Safety
   - Runtime field and type validation
   - Generic result types
   - Model-based metadata validation

2. Query Building
   - Structured query requests
   - Automatic field mapping
   - Built-in validation

3. Features
   - Dynamic field selection
   - WHERE clause support
   - ORDER BY support
   - Pagination
   - Automatic parameter binding

### Usage

#### Basic Query
```go
// Define your model
type Employee struct {
    ID         int64     `db:"id" json:"id"`
    FirstName  string    `db:"first_name" json:"first_name"`
    LastName   string    `db:"last_name" json:"last_name"`
    Email      string    `db:"email" json:"email"`
    Department string    `db:"department" json:"department"`
    IsActive   bool      `db:"is_active" json:"is_active"`
}

// Execute a query
resp, err := sqld.Execute[Employee](ctx, db, sqld.QueryRequest{
    Select: []string{"id", "first_name", "last_name", "email", "department"},
    Where: map[string]interface{}{
        "is_active": true,
        "department": "Engineering",
    },
    OrderBy: []sqld.OrderByClause{
        {Field: "first_name", Desc: false},
    },
})
```

#### With Pagination
```go
resp, err := sqld.Execute[Employee](ctx, db, sqld.QueryRequest{
    Select: []string{"id", "first_name", "last_name", "email"},
    Pagination: &sqld.PaginationRequest{
        Page: 1,
        PageSize: 10,
    },
})
```

## Raw Query System

### Overview
The Raw Query System allows writing raw SQL queries while maintaining full type safety through runtime validation and named parameters.

### Key Features
1. Type Safety
   - Parameter type validation against reference structs
   - Field mapping validation through struct tags
   - Automatic parameter binding

2. Query Building
   - Raw SQL with named parameters
   - Automatic parameter substitution
   - Type checking against reference structs

3. Features
   - Named parameter support using {{param}}
   - Full type validation against structs
   - Flexible SQL construction
   - Automatic field mapping

### Usage

#### Basic Query
```go
// Define parameter types
type QueryParams struct {
    DepartmentID int64  `db:"department_id"`
    MinSalary    int64  `db:"min_salary"`
}

// Define result type for type checking
type Result struct {
    ID        int64   `db:"id" json:"id"`
    FirstName string  `db:"first_name" json:"first_name"`
    Salary    float64 `db:"salary" json:"salary"`
}

// SQL query with type-safe parameters
query := `
    SELECT id, first_name, salary
    FROM employees
    WHERE department_id = {{department_id}}
    AND salary >= {{min_salary}}
    ORDER BY salary DESC;
`

params := map[string]interface{}{
    "department_id": 101,
    "min_salary":    50000,
}

// Using generic version - type safe parameters but flexible results
results, err := sqld.ExecuteRaw[QueryParams, Result](ctx, db, query, params)
if err != nil {
    log.Fatal(err)
}

// Results are maps with exactly what was selected
for _, emp := range results {
    fmt.Printf("Employee %v earns $%v\n", 
        emp["first_name"], 
        emp["salary"],
    )
}
```

## Safety Features

1. SQL Injection Prevention
   - Parameterized queries using $1, $2, etc.
   - Parameter binding through database driver
   - Type validation of parameters
   - No direct string concatenation or manual escaping

2. Type Safety
   - Runtime validation in both systems
     - Validates field existence
     - Checks value types against database schema
     - Ensures parameter types match before query execution
   
   - Model-based validation in Structured System
     - Uses Go generics for compile-time type checking
     - Validates all fields exist in model struct
     - Ensures query operations match field types
   
   - Full type checking in Raw System
     - Uses reference structs to validate parameter types
     - Ensures result columns map to expected types


## Choosing Between Systems

- Use Structured System for dynamic queries on flat tables -- if appropriate create views for flattening complex relationships
- Use Raw System for complex queries or specific SQL features


## Error Handling

Both systems provide detailed error messages for:
- Invalid field names
- Type mismatches
- Missing parameters
- SQL syntax errors
- Execution errors
