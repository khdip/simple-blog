package category

import (
	"context"

	cpb "blog/gunk/v1/category"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CategorySvc) DeleteCategory(ctx context.Context, req *cpb.DeleteCategoryRequest) (*cpb.DeleteCategoryResponse, error) {
	id := req.CategoryID

	err := s.core.DeleteCategory(context.Background(), id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to delete category")
	}

	return &cpb.DeleteCategoryResponse{}, nil
}
