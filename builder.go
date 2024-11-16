package sqld

import (
	"github.com/Masterminds/squirrel"
)

// BuildQuery converts a Query into a squirrel.SelectBuilder
func BuildQuery(q Query) (squirrel.SelectBuilder, error)

// private functions

// validateQuery checks if the query is valid
func validateQuery(q Query) error

// buildSelect processes the select fields
func buildSelect(fields []string) []string

// buildWhere processes the where conditions
func buildWhere(conditions map[string]interface{}) squirrel.Eq
