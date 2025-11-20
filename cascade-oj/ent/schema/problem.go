package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Problem struct {
	ent.Schema
}

func (Problem) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Immutable(),
		field.Int64("creator_id").
			Positive(),
		field.String("title").MaxLen(100).
			NotEmpty(),
		field.String("description").
			SchemaType(map[string]string{
				"mysql": "TEXT",
			}).NotEmpty(),
		field.Enum("problem_type").
			Values("OJ", "other").
			Default("OJ"),
		field.Int("time_limit").
			Positive().SchemaType(map[string]string{
			"mysql": "INT",
		}),
		field.Int("memory_limit").
			Positive().SchemaType(map[string]string{
			"mysql": "INT",
		}),
		field.Enum("use_status").
			Values("unavailable", "available", "using").
			Default("unavailable"),
	}
}

func (Problem) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "Problems"},
	}
}

func (Problem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("creator", User.Type).
			Field("creator_id").
			Ref("problems").
			Unique().
			Required(),
		edge.To("test_cases", TestCase.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("judge_records", JudgeRecord.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("submissions", SubmissionRecord.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("problem_set_problems", ProblemSet_Problem.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
