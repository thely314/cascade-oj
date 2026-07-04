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

// AlertReport holds AI-generated analysis reports after noise reduction.
type AlertReport struct {
	ent.Schema
}

func (AlertReport) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Immutable(),
		field.String("title").
			MaxLen(255).
			NotEmpty().
			Comment("report title summarizing the analysis period"),
		field.String("summary").
			SchemaType(map[string]string{
				dialect.MySQL: "TEXT",
			}).
			NotEmpty().
			Comment("executive summary of findings"),
		field.String("root_causes").
			SchemaType(map[string]string{
				dialect.MySQL: "JSON",
			}).
			Comment("JSON array of identified root causes"),
		field.String("suggested_actions").
			SchemaType(map[string]string{
				dialect.MySQL: "JSON",
			}).
			Comment("JSON array of suggested remediation actions"),
		field.Int("noise_count").
			Default(0).
			Comment("number of alerts classified as noise"),
		field.Int("real_risk_count").
			Default(0).
			Comment("number of alerts classified as real risks"),
		field.Enum("status").
			Values("draft", "published", "acknowledged").
			Default("draft"),
		field.String("full_report_markdown").
			SchemaType(map[string]string{
				dialect.MySQL: "TEXT",
			}).
			Comment("complete markdown report for OCEs"),
		field.Time("created_at").
			Default(func() time.Time { return time.Now() }).
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).
			Annotations(entsql.Annotation{
				Default: "CURRENT_TIMESTAMP",
			}).
			Immutable(),
		field.Time("acknowledged_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).
			Comment("time when an OCE acknowledged the report"),
	}
}

func (AlertReport) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
		index.Fields("created_at"),
	}
}

func (AlertReport) Edges() []ent.Edge {
	return nil
}

func (AlertReport) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "AlertReports"},
	}
}
