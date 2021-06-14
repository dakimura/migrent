// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// MigrationsColumns holds the columns for the "migrations" table.
	MigrationsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "applied_at", Type: field.TypeTime},
	}
	// MigrationsTable holds the schema information for the "migrations" table.
	MigrationsTable = &schema.Table{
		Name:        "migrations",
		Columns:     MigrationsColumns,
		PrimaryKey:  []*schema.Column{MigrationsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		MigrationsTable,
	}
)

func init() {
}