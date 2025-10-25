package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	pb "github.com/Nucleussss/hikayat-forum/post/api/post/v1"
	"github.com/Nucleussss/hikayat-forum/post/internal/models"
	"github.com/Nucleussss/hikayat-forum/post/internal/repository"
	"github.com/Nucleussss/hikayat-forum/post/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostRepo struct {
	pool *pgxpool.Pool
}

func NewPostRepo(pool *pgxpool.Pool) repository.PostRepoInterface {
	return &PostRepo{pool: pool}
}

func (r *PostRepo) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	query := `
		INSERT INTO posts(title, content, author_id, category, is_deleted) 
		VALUES($1, $2 ,$3 ,$4,false)
		RETURNING id, title, content, author_id, category, created_at, updated_at, is_deleted
	`

	uuidAuthorID := uuid.MustParse(req.AuthorId)

	post := &models.Post{}
	err := r.pool.QueryRow(ctx, query, req.Title, req.Content, uuidAuthorID, req.Category).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorID,
		&post.Category,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	pbPost := utils.PostModelToPB(post)

	return pbPost, nil
}

func (r *PostRepo) GetPost(ctx context.Context, id string) (*pb.Post, error) {
	query := `
		SELECT id, title, content, author_id, category, created_at, updated_at, is_deleted
		FROM posts 
		WHERE id = $1 AND is_deleted = false
	`
	uuidAuthorID := uuid.MustParse(id)

	post := &models.Post{}
	err := r.pool.QueryRow(ctx, query, uuidAuthorID).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorID,
		&post.Category,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.IsDeleted,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	log.Printf("Post retrieved successfully: %v", post)

	pbPost := utils.PostModelToPB(post)
	return pbPost, nil
}

func (r *PostRepo) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	query := `
		 SELECT id, title, content, author_id, category, created_at, updated_at, is_deleted
		 FROM posts
		 WHERE is_deleted = false
	`

	var args []interface{}
	argIndex := 1

	// optional: filter by ID
	if req.AuthorId != "" {
		query += fmt.Sprintf(" AND author_id = $%d", argIndex)
		uuidVal, err := uuid.Parse(req.AuthorId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid author ID: %s", err)
		}
		args = append(args, uuidVal)
		argIndex++
	}

	// optional: filter by category
	if req.Category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, req.Category)
		argIndex++
	}

	query += " ORDER BY created_at DESC"

	// pagination handling
	limitCheck := req.Limit
	if limitCheck == 0 || limitCheck > 100 {
		limitCheck = 10
	}
	fetchLimit := limitCheck + 1
	query += fmt.Sprintf(" LIMIT $%d", argIndex)
	args = append(args, fetchLimit)
	argIndex++

	// pagination handling
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Page > 1 {
		offset := (req.Page - 1) * req.Limit
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, offset)
		argIndex++
	}

	// do request to database
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// scan results into Post struct slice
	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.Category,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.IsDeleted,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("row iteration error: %w", rows.Err())
	}

	// calculate has_more based on limit and total number of posts
	hasMore := false
	if len(posts) > int(limitCheck) {
		posts = posts[:int(limitCheck)]
		hasMore = true
	}

	pbPosts := make([]*pb.Post, 0, len(posts))
	for _, post := range posts {
		pbPosts = append(pbPosts, utils.PostModelToPB(post))
	}

	response := &pb.ListPostsResponse{
		Posts:   pbPosts,
		HasMore: hasMore,
	}

	return response, nil
}

func (r *PostRepo) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.Post, error) {
	query := `
		UPDATE posts
		SET 
			title = COALESCE($1, title),
			content = COALESCE($2, content),
			category = COALESCE($3, category),
			updated_at = NOW()
		WHERE id = $4 AND is_deleted = false
		RETURNING id, title, content, author_id, category, created_at, updated_at, is_deleted
	`
	uuidAuthorID := uuid.MustParse(req.Id)

	var post models.Post
	err := r.pool.QueryRow(
		ctx,
		query,
		req.Post.Title,
		req.Post.Content,
		req.Post.Category,
		uuidAuthorID,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorID,
		&post.Category,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.IsDeleted,
	)

	posts := utils.PostModelToPB(&post)

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepo) DeletePost(ctx context.Context, id string) error {
	query := `
		DELETE FROM posts 
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowaffected := result.RowsAffected()
	if rowaffected == 0 {
		return errors.New("post not found")
	}

	return nil
}
