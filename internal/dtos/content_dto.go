package dtos

// Input para criação/update via Postman
type ContentInput struct {
	Type             string `json:"type" binding:"required,oneof=project post"`
	Title            string `json:"title" binding:"required,min=3"`
	ShortDescription string `json:"short_description"`
	Body             string `json:"body" binding:"required"`
	DemoURL          string `json:"demo_url" validate:"url"`
	RepoURL          string `json:"repo_url" validate:"url"`
}

// Output limpo para o Next.js
type ContentResponse struct {
	ID               uint   `json:"id"`
	Title            string `json:"title"`
	Slug             string `json:"slug"`
	ShortDescription string `json:"short_description"`
	Body             string `json:"body"`
	Type             string `json:"type"`
	DemoURL          string `json:"demo_url,omitempty"`
}
