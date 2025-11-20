package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Announcement struct {
	ent.Schema
}

func (Announcement) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Immutable(),
		field.Int64("publisher_id").Positive(),
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
