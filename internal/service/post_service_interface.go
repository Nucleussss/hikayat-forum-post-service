package service

import (
	"context"

	postpb "github.com/Nucleussss/hikayat-proto/gen/go/post/v1"
)

type PostServiceInterface interface {
	CreatePost(ctx context.Context, req *postpb.CreatePostRequest) (*postpb.Post, error)
	GetPost(ctx context.Context, req *postpb.GetPostRequest) (*postpb.Post, error)
	ListPosts(ctx context.Context, req *postpb.ListPostsRequest) (*postpb.ListPostsResponse, error)
	UpdatePost(ctx context.Context, req *postpb.UpdatePostRequest) (*postpb.Post, error)
	DeletePost(ctx context.Context, req *postpb.DeletePostRequest) error
}
