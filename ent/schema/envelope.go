package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

// Envelope holds the schema definition for the Envelope entity.
type Envelope struct {
	ent.Schema
}

// Fields of the Envelope.
func (Envelope) Fields() []ent.Field {
	return []ent.Field{
		field.String("to").NotEmpty(),
		field.String("from").NotEmpty(),
		field.String("subject").Default("no subject"),
		field.String("content").Default("no content"),
		field.Time("created_at").Default(time.Now),
	}
}

func (Envelope) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("to", "id"),
		index.Fields("created_at"),
	}
}

// Edges of the Envelope.
func (Envelope) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("attachments", Attachment.Type),
	}
}
