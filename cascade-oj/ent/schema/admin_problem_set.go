package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type AdminProblemSet struct {
	ent.Schema
}

func (AdminProblemSet) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Immutable(),
		field.Int64("admin_id").Positive(),
		field.Int64("problem_set_id").Positive(),
	}
}

func (AdminProblemSet) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "Admin_ProblemSet"},
	}
}

func (AdminProblemSet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("admin", User.Type).
			Field("admin_id").
			Ref("admin_problem_sets").
			Unique().
			Required(),

		edge.From("problem_set", ProblemSet.Type).
			Field("problem_set_id").
			Ref("admin_problem_sets").
			Unique().
			Required(),
	}
}
