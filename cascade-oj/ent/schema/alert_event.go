package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// AlertEvent holds raw alert events collected from metrics and logs.
type AlertEvent struct {
	ent.Schema
}

func (AlertEvent) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Immutable(),
		field.String("source").
			MaxLen(64).
			NotEmpty().
			Comment("service name that generated the alert"),
		field.Enum("severity").
			Values("critical", "warning", "info").
			Default("warning"),
		field.String("metric_name").
			MaxLen(128).
			Default("").
			Comment("e.g. http_server_requests_seconds_count"),
		field.Float("metric_value").
			Default(0),
		field.Float("threshold").
			Default(0).
			Comment("trigger threshold value"),
		field.String("raw_data").
			SchemaType(map[string]string{
				dialect.MySQL: "JSON",
			}).
			Comment("full alert payload as JSON"),
		field.Time("received_at").
			Default(func() time.Time { return time.Now() }).
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).
			Annotations(entsql.Annotation{
				Default: "CURRENT_TIMESTAMP",
			}),
		field.Int64("aggregated_into").
			Default(0).
			Comment("report ID that this alert was aggregated into, 0 = not yet aggregated"),
	}
}

func (AlertEvent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("source", "severity"),
		index.Fields("received_at"),
		index.Fields("aggregated_into"),
	}
}

func (AlertEvent) Edges() []ent.Edge {
	return nil
}

func (AlertEvent) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "AlertEvents"},
	}
}
