package utils

import (
	pb "github.com/Nucleussss/hikayat-forum/post/api/post/v1"
	"github.com/Nucleussss/hikayat-forum/post/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PostModelToPB(p *models.Post) *pb.Post {
	if p == nil {
		return nil
	}

	return &pb.Post{
		Id:       p.ID.String(),
		Title:    p.Title,
		Content:  p.Content,
		AuthorId: p.AuthorID.String(),
		Category: p.Category,

		CreatedAt: timestamppb.New(p.CreatedAt),
		UpdatedAt: timestamppb.New(p.UpdatedAt),
		IsDeleted: p.IsDeleted,
	}
}
