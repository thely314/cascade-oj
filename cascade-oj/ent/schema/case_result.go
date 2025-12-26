package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// judge result for each test case
type CaseResult struct {
	ent.Schema
}

func (CaseResult) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Immutable(),
		field.Int64("case_group_result_id"),
		// field.Enum("status").Values(
		// 	"pending", "judging", "accepted", "wrong_answer",
		// 	"time_limit_exceeded", "memory_limit_exceeded", "runtime_error", "compilation_error", "system_error",
		// ).Default("pending"),
		field.Int16("status").Default(0),
		field.Uint64("time_cost_ms"),
		field.Uint64("memory_cost_kb"),
		field.String("stdout").Default(""),
		field.String("stderr").Default(""),
		field.Int("score").Default(0).SchemaType(map[string]string{"mysql": "INT"}),
	}
}

func (CaseResult) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "Case_Results"},
	}
}

func (CaseResult) Edges() []ent.Edge {
	return []ent.Edge{
		// edge.To("problem", Problem.Type),
		edge.From("case_group_result", CaseGroupResult.Type).
			Field("case_group_result_id").
			Ref("case_results").
			Unique().
			Required(),
	}
}
