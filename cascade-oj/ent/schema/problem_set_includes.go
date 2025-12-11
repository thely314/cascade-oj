package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ProblemSet_Includes struct {
	ent.Schema
}

func (ProblemSet_Includes) Fields() []ent.Field {
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

func (ProblemSet_Includes) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "ProblemSet_Includes"},
	}
}

func (ProblemSet_Includes) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("problem_set", ProblemSet.Type).
			Field("problem_set_id").
			Ref("problem_set_includes").
			Unique().
			Required(),

		edge.From("problem", Problem.Type).
			Field("problem_id").
			Ref("problem_set_includes").
			Unique().
			Required(),
	}
}
