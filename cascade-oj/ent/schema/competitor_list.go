package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Competitor_List struct {
	ent.Schema
}

func (Competitor_List) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Immutable(),
		field.Int64("user_id").Positive(),
		field.Int64("problem_set_id").Positive(),
		field.Int("total_score").Default(0).SchemaType(map[string]string{"mysql": "INT"}).Optional(),
	}
}

func (Competitor_List) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "Competitor_List"},
	}
}

func (Competitor_List) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Field("user_id").
			Ref("competitor_list").
			Unique().
			Required(),

		edge.From("problem_set", ProblemSet.Type).
			Field("problem_set_id").
			Ref("competitor_list").
			Unique().
			Required(),
	}
}
