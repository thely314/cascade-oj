package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// describes the judge configuration for a problem
type ProblemJudgeConfig struct {
	ent.Schema
}

func (ProblemJudgeConfig) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Immutable(),
		field.String("config_name").Unique(),
		field.String("description").Default(""),
		field.String("submission_queue_name").Default("submission_queue"),
		field.String("self_test_queue_name").Default("self_test_queue"),
		field.String("judge_engine").Default("default").Comment("extensible judge engine identifier"),
	}
}

func (ProblemJudgeConfig) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "Problem_JudgeConfigs"},
	}
}

func (ProblemJudgeConfig) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("judge_config", Problem.Type),
	}
}
