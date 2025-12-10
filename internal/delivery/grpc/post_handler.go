package grpc

import (
	"context"
	"log"

	"github.com/Nucleussss/hikayat-forum/post/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	postpb "github.com/Nucleussss/hikayat-proto/gen/go/post/v1"
)

type PostHandler struct {
	postpb.UnimplementedPostServiceServer
	postService service.PostServiceInterface
}

func NewPostHandler(postService service.PostServiceInterface) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) CreatePost(ctx context.Context, req *postpb.CreatePostRequest) (*postpb.Post, error) {
	op := "grpc.PostHandler.CreatePost"
	post, err := h.postService.CreatePost(ctx, req)
	if err != nil {
		log.Printf("%s Error creating post: %v", op, err)
		return nil, status.Error(codes.Internal, "Failed to create post")
	}

	return post, nil
}

func (h *PostHandler) GetPost(ctx context.Context, req *postpb.GetPostRequest) (*postpb.Post, error) {
	op := "grpc.PostHandler.UpdatePost"
	post, err := h.postService.GetPost(ctx, req)
	if err != nil {
		log.Printf("%s Error updating post: %v", op, err)
		return nil, status.Error(codes.Internal, "Failed to update post")
	}

	return post, nil
}

func (h *PostHandler) ListPosts(ctx context.Context, req *postpb.ListPostsRequest) (*postpb.ListPostsResponse, error) {
	op := "grpc.PostHandler.ListPost"
	posts, err := h.postService.ListPosts(ctx, req)
	if err != nil {
		log.Printf("%s Error updating post: %v", op, err)
		return nil, status.Error(codes.Internal, "Failed to get list of posts")
	}

	return posts, nil
}

func (h *PostHandler) UpdatePost(ctx context.Context, req *postpb.UpdatePostRequest) (*postpb.Post, error) {
	op := "grpc.PostHandler.UpdatePost"

	post, err := h.postService.UpdatePost(ctx, req)
	if err != nil {
		log.Printf("%s Error updating post: %v", op, err)
		return nil, status.Error(codes.Internal, " Failed to update post")
	}

	return post, nil
}

func (h *PostHandler) DeletePost(ctx context.Context, req *postpb.DeletePostRequest) (*emptypb.Empty, error) {
	op := "grpc.PostHandler.DeletePost"

	

	err := h.postService.DeletePost(ctx, req)
	if err != nil {
		log.Printf("%s Error deleting post: %v", op, err)
		return nil, status.Error(codes.Internal, " Failed to delete post")
	}

	return nil, nil
}
