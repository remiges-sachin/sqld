# SQLD For Dynamic SQL

`sqld` is a package that enables dynamic SQL query generation and execution from JSON requests. It provides two distinct subsystems for handling different query needs:

1. **Structured Query System**: A high-level, type-safe abstraction for building queries using Go generics
2. **Safe Raw Query System**: A flexible system for writing raw SQL with safety guarantees

## Key Features

### Structured Query System
- Type-safe query building with Go generics
- Dynamic field selection and filtering
- Built-in pagination and ordering
- Automatic parameter binding

### Safe Raw Query System
- Raw SQL with named parameters
- Runtime type validation
- Safe parameter substitution
- Flexible query construction

## Usage

### Structured Query System
```go
// Register your model
if err := sqld.Register(Employee{}); err != nil {
    log.Fatal(err)
}

// Execute a structured query
resp, err := sqld.Execute[Employee](ctx, db, sqld.QueryRequest{
    Select: []string{"id", "name", "email"},
    Where: map[string]interface{}{
        "is_active": true,
    },
})
```

### Safe Raw Query System
```go
// Execute a raw query with named parameters
query := `
    SELECT id, name, email 
    FROM employees 
    WHERE department = {{dept}}
`
results, err := sqld.ExecuteRaw[EmployeeParams, Employee](
    ctx, 
    db, 
    query,
    map[string]interface{}{"dept": "Engineering"},
)
```

For more examples, check the `examples/` directory.

## Architecture

SQLD follows a clean design with components supporting both query systems:

### Core Components
1. Query Building
   - Structured query generation
   - Raw query processing
   - Type-safe parameter handling

2. Type System
   - Model metadata management
   - Custom type scanning
   - Null value handling

3. Execution
   - Query execution
   - Result mapping
   - Error handling

## Features
- Type-safe query building
- Automatic handling of null values
- Support for complex WHERE clauses
- ORDER BY support
- Custom type scanners
- Pagination support
- Named parameter support
- SQL injection prevention

## Documentation
Detailed documentation is available in the `doc/` directory:
- [Architecture Overview](doc/architecture.md)
- [Query Building Guide](doc/query-building.md)
