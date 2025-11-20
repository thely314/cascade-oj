package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type SystemLog struct {
	ent.Schema
}

func (SystemLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Immutable(),
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
