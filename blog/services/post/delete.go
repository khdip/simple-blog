package post

import (
	"context"

	ppb "blog/gunk/v1/post"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *PostSvc) DeletePost(ctx context.Context, req *ppb.DeletePostRequest) (*ppb.DeletePostResponse, error) {
	id := req.ID

	err := s.core.DeletePost(context.Background(), id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to delete post")
	}

	return &ppb.DeletePostResponse{}, nil
}
