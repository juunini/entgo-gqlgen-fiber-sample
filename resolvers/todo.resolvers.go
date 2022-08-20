package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"entgogqlgenfibersample/ent"
	"entgogqlgenfibersample/ent/todo"

	"github.com/google/uuid"
)

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input ent.CreateTodoInput) (*ent.Todo, error) {
	return ent.FromContext(ctx).Todo.Create().SetInput(input).Save(ctx)
}

// UpdateTodo is the resolver for the updateTodo field.
func (r *mutationResolver) UpdateTodo(ctx context.Context, id uuid.UUID, input ent.UpdateTodoInput) (*ent.Todo, error) {
	return ent.FromContext(ctx).Todo.UpdateOneID(id).SetInput(input).Save(ctx)
}

// UpdateTodos is the resolver for the updateTodos field.
func (r *mutationResolver) UpdateTodos(ctx context.Context, ids []uuid.UUID, input ent.UpdateTodoInput) ([]*ent.Todo, error) {
	client := ent.FromContext(ctx)
	if err := client.Todo.Update().Where(todo.IDIn(ids...)).SetInput(input).Exec(ctx); err != nil {
		return nil, err
	}
	return client.Todo.Query().Where(todo.IDIn(ids...)).All(ctx)
}

// Ping is the resolver for the ping field.
func (r *queryResolver) Ping(ctx context.Context) (string, error) {
	return "pong", nil
}

// Mutation returns ent.MutationResolver implementation.
func (r *Resolver) Mutation() ent.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
