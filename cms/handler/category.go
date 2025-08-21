package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	cpb "blog/gunk/v1/category"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

type Category struct {
	CategoryID   int64  `db:"category_id"`
	CategoryName string `db:"category_name"`
}

func (category *Category) Validate() error {
	return validation.ValidateStruct(category,
		validation.Field(&category.CategoryName, validation.Required.Error("Category name field can not be empty."), validation.Length(3, 50).Error("Category name field should have atleast 3 characters and atmost 50 characters")),
	)
}

func (h *Handler) listCategory(w http.ResponseWriter, r *http.Request) {
	categories, err := h.cc.GetAllCategories(r.Context(), &cpb.GetAllCategoriesRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.templates.ExecuteTemplate(w, "category-list.html", categories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) createCategory(w http.ResponseWriter, r *http.Request) {
	ErrorValue := map[string]string{}
	category := Category{}
	h.loadCreateCategoryForm(w, category, ErrorValue)
}

func (h *Handler) storeCategory(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var category Category
	err = h.decoder.Decode(&category, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = category.Validate()
	if err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			ErrorValue := make(map[string]string)
			for key, value := range vErrors {
				ErrorValue[strings.Title(key)] = value.Error()
			}
			h.loadCreateCategoryForm(w, category, ErrorValue)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.cc.CreateCategory(r.Context(), &cpb.CreateCategoryRequest{
		Category: &cpb.Category{
			CategoryName: category.CategoryName,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/categories", http.StatusTemporaryRedirect)
}

func (h *Handler) editCategory(w http.ResponseWriter, r *http.Request) {
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

	category, err := h.cc.GetCategory(r.Context(), &cpb.GetCategoryRequest{
		CategoryID: int64(int_id),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if category.Category.CategoryID == 0 {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}

	var c Category
	c.CategoryID = int64(int_id)
	c.CategoryName = category.Category.CategoryName

	h.loadEditCategoryForm(w, c, map[string]string{})
}

func (h *Handler) updateCategory(w http.ResponseWriter, r *http.Request) {
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

	var c Category
	c.CategoryID = int64(int_id)

	if c.CategoryID == 0 {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.decoder.Decode(&c, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = c.Validate()
	if err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			ErrorValue := make(map[string]string)
			for key, value := range vErrors {
				ErrorValue[strings.Title(key)] = value.Error()
			}
			h.loadEditCategoryForm(w, c, ErrorValue)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.cc.UpdateCategory(r.Context(), &cpb.UpdateCategoryRequest{
		Category: &cpb.Category{
			CategoryID:   c.CategoryID,
			CategoryName: c.CategoryName,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/categories", http.StatusTemporaryRedirect)
}

func (h *Handler) deleteCategory(w http.ResponseWriter, r *http.Request) {
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

	_, err = h.cc.DeleteCategory(r.Context(), &cpb.DeleteCategoryRequest{
		CategoryID: int64(int_id),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/categories", http.StatusTemporaryRedirect)
}

type SearchedCategoryFormData struct {
	SearchResult []Category
	SearchQuery  string
}

func (h *Handler) searchCategory(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sq := r.FormValue("SearchCategory")

	categories, err := h.cc.SearchCategory(context.Background(), &cpb.SearchCategoryRequest{SearchCategoryQuery: sq})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var sResult []Category
	for _, category := range categories.SearchCategoryResult {
		sResult = append(sResult, Category{
			CategoryID:   category.CategoryID,
			CategoryName: category.CategoryName,
		})
	}
	slt := SearchedCategoryFormData{
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
		err = h.templates.ExecuteTemplate(w, "search-result-category.html", slt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type CategoryFormData struct {
	Category Category
	Errors   map[string]string
}

func (h *Handler) loadCreateCategoryForm(w http.ResponseWriter, category Category, myErrors map[string]string) {
	form := CategoryFormData{
		Category: category,
		Errors:   myErrors,
	}

	err := h.templates.ExecuteTemplate(w, "create-category.html", form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) loadEditCategoryForm(w http.ResponseWriter, category Category, myErrors map[string]string) {
	form := CategoryFormData{
		Category: category,
		Errors:   myErrors,
	}

	err := h.templates.ExecuteTemplate(w, "edit-category.html", form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
