package post

import (
	"context"

	"blog/blog/storage"
	ppb "blog/gunk/v1/post"
)

type postCoreStore interface {
	CreatePost(context.Context, storage.Post) (int64, error)
	GetPost(context.Context, int64) (storage.Post, error)
	GetAllPosts(context.Context) ([]storage.Post, error)
	UpdatePost(context.Context, storage.Post) error
	DeletePost(context.Context, int64) error
	SearchPost(context.Context, string) ([]storage.Post, error)
}

type PostSvc struct {
	ppb.UnimplementedPostServiceServer
	core postCoreStore
}

func NewPostServer(c postCoreStore) *PostSvc {
	return &PostSvc{
		core: c,
	}
}
