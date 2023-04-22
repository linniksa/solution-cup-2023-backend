package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Transaction holds the schema definition for the Transaction entity.
type Transaction struct {
	ent.Schema
}

func (Transaction) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("operation").GoType(TransactionType("")),
		field.Time("date"),
		field.Float("amount").SchemaType(map[string]string{
			dialect.Postgres: "numeric", // Override Postgres.
		}),
	}
}

func (Transaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("account", Account.Type).Unique(),
	}
}

// Indexes of the Card.
func (Transaction) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("transaction_account", "date"),
	}
}

type TransactionType string

func (_ TransactionType) Values() []string {
	return []string{
		"deposit",
		"transfer",
	}
}

const (
	Deposit  TransactionType = "deposit"
	Transfer TransactionType = "transfer"
)
