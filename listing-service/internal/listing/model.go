package listing

import "time"

type Category string
type Status string

const (
	CatTextbook     Category = "TEXTBOOK"
	CatGadget       Category = "GADGET"
	CatEssential    Category = "ESSENTIAL"
	CatNonEssential Category = "NON-ESSENTIAL"
	CatOther        Category = "OTHER"

	StAvailable Status = "AVAILABLE"
	StPending   Status = "PENDING"
	StSold      Status = "SOLD"
	StArchived  Status = "ARCHIVED"
)

type Listing struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Price       int64     `json:"price"`
	Category    Category  `json:"category"`
	UserID      int64     `json:"user_id"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateParams struct {
	Title       string   `json:"title"`
	Description *string  `json:"description,omitempty"`
	Price       int64    `json:"price"`
	Category    Category `json:"category"`
}

type UpdateParams struct {
	Title       *string   `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`
	Price       *int64    `json:"price,omitempty"`
	Category    *Category `json:"category,omitempty"`
	Status      *Status   `json:"status,omitempty"`
}
