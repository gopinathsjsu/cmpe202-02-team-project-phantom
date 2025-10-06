package model

import "time"

// Listing struct remains the same
type Listing struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Category    string    `json:"category"`
	UserID      int       `json:"userId"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

// --- CHANGE: Keywords is now a slice of strings ---
type SearchParams struct {
	Category string   `json:"category"`
	Keywords []string `json:"keywords"` // Changed from string to []string
	MinPrice float64  `json:"min_price"`
	MaxPrice float64  `json:"max_price"`
}

// Gemini API structs remain the same
type GeminiRequest struct {
	Contents         []Content        `json:"contents"`
	GenerationConfig GenerationConfig `json:"generation_config"`
}

// ... (rest of the file is the same)
type Content struct {
	Parts []Part `json:"parts"`
}
type Part struct {
	Text string `json:"text"`
}
type GenerationConfig struct {
	ResponseMimeType string `json:"response_mime_type"`
}
type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}
type Candidate struct {
	Content Content `json:"content"`
}
