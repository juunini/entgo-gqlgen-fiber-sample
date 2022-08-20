package main

import (
	"context"
	"database/sql"
	"entgogqlgenfibersample/ent"
	"entgogqlgenfibersample/ent/migrate"
	"entgogqlgenfibersample/resolvers"
	"log"
	"os"
	"os/signal"
	"syscall"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"

	_ "github.com/jackc/pgx/v4/stdlib"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
)

func main() {
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().Database("postgres"))
	if err := postgres.Start(); err != nil {
		panic(err)
	}

	client := openPsql("postgresql://postgres:postgres@127.0.0.1/postgres")
	if err := client.Schema.Create(
		context.Background(),
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		panic(err)
	}

	go interruptEmbeddedPostgres(client, postgres)

	srv := handler.NewDefaultServer(resolvers.NewSchema(client))
	srv.Use(entgql.Transactioner{TxOpener: client})

	app := fiber.New()

	app.Get("/playground", adaptor.HTTPHandlerFunc(playground.Handler("Graphql Playground", "/query")))
	app.All("/query", adaptor.HTTPHandler(srv))

	app.Listen(":3000")
}

func openPsql(databaseUrl string) *ent.Client {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv))
}

func interruptEmbeddedPostgres(client *ent.Client, postgres *embeddedpostgres.EmbeddedPostgres) {
	sig := make(chan os.Signal, 1)
	signal.Notify(
		sig,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGINT,
		os.Interrupt,
	)

	<-sig

	client.Close()
	postgres.Stop()
	os.Exit(0)
}
