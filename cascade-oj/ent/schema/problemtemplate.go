package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ProblemTemplate struct {
	ent.Schema
}

func (ProblemTemplate) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Default("main.cpp").Comment("文件名"), // 为多文件预留
		field.Text("content").Comment("模板内容"),
	}
}

func (ProblemTemplate) Edges() []ent.Edge {
	return []ent.Edge{
		// 建立与 Problem 的反向关联
		edge.From("problem", Problem.Type).
			Ref("templates").
			Unique(),
	}
}
