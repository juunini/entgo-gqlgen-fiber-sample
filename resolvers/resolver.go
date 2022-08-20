package resolvers

import (
	"entgogqlgenfibersample/ent"

	"github.com/99designs/gqlgen/graphql"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{ client *ent.Client }

func NewSchema(client *ent.Client) graphql.ExecutableSchema {
	return ent.NewExecutableSchema(ent.Config{
		Resolvers: &Resolver{client},
	})
}
