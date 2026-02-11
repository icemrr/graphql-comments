package gql

import (
	"graphql-comments/internal/storage"

	"github.com/graphql-go/graphql"
)

func BuildSchema(store storage.Storage) (*graphql.Schema, error) {
	resolverContext := &ResolverContext{
		Storage: store,
	}

	// Comment тип
	commentType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"id":       &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"postId":   &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"parentId": &graphql.Field{Type: graphql.String},
			"content":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		},
	})

	// Добавляем replies рекурсивно
	commentType.AddFieldConfig("replies", &graphql.Field{
		Type:    graphql.NewList(commentType),
		Resolve: resolverContext.RepliesResolver,
	})

	// Post тип
	postType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Post",
		Fields: graphql.Fields{
			"id":      &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"title":   &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"content": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"comments": &graphql.Field{
				Type:    graphql.NewList(commentType),
				Resolve: resolverContext.CommentsResolver,
			},
		},
	})

	// Query
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"posts": &graphql.Field{
				Type:    graphql.NewList(postType),
				Resolve: resolverContext.PostsResolver,
			},
		},
	})

	// Mutation
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createPost": &graphql.Field{
				Type: postType,
				Args: graphql.FieldConfigArgument{
					"title":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"content": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolverContext.CreatePostResolver,
			},
			"createComment": &graphql.Field{
				Type: commentType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewInputObject(graphql.InputObjectConfig{
							Name: "CreateCommentInput",
							Fields: graphql.InputObjectConfigFieldMap{
								"postId":   &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
								"parentId": &graphql.InputObjectFieldConfig{Type: graphql.String},
								"content":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
							},
						}),
					},
				},
				Resolve: resolverContext.CreateCommentResolver,
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	return &schema, err
}
