package service

import (
	"context"
	"log"

	pb "github.com/Nucleussss/hikayat-forum/post/api/post/v1"
	"github.com/Nucleussss/hikayat-forum/post/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostService struct {
	postRepo repository.PostRepoInterface
}

func NewPostService(postRepo repository.PostRepoInterface) PostServiceInterface {
	return &PostService{postRepo: postRepo}
}

func (s *PostService) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	op := "service.CreatePost"
	post, err := s.postRepo.CreatePost(ctx, req)
	if err != nil {
		log.Printf("%s Error creating post: %v", op, err)
		return nil, status.Error(codes.Internal, "Failed to create post")
	}

	log.Printf("%s succeed create post", op)
	return post, nil
}

func (s *PostService) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.Post, error) {
	op := "service.GetPost"
	post, err := s.postRepo.GetPost(ctx, req.Id)
	if err != nil {
		log.Printf("%s Error getting post: %v", op, err)
		return nil, status.Error(codes.Internal, "Failed to get post")
	}

	log.Printf("%s succeed get post", op)
	return post, nil
}

func (s *PostService) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	op := "service.ListPosts"
	list, err := s.postRepo.ListPosts(ctx, req)
	if err != nil {
		log.Printf("%s Error listing posts: %v", op, err)
		return nil, status.Error(codes.Internal, "Failed to list posts")
	}

	if len(list.Posts) == 0 {
		log.Printf("%s Failed to get list of posts, no post match the paramaeter: %v", op, err)
	}

	return list, nil
}

func (s *PostService) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.Post, error) {
	op := "service.UpdatePost"

	// validate fieldmask
	if !req.UpdateMask.IsValid(req.Post) {
		log.Printf("Invalid fieldmask for update request")
		return nil, status.Error(codes.InvalidArgument, "Invalid fieldmask for update request")
	}

	post, err := s.postRepo.GetPost(ctx, req.Id)
	if err != nil {
		log.Printf("%s Error getting post to update: %v", op, err)
		return nil, status.Error(codes.Internal, " Failed to get post")
	}

	// update post fields based on fieldmask
	for _, path := range req.UpdateMask.Paths {
		switch path {
		case "title":
			post.Title = req.Post.Title
		case "content":
			post.Content = req.Post.Content
		case "category":
			post.Category = req.Post.Category
		}
	}

	// create new update request with updated fields and fieldmask
	newReq := &pb.UpdatePostRequest{
		Id:         req.Id,
		Post:       post,
		UpdateMask: req.UpdateMask,
	}

	// update post base on new request
	postUpdate, err := s.postRepo.UpdatePost(ctx, newReq)
	if err != nil {
		log.Printf("%s Error getting post to update: %v", op, err)
		return nil, status.Error(codes.Internal, " Failed to get post")
	}

	return postUpdate, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *pb.DeletePostRequest) error {
	op := "service.DeletePost"

	err := s.postRepo.DeletePost(ctx, req.Id)
	if err != nil {
		log.Printf("%s Error deleting post: %v", op, err)
		return status.Error(codes.Internal, "Failed to delete post")
	}

	return nil
}
