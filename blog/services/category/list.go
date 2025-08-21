package category

import (
	"context"

	cpb "blog/gunk/v1/category"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CategorySvc) GetAllCategories(ctx context.Context, req *cpb.GetAllCategoriesRequest) (*cpb.GetAllCategoriesResponse, error) {
	categories, err := s.core.GetAllCategories(context.Background())
	var c []*cpb.Category
	for _, category := range categories {
		c = append(c, &cpb.Category{
			CategoryID:   category.CategoryID,
			CategoryName: category.CategoryName,
		})
	}

	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get category")
	}

	return &cpb.GetAllCategoriesResponse{
		Categories: c,
	}, nil
}
