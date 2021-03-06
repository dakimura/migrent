// Code generated by entc, DO NOT EDIT.

package migration

import (
	"time"
)

const (
	// Label holds the string label denoting the migration type in the database.
	Label = "migration"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldAppliedAt holds the string denoting the applied_at field in the database.
	FieldAppliedAt = "applied_at"
	// Table holds the table name of the migration in the database.
	Table = "migrations"
)

// Columns holds all SQL columns for migration fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldAppliedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultAppliedAt holds the default value on creation for the "applied_at" field.
	DefaultAppliedAt func() time.Time
	// UpdateDefaultAppliedAt holds the default value on update for the "applied_at" field.
	UpdateDefaultAppliedAt func() time.Time
)
