package listing

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/listings-service/internal/platform"
)

type Handlers struct{ S *Store }

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := platform.UserIDFromHeader(r)
	if !ok {
		platform.Error(w, http.StatusUnauthorized, "missing user")
		return
	}

	var p CreateParams
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		platform.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	if p.Title == "" || p.Price <= 0 || p.Category == "" {
		platform.Error(w, http.StatusBadRequest, "title, price, category required")
		return
	}

	l, err := h.S.Create(r.Context(), userID, p)
	if err != nil {
		platform.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	platform.JSON(w, http.StatusCreated, l)
}

func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	l, err := h.S.Get(r.Context(), id)
	if err != nil {
		platform.Error(w, http.StatusNotFound, "not found")
		return
	}
	platform.JSON(w, http.StatusOK, l)
}

func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := ListFilters{
		Limit:  parseInt(q.Get("limit"), 20),
		Offset: parseInt(q.Get("offset"), 0),
		Sort:   q.Get("sort"),
	}
	if s := q.Get("q"); s != "" {
		f.Q = &s
	}
	if s := q.Get("category"); s != "" {
		c := Category(s)
		f.Category = &c
	}
	if s := q.Get("status"); s != "" {
		st := Status(s)
		f.Status = &st
	}
	if s := q.Get("min_price"); s != "" {
		v := parseInt64(s, 0)
		f.MinPrice = &v
	}
	if s := q.Get("max_price"); s != "" {
		v := parseInt64(s, 0)
		f.MaxPrice = &v
	}

	items, err := h.S.List(r.Context(), f)
	if err != nil {
		platform.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	platform.JSON(w, http.StatusOK, map[string]any{"items": items, "count": len(items)})
}

func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := platform.UserIDFromHeader(r)
	if !ok {
		platform.Error(w, http.StatusUnauthorized, "missing user")
		return
	}

	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	var p UpdateParams
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		platform.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	l, err := h.S.Update(r.Context(), id, userID, p)
	if err != nil {
		platform.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	platform.JSON(w, http.StatusOK, l)
}

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if r.URL.Query().Get("hard") == "true" {
		if err := h.S.Delete(r.Context(), id); err != nil {
			platform.Error(w, 500, err.Error())
			return
		}
	} else {
		if err := h.S.Archive(r.Context(), id); err != nil {
			platform.Error(w, 500, err.Error())
			return
		}
	}
	platform.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func parseInt(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
func parseInt64(s string, def int64) int64 {
	if s == "" {
		return def
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return def
	}
	return v
}
