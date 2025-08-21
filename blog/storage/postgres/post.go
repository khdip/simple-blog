package postgres

import (
	"blog/blog/storage"
	"context"
)

const insertPost = `
	INSERT INTO posts(
		title,
		author,
		description,
		category
	) VALUES(
		:title,
		:author,
		:description,
		:category
	)
	RETURNING id;
`
const getPost = `
	SELECT * FROM posts
	WHERE id=$1;
`
const getAllPosts = `
	SELECT * FROM posts;
`
const updatePost = `
	UPDATE posts 
	SET title=:title,
	author=:author,
	description=:description,
	category=:category
	WHERE id=:id
	RETURNING id;
`
const deletePost = `
	DELETE FROM posts 
	WHERE id=$1;
`
const searchPost = `
	SELECT * FROM posts 
	WHERE title ILIKE '%%' || $1 || '%%';
`

func (s *Storage) CreatePost(ctx context.Context, t storage.Post) (int64, error) {
	stmt, err := s.db.PrepareNamed(insertPost)
	if err != nil {
		return 0, err
	}

	var id int64
	if err := stmt.Get(&id, t); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) GetPost(ctx context.Context, id int64) (storage.Post, error) {
	var post storage.Post
	if err := s.db.Get(&post, getPost, id); err != nil {
		return storage.Post{}, err
	}
	return post, nil
}

func (s *Storage) GetAllPosts(ctx context.Context) ([]storage.Post, error) {
	var post []storage.Post
	if err := s.db.Select(&post, getAllPosts); err != nil {
		return []storage.Post{}, err
	}
	return post, nil
}

func (s *Storage) UpdatePost(ctx context.Context, t storage.Post) error {
	stmt, err := s.db.PrepareNamed(updatePost)
	if err != nil {
		return err
	}

	var id int64
	if err := stmt.Get(&id, t); err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeletePost(ctx context.Context, id int64) error {
	res := s.db.MustExec(deletePost, id)
	ok, err := res.RowsAffected()
	if err != nil || ok == 0 {
		return err
	}
	return nil
}

func (s *Storage) SearchPost(ctx context.Context, sq string) ([]storage.Post, error) {
	var post []storage.Post
	if err := s.db.Select(&post, searchPost, sq); err != nil {
		return []storage.Post{}, err
	}
	return post, nil
}
