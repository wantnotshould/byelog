package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// VisitLog holds the schema definition for the VisitLog entity.
type VisitLog struct {
	ent.Schema
}

func (VisitLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "visit_log"},
		schema.Comment("Access log storage."),
		entsql.WithComments(true),
	}
}

// Fields of the VisitLog.
func (VisitLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("ip").
			MaxLen(45).
			Comment("Client IP address"),

		field.String("method").
			MaxLen(10).
			Comment("HTTP request method"),

		field.String("path").
			MaxLen(512).
			Comment("Request URL path"),

		field.Text("query").
			Optional().
			Comment("Query string"),

		field.String("title").
			MaxLen(255).
			Optional().
			Comment("Page title"),

		field.Text("referer").
			Optional().
			Comment("Referrer URL"),

		field.String("ua").
			MaxLen(512).
			Optional().
			Comment("User-Agent string"),

		field.Bool("is_bot").
			Default(false).
			Comment("Is bot/crawler"),

		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("Log timestamp"),
	}
}

func (VisitLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_at"),
		index.Fields("ip"),
	}
}

// Edges of the VisitLog.
func (VisitLog) Edges() []ent.Edge {
	return nil
}
