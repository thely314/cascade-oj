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

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Unique().
			Immutable().SchemaType(map[string]string{
			"mysql": "INT AUTO_INCREMENT",
		}),
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
		edge.To("admin_problem_sets", AdminProblemSet.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("problem_set_users", ProblemSetUser.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

type Problem struct {
	ent.Schema
}

func (Problem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Unique().
			Immutable().SchemaType(map[string]string{
			"mysql": "INT AUTO_INCREMENT",
		}),
		field.Int("creator_id").
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
		edge.To("problem_set_problems", ProblemSetProblem.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

type TestCase struct {
	ent.Schema
}

func (TestCase) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Unique().
			Immutable().
			SchemaType(map[string]string{"mysql": "INT AUTO_INCREMENT"}),
		field.Int("problem_id").Positive(),
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

type ProblemSet struct {
	ent.Schema
}

func (ProblemSet) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Unique().
			Immutable().
			SchemaType(map[string]string{"mysql": "INT AUTO_INCREMENT"}),
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
		edge.To("admin_problem_sets", AdminProblemSet.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("problem_set_problems", ProblemSetProblem.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("problem_set_users", ProblemSetUser.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

type JudgeRecord struct {
	ent.Schema
}

func (JudgeRecord) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Unique().
			Immutable().
			SchemaType(map[string]string{"mysql": "INT AUTO_INCREMENT"}),
		field.Int("problem_id").Positive(),
		field.Int("user_id").Positive(),
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

type SubmissionRecord struct {
	ent.Schema
}

func (SubmissionRecord) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Unique().
			Immutable().
			SchemaType(map[string]string{"mysql": "INT AUTO_INCREMENT"}),
		field.Int("judge_id").Positive(),
		field.Int("problem_id").Positive(),
		field.Int("problem_set_id").Optional(),
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

type Announcement struct {
	ent.Schema
}

func (Announcement) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Unique().
			Immutable().
			SchemaType(map[string]string{"mysql": "INT AUTO_INCREMENT"}),
		field.Int("publisher_id").Positive(),
		field.String("title").MaxLen(100).NotEmpty(),
		field.String("content").SchemaType(map[string]string{"mysql": "TEXT"}).NotEmpty(),
	}
}

func (Announcement) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "Announcements"},
	}
}

func (Announcement) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("publisher", User.Type).
			Field("publisher_id").
			Ref("announcements").
			Unique().
			Required(),
	}
}

type SystemLog struct {
	ent.Schema
}

func (SystemLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Unique().
			Immutable().
			SchemaType(map[string]string{"mysql": "INT AUTO_INCREMENT"}),
		field.Time("log_time").Default(func() time.Time { return time.Now() }).SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}).
			Annotations(entsql.Annotation{
				Default: "CURRENT_TIMESTAMP",
			}),
		field.String("log_info").SchemaType(map[string]string{"mysql": "TEXT"}).NotEmpty(),
	}
}

func (SystemLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "Logs"},
	}
}

func (SystemLog) Edges() []ent.Edge {
	return nil
}

type AdminProblemSet struct {
	ent.Schema
}

func (AdminProblemSet) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Unique().
			Immutable().
			SchemaType(map[string]string{"mysql": "INT AUTO_INCREMENT"}),
		field.Int("admin_id").Positive(),
		field.Int("problem_set_id").Positive(),
	}
}

func (AdminProblemSet) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "Admin_ProblemSet"},
	}
}

func (AdminProblemSet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("admin", User.Type).
			Field("admin_id").
			Ref("admin_problem_sets").
			Unique().
			Required(),

		edge.From("problem_set", ProblemSet.Type).
			Field("problem_set_id").
			Ref("admin_problem_sets").
			Unique().
			Required(),
	}
}

type ProblemSetProblem struct {
	ent.Schema
}

func (ProblemSetProblem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Unique().
			Immutable().
			SchemaType(map[string]string{"mysql": "INT AUTO_INCREMENT"}),
		field.Int("problem_set_id").Positive(),
		field.Int("problem_id").Positive(),
		field.Int("problem_order").Positive().
			SchemaType(map[string]string{"mysql": "INT"}),
	}
}

func (ProblemSetProblem) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "ProblemSet_Problems"},
	}
}

func (ProblemSetProblem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("problem_set", ProblemSet.Type).
			Field("problem_set_id").
			Ref("problem_set_problems").
			Unique().
			Required(),

		edge.From("problem", Problem.Type).
			Field("problem_id").
			Ref("problem_set_problems").
			Unique().
			Required(),
	}
}

type ProblemSetUser struct {
	ent.Schema
}

func (ProblemSetUser) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Unique().
			Immutable().
			SchemaType(map[string]string{"mysql": "INT AUTO_INCREMENT"}),
		field.Int("user_id").Positive(),
		field.Int("problem_set_id").Positive(),
		field.Int("total_score").Default(0).SchemaType(map[string]string{"mysql": "INT"}).Optional(),
	}
}

func (ProblemSetUser) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "ProblemSet_Users"},
	}
}

func (ProblemSetUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Field("user_id").
			Ref("problem_set_users").
			Unique().
			Required(),

		edge.From("problem_set", ProblemSet.Type).
			Field("problem_set_id").
			Ref("problem_set_users").
			Unique().
			Required(),
	}
}
