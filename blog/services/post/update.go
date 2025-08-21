package post

import (
	"context"

	"blog/blog/storage"
	ppb "blog/gunk/v1/post"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *PostSvc) UpdatePost(ctx context.Context, req *ppb.UpdatePostRequest) (*ppb.UpdatePostResponse, error) {
	post := storage.Post{
		ID:          req.Post.ID,
		Title:       req.Post.Title,
		Author:      req.Post.Author,
		Description: req.Post.Description,
		Category:    req.Post.Category,
	}

	err := s.core.UpdatePost(context.Background(), post)

	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to update post")
	}

	return &ppb.UpdatePostResponse{}, nil
}
