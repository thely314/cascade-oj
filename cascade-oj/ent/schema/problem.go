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
		field.Int64("judge_config_id"),
		field.Int16("case_version").Default(1),
		field.Int("time_limit_ms").
			Positive().SchemaType(map[string]string{
			"mysql": "INT",
		}).Comment("milliseconds"),
		field.Int("memory_limit_kb").
			Positive().SchemaType(map[string]string{
			"mysql": "INT",
		}).Comment("kilobytes"),
		field.Enum("use_status").
			Values("unavailable", "available", "using", "deleted").
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
		edge.From("judge_config", ProblemJudgeConfig.Type).
			Field("judge_config_id").
			Ref("judge_config").
			Unique().
			Required(),
		edge.To("judge_records", JudgeRecord.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("submissions", SubmissionRecord.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("problem_set_includes", ProblemSet_Includes.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("templates", ProblemTemplate.Type),
	}
}
