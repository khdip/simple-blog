package category

import (
	"context"

	"blog/blog/storage"
	cpb "blog/gunk/v1/category"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CategorySvc) CreateCategory(ctx context.Context, req *cpb.CreateCategoryRequest) (*cpb.CreateCategoryResponse, error) {
	category := storage.Category{
		CategoryID:   req.Category.CategoryID,
		CategoryName: req.Category.CategoryName,
	}

	id, err := s.core.CreateCategory(context.Background(), category)

	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create category")
	}

	return &cpb.CreateCategoryResponse{
		CategoryID: id,
	}, nil
}
