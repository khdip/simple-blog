package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	cpb "blog/gunk/v1/category"
	ppb "blog/gunk/v1/post"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

type Post struct {
	ID          int64  `db:"id"`
	Title       string `db:"title"`
	Author      string `db:"author"`
	Description string `db:"description"`
	Category    string `db:"category"`
}

func (post *Post) Validate() error {
	return validation.ValidateStruct(post,
		validation.Field(&post.Title, validation.Required.Error("Title field can not be empty."), validation.Length(3, 50).Error("Title field should have atleast 3 characters and atmost 50 characters")),
		validation.Field(&post.Author, validation.Required.Error("Author field can not be empty.")),
		validation.Field(&post.Description, validation.Required.Error("Description field can not be empty."), validation.Length(10, 0).Error("Description field should have atleast 3 characters")),
	)
}

func (h *Handler) listPost(w http.ResponseWriter, r *http.Request) {
	posts, err := h.pc.GetAllPosts(r.Context(), &ppb.GetAllPostsRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.templates.ExecuteTemplate(w, "home.html", posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}

	int_id, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}

	post, err := h.pc.GetPost(r.Context(), &ppb.GetPostRequest{
		ID: int64(int_id),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if post.Post.ID == 0 {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}

	var p Post
	p.ID = int64(int_id)
	p.Title = post.Post.Title
	p.Author = post.Post.Author
	p.Description = post.Post.Description
	p.Category = post.Post.Category

	err = h.templates.ExecuteTemplate(w, "post.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	ErrorValue := map[string]string{}
	categories, err := h.cc.GetAllCategories(r.Context(), &cpb.GetAllCategoriesRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	category := []Category{}
	for _, v := range categories.Categories {
		category = append(category, Category{
			CategoryName: v.CategoryName,
		})
	}

	post := Post{}
	h.loadCreateForm(w, post, category, ErrorValue)
}

func (h *Handler) storePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var post Post
	err = h.decoder.Decode(&post, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories, err := h.cc.GetAllCategories(r.Context(), &cpb.GetAllCategoriesRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	category := []Category{}
	for _, v := range categories.Categories {
		category = append(category, Category{
			CategoryName: v.CategoryName,
		})
	}

	err = post.Validate()
	if err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			ErrorValue := make(map[string]string)
			for key, value := range vErrors {
				ErrorValue[strings.Title(key)] = value.Error()
			}
			h.loadCreateForm(w, post, category, ErrorValue)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.pc.CreatePost(r.Context(), &ppb.CreatePostRequest{
		Post: &ppb.Post{
			Title:       post.Title,
			Author:      post.Author,
			Description: post.Description,
			Category:    post.Category,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) editPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}

	int_id, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}

	post, err := h.pc.GetPost(r.Context(), &ppb.GetPostRequest{
		ID: int64(int_id),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if post.Post.ID == 0 {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}

	var p Post
	p.ID = int64(int_id)
	p.Title = post.Post.Title
	p.Author = post.Post.Author
	p.Description = post.Post.Description
	p.Category = post.Post.Category

	categories, err := h.cc.GetAllCategories(r.Context(), &cpb.GetAllCategoriesRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	category := []Category{}
	for _, v := range categories.Categories {
		category = append(category, Category{
			CategoryName: v.CategoryName,
		})
	}

	h.loadEditForm(w, p, category, map[string]string{})
}

func (h *Handler) updatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}
	int_id, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}

	var p Post
	p.ID = int64(int_id)

	if p.ID == 0 {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.decoder.Decode(&p, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories, err := h.cc.GetAllCategories(r.Context(), &cpb.GetAllCategoriesRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	category := []Category{}
	for _, v := range categories.Categories {
		category = append(category, Category{
			CategoryName: v.CategoryName,
		})
	}

	err = p.Validate()
	if err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			ErrorValue := make(map[string]string)
			for key, value := range vErrors {
				ErrorValue[strings.Title(key)] = value.Error()
			}
			h.loadEditForm(w, p, category, ErrorValue)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.pc.UpdatePost(r.Context(), &ppb.UpdatePostRequest{
		Post: &ppb.Post{
			ID:          p.ID,
			Title:       p.Title,
			Author:      p.Author,
			Description: p.Description,
			Category:    p.Category,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}

	int_id, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}

	_, err = h.pc.DeletePost(r.Context(), &ppb.DeletePostRequest{
		ID: int64(int_id),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

type SearchedFormData struct {
	SearchResult []Post
	SearchQuery  string
}

func (h *Handler) searchPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sq := r.FormValue("SearchPost")

	posts, err := h.pc.SearchPost(context.Background(), &ppb.SearchPostRequest{SearchPostQuery: sq})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var sResult []Post
	for _, post := range posts.SearchPostResult {
		sResult = append(sResult, Post{
			ID:          post.ID,
			Title:       post.Title,
			Author:      post.Author,
			Description: post.Description,
			Category:    post.Category,
		})
	}
	slt := SearchedFormData{
		SearchResult: sResult,
		SearchQuery:  sq,
	}
	if len(sResult) == 0 {
		err = h.templates.ExecuteTemplate(w, "no-search-result.html", slt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err = h.templates.ExecuteTemplate(w, "search-result-post.html", slt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type FormData struct {
	Post     Post
	Category []Category
	Errors   map[string]string
}

func (h *Handler) loadCreateForm(w http.ResponseWriter, post Post, category []Category, myErrors map[string]string) {
	form := FormData{
		Post:     post,
		Category: category,
		Errors:   myErrors,
	}

	err := h.templates.ExecuteTemplate(w, "create-post.html", form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) loadEditForm(w http.ResponseWriter, post Post, category []Category, myErrors map[string]string) {
	form := FormData{
		Post:     post,
		Category: category,
		Errors:   myErrors,
	}

	err := h.templates.ExecuteTemplate(w, "edit-post.html", form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
