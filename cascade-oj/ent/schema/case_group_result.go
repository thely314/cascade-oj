package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// a group of case results
type CaseGroupResult struct {
	ent.Schema
}

func (CaseGroupResult) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Immutable(),
		// field.Enum("status").Values(
		// 	"pending", "judging", "accepted", "wrong_answer",
		// 	"time_limit_exceeded", "memory_limit_exceeded", "runtime_error", "compilation_error", "system_error",
		// ).Default("pending"),
		field.Int16("status").Default(0),
		field.Uint64("total_time_cost_ms"),
		field.Uint64("max_memory_cost_kb"),
		field.Int("score").Default(0).SchemaType(map[string]string{"mysql": "INT"}),
	}
}

func (CaseGroupResult) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "CaseGroup_Results"},
	}
}

func (CaseGroupResult) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("problem", Problem.Type),
		edge.To("case_results", CaseResult.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
