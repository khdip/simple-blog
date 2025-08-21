package post

import (
	"context"

	ppb "blog/gunk/v1/post"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *PostSvc) GetAllPosts(ctx context.Context, req *ppb.GetAllPostsRequest) (*ppb.GetAllPostsResponse, error) {
	posts, err := s.core.GetAllPosts(context.Background())
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
		return nil, status.Error(codes.Internal, "Failed to get posts")
	}

	return &ppb.GetAllPostsResponse{
		Posts: p,
	}, nil
}
