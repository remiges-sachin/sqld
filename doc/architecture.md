# SQLD Architecture

## Overview

SQLD is a Go package that provides type-safe, dynamic SQL query capabilities through a set of components to handle query building, execution, and result mapping. It offers two complementary subsystems for different use cases:

1. **Structured Query System**: A high-level abstraction with runtime type validation and generic result types
2. **Raw Query System**: A flexible system for type-safe raw SQL queries with named parameters

## Query Subsystems

SQLD provides two subsystems for handling queries:

### 1. Structured Query System
- High-level, runtime type-safe query building
- Uses structured QueryRequest structure
- Runtime field and type validation
- Generic result types for type-safe scanning
- Built-in support for WHERE, ORDER BY, and pagination
- Best for standard CRUD operations and API endpoints

### 2. Raw Query System
- Works with hand-written SQL queries
- Named parameter support with {{param}} syntax
- Runtime type validation against reference structs
- Automatic field mapping through struct tags
- Best for complex queries and joins
- Type safety through parameter validation

## Core Components

The package consists of the following components: 

### Query Building

#### Structured Query System
- Query construction through structured requests
- Runtime field and type validation
- Generic result types for type-safe scanning
- Built-in support for WHERE, ORDER BY, and pagination

#### Raw Query System
- Raw SQL query support with named parameters ({{param}} syntax)
- Runtime type validation against reference structs
- Automatic parameter substitution and validation
- Flexible SQL construction with safety checks

### Type System
- Model interface and metadata management
- Runtime type validation and checking
- Null value handling
- JSON field mapping

### Query Execution
- Type-safe query execution with context support
- Result mapping and scanning
- Database type abstraction (supports sql.DB and pgx)
- Error handling and reporting

### Additional Features
- Pagination support with offset/limit and page-based options
- Model registry for type information
- Validation system for query parameters
- Type-safe named parameter handling

## Flow Diagram

### Structured Query System
```
[QueryRequest] -> Runtime Validation -> Query Building -> Execution -> Result Mapping -> [Response]
```

### Raw Query System
```
[Raw SQL + Params] -> Parameter Extraction -> Runtime Type Validation -> Parameter Binding -> Execution -> [Response]
```

## Design Decisions

1. Type Safety Approach:
   - Runtime validation for query parameters and fields
   - Generic result types for type-safe scanning
   - Model-based metadata validation

2. Dual System Design:
   - Structured System for standard queries with runtime safety
   - Raw System for flexible queries with type-safe parameters

3. Integration:
   - Components work together seamlessly
   - Shared type system and execution layer

4. Safety:
   - Built-in protection against SQL injection through named parameters
   - Runtime type validation for all parameters

5. Flexibility:
   - Choice between structured and raw query APIs
   - Support for both simple and complex queries
   - Named parameters for better readability and safety
