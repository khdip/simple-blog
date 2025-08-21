package category

import (
	"context"

	"blog/blog/storage"
	cpb "blog/gunk/v1/category"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CategorySvc) UpdateCategory(ctx context.Context, req *cpb.UpdateCategoryRequest) (*cpb.UpdateCategoryResponse, error) {
	category := storage.Category{
		CategoryID:   req.Category.CategoryID,
		CategoryName: req.Category.CategoryName,
	}

	err := s.core.UpdateCategory(context.Background(), category)

	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to update category")
	}

	return &cpb.UpdateCategoryResponse{}, nil
}
