# entgo, gqlgen, fiber sample

## How to?

1. Insert under commands

```
$ go get -d entgo.io/ent/cmd/ent
$ go run entgo.io/ent/cmd/ent init Todo
$ go get -u entgo.io/contrib/entgql
```

2. Add `ent/schema/todo.go` Field and Annotation

3. Add `ent/schema/todo.graphql` and `ent/gqlgen.yml`

4. Add `ent/entc.go`

5. Fix `ent/generate.go` content

6. Do generate
  ```sh
  $ go generate ./...
  ```

7. Fill `resolvers/*` contents

8. Install fiber
  ```sh
  $ go get -u github.com/gofiber/fiber/v2
  $ go get -u github.com/gofiber/adaptor/v2
  $ go get -u github.com/99designs/gqlgen/graphql/playground
  $ go get -u github.com/mattn/go-sqlite3
  ```

9. Add `main.go` and Start server
  ```go
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
    client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
    if err != nil {
      panic(err)
    }
    
    if err := client.Schema.Create(
      context.Background(),
      migrate.WithGlobalUniqueID(true),
    ); err != nil {
      panic(err)
    }

    srv := handler.NewDefaultServer(resolvers.NewSchema(client))
    srv.Use(entgql.Transactioner{TxOpener: client})
    
    app := fiber.New()
    app.Get("/playground", adaptor.HTTPHandlerFunc(playground.Handler("Graphql Playground", "/query")))
    app.All("/query", adaptor.HTTPHandler(srv))
    app.Listen(":3000")
  }
  ```

  ```sh
  $ go run main.go
  ```

10. Enter http://localhost:3000/playground

## If you want embedded postgresql

1. Add dependencies
  ```sh
  $ go get -u github.com/fergusstrange/embedded-postgres
  $ go get -u github.com/jackc/pgx/v4/stdlib
  ```

2. Write `main.go`

```go
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
```
