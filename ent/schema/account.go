package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"time"
)

// Account holds the schema definition for the Account entity.
type Account struct {
	ent.Schema
}

// Fields of the Account.
func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.String("first_name"),
		field.String("last_name"),
		field.String("country"),
		field.Time("birthDay"),

		field.Enum("currency").
			GoType(Currency("")),

		field.Float("amount").SchemaType(map[string]string{
			dialect.Postgres: "numeric", // Override Postgres.
		}),

		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Site.
func (Account) Edges() []ent.Edge {
	return nil
}
