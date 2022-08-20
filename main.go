package main

import (
	"context"
	"entgogqlgenfibersample/ent"
	"entgogqlgenfibersample/ent/migrate"
	"entgogqlgenfibersample/resolvers"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/dialect"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	app := fiber.New()

	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		panic(err)
	}
	if client.Schema.Create(
		context.Background(),
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		panic(err)
	}

	srv := handler.NewDefaultServer(resolvers.NewSchema(client))
	srv.Use(entgql.Transactioner{TxOpener: client})

	app.Get("/playground", adaptor.HTTPHandlerFunc(playground.Handler("Graphql Playground", "/query")))
	app.All("/query", adaptor.HTTPHandler(srv))

	app.Listen(":3000")
}
