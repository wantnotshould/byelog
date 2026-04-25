package schema

import "entgo.io/ent"

// VisitLog holds the schema definition for the VisitLog entity.
type VisitLog struct {
	ent.Schema
}

// Fields of the VisitLog.
func (VisitLog) Fields() []ent.Field {
	return nil
}

// Edges of the VisitLog.
func (VisitLog) Edges() []ent.Edge {
	return nil
}
