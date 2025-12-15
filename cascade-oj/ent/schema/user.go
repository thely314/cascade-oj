package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Immutable(),
		field.String("username").MaxLen(50).NotEmpty().Unique(),
		field.String("email").MaxLen(100).NotEmpty().Unique(),
		field.String("password_hash").MaxLen(255).NotEmpty(),
		field.Enum("role").Values("competitor", "creator", "admin").Default("competitor"),
	}
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "Users"},
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("problems", Problem.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("judge_records", JudgeRecord.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("announcements", Announcement.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("problem_set_manager", ProblemSetManager.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("competitor_list", Competitor_List.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
