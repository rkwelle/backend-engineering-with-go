package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rkwelle/social-app/internal/store"
)

type postKey string

const postCtxKey postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		// writeJSONError(w, http.StatusBadRequest, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	// validate for required content field
	// if payload.Content == "" {
	// 	app.badRequestResponse(w, r, fmt.Errorf("content is required"))
	// 	return
	// }
	if err := validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// TODO: Change after authentication
		UserID: 1,
	}

	ctx := r.Context()
	if err := app.store.Posts.Create(ctx, post); err != nil {
		// writeJSONError(w, http.StatusInternalServerError, err.Error())
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		// writeJSONError(w, http.StatusInternalServerError, err.Error())
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	// Now done on the middleware context
	//
	// idParam := chi.URLParam(r, "postID")
	// id, err := strconv.ParseInt(idParam, 10, 64)
	// if err != nil {
	// 	// writeJSONError(w, http.StatusInternalServerError, err.Error())
	// 	app.internalServerError(w, r, err)
	// 	return
	// }
	// ctx := r.Context()

	// post, err := app.store.Posts.GetByID(ctx, id)

	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, store.ErrNotFound):
	// 		// writeJSONError(w, http.StatusNotFound, err.Error())
	// 		app.notFoundResponse(w, r, err)
	// 	default:
	// 		// writeJSONError(w, http.StatusInternalServerError, err.Error())
	// 		app.internalServerError(w, r, err)
	// 	}
	// 	return
	// }

	post := getPostFromCtx(r)

	// comments, err := app.store.Comments.GetByPostID(ctx, id)
	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		// writeJSONError(w, http.StatusInternalServerError, err.Error())
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		// writeJSONError(w, http.StatusInternalServerError, err.Error())
		app.internalServerError(w, r, err)
		return
	}
	ctx := r.Context()

	err = app.store.Posts.Delete(ctx, id)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			// writeJSONError(w, http.StatusNotFound, err.Error())
			app.notFoundResponse(w, r, err)
		default:
			// writeJSONError(w, http.StatusInternalServerError, err.Error())
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		// writeJSONError(w, http.StatusInternalServerError, err.Error())
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			// writeJSONError(w, http.StatusInternalServerError, err.Error())
			app.internalServerError(w, r, err)
			return
		}
		ctx := r.Context()

		post, err := app.store.Posts.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				// writeJSONError(w, http.StatusNotFound, err.Error())
				app.notFoundResponse(w, r, err)
			default:
				// writeJSONError(w, http.StatusInternalServerError, err.Error())
				app.internalServerError(w, r, err)
			}
			return
		}

		// instead of using string as key which is error prone
		// ctx = context.WithValue(ctx, "post", post)
		ctx = context.WithValue(ctx, postCtxKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	// use the postCtxKey instead of string
	// post, _ := r.Context().Value("post").(*store.Post)
	post, _ := r.Context().Value(postCtxKey).(*store.Post)
	return post
}
