package post

import (
	"context"

	"blog/blog/storage"
)

type postStore interface {
	CreatePost(context.Context, storage.Post) (int64, error)
	GetPost(context.Context, int64) (storage.Post, error)
	GetAllPosts(context.Context) ([]storage.Post, error)
	UpdatePost(context.Context, storage.Post) error
	DeletePost(context.Context, int64) error
	SearchPost(context.Context, string) ([]storage.Post, error)
}

type CorePostSvc struct {
	store postStore
}

func NewCorePostSvc(s postStore) *CorePostSvc {
	return &CorePostSvc{
		store: s,
	}
}

func (cs CorePostSvc) CreatePost(ctx context.Context, t storage.Post) (int64, error) {
	return cs.store.CreatePost(ctx, t)
}

func (cs CorePostSvc) GetPost(ctx context.Context, id int64) (storage.Post, error) {
	return cs.store.GetPost(ctx, id)
}

func (cs CorePostSvc) GetAllPosts(ctx context.Context) ([]storage.Post, error) {
	return cs.store.GetAllPosts(ctx)
}

func (cs CorePostSvc) UpdatePost(ctx context.Context, t storage.Post) error {
	return cs.store.UpdatePost(ctx, t)
}

func (cs CorePostSvc) DeletePost(ctx context.Context, id int64) error {
	return cs.store.DeletePost(ctx, id)
}

func (cs CorePostSvc) SearchPost(ctx context.Context, q string) ([]storage.Post, error) {
	return cs.store.SearchPost(ctx, q)
}
