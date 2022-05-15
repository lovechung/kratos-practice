package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// Car holds the schema definition for the Car entity.
type Car struct {
	ent.Schema
}

func (Car) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "car"},
	}
}

// Fields of the Car.
func (Car) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Int64("user_id").
			Optional(),
		field.String("model").
			Optional(),
		field.Time("registered_at").
			Default(time.Now).
			SchemaType(map[string]string{dialect.MySQL: "datetime"}),
	}
}

// Edges of the Car.
func (Car) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("cars").
			Field("user_id").
			Unique(),
	}
}
