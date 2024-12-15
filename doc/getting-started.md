# Getting Started with SQLD

## Installation

```bash
go get github.com/remiges-tech/sqld
```

## Basic Usage

1. First, define your model:

```go
type Employee struct {
    ID        int64  `db:"id" json:"id"`
    FirstName string `db:"first_name" json:"first_name"`
    LastName  string `db:"last_name" json:"last_name"`
    Email     string `db:"email" json:"email"`
}

// Implement the Model interface
func (e Employee) TableName() string {
    return "employees"
}
```

2. Register your model:

```go
if err := sqld.Register(Employee{}); err != nil {
    log.Fatal(err)
}
```

3. Execute queries:

```go
resp, err := sqld.Execute[Employee](ctx, db, sqld.QueryRequest{
    Select: []string{"id", "first_name", "email"},
    Where: map[string]interface{}{
        "is_active": true,
    },
    Pagination: &sqld.PaginationRequest{
        Page:     1,
        PageSize: 10,
    },
})
```

## Next Steps

- Check out the [Examples](./examples.md) for more usage patterns
- Read the [Architecture](./architecture.md) document for deeper understanding
