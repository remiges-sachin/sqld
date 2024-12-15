// Package sqld provides a type-safe, flexible SQL query builder and executor for dynamic database operations.
//
// sqld simplifies building and executing database queries where the fields and filters are not known at compile time,
// while maintaining type safety and secure scanning. It is particularly useful for building REST APIs that need to
// handle dynamic requests for field selection, filtering, and pagination.
//
// Key Features:
//   - Type-safe query building with compile-time checks
//   - Dynamic field selection and filtering
//   - Built-in pagination support
//   - Automatic null value handling
//   - Custom type scanner support
//   - Protection against SQL injection
//
// Architecture:
//
// The package follows a clean separation of concerns with four main components:
//   1. Registry: Manages model metadata and type scanners
//   2. Builder: Constructs type-safe SQL queries
//   3. Scanner: Handles type conversions and null values
//   4. Executor: Executes queries and maps results
//
// Basic Usage:
//
//	// Register your model
//	if err := sqld.Register(Employee{}); err != nil {
//	    log.Fatal(err)
//	}
//
//	// Execute a query
//	resp, err := sqld.Execute[Employee](ctx, db, sqld.QueryRequest{
//	    Select: []string{"id", "name", "email"},
//	    Where: map[string]interface{}{
//	        "is_active": true,
//	    },
//	})
//
// For more examples and detailed documentation, visit the examples directory
// or refer to the package documentation.
package sqld
