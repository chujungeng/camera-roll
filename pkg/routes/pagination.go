package routes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

const (
	PaginationDefaultOffset uint64 = 0
	PaginationDefaultLimit  uint64 = 12
)

const (
	ParamNumberBase = 10
	ParamNumberBit  = 64
)

const (
	ParamOffset = "offset"
	ParamLimit  = "limit"
	ParamPageID = "page"
)

// Pagination middleware is used to extract the next page id from the url query
func Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		PageID := r.URL.Query().Get(ParamPageID)
		intPageID := 0
		var err error
		if PageID != "" {
			intPageID, err = strconv.Atoi(PageID)
			if err != nil {
				_ = render.Render(w, r, ErrInvalidRequest(fmt.Errorf("couldn't read %s: %w", ParamPageID, err)))
				return
			}
		}
		ctx := context.WithValue(r.Context(), pageIDKey, intPageID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
