package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Nikhil1169/search-agent/internal/model"
)

// The interface is now more generic.
type AIClient interface {
	GetSearchParams(ctx context.Context, userQuery string) (*model.SearchParams, error)
}

type DBStore interface {
	FindListingsByParams(ctx context.Context, params *model.SearchParams) ([]model.Listing, error)
}

type SearchHandler struct {
	AI    AIClient // Uses the generic interface
	Store DBStore
}

// ... the rest of the file (searchRequest struct and HandleSearch function) is exactly the same ...
type searchRequest struct {
	Query string `json:"query"`
}

func (h *SearchHandler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	var req searchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Query == "" {
		http.Error(w, "Query cannot be empty", http.StatusBadRequest)
		return
	}

	log.Printf("Received search query: %s", req.Query)

	searchParams, err := h.AI.GetSearchParams(r.Context(), req.Query)
	if err != nil {
		log.Printf("ERROR getting search params from AI: %v", err)
		http.Error(w, "Failed to understand query", http.StatusInternalServerError)
		return
	}

	listings, err := h.Store.FindListingsByParams(r.Context(), searchParams)
	if err != nil {
		log.Printf("ERROR finding listings in database: %v", err)
		http.Error(w, "Failed to retrieve listings", http.StatusInternalServerError)
		return
	}

	log.Printf("Found %d listings for query '%s'", len(listings), req.Query)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(listings)
}
