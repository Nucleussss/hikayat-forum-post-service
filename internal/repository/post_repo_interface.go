package repository

import (
	"context"

	pb "github.com/Nucleussss/hikayat-forum/post/api/post/v1"
)

type PostRepoInterface interface {
	CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error)
	GetPost(ctx context.Context, id string) (*pb.Post, error)
	ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error)
	UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.Post, error)
	DeletePost(ctx context.Context, id string) error
}
