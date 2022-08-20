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
  ```sh
  $ go run main.go
  ```

10. Enter http://localhost:3000/playground
