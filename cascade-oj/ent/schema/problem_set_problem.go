package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ProblemSet_Problem struct {
	ent.Schema
}

func (ProblemSet_Problem) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Immutable(),
		field.Int64("problem_set_id").Positive(),
		field.Int64("problem_id").Positive(),
		field.Int("problem_order").Positive().
			SchemaType(map[string]string{"mysql": "INT"}),
	}
}

func (ProblemSet_Problem) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "ProblemSet_Problems"},
	}
}

func (ProblemSet_Problem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("problem_set", ProblemSet.Type).
			Field("problem_set_id").
			Ref("problem_set_problems").
			Unique().
			Required(),

		edge.From("problem", Problem.Type).
			Field("problem_id").
			Ref("problem_set_problems").
			Unique().
			Required(),
	}
}
