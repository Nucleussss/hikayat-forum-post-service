package repository

import (
	"context"

	postpb "github.com/Nucleussss/hikayat-proto/gen/go/post/v1"
)

type PostRepoInterface interface {
	CreatePost(ctx context.Context, req *postpb.CreatePostRequest) (*postpb.Post, error)
	GetPost(ctx context.Context, id string) (*postpb.Post, error)
	ListPosts(ctx context.Context, req *postpb.ListPostsRequest) (*postpb.ListPostsResponse, error)
	UpdatePost(ctx context.Context, req *postpb.UpdatePostRequest) (*postpb.Post, error)
	DeletePost(ctx context.Context, id string) error
}
