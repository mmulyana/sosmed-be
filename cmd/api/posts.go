package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mmulyana/sosmed-be/internal/store"
)

type CreatePostPayload struct {
	Title   string `json:"title" validate:"required,max=100"`
	Content string `json:"content" validate:"required"`
	UserId  int    `json:"userId" validate:"required"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	err := ReadJSON(w, r, &payload)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = Validate.Struct(payload)
	if err != nil {
		validationErrors := formatValidationError(err)
		WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errors": validationErrors,
		})
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserId:  int64(payload.UserId),
	}

	ctx := r.Context()

	err = app.store.Posts.Create(ctx, post)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = WriteJSON(w, http.StatusCreated, post)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := r.Context()
	post, err := app.store.Posts.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			WriteJSONError(w, http.StatusNotFound, err.Error())
			return
		default:
			WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	err = WriteJSON(w, http.StatusCreated, post)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	posts, err := app.store.Posts.GetAll(ctx)
	if err != nil {
		WriteJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	err = WriteJSON(w, http.StatusOK, posts)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
