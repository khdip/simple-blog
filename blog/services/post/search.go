package post

import (
	"context"

	ppb "blog/gunk/v1/post"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *PostSvc) SearchPost(ctx context.Context, req *ppb.SearchPostRequest) (*ppb.SearchPostResponse, error) {
	query := req.GetSearchPostQuery()
	posts, err := s.core.SearchPost(context.Background(), query)
	var p []*ppb.Post
	for _, post := range posts {
		p = append(p, &ppb.Post{
			ID:          post.ID,
			Title:       post.Title,
			Author:      post.Author,
			Description: post.Description,
			Category:    post.Category,
		})
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to search post")
	}

	return &ppb.SearchPostResponse{
		SearchPostResult: p,
	}, nil
}
