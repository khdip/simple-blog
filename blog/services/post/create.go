package post

import (
	"context"

	"blog/blog/storage"
	ppb "blog/gunk/v1/post"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *PostSvc) CreatePost(ctx context.Context, req *ppb.CreatePostRequest) (*ppb.CreatePostResponse, error) {
	post := storage.Post{
		ID:          req.GetPost().ID,
		Title:       req.GetPost().Title,
		Author:      req.GetPost().Author,
		Description: req.GetPost().Description,
		Category:    req.GetPost().Category,
	}
	id, err := s.core.CreatePost(context.Background(), post)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create post")
	}

	return &ppb.CreatePostResponse{
		ID: id,
	}, nil
}
