// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// EnvelopesColumns holds the columns for the "envelopes" table.
	EnvelopesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "to", Type: field.TypeString},
		{Name: "from", Type: field.TypeString},
		{Name: "subject", Type: field.TypeString},
		{Name: "content", Type: field.TypeString},
		{Name: "created_at", Type: field.TypeTime},
	}
	// EnvelopesTable holds the schema information for the "envelopes" table.
	EnvelopesTable = &schema.Table{
		Name:       "envelopes",
		Columns:    EnvelopesColumns,
		PrimaryKey: []*schema.Column{EnvelopesColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "envelope_to",
				Unique:  false,
				Columns: []*schema.Column{EnvelopesColumns[1]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		EnvelopesTable,
	}
)

func init() {
}
