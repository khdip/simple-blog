package category

import (
	"context"

	cpb "blog/gunk/v1/category"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CategorySvc) GetCategory(ctx context.Context, req *cpb.GetCategoryRequest) (*cpb.GetCategoryResponse, error) {
	id := req.CategoryID
	category, err := s.core.GetCategory(context.Background(), id)

	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get category")
	}

	return &cpb.GetCategoryResponse{
		Category: &cpb.Category{
			CategoryID:   category.CategoryID,
			CategoryName: category.CategoryName,
		},
	}, nil
}
