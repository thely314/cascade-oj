package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ProblemSetManager struct {
	ent.Schema
}

func (ProblemSetManager) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Immutable(),
		field.Int64("admin_id").Positive(),
		field.Int64("problem_set_id").Positive(),
	}
}

func (ProblemSetManager) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "ProblemSet_Managers"},
	}
}

func (ProblemSetManager) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("admin", User.Type).
			Field("admin_id").
			Ref("problem_set_manager").
			Unique().
			Required(),

		edge.From("problem_set", ProblemSet.Type).
			Field("problem_set_id").
			Ref("problem_set_manager").
			Unique().
			Required(),
	}
}
