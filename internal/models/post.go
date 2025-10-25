package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	AuthorID  uuid.UUID `db:"author_id"`
	Category  string    `db:"category"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	IsDeleted bool      `db:"is_deleted"`
}

// // Request
// type CreatePostRequest struct {
// 	Title    string
// 	Content  string
// 	AuthorID string
// 	Category string
// }

// type GetPostRequest struct {
// 	ID uuid.UUID
// }

// type ListPostsRequest struct {
// 	AuthorID string
// 	Category string
// 	Page     int32
// 	Limit    int32
// }

// type UpdatePostRequest struct {
// 	ID       uuid.UUID
// 	Title    *string
// 	Content  *string
// 	Category *string
// }

// type DeletePostRequest struct {
// 	ID uuid.UUID
// }

// // Response
// type ListPostsResponse struct {
// 	Posts    []*Post
// 	Has_more bool
// }
