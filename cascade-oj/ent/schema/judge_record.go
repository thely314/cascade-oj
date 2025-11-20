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

type JudgeRecord struct {
	ent.Schema
}

func (JudgeRecord) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Immutable(),
		field.Int64("problem_id").Positive(),
		field.Int64("user_id").Positive(),
		field.Enum("status").Values(
			"pending", "judging", "accepted", "wrong_answer",
			"time_limit_exceeded", "memory_limit_exceeded", "runtime_error", "compilation_error",
		).Default("pending"),
		field.Time("judge_start_time").Default(time.Now()).SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}).
			Annotations(entsql.Annotation{
				Default: "CURRENT_TIMESTAMP",
			}),
		field.String("result").SchemaType(map[string]string{"mysql": "TEXT"}).Optional(),
		field.String("code").SchemaType(map[string]string{"mysql": "TEXT"}).NotEmpty(),
		field.Enum("language").Values("c", "cpp", "python", "rust").Default("c"),
		field.Enum("judge_type").Values("test_case", "custom_test_case").Default("custom_test_case"),
	}
}

func (JudgeRecord) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "JudgeRecords"},
	}
}

func (JudgeRecord) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("problem", Problem.Type).
			Field("problem_id").
			Ref("judge_records").
			Unique().
			Required(),
		edge.From("user", User.Type).
			Field("user_id").
			Ref("judge_records").
			Unique().
			Required(),
		edge.To("submissions", SubmissionRecord.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
