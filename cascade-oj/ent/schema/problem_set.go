package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ProblemSet struct {
	ent.Schema
}

func (ProblemSet) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Immutable(),
		field.String("name").MaxLen(100).NotEmpty(),
		field.String("description").SchemaType(map[string]string{"mysql": "TEXT"}).Optional(),
		field.Time("start_time").SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}),
		field.Time("end_time").SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}),
		field.Enum("status").Values("upcoming", "ongoing", "ended").Default("upcoming"),
	}
}

func (ProblemSet) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "ProblemSets"},
	}
}

func (ProblemSet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("submissions", SubmissionRecord.Type).
			Annotations(entsql.OnDelete(entsql.SetNull)),
		edge.To("problem_set_manager", ProblemSetManager.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("problem_set_includes", ProblemSet_Includes.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("competitor_list", Competitor_List.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
