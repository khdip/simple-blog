package postgres

import (
	"blog/blog/storage"
	"context"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreatePost(t *testing.T) {
	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      storage.Post
		want    int64
		wantErr bool
	}{
		{
			name: "CREATE_POST_SUCCESS",
			in: storage.Post{
				Title:       "A Post Title",
				Author:      "John Doe",
				Description: "Bla Bla Bla",
				Category:    "Uncategorized",
			},
			want: 1,
		},
		{
			name: "FAILED_DUPLICATE_TITLE",
			in: storage.Post{
				Title:       "A Post Title",
				Author:      "John Doe",
				Description: "Bla Bla Bla",
				Category:    "Uncategorized",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.CreatePost(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.CreatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Storage.CreatePost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPost(t *testing.T) {
	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      int64
		want    storage.Post
		wantErr bool
	}{
		{
			name: "GET_POST_SUCCESS",
			in:   1,
			want: storage.Post{
				ID:          1,
				Title:       "A Post Title",
				Author:      "John Doe",
				Description: "Bla Bla Bla",
				Category:    "Uncategorized",
			},
		},
		{
			name:    "GET_POST_INVALID",
			in:      100,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetPost(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Diff = %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetAllPosts(t *testing.T) {
	s := newTestStorage(t)
	tests := []struct {
		name    string
		want    []storage.Post
		wantErr bool
	}{
		{
			name: "GET_ALL_POST_SUCCESS",
			want: []storage.Post{
				{
					ID:          1,
					Title:       "A Post Title",
					Author:      "John Doe",
					Description: "Bla Bla Bla",
					Category:    "Uncategorized",
				},
				{
					ID:          2,
					Title:       "Another Post Title",
					Author:      "Melinda Doe",
					Description: "Bla Bla Bla Bla Bla",
					Category:    "Uncategorized 2",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotList, err := s.GetAllPosts(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i].ID < tt.want[j].ID
			})

			sort.Slice(gotList, func(i, j int) bool {
				return gotList[i].ID < gotList[j].ID
			})

			for i, got := range gotList {

				if !cmp.Equal(got, tt.want[i]) {
					t.Errorf("Diff: got -, want += %v", cmp.Diff(got, tt.want[i]))
				}

			}

		})
	}
}

func TestUpdatePost(t *testing.T) {
	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      storage.Post
		wantErr bool
	}{
		{
			name: "UPDATE_POST_SUCCESS",
			in: storage.Post{
				ID:          1,
				Title:       "A Post Title",
				Author:      "John Doe",
				Description: "Bla Bla Bla",
				Category:    "Uncategorized",
			},
			wantErr: false,
		},
		{
			name:    "UPDATE_POST_FAILED",
			in:      storage.Post{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := s.UpdatePost(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.UpdatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDeletePost(t *testing.T) {
	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      int64
		wantErr bool
	}{
		{
			name:    "DELETE_POST_SUCCESS",
			in:      1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := s.DeletePost(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.DeletePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSearchPost(t *testing.T) {
	s := newTestStorage(t)
	tests := []struct {
		name    string
		in      string
		want    []storage.Post
		wantErr bool
	}{
		{
			name: "SEARCH_POST_SUCCESS",
			in:   "Title",
			want: []storage.Post{
				{
					ID:          1,
					Title:       "A Post Title",
					Author:      "John Doe",
					Description: "Bla Bla Bla",
					Category:    "Uncategorized",
				},
				{
					ID:          2,
					Title:       "Another Post Title",
					Author:      "Melinda Doe",
					Description: "Bla Bla Bla Bla Bla",
					Category:    "Uncategorized 2",
				},
			},
		},
		{
			name: "SEARCH_POST_NOT_SUCCESS",
			in:   "Title",
			want: []storage.Post{
				{
					ID:          3,
					Title:       "A Post",
					Author:      "John Doe",
					Description: "Bla Bla Bla",
					Category:    "Uncategorized",
				},
				{
					ID:          4,
					Title:       "Another Post",
					Author:      "Melinda Doe",
					Description: "Bla Bla Bla Bla Bla",
					Category:    "Uncategorized 2",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotList, err := s.SearchPost(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i].ID < tt.want[j].ID
			})

			sort.Slice(gotList, func(i, j int) bool {
				return gotList[i].ID < gotList[j].ID
			})

			for i, got := range gotList {

				if !cmp.Equal(got, tt.want[i]) {
					t.Errorf("Diff: got -, want += %v", cmp.Diff(got, tt.want[i]))
				}

			}

		})
	}
}
