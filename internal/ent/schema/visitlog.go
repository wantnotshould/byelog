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

		field.String("query").
			MaxLen(1024).
			Optional().
			Comment("Query string"),

		field.String("title").
			MaxLen(255).
			Optional().
			Comment("Page title"),

		field.String("referer").
			MaxLen(1024).
			Optional().
			Comment("Referrer URL"),

		field.String("os").
			MaxLen(128).
			Optional().
			Comment("OS name"),

		field.String("browser").
			MaxLen(128).
			Optional().
			Comment("Browser name"),

		field.String("browser_version").
			MaxLen(64).
			Optional().
			Comment("Full version"),

		field.String("device_type").
			MaxLen(32).
			Optional().
			Comment("General category"),

		field.String("device_model").
			MaxLen(128).
			Optional().
			Comment("Hardware info"),

		field.String("engine").
			MaxLen(64).
			Optional().
			Comment("Rendering engine"),

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
		index.Fields("ip"),
		index.Fields("created_at"),

		index.Fields("created_at", "path").
			Annotations(entsql.PrefixColumn("path", 128)),

		index.Fields("created_at", "referer").
			Annotations(entsql.PrefixColumn("referer", 128)),

		index.Fields("created_at", "device_type", "is_bot"),

		index.Fields("ip", "created_at"),
	}
}

// Edges of the VisitLog.
func (VisitLog) Edges() []ent.Edge {
	return nil
}
