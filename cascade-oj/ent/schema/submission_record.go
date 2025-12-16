package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type SubmissionRecord struct {
	ent.Schema
}

func (SubmissionRecord) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Immutable(),
		field.Int64("judge_id").Positive(),
		field.Int64("problem_id").Positive(),
		field.Int64("problem_set_id").Optional(),
		field.Enum("result").Values("Pass", "Fail"),
		field.Time("submission_time").Default(time.Now()).SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}).
			Annotations(entsql.Annotation{
				Default: "CURRENT_TIMESTAMP",
			}),
		field.Int("score").Default(0).SchemaType(map[string]string{"mysql": "INT"}),
	}
}

func (SubmissionRecord) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "SubmissionRecords"},
	}
}

func (SubmissionRecord) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("judge", JudgeRecord.Type).
			Field("judge_id").
			Ref("submissions").
			Unique().
			Required(),

		edge.From("problem", Problem.Type).
			Field("problem_id").
			Ref("submissions").
			Unique().
			Required(),

		edge.From("problem_set", ProblemSet.Type).
			Field("problem_set_id").
			Ref("submissions").
			Unique(),
	}
}
