package main

import (
	"fmt"
	"net/http"

	"github.com/mmulyana/sosmed-be/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required"`
	PostId  int    `json:"postId" validate:"required"`
	UserId  int    `json:"userId" validate:"required"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentPayload
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

	fmt.Println("userId -> ", payload.UserId)

	comment := &store.Comment{
		Content: payload.Content,
		UserId:  int64(payload.UserId),
		PostId:  int64(payload.PostId),
	}

	fmt.Println("comment -> ", comment)

	ctx := r.Context()

	err = app.store.Comments.Create(ctx, comment)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = WriteJSON(w, http.StatusCreated, comment)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
