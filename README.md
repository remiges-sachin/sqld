# sqld - Dynamic SQL Query Builder

`sqld` is a lightweight Go package that enables dynamic SQL query generation from JSON requests. It simplifies building flexible database queries where the fields and filters are not known at compile time.

## Overview

Modern APIs often need to be flexible in how they return data, allowing clients to:
- Request specific fields
- Apply dynamic filters
- Control data shape

`sqld` addresses these needs by providing a simple JSON interface that maps to SQL queries.

## Core Features

### Dynamic Field Selection
Clients can specify which fields they want to retrieve:

```json
{
	"select": ["id", "name", "email"],
	"from": "users"
}
```



### Dynamic Filtering
Apply filters using a simple key-value structure:

```json
{
    "select": ["id", "name"],
    "from": "users",
    "where": {
        "status": "active",
    "role": "admin"
    }
}
```


## Design Philosophy

- **Simplicity First**: Start with basic operations and grow based on needs
- **Security Focused**: Built-in validation and sanitization
- **Type Safety**: Leverages Go's type system for reliable query building
- **Extensible**: Designed to grow with additional features
- **Performance**: Efficient query generation using [Squirrel](https://github.com/Masterminds/squirrel)

## Use Cases

- REST APIs with flexible response fields
- GraphQL-like query capabilities without GraphQL complexity
- Admin interfaces with dynamic data requirements
- Report generators with configurable outputs

## Current Limitations

- Supports only equality comparisons in WHERE clause
- Single table queries only
- No support for JOINs yet
- No aggregations or GROUP BY
- No sorting or pagination

## Roadmap

Future versions may include:
- Additional comparison operators (>, <, LIKE, etc.)
- Multi-table JOIN support
- ORDER BY and GROUP BY
- Pagination
- Field and table aliases
- Complex WHERE conditions (AND, OR combinations)
- Aggregation functions
- Query optimization hints

## Dependencies

- [Squirrel](https://github.com/Masterminds/squirrel) for SQL query building

## Contributing

This package is designed to grow based on real-world usage patterns. Contributions are welcome, especially in the following areas:
- Additional operators
- Performance improvements
- Security enhancements
- Documentation
- Test cases

## License

