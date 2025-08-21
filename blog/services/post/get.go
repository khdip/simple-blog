package post

import (
	"context"

	ppb "blog/gunk/v1/post"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *PostSvc) GetPost(ctx context.Context, req *ppb.GetPostRequest) (*ppb.GetPostResponse, error) {
	id := req.ID
	post, err := s.core.GetPost(context.Background(), id)

	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get post")
	}

	return &ppb.GetPostResponse{
		Post: &ppb.Post{
			ID:          post.ID,
			Title:       post.Title,
			Author:      post.Author,
			Description: post.Description,
			Category:    post.Category,
		},
	}, nil
}
