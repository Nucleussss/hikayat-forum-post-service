package utils

import (
	"github.com/Nucleussss/hikayat-forum/post/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"

	postpb "github.com/Nucleussss/hikayat-proto/gen/go/post/v1"
)

func PostModelToPB(p *models.Post) *postpb.Post {
	if p == nil {
		return nil
	}

	return &postpb.Post{
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
