package sqld

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/jackc/pgx/v5"
)

type fieldInfo struct {
	jsonKey string
	goType  reflect.Type
	fieldName string
}

// BuildMetadataMap uses reflection on the model struct to map db tags to fieldInfo.
// It extracts the 'db' and 'json' tags from the struct fields and creates a map
// where the key is the 'db' tag and the value is a fieldInfo struct containing the
// 'json' tag and the Go type of the field. This map is used later in the ExecuteRaw
// function to map database column names to JSON keys in the result.
func BuildMetadataMap[T any]() (map[string]fieldInfo, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("model must be a struct")
	}

	metaMap := make(map[string]fieldInfo)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		jsonTag := field.Tag.Get("json")
		if dbTag != "" && jsonTag != "" {
			metaMap[dbTag] = fieldInfo{
				jsonKey: jsonTag,
				goType:  field.Type,
				fieldName: field.Name,
			}
		}
	}
	return metaMap, nil
}

// isTypeCompatible checks if the runtime type of a value matches the expected type.
// It returns true if the value's type is compatible with the expected type,
// and false otherwise. It also handles the case where the expected type is an
// empty interface, in which case any type is considered compatible.
func isTypeCompatible(valType, expectedType reflect.Type) bool {
	if valType == nil || expectedType == nil {
		return false
	}

	// If the expected type is an empty interface, accept any type.
	if expectedType.Kind() == reflect.Interface && expectedType.NumMethod() == 0 {
		// This means expectedType is `interface{}`
		return true
	}

	return valType == expectedType
}

func typeNameOrNil(t reflect.Type) string {
	if t == nil {
		return "nil"
	}
	return t.String()
}

// Named parameter regex to find patterns like {{param_name}}
var namedParamRegex = regexp.MustCompile(`\{\{([a-zA-Z0-9_]+)\}\}`)

// ExtractNamedPlaceholders finds all named parameters in the {{param_name}} format.
func ExtractNamedPlaceholders(query string) ([]string, error) {
	matches := namedParamRegex.FindAllStringSubmatch(query, -1)
	var params []string
	seen := make(map[string]bool)
	for _, match := range matches {
		paramName := match[1]
		if !seen[paramName] {
			seen[paramName] = true
			params = append(params, paramName)
		}
	}
	return params, nil
}

// ReplaceNamedWithDollarPlaceholders replaces {{param_name}} with $1, $2, ...
func ReplaceNamedWithDollarPlaceholders(query string, queryParams []string) (string, error) {
	for i, p := range queryParams {
		placeholder := fmt.Sprintf("{{%s}}", p)
		newPlaceholder := fmt.Sprintf("$%d", i+1)
		query = strings.ReplaceAll(query, placeholder, newPlaceholder)
	}
	return query, nil
}

// ValidateMapParamsAgainstStructNamed ensures the params map matches the expected types from P.
// It uses the isTypeCompatible function to check if the type of each parameter in the map
// matches the expected type from P. This is primarily to prevent runtime errors due to type mismatches.
func ValidateMapParamsAgainstStructNamed[P any](
	paramMap map[string]interface{},
	queryParams []string,
) ([]interface{}, error) {
	t := reflect.TypeOf((*P)(nil)).Elem()
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("model must be a struct")
	}

	typeByName := make(map[string]reflect.Type)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		jsonTag := field.Tag.Get("json")
		
		// Validate that all fields with db tag must have json tag
		if dbTag != "" && jsonTag == "" {
			return nil, fmt.Errorf("field %s has db tag but missing json tag", field.Name)
		}
		
		if dbTag != "" {
			typeByName[dbTag] = field.Type
		}
	}

	args := make([]interface{}, 0, len(queryParams))
	for _, p := range queryParams {
		expectedType, found := typeByName[p]
		if !found {
			return nil, fmt.Errorf("no type info for param %s", p)
		}

		val, present := paramMap[p]
		if !present {
			// If the parameter is optional and not present, append nil or handle as needed
			args = append(args, nil)
			continue
		}

		valType := reflect.TypeOf(val)
		if !isTypeCompatible(valType, expectedType) {
			return nil, fmt.Errorf("parameter %s type mismatch: got %s, want %s",
				p, typeNameOrNil(valType), typeNameOrNil(expectedType))
		}

		args = append(args, val)
	}

	return args, nil
}

// ExecuteRaw takes a query with {{param_name}} placeholders and executes it.
// P is the type that defines parameter structure (with `db` tags)
// R is the type that defines result structure (with `db` and `json` tags)
func ExecuteRaw[P, R any](
	ctx context.Context,
	db interface{},
	query string,
	params map[string]interface{},
) ([]map[string]interface{}, error) {
	// 1. Extract named placeholders
	queryParams, err := ExtractNamedPlaceholders(query)
	if err != nil {
		return nil, fmt.Errorf("failed to extract named placeholders: %w", err)
	}

	// 2. Validate and convert map params to arguments in correct order
	args, err := ValidateMapParamsAgainstStructNamed[P](params, queryParams)
	if err != nil {
		return nil, fmt.Errorf("parameter validation failed: %w", err)
	}

	// 3. Replace named placeholders with $N placeholders
	finalQuery, err := ReplaceNamedWithDollarPlaceholders(query, queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to replace named placeholders: %w", err)
	}

	// 4. Build metadata map for results (no instance needed)
	metaMap, err := BuildMetadataMap[R]()
	if err != nil {
		return nil, fmt.Errorf("failed to build metadata map: %w", err)
	}

	// 5. Execute query and scan into slice of structs first to handle custom types
	var structResults []R
	switch db := db.(type) {
	case *sql.DB:
		if err := sqlscan.Select(ctx, db, &structResults, finalQuery, args...); err != nil {
			return nil, fmt.Errorf("failed to execute query: %w", err)
		}
	case *pgx.Conn:
		if err := pgxscan.Select(ctx, db, &structResults, finalQuery, args...); err != nil {
			return nil, fmt.Errorf("failed to execute query: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported database type: %T", db)
	}

	// 6. Convert struct results to maps with only requested fields
	results := make([]map[string]interface{}, len(structResults))
	for i, row := range structResults {
		val := reflect.ValueOf(row)
		typ := val.Type()
		resultMap := make(map[string]interface{})

		// Only include fields that were in the original query's SELECT clause
		for _, info := range metaMap {
			if field, ok := typ.FieldByName(info.fieldName); ok {
				fieldVal := val.FieldByName(field.Name)
				if fieldVal.IsValid() {
					resultMap[info.jsonKey] = fieldVal.Interface()
				}
			}
		}
		results[i] = resultMap
	}

	return results, nil
}
