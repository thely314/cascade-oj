package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type TestCase struct {
	ent.Schema
}

func (TestCase) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Immutable(),
		field.Int64("problem_id").Positive(),
		field.String("input").SchemaType(map[string]string{"mysql": "TEXT"}).NotEmpty(),
		field.String("output").SchemaType(map[string]string{"mysql": "TEXT"}).NotEmpty(),
	}
}

func (TestCase) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "TestCases"},
	}
}

func (TestCase) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("problem", Problem.Type).
			Field("problem_id").
			Ref("test_cases").
			Unique().
			Required(),
	}
}
